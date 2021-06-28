// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
//
// SPDX-License-Identifier: Apache-2.0

package stages

import (
	"fmt"
	"init/config"
	"init/shell"
	"os"
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

	commands := []shell.Command{}
	commands = append(commands, shell.Command{Executable: shell.IP, Arguments: []string{"link", "set", "dev", "lo", "up"}})
	for _, nd := range c.Secondary.Networking.Networks {

		if nd.Type != "" {
			// If type not default, create as specified
			commands = append(commands, shell.Command{Executable: shell.IP, Arguments: []string{"link", "add", "dev", nd.Device, "type", nd.Type}})
		}

		if nd.DHCP {
			if nd.IPV6 {
				commands = append(commands, shell.Command{Executable: shell.IP, Arguments: []string{"link", "set", "dev", nd.Device, "up"}})
				commands = append(commands, shell.Command{Executable: shell.Udhcp, Arguments: []string{"-b", "-i", nd.Device, "-p", "/var/run/udhcpc.pid"}})
			} else {
				commands = append(commands, shell.Command{Executable: shell.IP, Arguments: []string{"link", "set", "dev", nd.Device, "up"}})
				commands = append(commands, shell.Command{Executable: shell.Udhcp, Arguments: []string{"-b", "-i", nd.Device, "-p", "/var/run/udhcpc.pid"}})
			}
		} else {

			commands = append(commands, shell.Command{Executable: shell.IP, Arguments: []string{"link", "set", "dev", nd.Device, "up"}})

			for _, v := range nd.Addresses {
				commands = append(commands, shell.Command{Executable: shell.IP, Arguments: []string{"addr", "add", v, "dev", nd.Device}})
			}

			if nd.DefaultGateway != "" {
				commands = append(commands, shell.Command{Executable: shell.IP, Arguments: []string{"route", "add", "default", "via", nd.DefaultGateway, "dev", nd.Device}})
			}
		}
	}

	err := shell.Executor.Execute(commands)
	if err != nil {
		return err
	}

	commands = []shell.Command{}
	for _, rt := range c.Secondary.Networking.Routes {
		commands = append(commands, shell.Command{Executable: shell.IP, Arguments: []string{"route", "add", rt.Address, "dev", rt.Device}})
	}

	err = shell.Executor.Execute(commands)
	if err != nil {
		return err
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
