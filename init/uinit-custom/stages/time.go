package stages

import (
	"fmt"
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

	_, err := executeOne(fmt.Sprintf("/bbin/ntpdate"), "")
	if err != nil {
		log.Printf("Error updating NTP: %v", err)
	}

	_, err = executeOne(fmt.Sprintf("/bbin/hwclock -w"), "")
	if err != nil {
		log.Printf("Error setting HW Clock: %v", err)
	}

	return nil
}
