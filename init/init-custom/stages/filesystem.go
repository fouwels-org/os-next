package stages

import (
	"fmt"
	"init-custom/config"
)

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

	commands := []command{}

	for _, v := range c.Primary.Filesystem.Devices {
		commands = append(commands, command{command: "/bin/mkdir", arguments: []string{"-p", v.MountPoint}})
		commands = append(commands, command{command: "/bin/mount", arguments: []string{"-t", v.FileSystem, v.ID, v.MountPoint}})
	}

	err := execute(commands)
	if err != nil {
		return fmt.Errorf("Error mounting: %w", err)
	}

	return nil
}
