package stages

import (
	"fmt"
	"uinit-custom/config"
)

//FS implementes IStage
type FS struct {
	finals []string
}

//String ..
func (f FS) String() string {
	return "FS"
}

//Finalise ..
func (f FS) Finalise() []string {
	return f.finals
}

//Run ..
func (f FS) Run(c config.Config, s config.Secrets) error {
	commands := []string{}

	//Check and set up partitions
	_, err := executeOne(fmt.Sprintf("ls %v", c.FileSystem.Partitions.Data), "")
	if err != nil { //if does not exist
		commands = append(commands, fmt.Sprintf("parted %v --script mkpart primary 0% 100%", c.FileSystem.Device))
	}

	err = execute(commands)
	if err != nil {
		return err
	}

	return nil
}
