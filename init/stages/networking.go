// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package stages

import (
	"fmt"
	"init/config"
	"init/shell"
	"net"
	"os"

	"github.com/vishvananda/netlink"
)

//Networking implements IStage
type Networking struct {
	finals []string
}

//String ..
func (n *Networking) String() string {
	return "networking"
}

//Policy ..
func (m *Networking) Policy() Policy {
	return PolicyHard
}

//Finalise ..
func (n *Networking) Finalise() []string {
	return n.finals
}

//Run ..
func (n *Networking) Run(c config.Config) (e error) {

	lo, err := netlink.LinkByName("lo")
	if err != nil {
		return fmt.Errorf("failed to get link lo: %w", err)
	}
	err = netlink.LinkSetUp(lo)
	if err != nil {
		return fmt.Errorf("failed to set link lo up: %w", err)
	}

	for _, nd := range c.Secondary.Networking.Networks {

		if nd.Type != "" {

			la := netlink.NewLinkAttrs()
			la.Name = nd.Device

			err := netlink.LinkAdd(&netlink.GenericLink{
				LinkAttrs: la,
				LinkType:  nd.Type,
			})
			if err != nil {
				return fmt.Errorf("failed to create link %v or type %v: %w", nd.Device, nd.Type, err)
			}
		}

		link, err := netlink.LinkByName(nd.Device)
		if err != nil {
			return fmt.Errorf("failed to get link %v: %w", nd.Device, err)
		}
		err = netlink.LinkSetUp(link)
		if err != nil {
			return fmt.Errorf("failed to set link %v up: %w", nd.Device, err)
		}

		if nd.DHCP {

			com := []shell.Command{
				// -i: interface
				// -t: send up to n discover packets
				// -T: pause between packets
				// -f: run in foreground
				// -n: exit if no lease
				// -q: exit if lease
				{Executable: shell.DHCP, Arguments: []string{"-i", nd.Device, "-t", "5", "-T", "3", "-f", "-n", "-q", "-p", "/var/run/dhcp.pid"}},
			}

			err = shell.Executor.Execute(com)
			if err != nil {
				return fmt.Errorf("failed to start udhcpc: %w", err)
			}

		} else {

			for _, v := range nd.Addresses {

				addr, err := netlink.ParseAddr(v)
				if err != nil {
					return fmt.Errorf("failed to parse address %v: %w", v, err)
				}
				err = netlink.AddrAdd(link, addr)
				if err != nil {
					return fmt.Errorf("failed to add address %v to %v: %w", v, nd.Device, err)
				}
			}

			if nd.DefaultGateway != "" {

				gatewayIP, err := netlink.ParseAddr(nd.DefaultGateway)
				if err != nil {
					return fmt.Errorf("failed to parse default gateway %v: %w", nd.DefaultGateway, err)
				}

				err = netlink.RouteAdd(&netlink.Route{
					Scope:     netlink.SCOPE_UNIVERSE,
					LinkIndex: link.Attrs().Index,
					Dst:       &net.IPNet{IP: gatewayIP.IP, Mask: gatewayIP.Mask},
				})
				if err != nil {
					return fmt.Errorf("failed to set default gateway for %v: %w", nd.Device, err)
				}
			}
		}
	}

	for _, rt := range c.Secondary.Networking.Routes {

		link, err := netlink.LinkByName(rt.Device)
		if err != nil {
			return fmt.Errorf("failed to get link %v: %w", rt.Device, err)
		}

		ip, err := netlink.ParseAddr(rt.Address)
		if err != nil {
			return fmt.Errorf("failed to parse link address %v: %w", rt, err)
		}

		err = netlink.RouteAdd(&netlink.Route{
			Scope:     netlink.SCOPE_UNIVERSE,
			LinkIndex: link.Attrs().Index,
			Dst:       &net.IPNet{IP: ip.IP, Mask: ip.Mask},
		})

		if err != nil {
			return fmt.Errorf("failed to set link %v: %w", rt, err)
		}
	}

	if len(c.Secondary.Networking.Nameservers) != 0 {
		f, err := os.OpenFile("/etc/resolv.conf", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file to write nameservers: %v", err)
		}
		defer f.Close()

		for _, ns := range c.Secondary.Networking.Nameservers {
			_, err = fmt.Fprintf(f, "nameserver %v\n", ns)
			if err != nil {
				return fmt.Errorf("failed to write nameserver: %v", err)
			}
		}

		err = f.Sync()
		if err != nil {
			return fmt.Errorf("failed to sync on %v: %v", f.Name(), err)
		}

		err = f.Close()
		if err != nil {
			return fmt.Errorf("failed to close on %v: %v", f.Name(), err)
		}

		n.finals = append(n.finals, fmt.Sprintf("nameservers configured to %v", c.Secondary.Networking.Nameservers))
	}
	return nil
}
