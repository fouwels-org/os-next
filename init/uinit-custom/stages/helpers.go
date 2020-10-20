package stages

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func executeOne(command string, stdin string) (string, error) {

	cmdSplit := strings.Split(command, " ")
	if len(cmdSplit) == 0 {
		return "", fmt.Errorf("Empty command provided")
	}

	buffer := bytes.Buffer{}
	buffer.Write([]byte(stdin))

	cmd := exec.Command(cmdSplit[0], cmdSplit[1:]...)
	cmd.Stdin = &buffer
	out, err := cmd.CombinedOutput()

	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

func execute(command []string) error {
	for _, c := range command {
		out, err := executeOne(c, "")
		if err != nil {
			return fmt.Errorf("%v failed: %v %w", command, string(out), err)
		}
	}

	return nil
}

func logf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Printf("[uinit] %v", message)
}
