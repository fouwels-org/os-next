package stages

import (
	"fmt"
	"init/config"
	"init/shell"
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

	commands := []shell.Command{}

	for _, vloop := range c.Primary.Filesystem.Devices {

		v := vloop

		commands = append(commands, shell.Command{Executable: shell.Mkdir, Arguments: []string{"-p", v.MountPoint}})
		commands = append(commands, shell.Command{Executable: shell.Mount, Arguments: []string{"-t", v.FileSystem, v.ID, v.MountPoint}})
	}

	err := shell.Executor.Execute(commands)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	return nil
}
