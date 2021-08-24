// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package stages

import (
	"fmt"
	"init/config"
	"init/shell"
	"time"
)

//Time implements IStage
type Time struct {
	finals []string
}

//String ..
func (n *Time) String() string {
	return "time"
}

//Policy ..
func (n *Time) Policy() Policy {
	return PolicyHard
}

//Finalise ..
func (n *Time) Finalise() []string {
	return n.finals
}

//Run ..
func (n *Time) Run(c config.Config) (e error) {

	if c.Secondary.Time.NTP {

		if len(c.Secondary.Time.Servers) == 0 {
			return fmt.Errorf("NTP is enabled, but no NTP servers have been specified")
		}

		commands := []shell.Command{}
		commands = append(commands, shell.Command{Executable: shell.Sntp, Arguments: []string{"-q", c.Secondary.Time.Servers[0]}})

		err := shell.Executor.Execute(commands)
		if err != nil {
			n.finals = append(n.finals, fmt.Sprintf("warning: failed to update NTP from %v: %v", c.Secondary.Time.Servers[0], err))
		}
	} else {
		n.finals = append(n.finals, "notice: NTP is not enabled")
	}
	if c.Secondary.Time.HWClock {

		commands := []shell.Command{}
		commands = append(commands, shell.Command{Executable: shell.Hwclock, Arguments: []string{"-w"}})

		err := shell.Executor.Execute(commands)
		if err != nil {
			n.finals = append(n.finals, fmt.Sprintf("warning: failed to set HW clock: %v", err))
		}
	} else {
		n.finals = append(n.finals, "notice: HW Clock Disabled")
	}

	n.finals = append(n.finals, fmt.Sprintf("time is now: %v", time.Now().UTC()))

	return nil
}
