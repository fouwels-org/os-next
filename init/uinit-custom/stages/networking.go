package stages

import (
	"fmt"
	"uinit-custom/config"
)

//Networking implements IStage
type Networking struct {
	finals []string
}

//String ..
func (n Networking) String() string {
	return "Networking"
}

//Finalise ..
func (n Networking) Finalise() []string {
	return n.finals
}

//Run ..
func (n Networking) Run(c config.Config) error {

	commands := []command{}

	for _, n := range c.Networking.Networks {
		if n.DHCP {
			if n.IPV6 {
				commands = append(commands, command{command: "/bbin/dhclient", arguments: []string{n.Device}})
			} else {
				commands = append(commands, command{command: "/bbin/dhclient", arguments: []string{"--ipv6=false", n.Device}})
			}
		} else {
			return fmt.Errorf("NOTIMPLEMENTED: Static Addressing")
		}
	}

	err := execute(commands)
	if err != nil {
		return err
	}
	return nil
}
