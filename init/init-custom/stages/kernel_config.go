package stages

import (
	"fmt"
	"init-custom/config"
)

//KernelConfig implements IStage
type KernelConfig struct {
	finals []string
}

//String ..
func (n KernelConfig) String() string {
	return "Kernel Config"
}

//Finalise ..
func (n KernelConfig) Finalise() []string {
	return n.finals
}

//Run ..
func (n KernelConfig) Run(c config.Config) error {

	err := setFile("/sys/fs/cgroup/memory/memory.use_hierarchy", "1", 0644)
	if err != nil {
		return fmt.Errorf("Failed to set file: %w", err)
	}
	return nil
}
