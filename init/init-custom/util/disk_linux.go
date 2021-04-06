// +build linux

package util

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/anatol/luks.go"
)

//DiskUtil ..
type DiskUtil struct {
	MapperIndex int
}

//Decrypt decrypts and maps the device, if found to be encrypted
//device: /dev/id, dmlabel: LABEL of matching /dev/dm- device, key: LUKS key, returns: [mapped /dev/id, error]
func (d *DiskUtil) Decrypt(device string, dmlabel string, key string) (string, error) {

	devices, err := d.GetBLKDevices()
	if err != nil {
		return "", fmt.Errorf("failed to get BLKID (labelled) devices: %w", err)
	}

	for _, dev := range devices {

		if dev.Device == device {

			if strings.ToUpper(dev.FsType) != "CRYPTO_LUKS" {
				log.Printf("%v is not of type CRYPTO_LUKS, skipped", dev.Device)
				return device, nil
			}

			lkd, err := luks.Open(device)
			if err != nil {
				return "", fmt.Errorf("failed to LUKS open %v: %w", device, err)
			}
			defer lkd.Close()

			// equivalent of `cryptsetup open /dev/sda1 volumename`
			d.MapperIndex++
			err = lkd.Unlock(0, []byte(key), fmt.Sprintf("luks%v", d.MapperIndex))
			mapper := fmt.Sprintf("/dev/mapper/luks%v", d.MapperIndex)

			if err == luks.ErrPassphraseDoesNotMatch {
				return "", fmt.Errorf("failed to unlock %v, password is incorrect", device)
			}
			if err != nil {
				return "", fmt.Errorf("unexpected error unlocking %v: %w", device, err)
			}

			log.Printf("%v mapped to %v", device, mapper)
			return mapper, nil
		}
	}

	return "", fmt.Errorf("device %v not found in BLKID", device)
}

// GetBLKDevices returns the decoded block devices
func (d *DiskUtil) GetBLKDevices() ([]BLKDevice, error) {

	var re1 = regexp.MustCompile(`"(.*?)"`)

	devices := []BLKDevice{}

	cmd := exec.Command("/sbin/blkid")
	stdout, err := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		return []BLKDevice{}, fmt.Errorf("Failed to get blkid: %w", err)
	}

	r := bufio.NewReader(stdout)
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}

		data := BLKDevice{}

		strLine := string(line)
		partLine := strings.Fields(strLine)

		devStr := strings.Split(strLine, ":")
		if len(devStr) > 0 {
			data.Device = devStr[0]
		}

		for i := 0; i < len(partLine); i++ {
			col := partLine[i]
			if strings.HasPrefix(col, "UUID") {

				ms := re1.FindString(col)
				data.UUID = strings.Trim(ms, "\"")

			} else if strings.HasPrefix(col, "LABEL") {

				ms := re1.FindString(col)
				data.Label = strings.Trim(ms, "\"")

			} else if strings.HasPrefix(col, "TYPE") {

				ms := re1.FindString(col)
				data.FsType = strings.Trim(ms, "\"")

			} else if strings.HasPrefix(col, "PARTUUID") {

				ms := re1.FindString(col)
				data.PartUUID = strings.Trim(ms, "\"")
			}
		}
		devices = append(devices, data)
	}

	return devices, nil
}
