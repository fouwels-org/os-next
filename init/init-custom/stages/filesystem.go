package stages

import (
	"fmt"
	"init-custom/config"
	"init-custom/util"
)

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

	commands := []util.Command{}

	for _, vloop := range c.Primary.Filesystem.Devices {

		v := vloop

		commands = append(commands, util.Command{Target: "/bin/mkdir", Arguments: []string{"-p", v.MountPoint}})
		commands = append(commands, util.Command{Target: "/bin/mount", Arguments: []string{"-t", v.FileSystem, v.ID, v.MountPoint}})
	}

	err := util.Shell.Execute(commands)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	return nil
}
