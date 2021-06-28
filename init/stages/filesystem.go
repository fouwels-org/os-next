// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 K. Fouwels <k@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package stages

import (
	"errors"
	"fmt"
	"init/config"
	"init/disks"
	"init/shell"
	"io/ioutil"
	"log"
	"os"
)

//Filesystem implements IStage
type Filesystem struct {
	finals []string
}

//String ..
func (n *Filesystem) String() string {
	return "filesystem"
}

//Policy ..
func (n *Filesystem) Policy() Policy {
	return PolicyHard
}

//Finalise ..
func (n *Filesystem) Finalise() []string {
	return n.finals
}

//Run ..
func (n *Filesystem) Run(c config.Config) error {

	blklist, err := disks.GetBlkid("")
	if err != nil {
		return fmt.Errorf("failed to get blkid: %w", err)
	}

	for _, v := range c.Primary.Filesystem.Devices {

		if v.Label == "" {
			return fmt.Errorf("label missing for %+v, aborting", v)
		}

		blk := disks.Blkid{}
		for _, j := range blklist {
			if j.LABEL == v.Label {
				blk = j
			}
		}

		if blk.Device == "" || blk.LABEL == "" {
			return fmt.Errorf("device %v for label %v not found in blkid, cannot mount filesystem", v.Label, blk.Device)
		}

		if v.Label != blk.LABEL {
			return fmt.Errorf("device %v label %v does not match expected %v, will not mount filesystem", blk.Device, v.Label, blk.LABEL)
		}

		// Mount it
		commands := []shell.Command{
			{Executable: shell.Mkdir, Arguments: []string{"-p", v.MountPoint}},
			{Executable: shell.Mount, Arguments: []string{"-t", v.FileSystem, blk.Device, v.MountPoint}},
		}

		// If cannot mount, return with err
		err := shell.Executor.Execute(commands)
		if err != nil {
			log.Printf("failed to mount: %v", err)
		}
	}

	// Deploy default secondary config if one does not exist
	_, err = os.Stat("/var/config/secondary.yml")
	if errors.Is(err, os.ErrNotExist) {

		secondary, err := ioutil.ReadFile("/config/secondary.yml")
		if err != nil {
			return fmt.Errorf("failed to copy secondary config: %w", err)
		}
		err = ioutil.WriteFile("/var/config/secondary.yml", secondary, 0644)
		if err != nil {
			return fmt.Errorf("failed to copy secondary config: %w", err)
		}
		if err != nil {
			return fmt.Errorf("failed to symlink default secondary config: %w", err)
		}

		n.finals = append(n.finals, "default secondary configuration was written to /var/config/secondary.yml")

	} else if err != nil {
		return fmt.Errorf("error checking if secondary config file exists: %w", err)
	}

	return nil
}
