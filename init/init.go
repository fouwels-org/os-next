// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"os-next/init/config"
	"os-next/init/console"
	"os-next/init/external/u-root/libinit"
	"os-next/init/journal"
	"os-next/init/kernel"
	"os-next/init/shell"
	"os-next/init/stages"
	"syscall"
	"time"
)

const _configPrimaryPath = "/config/primary.yml"
const _configSecondaryPath = "/var/config/secondary.yml"

func main() {

	err := run()
	if err != nil {
		journal.Logfln("exit with err: %v", err)
	} else {
		journal.Logfln("exit without error")
	}

	// Sync file system
	_, _, serr := syscall.Syscall(syscall.SYS_SYNC, 0, 0, 0)
	if serr != 0 {
		journal.Logfln("failed to sync file system: %v", err)
	}

	journal.Logfln("rebooting in 5 seconds")
	time.Sleep(5 * time.Second)

	// Reboot
	err = syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
	journal.Logfln("reboot syscall failed, exiting to kernel in 5 seconds, good luck: %v", err)
	time.Sleep(5 * time.Second)

	os.Exit(1)

}

func run() error {

	err := kernel.SetLogLevel(kernel.KLogCritical)
	if err != nil {
		return fmt.Errorf("failed to set kernel log level: %w", err)
	}

	fmt.Printf("\033[2J") // Clear console

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
	err = console.Login(cfg.Secondary.Authenticators)
	if err != nil {
		return fmt.Errorf("failed to run login: %w", err)
	}

	err = console.Shell()
	if err != nil {
		return fmt.Errorf("failed to run shell: %w", err)
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

	journal.Logf("primary config: ")

	configPrimary := config.PrimaryFile{}

	err := config.LoadConfig(_configPrimaryPath, &configPrimary)
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to load primary config from %v: %v", _configPrimaryPath, err)
	}
	c.Primary = configPrimary.Primary

	journal.Logfln("✔️")

	err = executeStages(c, primary)
	if err != nil {
		return config.Config{}, fmt.Errorf("primary: %w", err)
	}

	journal.Logf("secondary config: ")

	configSecondary := config.SecondaryFile{}
	err = config.LoadConfig(_configSecondaryPath, &configSecondary)
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to find secondary config, not running second stage: %v", err)
	}

	c.Secondary = configSecondary.Secondary

	journal.Logfln("✔️")

	err = executeStages(c, secondary)
	if err != nil {
		return config.Config{}, fmt.Errorf("secondary: %w", err)
	}

	journal.Logfln("messages:")

	for _, st := range primary {

		finals := st.Finalise()
		if len(finals) == 0 {
			continue
		}

		for _, f := range finals {
			journal.Logfln("| %v: %v", st, f)
		}
	}

	for _, st := range secondary {

		finals := st.Finalise()
		if len(finals) == 0 {
			continue
		}

		for _, f := range finals {
			journal.Logfln("| %v: %v", st, f)
		}
	}

	return c, nil
}

func executeStages(c config.Config, sts []stages.IStage) error {

	for _, st := range sts {

		journal.Logf("%v: ", st)

		err := st.Run(c)
		if err != nil {

			switch st.Policy() {
			case stages.PolicyHard:
				journal.Logfln("❌ hard fail: %v", err)
				//return fmt.Errorf("%v failed", st)
				console.Shell()
			case stages.PolicySoft:
				journal.Logfln("❗ soft fail: %v", err)
			}

		} else {
			journal.Logfln("✔️")
		}

		// Sync file system
		_, _, serr := syscall.Syscall(syscall.SYS_SYNC, 0, 0, 0)
		if serr != 0 {
			journal.Logfln("failed to sync file system after %v: %v", st, err)
		}
	}

	return nil
}
