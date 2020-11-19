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

	"log"
	"os/exec"
	"syscall"

	"github.com/u-root/u-root/pkg/libinit"
	"github.com/u-root/u-root/pkg/ulog"
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

const _configPath = "/etc/init/config.json"

func main() {
	flag.Parse()

	log.Printf("Welcome to Mjolnir - IIoT OS!")
	log.SetPrefix("init: ")

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
	libinit.SetEnv()
	// creates the rootfs
	libinit.CreateRootfs()

	// Turn off job control when test mode is on.
	ctty := libinit.WithTTYControl(!*test)

	cmdList := []*exec.Cmd{
		libinit.Command("/bin/ash", ctty), // start ash when all user defined init is done
	}
	// run the user-defined init taskes
	err := uinit()
	if err != nil {
		logf("failed loading the staged init services: %v", err)
	}
	// finally run the list of commands
	cmdCount := libinit.RunCommands(debug, cmdList...)
	if cmdCount == 0 {
		log.Printf("No suitable executable found in %v", cmdList)
	}

	// We need to reap all children before exiting.
	log.Printf("Waiting for orphaned children")
	libinit.WaitOrphans()
	log.Printf("All commands exited")
	log.Printf("Syncing filesystems")
	syscall.Sync()
	log.Printf("Exiting...")

}

func uinit() error {

	logf("loading config")
	c, err := config.LoadConfig(_configPath)
	if err != nil {
		return fmt.Errorf("failed to load config from %v: %v", _configPath, err)
	}

	stageList := []stages.IStage{
		&stages.Modules{},
		&stages.Housekeeping{},
		&stages.Networking{},
		&stages.Time{},
		&stages.Docker{},
		&stages.Systeminfo{},
	}

	logf("executing stages")

	for _, st := range stageList {

		logf("[%v] starting", st)

		err := st.Run(c)
		if err != nil {
			logf("[%v] failed: %v/n", st, err)
		} else {
			logf("[%v] succeeded", st)
		}
	}

	logf("stage information")

	for _, st := range stageList {

		finals := st.Finalise()
		if len(finals) == 0 {
			continue
		}

		for _, f := range finals {
			logf("[%v] %v", st, f)
		}
	}

	return nil
}

func logf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Printf("[uinit] %v", message)
}
