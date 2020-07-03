package stages

import (
	"fmt"
	"uinit-custom/config"
)

//Networking implements IStage
type Networking struct {
}

//String ..
func (n Networking) String() string {
	return "Networking"
}

//Run ..
func (n Networking) Run(c config.Config, s config.Secrets) error {

	commands := []string{}

	for _, n := range c.Networking.Networks {
		if n.DHCP {
			if n.IPV6 {
				commands = append(commands, fmt.Sprintf("/bbin/dhclient %v", n.Device))
			} else {
				commands = append(commands, fmt.Sprintf("/bbin/dhclient -ipv6=false %v", n.Device))
			}
		} else {
			return fmt.Errorf("NOTIMPLEMENTED: Static Addressing")
		}
	}

	commands = append(commands, "/bbin/ip a")

	err := execute(commands)
	if err != nil {
		return err
	}
	return nil
}
