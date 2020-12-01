package stages

import (
	"fmt"
	"init-custom/config"
	"log"
	"os"
)

//Time implements IStage
type Time struct {
	finals []string
}

//String ..
func (n Time) String() string {
	return "Time"
}

//Finalise ..
func (n Time) Finalise() []string {
	return n.finals
}

//Run ..
func (n Time) Run(c config.Config) (e error) {

	// Configure NTP
	err := func() error {
		// #nosec G302 (CWE-276). 644 is intentional.
		f, err := os.OpenFile("/etc/ntp.conf", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 644)
		if err != nil {
			return fmt.Errorf("Failed to open file to write ntp settings: %v", err)
		}
		// #nosec G307. Double defer is safe for file.Writer
		defer f.Close()

		for _, ns := range c.Secondary.Time.Servers {
			_, err = fmt.Fprintf(f, "server %v\n", ns)
			if err != nil {
				return fmt.Errorf("Failed to write ntp server: %v", err)
			}
		}

		err = f.Sync()
		if err != nil {
			return fmt.Errorf("Failed to sync on %v: %v", f.Name(), err)
		}

		ferr := f.Close()
		if ferr != nil {
			e = fmt.Errorf("Failed to close on %v: %w", f.Name(), ferr)
		}

		return nil
	}()
	if err != nil {
		return err
	}

	if len(c.Secondary.Time.Servers) != 0 {
		n.finals = append(n.finals, fmt.Sprintf("NTP servers set to %v", c.Secondary.Time.Servers))
	} else {
		if c.Secondary.Time.NTP {
			n.finals = append(n.finals, fmt.Sprintf("[warn] NTP enabled, but no NTP servers have been configured"))
		}
	}

	// Run command set
	if c.Secondary.Time.NTP {
		_, err = executeOne(command{command: "/sbin/ntpd", arguments: []string{}}, "")
		if err != nil {
			log.Printf("Error updating NTP: %v", err)
		}
	} else {
		n.finals = append(n.finals, fmt.Sprintf("[Notice] NTP Disabled"))
	}

	if c.Secondary.Time.HWClock {
		_, err = executeOne(command{command: "/sbin/hwclock", arguments: []string{"-w"}}, "")
		if err != nil {
			log.Printf("Error setting HW Clock: %v", err)
		}
	} else {
		n.finals = append(n.finals, fmt.Sprintf("[Notice] HW Clock Disabled"))
	}

	return nil
}
