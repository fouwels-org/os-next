package stages

import (
	"fmt"
	"init-custom/config"
)

//Housekeeping implements IStage
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

	err := setFile("/sys/fs/cgroup/memory/memory.use_hierarchy", "1", 0644)
	if err != nil {
		return fmt.Errorf("Failed to set file: %w", err)
	}

	commands := []command{
		{command: "/bin/mount", arguments: []string{"-t", "ext4", "/dev/sda2", "/var/lib/docker"}},
	}
	err = execute(commands)
	if err != nil {
		return fmt.Errorf("Error mounting: %w", err)
	}

	return nil
}
