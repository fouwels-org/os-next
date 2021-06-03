package main

import (
	"fmt"
	"init/config"
	"init/contrib/u-root/libinit"
	"init/stages"
	"init/static"
	"init/util"
	"log"
	"os"

	"syscall"
	"time"
)

var _banner string = static.Splash

const _configPrimaryPath = "/config/primary.json"
const _configSecondaryPath = "/var/config/secondary.json"

func main() {

	log.SetFlags(0)

	err := run()
	if err != nil {
		log.Printf("exit with err: %v", err)
	} else {
		log.Printf("exit without error?")
	}

	// We need to reap all children before exiting.
	log.Printf("all commands exited, syncing filesystems")
	syscall.Sync()

	log.Printf("rebooting in 5 seconds")
	time.Sleep(5 * time.Second)

	os.Exit(1)
}

func run() error {

	err := util.System.SetConsoleLogLevel(util.KLogCritical)

	if err != nil {
		return fmt.Errorf("failed to set kernel log level: %w", err)
	}

	log.Printf("\033[2J") // Clear console
	log.Printf("%v", _banner)

	// creates the rootfs
	libinit.CreateRootfs()

	// run the user-defined init tasks
	err = uinit()
	if err != nil {
		log.Printf("init failed: %v", err)
	}

	for {
		commands := []util.Command{}

		commands = append(commands, util.Command{Target: "/bin/login", Arguments: []string{}})

		err := util.Shell.ExecuteInteractive(commands)
		if err != nil {
			log.Printf("error returned from shell, restarting: %v", err)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// uinit loads the primary and secondary boot stages after the kernel hands over to the init process
// Returns nil if both stages loads successfully, otherwise error
func uinit() error {
	c := config.Config{}

	primary := []stages.IStage{
		&stages.Modules{},
		&stages.KernelConfig{},
		&stages.Filesystem{},
		&stages.Microcode{},
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
		return fmt.Errorf("failed to load primary config from %v: %v", _configPrimaryPath, err)
	}
	c.Primary = configPrimary.Primary

	log.Printf("[uinit] running primary stage")

	_ = executeStages(c, primary)

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
		return fmt.Errorf("failed to find secondary config, not running second stage: %v", err)
	}

	c.Secondary = configSecondary.Secondary

	_ = executeStages(c, secondary)

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

	return nil
}

// function to execute the stages, returns an error if any of the stages fail
func executeStages(c config.Config, stages []stages.IStage) []error {

	errors := []error{}
	for _, st := range stages {

		log.Printf("[%v] starting", st)

		err := st.Run(c)
		if err != nil {
			errors = append(errors, fmt.Errorf("%v failed: %w", st, err))
			log.Printf("[%v] failed: %v", st, err)
		} else {
			log.Printf("[%v] succeeded", st)
		}
	}
	return errors
}
