// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
//
// SPDX-License-Identifier: Apache-2.0

package stages

import (
	"fmt"
	"init/config"
	"io/ioutil"
)

//Microcode implementes IStage
type Microcode struct {
	finals []string
}

//String ..
func (m *Microcode) String() string {
	return "microcode"
}

//Policy ..
func (m *Microcode) Policy() Policy {
	return PolicySoft
}

//Finalise ..
func (m *Microcode) Finalise() []string {
	return m.finals
}

//Run ..
func (m *Microcode) Run(c config.Config) error {

	err := ioutil.WriteFile("/sys/devices/system/cpu/microcode/reload", []byte("1"), 0644)
	if err != nil {
		return fmt.Errorf("failed to trigger microcode load: %w", err)
	}

	return nil
}
