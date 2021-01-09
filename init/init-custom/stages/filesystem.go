package stages

import (
	"bufio"
	"fmt"
	"init-custom/config"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/anatol/luks.go"
)

const (
	bootPartition   = "BOOT"
	dataPartition   = "DATA"
	configPartition = "CONFIG"

	bootFile = "/tmp/vfat.txt"

	defaultDevConfig = "/dev/sda2"
	defaultDevData   = "/dev/sda3"
)

// DeviceType ..
type DeviceType struct {
	DEV      string
	LABEL    string
	UUID     string
	FSTYPE   string
	PARTUUID string
}

//Filesystem implements IStage
type Filesystem struct {
	finals []string
}

//String ..
func (n *Filesystem) String() string {
	return "Filesystem"
}

//Finalise ..
func (n *Filesystem) Finalise() []string {
	return n.finals
}

//Run ..
func (n *Filesystem) Run(c config.Config) error {

	// Find LUKS volumes and open them
	err := openLUKSvolumes()
	if err != nil {
		logf("Opening LUKS volumes failed: %v ", err)
	}

	// Find the Labelled device points from the linux blkid
	configDev, dataDev, err := findLabelledDevices()
	if err != nil {
		logf("finding lablled mount points failed: %v ", err)
	}

	logf("Config : %s  --- Data : %s", configDev, dataDev)

	commands := []command{}

	for _, v := range c.Primary.Filesystem.Devices {
		commands = append(commands, command{command: "/bin/mkdir", arguments: []string{"-p", v.MountPoint}})

		// set the mount point based on the tag in the primary.json config file
		dev := v.ID
		if strings.ToUpper(v.LABEL) == dataPartition {
			dev = dataDev
		} else if strings.ToUpper(v.LABEL) == configPartition {
			dev = configDev
		}
		commands = append(commands, command{command: "/bin/mount", arguments: []string{"-t", v.FileSystem, dev, v.MountPoint}})
	}

	err = execute(commands)
	if err != nil {
		return fmt.Errorf("Error mounting: %w", err)
	}

	return nil
}

// returns the available blockId as a deviceType from the underlying OS
func formatBlkid() ([]DeviceType, error) {
	var re1 = regexp.MustCompile(`"(.*?)"`)

	devices := []DeviceType{}

	cmd := exec.Command("blkid")
	stdout, err := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		logf("Command failed: %v ", err)
		return nil, fmt.Errorf("Command failed: %v", err)
	}

	r := bufio.NewReader(stdout)

	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}

		data := DeviceType{}

		strLine := string(line)
		partLine := strings.Fields(strLine)

		devStr := strings.Split(strLine, ":")
		if len(devStr) > 0 {
			data.DEV = devStr[0]
		}

		for i := 0; i < len(partLine); i++ {
			col := partLine[i]
			if strings.HasPrefix(col, "UUID") {
				ms := re1.FindString(col)
				data.UUID = strings.Trim(ms, "\"")
			} else if strings.HasPrefix(col, "LABEL") {
				ms := re1.FindString(col)
				data.LABEL = strings.Trim(ms, "\"")
			} else if strings.HasPrefix(col, "TYPE") {
				ms := re1.FindString(col)
				data.FSTYPE = strings.Trim(ms, "\"")
			} else if strings.HasPrefix(col, "PARTUUID") {
				ms := re1.FindString(col)
				data.PARTUUID = strings.Trim(ms, "\"")
			}
		}
		devices = append(devices, data)
	}
	return devices, nil

}

// Decrypts the volume or returns an error if decryption fails
func openLUKSvolumes() error {

	devices, err := formatBlkid()
	if err != nil {
		logf("Faild to format BLKID: %v ", err)
		return err
	}

	luksFound := 0

	for index, deviceBlk := range devices {
		blkFSType := strings.ToUpper(deviceBlk.FSTYPE)
		if blkFSType == "CRYPTO_LUKS" {
			luksFound++
			// LOAD KEY FROM TPM
			key := "pass0"
			dev, err := luks.Open(deviceBlk.DEV)
			if err != nil {
				logf("Faild to format BLKID: %v ", err)
				return err
			}
			defer dev.Close()

			// equivalent of `cryptsetup open /dev/sda1 volumename`
			err = dev.Unlock( /* slot */ 0, []byte(key), ("luks" + strconv.Itoa(index)))

			if err == luks.ErrPassphraseDoesNotMatch {
				log.Printf("The password is incorrect")
				return err
			} else if err != nil {
				log.Print(err)
				return err
			} else {
				logf("LUKS opened %s", deviceBlk.DEV)
				logf("mapped to /dev/mapper/%s", "luks"+strconv.Itoa(index))

				// at this point system should have a file `/dev/mapper/volumename`.
			}
		}
	}

	if luksFound == 0 {
		logf("No encrypted LUKS volumes found")
	} else {
		logf("%d encrypted LUKS volumes found", luksFound)
	}

	return nil
}

// returns the config and data partitions based on the formatted partitions.
// string: Config partition e.g /dev/sda2
// string: Data partition e.g /dev/sda3
// error: nil if partitons are found otherwise !nil
func findLabelledDevices() (string, string, error) {
	devices, err := formatBlkid()
	if err != nil {
		logf("Faild to format BLKID: %v ", err)
		return "", "", err
	}
	// assign the default return values
	config, data := defaultDevConfig, defaultDevData
	for _, dev := range devices {
		devLabel := strings.ToUpper(dev.LABEL)
		if devLabel == configPartition {
			config = dev.DEV
			logf("Config label found %s", dev.DEV)
		} else if devLabel == dataPartition {
			data = dev.DEV
			logf("Data label found %s", dev.DEV)
		} else if devLabel == bootPartition {
			err := ioutil.WriteFile(bootFile, []byte(dev.DEV), 0600)
			if err != nil {
				logf("Could not write the file %s", dev.DEV)
			}
		}
	}

	return config, data, nil
}
