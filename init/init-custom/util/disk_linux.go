// +build linux

package util

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/anatol/luks.go"
)

//DiskUtil ..
type DiskUtil struct {
}

//OpenLUKSvolumes Decrypts the volume or returns an error if decryption fails
func (d *DiskUtil) OpenLUKSvolumes() error {

	devices, err := d.FormatBlkid()
	if err != nil {
		log.Printf("Faild to format BLKID: %v ", err)
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
				log.Printf("Faild to format BLKID: %v ", err)
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
				log.Printf("LUKS opened %s", deviceBlk.DEV)
				log.Printf("mapped to /dev/mapper/%s", "luks"+strconv.Itoa(index))

				// at this point system should have a file `/dev/mapper/volumename`.
			}
		}
	}

	if luksFound == 0 {
		log.Printf("No encrypted LUKS volumes found")
	} else {
		log.Printf("%d encrypted LUKS volumes found", luksFound)
	}

	return nil
}

//FormatBlkid returns the available blockId as a deviceType from the underlying OS
func (d *DiskUtil) FormatBlkid() ([]DeviceType, error) {
	var re1 = regexp.MustCompile(`"(.*?)"`)

	devices := []DeviceType{}

	cmd := exec.Command("/sbin/blkid")
	stdout, err := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		log.Printf("Command failed: %v ", err)
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

// FindLabelledDevices returns the config and data partitions based on the formatted partitions.
// string: Config partition e.g /dev/sda2
// string: Data partition e.g /dev/sda3
// error: nil if partitons are found otherwise !nil
func (d *DiskUtil) FindLabelledDevices(partitions Partitions) (string, string, error) {
	devices, err := d.FormatBlkid()
	if err != nil {
		log.Printf("Faild to format BLKID: %v ", err)
		return "", "", err
	}
	// assign the default return values
	config, data := partitions.DefaultDevConfig, partitions.DefaultDevData
	for _, dev := range devices {
		devLabel := strings.ToUpper(dev.LABEL)
		if devLabel == partitions.ConfigPartition {
			config = dev.DEV
			log.Printf("Config label found %s", dev.DEV)
		} else if devLabel == partitions.DataPartition {
			data = dev.DEV
			log.Printf("Data label found %s", dev.DEV)
		} else if devLabel == partitions.BootPartition {
			err := ioutil.WriteFile(partitions.BootFile, []byte(dev.DEV), 0600)
			if err != nil {
				log.Printf("Could not write the file %s", dev.DEV)
			}
		}
	}

	return config, data, nil
}
