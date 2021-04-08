package stages

import (
	"fmt"
	"init-custom/config"
	"init-custom/util"
)

//KernelConfig implements IStage
type KernelConfig struct {
	finals []string
}

//String ..
func (n *KernelConfig) String() string {
	return "kernel config"
}

//Finalise ..
func (n *KernelConfig) Finalise() []string {
	return n.finals
}

//Run ..
func (n *KernelConfig) Run(c config.Config) error {

	err := util.File.SetFile("/sys/fs/cgroup/memory/memory.use_hierarchy", "1", 0644)
	if err != nil {
		return fmt.Errorf("failed to set file: %w", err)
	}
	return nil
}
