package console

import (
	"bufio"
	"fmt"
	"init/shell"
	"log"
	"os"
	"strings"
	"time"
)

const _authenticator = "LagoniP13"

func Start() error {

	err := login()
	if err != nil {
		return err
	}

	err = bash()
	if err != nil {
		return err
	}
	return nil
}

func login() error {

	success := false
	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read stdin: %w", err)
	}

	for !success {
		fmt.Printf("enter authenticator for shell\n> ")
		text, err := reader.ReadString('\n')

		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}

		if strings.TrimSuffix(text, "\n") == _authenticator {
			success = true
			log.Printf("user succeeded to authenticate")
		} else {
			fmt.Printf("authenticator incorrect\n")
			log.Printf("user failed to authenticate")
			time.Sleep(2 * time.Second)
		}
	}

	return nil
}

func bash() error {
	commands := []shell.Command{
		{Executable: shell.Ash, Arguments: []string{}},
	}

	err := shell.Executor.ExecuteInteractive(commands)
	if err != nil {
		return err
	}

	return nil
}
