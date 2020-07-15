package stages

import (
	"fmt"
	"log"
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

	err := execute(commands)
	if err != nil {
		log.Printf("")
	}
	return nil
}
