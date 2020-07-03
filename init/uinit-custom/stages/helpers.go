package stages

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func executeOne(command string) error {

	cmdSplit := strings.Split(command, " ")
	if len(cmdSplit) == 0 {
		return fmt.Errorf("Empty command provided")
	}

	cmd := exec.Command(cmdSplit[0], cmdSplit[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func execute(command []string) error {
	for _, c := range command {
		err := executeOne(c)
		if err != nil {
			return err
		}
	}

	return nil
}
