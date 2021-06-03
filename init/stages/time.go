package stages

import (
	"fmt"
	"init/config"
	"init/util"
	"log"
	"os"
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

		commands := []util.Command{}
		commands = append(commands, util.Command{Target: "/sbin/ntpd", Arguments: []string{"-q"}})

		err := util.Shell.Execute(commands)
		if err != nil {
			log.Printf("Error updating NTP: %v", err)
		}
	} else {
		n.finals = append(n.finals, "notice: NTP Disabled")
	}

	if c.Secondary.Time.HWClock {

		commands := []util.Command{}
		commands = append(commands, util.Command{Target: "/sbin/hwclock", Arguments: []string{"-w"}})

		err := util.Shell.Execute(commands)
		if err != nil {
			log.Printf("Error setting HW Clock: %v", err)
		}
	} else {
		n.finals = append(n.finals, "notice: HW Clock Disabled")
	}

	n.finals = append(n.finals, fmt.Sprintf("time is now: %v", time.Now().UTC()))

	return nil
}
