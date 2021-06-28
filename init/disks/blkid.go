// SPDX-FileCopyrightText: 2020 Lagoni Engineering
// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
//
// SPDX-License-Identifier: Apache-2.0

package disks

import (
	"fmt"
	"init/shell"
	"log"
	"strings"
)

type Blkid struct {
	Device   string
	UUID     string
	TYPE     string
	LABEL    string
	PARTUUID string
}

//GetBlkids get BLKID of target device, or all devices is passed blank
func GetBlkid(target string) ([]Blkid, error) {

	blkidList := []Blkid{}

	args := []string{}
	if target != "" {
		args = append(args, target)
	}

	command := shell.Command{
		Executable: shell.Blkid,
		Arguments:  args,
		Env:        []string{},
	}

	output, err := shell.Executor.ExecuteOne(command)
	if err != nil {
		return blkidList, fmt.Errorf("failed to call blkid executable: %w", err)
	}

	// parse each blkid line
	for _, l := range strings.Split(output, "\n") {

		c := strings.Split(l, " ")
		if len(c) < 1 {
			log.Printf("length of %v < 1, skipped", c)
			continue
		}

		b := Blkid{}
		b.Device = strings.Trim(c[0], ":") // add device identifier

		// parse remaining fields
		for _, r := range c[1:] {

			k := strings.Split(r, "=")
			if len(k) != 2 {
				log.Printf("length of split %v != 2, skipped", k)
				continue
			}

			key := k[0]
			value := k[1]

			switch key {
			case "UUID":
				b.UUID = strings.Trim(value, "\"")
			case "TYPE":
				b.TYPE = strings.Trim(value, "\"")
			case "LABEL":
				b.LABEL = strings.Trim(value, "\"")
			case "PARTUUID":
				b.PARTUUID = strings.Trim(value, "\"")
			default:
				log.Printf("unrecognised BLKID key %v, ignored", key)
			}
		}
		blkidList = append(blkidList, b)
	}

	return blkidList, nil
}
