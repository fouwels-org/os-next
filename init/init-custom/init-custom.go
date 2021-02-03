package main

import (
	"fmt"
	"init-custom/config"
	"init-custom/contrib/u-root/libinit"
	"init-custom/stages"
	"init-custom/static"
	"init-custom/util"
	"os"
	"time"

	"log"
	"syscall"
)

var _banner string = static.Splash

const _configPrimaryPath = "/etc/init/primary.json"
const _configSecondaryPath = "/var/config/secondary.json"

func main() {

	log.SetFlags(0)

	err := run()
	if err != nil {
		log.Printf("Exit with err: %v", err)
	} else {
		log.Printf("Exit without error?")
	}

	// We need to reap all children before exiting.
	log.Printf("all commands exited, syncing filesystems")
	syscall.Sync()

	os.Exit(1)
}

func run() error {

	err := util.System.SetConsoleLogLevel(util.KLogCritical)
	if err != nil {
		return fmt.Errorf("Failed to set kernel log level: %v", err)
	}

	log.Printf("\033[2J") // Clear console
	log.Printf(_banner)

	// creates the rootfs
	libinit.CreateRootfs()

	// run the user-defined init tasks
	_, secondaryLoaded, err := uinit()
	if err != nil {
		return fmt.Errorf("Init failed : %v", err)
	}

	for {
		commands := []util.Command{}
		if !secondaryLoaded {
			commands = append(commands, util.Command{Target: "/bin/ash", Arguments: []string{}})
		} else {
			// Loading the full operational OS
			commands = append(commands, util.Command{Target: "/bin/login", Arguments: []string{}})
		}

		err := util.Shell.ExecuteInteractive(commands)
		if err != nil {
			log.Printf("Error returned from shell, restarting: %v", err)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// This function loads the primary and secondary boot stages after the kernel hands over to the init process
// return parameters:
// 		bool: true if the primary stage loads successfully otherwise false
// 		bool: true if the secondary stage loads successfully otherwise false
//		error: nil if both stages loads successfully, otherwise an error is returned
func uinit() (bool, bool, error) {
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

	log.Printf("[uinit] loading primary config")
	configPrimary := config.PrimaryFile{}

	err := config.LoadConfig(_configPrimaryPath, &configPrimary)
	if err != nil {
		return false, false, fmt.Errorf("failed to load primary config from %v: %v", _configPrimaryPath, err)
	}
	c.Primary = configPrimary.Primary

	log.Printf("[uinit] running primary stage")
	err = executeStages(c, primary)
	if err != nil {
		return false, false, err
	}

	log.Printf("[uinit] loading secondary config")
	secondLoaded := false

	configSecondary := config.SecondaryFile{}
	err = config.LoadConfig(_configSecondaryPath, &configSecondary)
	if err != nil {
		log.Printf("[uinit] failed to find secondary config, not running second stage: %v", err)
		secondLoaded = false
	} else {
		secondLoaded = true
	}

	if secondLoaded {
		c.Secondary = configSecondary.Secondary
		err = executeStages(c, secondary)
		if err != nil {
			return true, secondLoaded, err
		}
	}

	log.Printf("[uinit] init complete")
	log.Printf("[uinit] messages:")

	for _, st := range primary {

		finals := st.Finalise()
		if len(finals) == 0 {
			continue
		}

		for _, f := range finals {
			log.Printf("[uinit] %v: %v", st, f)
		}
	}

	if secondLoaded {
		for _, st := range secondary {

			finals := st.Finalise()
			if len(finals) == 0 {
				continue
			}

			for _, f := range finals {
				log.Printf("[uinit] %v: %v", st, f)
			}
		}
	}

	return true, true, nil
}

// function to execute the stages, returns an error if any of the stages fail
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
