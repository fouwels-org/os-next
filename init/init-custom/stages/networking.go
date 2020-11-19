package stages

import (
	"fmt"
	"init-custom/config"
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
	commands = append(commands, command{command: "/sbin/ip", arguments: []string{"link", "set", "dev", "lo", "up"}})
	for _, n := range c.Networking.Networks {
		if n.DHCP {
			if n.IPV6 {
				commands = append(commands, command{command: "/sbin/ip", arguments: []string{"link", "set", "dev", n.Device, "up"}})
				commands = append(commands, command{command: "/sbin/udhcpc", arguments: []string{"-b", "-i", n.Device, "-p", "/var/run/udhcpc.pid"}})
			} else {
				commands = append(commands, command{command: "/sbin/ip", arguments: []string{"link", "set", "dev", n.Device, "up"}})
				commands = append(commands, command{command: "/sbin/udhcpc", arguments: []string{"-b", "-i", n.Device, "-p", "/var/run/udhcpc.pid"}})
			}
		} else {
			return fmt.Errorf("NOTIMPLEMENTED: Static Addressing")
		}
	}

	// wireguard test
	/*
		ip link add dev wg0 type wireguard
		wg set wg0 listen-port 51820 private-key /root/wg.key peer vezQ++zg/pvTjZ73XAXHtTnYi618BvllGHQ37a74tgc= allowed-ips 10.200.4.0/24 persistent-keepalive 5 endpoint 81.201.135.86:51820
		ip address add dev wg0 10.200.4.99/32
		ip link set up dev wg0
		ip route add 10.200.4.0/24 dev wg0
	*/

	err := setFile("/root/wg.key", string("0JD938XTDYNfszGp8EoMOoT1eq710ryzJm6a0JPEkEs="), 0600)
	if err != nil {
		return err
	}

	commands = append(commands, command{command: "/sbin/ip", arguments: []string{"link", "add", "dev", "wg0", "type", "wireguard"}})
	commands = append(commands, command{command: "/usr/sbin/wg", arguments: []string{"set", "wg0", "listen-port", "51820", "private-key", "/root/wg.key", "peer", "vezQ++zg/pvTjZ73XAXHtTnYi618BvllGHQ37a74tgc=", "allowed-ips", "10.200.4.0/24", "persistent-keepalive", "5", "endpoint", "concentrator.lagoni.co.uk:51820"}})
	commands = append(commands, command{command: "/sbin/ip", arguments: []string{"address", "add", "dev", "wg0", "10.200.4.99/32"}})
	commands = append(commands, command{command: "/sbin/ip", arguments: []string{"link", "set", "up", "dev", "wg0"}})
	commands = append(commands, command{command: "/sbin/ip", arguments: []string{"route", "add", "10.200.4.0/24", "dev", "wg0"}})

	err = execute(commands)
	if err != nil {
		return err
	}
	return nil
}
