// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
//
// SPDX-License-Identifier: Apache-2.0

package stages

import (
	"fmt"
	"init/config"
	"init/filesystem"
)

//KernelConfig implements IStage
type KernelConfig struct {
	finals []string
}

//String ..
func (n *KernelConfig) String() string {
	return "kernel config"
}

//Policy ..
func (n *KernelConfig) Policy() Policy {
	return PolicyHard
}

//Finalise ..
func (n *KernelConfig) Finalise() []string {
	return n.finals
}

//Run ..
func (n *KernelConfig) Run(c config.Config) error {

	err := filesystem.WriteSync("/sys/fs/cgroup/memory/memory.use_hierarchy", []byte("1"))
	if err != nil {
		return fmt.Errorf("failed to set file: %w", err)
	}
	return nil
}
