// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package stages

import (
	"fmt"
	"init/config"
	"init/tpm"
)

//TPM implements IStage
type TPM struct {
	finals []string
}

//String ..
func (n *TPM) String() string {
	return "TPM"
}

//Policy ..
func (n *TPM) Policy() Policy {
	return PolicySoft
}

//Finalise ..
func (n *TPM) Finalise() []string {
	return n.finals
}

//Run ..
func (n *TPM) Run(c config.Config) (e error) {

	values, err := tpm.ReadPCRs()
	if err != nil {
		return fmt.Errorf("failed to read TPM PCRs: %w", err)
	}

	str := ""

	for _, v := range values {
		str = str + fmt.Sprintf("%X", v.Value[len(v.Value)-2:]) + "-"
	}

	n.finals = append(n.finals, str[:len(str)-1])

	return nil
}
