package stages

import (
	"log"
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

	_, err := executeOne(command{command: "/bbin/ntpdate", arguments: []string{}}, "")
	if err != nil {
		log.Printf("Error updating NTP: %v", err)
	}

	_, err = executeOne(command{command: "/bbin/hwclock", arguments: []string{"-w"}}, "")
	if err != nil {
		log.Printf("Error setting HW Clock: %v", err)
	}

	return nil
}
