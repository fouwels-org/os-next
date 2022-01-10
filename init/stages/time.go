// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package stages

import (
	"fmt"
	"os"
	"time"

	"github.com/fouwels/os-next/init/config"
	"github.com/fouwels/os-next/init/journal"
	"github.com/fouwels/os-next/init/shell"
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

	// Configure NTP
	f, err := os.OpenFile("/etc/ntp.conf", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file to write ntp settings: %v", err)
	}
	defer f.Close()

	for _, ns := range c.Secondary.Time.Servers {
		_, err = fmt.Fprintf(f, "server %v\n", ns)
		if err != nil {
			return fmt.Errorf("failed to write ntp server: %v", err)
		}
	}

	err = f.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync on %v: %v", f.Name(), err)
	}

	ferr := f.Close()
	if ferr != nil {
		e = fmt.Errorf("failed to close on %v: %w", f.Name(), ferr)
	}

	if len(c.Secondary.Time.Servers) != 0 {
		n.finals = append(n.finals, fmt.Sprintf("NTP servers set to %v", c.Secondary.Time.Servers))
	} else {
		if c.Secondary.Time.NTP {
			n.finals = append(n.finals, "warning: NTP enabled, but no NTP servers have been configured")
		}
	}

	// Run command set
	if c.Secondary.Time.NTP {

		commands := []shell.Command{}
		commands = append(commands, shell.Command{Executable: shell.Ntpd, Arguments: []string{"-q"}})

		err := shell.Executor.Execute(commands)
		if err != nil {
			journal.Logfln("Error updating NTP: %v", err)
		}
	} else {
		n.finals = append(n.finals, "notice: NTP Disabled")
	}

	if c.Secondary.Time.HWClock {

		commands := []shell.Command{}
		commands = append(commands, shell.Command{Executable: shell.Hwclock, Arguments: []string{"-w"}})

		err := shell.Executor.Execute(commands)
		if err != nil {
			journal.Logfln("Error setting HW Clock: %v", err)
		}
	} else {
		n.finals = append(n.finals, "notice: HW Clock Disabled")
	}

	n.finals = append(n.finals, fmt.Sprintf("time is now: %v", time.Now().UTC()))

	return nil
}
