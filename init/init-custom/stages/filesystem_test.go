package stages_test

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

const (
	dataPartition   = "DATA"
	configPartition = "CONFIG"

	defaultDevConfig = "/dev/sda2"
	defaultDevData   = "/dev/sda3"
)

// DeviceType ..
type DeviceType struct {
	dev      string
	label    string
	uuid     string
	fstype   string
	partuuid string
}

func TestPartition(t *testing.T) {
	a, b := findMountPoints()
	t.Logf("Config : %s       Data : %s", a, b)
}

func findMountPoints() (string, string) {
	var re1 = regexp.MustCompile(`"(.*?)"`)

	cmd := exec.Command("blkid")
	stdout, err := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Command failed: %v ", err)
		return "", ""
	}

	r := bufio.NewReader(stdout)

	devices := []DeviceType{}

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
			data.dev = devStr[0]
		}

		for i := 0; i < len(partLine); i++ {
			col := partLine[i]
			if strings.HasPrefix(col, "UUID") {
				ms := re1.FindString(col)
				data.uuid = strings.Trim(ms, "\"")
			} else if strings.HasPrefix(col, "LABEL") {
				ms := re1.FindString(col)
				data.label = strings.Trim(ms, "\"")
			} else if strings.HasPrefix(col, "TYPE") {
				ms := re1.FindString(col)
				data.fstype = strings.Trim(ms, "\"")
			} else if strings.HasPrefix(col, "PARTUUID") {
				ms := re1.FindString(col)
				data.partuuid = strings.Trim(ms, "\"")
			}
		}
		devices = append(devices, data)
	}
	config, data := defaultDevConfig, defaultDevData
	for _, dev := range devices {
		upper := strings.ToUpper(dev.label)
		if upper == configPartition {
			config = dev.dev
		} else if upper == dataPartition {
			data = dev.dev
		}
	}
	return config, data
}
