package stages

import (
	"fmt"
	"uinit-custom/config"
)

//Networking implements IStage
type Housekeeping struct {
	finals []string
}

//String ..
func (n Housekeeping) String() string {
	return "House Keeping"
}

//Finalise ..
func (n Housekeeping) Finalise() []string {
	return n.finals
}

//Run ..
func (n Housekeeping) Run(c config.Config) error {

	err := writeLines("/sys/fs/cgroup/memory/memory.use_hierarchy", "1")
	if err != nil {
		logf("Error setting memory use_hierarchy:  " + err.Error())
		return err
	}

	commands := []string{}
	v := "/dev/sda2"
	commands = append(commands, fmt.Sprintf("mount -t ext4 %v /var/lib/docker", v))

	err = execute(commands)
	if err != nil {
		logf("Error setting memory use_hierarchy:  " + err.Error())
		return err
	}

	return nil
}
