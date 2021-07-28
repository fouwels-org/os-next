// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"init/config"
	"init/console"
	"init/contrib/u-root/libinit"
	"init/kernel"
	"init/shell"
	"init/stages"
	"log"
	"os"
	"syscall"
	"time"
)

const _configPrimaryPath = "/config/primary.yml"
const _configSecondaryPath = "/var/config/secondary.yml"

func main() {

	log.SetFlags(log.Lmicroseconds | log.LUTC)

	err := run()
	if err != nil {
		log.Printf("exit with err: %v", err)
	} else {
		log.Printf("exit without error")
	}

	// Sync file system
	_, _, serr := syscall.Syscall(syscall.SYS_SYNC, 0, 0, 0)
	if serr != 0 {
		log.Printf("failed to sync file system: %v", err)
	}

	log.Printf("rebooting in 5 seconds")
	time.Sleep(5 * time.Second)

	// Reboot
	err = syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
	log.Printf("reboot syscall failed, exiting to kernel in 5 seconds, good luck: %v", err)
	time.Sleep(5 * time.Second)

	os.Exit(1)

}

func run() error {

	err := kernel.SetLogLevel(kernel.KLogCritical)
	if err != nil {
		return fmt.Errorf("failed to set kernel log level: %w", err)
	}

	log.Printf("\033[2J") // Clear console

	// run self tests
	err = shell.SelfTest()
	if err != nil {
		return fmt.Errorf("failed shell self test: %w", err)
	}

	// creates the rootfs
	libinit.CreateRootfs()

	// run the user defined init tasks
	cfg, err := uinit()
	if err != nil {

		return fmt.Errorf("init failed: %v", err)
	}

	// start the console
	err = console.Start(cfg.Secondary.Authenticators)
	if err != nil {
		return fmt.Errorf("failed to run console: %w", err)
	}

	return nil
}

// uinit loads the primary and secondary boot stages after the kernel hands over to the init process
// Returns nil if both stages loads successfully, otherwise error
func uinit() (config.Config, error) {
	c := config.Config{}

	primary := []stages.IStage{
		&stages.Modules{},
		&stages.KernelConfig{},
		&stages.TPM{},
		&stages.Filesystem{},
		&stages.Microcode{},
	}

	secondary := []stages.IStage{
		&stages.Modules{},
		&stages.Networking{},
		&stages.Wireguard{},
		&stages.Time{},
		&stages.Docker{},
	}

	log.Printf("[uinit] loading primary config")

	configPrimary := config.PrimaryFile{}

	err := config.LoadConfig(_configPrimaryPath, &configPrimary)
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to load primary config from %v: %v", _configPrimaryPath, err)
	}
	c.Primary = configPrimary.Primary

	log.Printf("[uinit] running primary stage")

	err = executeStages(c, primary)
	if err != nil {
		return config.Config{}, fmt.Errorf("primary: %w", err)
	}

	log.Printf("[uinit] messages (primary):")

	for _, st := range primary {

		finals := st.Finalise()
		if len(finals) == 0 {
			continue
		}

		for _, f := range finals {
			log.Printf("[uinit] %v: %v", st, f)
		}
	}

	log.Printf("[uinit] loading secondary config")

	configSecondary := config.SecondaryFile{}
	err = config.LoadConfig(_configSecondaryPath, &configSecondary)
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to find secondary config, not running second stage: %v", err)
	}

	c.Secondary = configSecondary.Secondary

	log.Printf("[uinit] running secondary stage")

	err = executeStages(c, secondary)
	if err != nil {
		return config.Config{}, fmt.Errorf("secondary: %w", err)
	}

	log.Printf("[uinit] messages (secondary):")

	for _, st := range secondary {

		finals := st.Finalise()
		if len(finals) == 0 {
			continue
		}

		for _, f := range finals {
			log.Printf("[uinit] %v: %v", st, f)
		}
	}

	log.Printf("[uinit] initialisation complete")

	return c, nil
}

func executeStages(c config.Config, sts []stages.IStage) error {

	for _, st := range sts {

		log.Printf("[%v] starting", st)

		err := st.Run(c)
		if err != nil {

			switch st.Policy() {
			case stages.PolicyHard:
				log.Printf("[%v] failed (hard): %v", st, err)
				return fmt.Errorf("%v failed", st)
			case stages.PolicySoft:
				log.Printf("[%v] failed (soft): %v", st, err)
			}

		} else {
			log.Printf("[%v] succeeded", st)
		}

		// Sync file system
		_, _, serr := syscall.Syscall(syscall.SYS_SYNC, 0, 0, 0)
		if serr != 0 {
			log.Printf("failed to sync file system after %v: %v", st, err)
		}
	}

	return nil
}
