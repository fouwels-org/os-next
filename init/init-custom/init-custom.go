// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is a basic init script.
package main

import (
	"flag"
	"fmt"
	"init-custom/config"
	"init-custom/stages"
	"os"
	"time"

	"log"
	"os/exec"
	"syscall"

	"init-custom/contrib/u-root/libinit"
	"init-custom/contrib/u-root/ulog"
)

// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// init is u-root's standard userspace init process.
//
// init is intended to be the first process run by the kernel when it boots up.
// init does some basic initialization (mount file systems, turn on loopback)
// and then tries to execute, in order, /inito, a uinit (either in /bin, /bbin,
// or /ubin), and then a shell (/bin/defaultsh and /bin/sh).

var (
	verbose  = flag.Bool("v", false, "print all build commands")
	test     = flag.Bool("test", false, "Test mode: don't try to set control tty")
	debug    = func(string, ...interface{}) {}
	osInitGo = func() {}
)

const _configPrimaryPath = "/etc/init/primary.json"
const _configSecondaryPath = "/var/config/secondary.json"

func main() {
	err := run()
	if err != nil {
		log.Printf("Exit with err: %v", err)
	} else {
		log.Printf("Exit without error?")
	}

	// We need to reap all children before exiting.
	log.Printf("All commands exited")
	log.Printf("Syncing filesystems")
	syscall.Sync()
	log.Printf("Exiting...")

	os.Exit(1)
}

func run() error {

	flag.Parse()
	log.Printf("Welcome to Mjolnir - IIoT OS!")
	log.SetPrefix("[init]: ")

	if *verbose {
		debug = log.Printf
	}

	// Before entering an interactive shell, decrease the loglevel because
	// spamming non-critical logs onto the shell frustrates users. The logs
	// are still accessible through dmesg.
	if !*verbose {
		// Only messages more severe than "notice" are printed.
		if err := ulog.KernelLog.SetConsoleLogLevel(ulog.KLogNotice); err != nil {
			log.Printf("Could not set log level: %v", err)
		}
	}

	// sets the system environmental variables
	err := libinit.SetEnv()
	if err != nil {
		return fmt.Errorf("Failed to set system environment variables: %w", err)
	}

	// creates the rootfs
	libinit.CreateRootfs()

	// run the user-defined init taskes
	err = uinit()
	if err != nil {
		return fmt.Errorf("failed loading the staged init services: %v", err)
	}

	// sleep to allow output from uinit to be read
	log.Printf("Waiting 3 seconds to start TTY")
	time.Sleep(3 * time.Second)

	for {
		// Turn off job control when test mode is on.
		ctty := libinit.WithTTYControl(!*test)

		cmdList := []*exec.Cmd{
			libinit.Command("/bin/top", ctty),   // display the processes running and memory usage
			libinit.Command("/bin/login", ctty), // start login so there is no direct access to the shell, if not logged in within 60sec of existing top then it will exit and show top again.
		}
		// finally run the list of commands
		cmdCount := libinit.RunCommands(debug, cmdList...)
		if cmdCount == 0 {
			return fmt.Errorf("No suitable executable found in %v", cmdList)
		}
	}
}

func uinit() error {
	c := config.Config{}

	primary := []stages.IStage{
		&stages.Modules{},
		&stages.KernelConfig{},
		&stages.Filesystem{},
		&stages.Systeminfo{},
	}

	secondary := []stages.IStage{
		&stages.Networking{},
		&stages.Wireguard{},
		&stages.Time{},
		&stages.Docker{},
	}

	log.Printf("loading primary config")
	configPrimary := config.PrimaryFile{}

	err := config.LoadConfig(_configPrimaryPath, &configPrimary)
	if err != nil {
		return fmt.Errorf("failed to load primary config from %v: %v", _configPrimaryPath, err)
	}
	c.Primary = configPrimary.Primary

	log.Printf("running primary stage")
	err = executeStages(c, primary)
	if err != nil {
		return err
	}

	log.Printf("loading secondary config")
	secondLoaded := false

	configSecondary := config.SecondaryFile{}
	err = config.LoadConfig(_configSecondaryPath, &configSecondary)
	if err != nil {
		log.Printf("Failed to load secondary config, not running second stage: %v", err)
		secondLoaded = false
	} else {
		secondLoaded = true
	}

	if secondLoaded {
		c.Secondary = configSecondary.Secondary
		err = executeStages(c, secondary)
		if err != nil {
			return err
		}
	}

	log.Printf("finalised:")

	for _, st := range primary {

		finals := st.Finalise()
		if len(finals) == 0 {
			continue
		}

		for _, f := range finals {
			log.Printf("[%v] %v", st, f)
		}
	}

	if secondLoaded {
		for _, st := range secondary {

			finals := st.Finalise()
			if len(finals) == 0 {
				continue
			}

			for _, f := range finals {
				log.Printf("[%v] %v", st, f)
			}
		}
	}

	return nil
}

func executeStages(c config.Config, stages []stages.IStage) error {

	for _, st := range stages {

		log.Printf("[%v] starting", st)

		err := st.Run(c)
		if err != nil {
			log.Printf("[%v] failed: %v/n", st, err)
		} else {
			log.Printf("[%v] succeeded", st)
		}
	}
	return nil
}
