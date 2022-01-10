// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package stages

import "github.com/fouwels/os-next/init/config"

type Policy int

const (
	PolicyHard Policy = iota // Hard fail, failure of stage aborts bootup
	PolicySoft               // Soft fail, failure of stage continued boot
)

//IStage ..
type IStage interface {
	//Called during the sequental stage init
	Run(config.Config) error
	//Friendly stage name
	String() string
	//Final strings, called after all stages have been initialized. (eg. to render the stage acquired DHCP IP, or generated SSH public key)
	Finalise() []string
	//Policy
	Policy() Policy
}
