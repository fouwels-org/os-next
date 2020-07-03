package stages

import (
	"fmt"
	"uinit-custom/config"
)

//Modules implementes IStage
type Modules struct {
}

//String ..
func (m Modules) String() string {
	return "Modules"
}

//Run ..
func (m Modules) Run(c config.Config, s config.Secrets) error {

	commands := []string{}

	for _, v := range c.Modules {
		commands = append(commands, fmt.Sprintf("modprobe %v", v))
	}

	err := execute(commands)
	if err != nil {
		return err
	}

	return nil
}
