// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is a basic init script.
package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	commands = []string{
		"modprobe igb",      // K300 Ethernet Driver
		"modprobe btrfs",    // File system needed by Docker
		"modprobe usbnet",   // Not sure if this is needed
		"modprobe qmi_wwan", // Not sure what this is for
		"/bbin/date",
		"/bbin/dhclient -ipv6=false eth0",
		"/bbin/ip a",
	}
)

func main() {
	for _, line := range commands {
		log.Printf("Executing Command: %v", line)
		cmdSplit := strings.Split(line, " ")
		if len(cmdSplit) == 0 {
			continue
		}

		cmd := exec.Command(cmdSplit[0], cmdSplit[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Print(err)
		}

	}
	log.Print("Uinit Done!")
}
