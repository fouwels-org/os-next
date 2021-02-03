package stages

import (
	"fmt"
	"init-custom/config"
	"init-custom/util"
	"strings"
)

var _partitions = util.Partitions{
	BootPartition:   "BOOT",
	DataPartition:   "DATA",
	ConfigPartition: "CONFIG",

	BootFile: "/tmp/vfat.txt",

	DefaultDevConfig: "/dev/sda2",
	DefaultDevData:   "/dev/sda3",
}

//Filesystem implements IStage
type Filesystem struct {
	finals []string
}

//String ..
func (n *Filesystem) String() string {
	return "filesystem"
}

//Finalise ..
func (n *Filesystem) Finalise() []string {
	return n.finals
}

//Run ..
func (n *Filesystem) Run(c config.Config) error {

	// Find LUKS volumes and open them
	err := util.Disk.OpenLUKSvolumes()
	if err != nil {
		return fmt.Errorf("Opening LUKS volumes failed: %v ", err)
	}

	// Find the Labelled device points from the linux blkid
	configDev, dataDev, err := util.Disk.FindLabelledDevices(_partitions)
	if err != nil {
		return fmt.Errorf("Finding lablled mount points failed: %v ", err)
	}

	commands := []util.Command{}

	for _, v := range c.Primary.Filesystem.Devices {
		commands = append(commands, util.Command{Target: "/bin/mkdir", Arguments: []string{"-p", v.MountPoint}})

		// set the mount point based on the tag in the primary.json config file
		dev := v.ID
		if strings.ToUpper(v.LABEL) == _partitions.DataPartition {
			dev = dataDev
		} else if strings.ToUpper(v.LABEL) == _partitions.ConfigPartition {
			dev = configDev
		}
		commands = append(commands, util.Command{Target: "/bin/mount", Arguments: []string{"-t", v.FileSystem, dev, v.MountPoint}})
	}

	err = util.Shell.Execute(commands)
	if err != nil {
		return fmt.Errorf("Error mounting: %w", err)
	}

	n.finals = append(n.finals, fmt.Sprintf("initialized as config %v data %v", configDev, dataDev))

	return nil
}
