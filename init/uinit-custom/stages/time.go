package stages

import (
	"fmt"
	"uinit-custom/config"
)

//Time implements IStage
type Time struct {
	finals []string
}

//String ..
func (n Time) String() string {
	return "Time"
}

//Finalise ..
func (n Time) Finalise() []string {
	return n.finals
}

//Run ..
func (n Time) Run(c config.Config) error {

	commands := []string{}

	commands = append(commands, fmt.Sprintf("/bbin/ntpdate"))

	err := execute(commands)
	if err != nil {
		return err
	}
	return nil
}
