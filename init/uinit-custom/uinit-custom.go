package main

import (
	"fmt"
	"log"
	"os"
	"uinit-custom/config"
	"uinit-custom/stages"
)

var _configPath = "/boot/config.json"
var _secretsPath = "/boot/secrets.json"

func main() {
	err := run()
	if err != nil {
		logf("%v", err)
	} else {
		logf("Exit without error - this is unexpected")
	}

	os.Exit(-1)
}

func run() error {

	logf("Loading config")
	c, err := config.LoadConfig(_configPath)
	if err != nil {
		return fmt.Errorf("Failed to load config from %v: %v", _configPath, err)
	}
	s, err := config.LoadSecrets(_secretsPath)
	if err != nil {
		return fmt.Errorf("Failed to load secrets from %v: %v", _secretsPath, err)
	}

	stageList := []stages.IStage{
		stages.Modules{},
		stages.Networking{},
	}

	logf("Executing stages")

	for _, st := range stageList {

		logf("[%v] starting", st)

		err := st.Run(c, s)
		if err != nil {
			return fmt.Errorf("[%v] failed: %v", st, err)
		}
		logf("[%v] succeeded", st)
	}

	return nil
}

func logf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Printf("%v", message)
}
