package main

import (
	"fmt"
	"log"
	"os"
	"uinit-custom/config"
	"uinit-custom/stages"
)

var _configPath = "/etc/uinit/config.json"

func main() {
	err := run()
	if err != nil {
		logf("%v", err)
	}

	fmt.Printf("press enter to drop to shell")
	fmt.Scanln()

	os.Exit(-1)
}

func run() error {

	logf("loading config")
	c, err := config.LoadConfig(_configPath)
	if err != nil {
		return fmt.Errorf("failed to load config from %v: %v", _configPath, err)
	}

	stageList := []stages.IStage{
		&stages.Modules{},
		&stages.Networking{},
		&stages.Docker{},
	}

	logf("executing stages")

	for _, st := range stageList {

		logf("[%v] starting", st)

		err := st.Run(c)
		if err != nil {
			return fmt.Errorf("[%v] failed: %v", st, err)
		}
		logf("[%v] succeeded", st)
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
