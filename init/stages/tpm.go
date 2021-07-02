// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: LicenseRef-MS-RSL

package stages

import (
	"fmt"
	"init/config"
	"init/tpm"
	"log"
)

//TPM implements IStage
type TPM struct {
	finals []string
}

//String ..
func (n *TPM) String() string {
	return "time"
}

//Policy ..
func (n *TPM) Policy() Policy {
	return PolicyHard
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

	for _, v := range values {
		log.Printf("%v 0x%X (%X)", v.ID, v.Value, v.Value[len(v.Value)-3:])
	}
	return nil
}
