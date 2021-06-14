package stages

import (
	"fmt"
	"init/config"
	"init/disks"
	"init/shell"
	"log"
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

	blkmap := blklist.LabelMap()

	for _, v := range c.Primary.Filesystem.Devices {

		// Check if device exists for specified label
		b, ok := blkmap[v.Label]

		// If does not exist, return with err
		if !ok {
			return fmt.Errorf("BLKID %v missing, cannot mount filesystem\n: %+v", v.Label, blklist)
		}

		// Mount it
		commands := []shell.Command{
			{Executable: shell.Mkdir, Arguments: []string{"-p", v.MountPoint}},
			{Executable: shell.Mount, Arguments: []string{"-t", v.FileSystem, b.Device, v.MountPoint}},
		}

		// If cannot mount, return with err
		err := shell.Executor.Execute(commands)
		if err != nil {
			log.Printf("failed to mount: %v", err)
			//return fmt.Errorf("failed to mount: %w", err)
		}
	}

	return nil
}
