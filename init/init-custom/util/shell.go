package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

//Shell static instance of ShellUtil{}
var Shell ShellUtil = ShellUtil{}

//ShellUtil ..
type ShellUtil struct {
}

//ExecuteOne ..
func (s *ShellUtil) executeOne(c Command) (string, error) {

	// #nosec G204 (CWE-78).
	// N/A: subprocesses are executed dynamically by design, based on fixed
	// configuration within calling functions.
	cmd := exec.Command(c.Target, c.Arguments...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

//Execute ..
func (s *ShellUtil) Execute(commands []Command) error {
	for _, c := range commands {
		out, err := s.executeOne(c)
		if err != nil {
			return fmt.Errorf("%v failed: %v %w", c, string(out), err)
		}
	}

	return nil
}

//ExecuteDaemon execute in daemon mode
func (s *ShellUtil) ExecuteDaemon(c Command, writer io.Writer) error {

	// #nosec G204 (CWE-78).
	// N/A: subprocesses are executed dynamically by design, based on fixed
	// configuration within calling functions.
	cmd := exec.Command(c.Target, c.Arguments...)

	cmd.Env = append(cmd.Env, c.Env...)

	cmd.Stderr = writer
	cmd.Stdout = writer

	go func(cmd *exec.Cmd, writer io.Writer) {
		for {
			err := cmd.Run()
			_, err = writer.Write([]byte(fmt.Sprintf("Daemon exit with err! Restarting after 5s!: %v\n", err)))
			log.Printf("Failed to write to writer for daemon %v: %v", c, err)
			time.Sleep(5 * time.Second)
		}
	}(cmd, writer)

	return nil
}

//ExecuteInteractive Execute with attached TTY
func (s *ShellUtil) ExecuteInteractive(commands []Command) error {
	for _, c := range commands {
		err := s.executeOneInteractive(c)
		if err != nil {
			return fmt.Errorf("%v failed: %w", c, err)
		}
	}

	return nil
}

func (s *ShellUtil) executeOneInteractive(c Command) error {

	// #nosec G204 (CWE-78).
	// N/A: subprocesses are executed dynamically by design, based on fixed
	// configuration within calling functions.
	cmd := exec.Command(c.Target, c.Arguments...)

	cmd.Env = append(cmd.Env, c.Env...)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil

}

//Command ..
type Command struct {
	Target    string
	Arguments []string
	Env       []string
}

func (c Command) String() string {

	combo := []string{}
	combo = append(combo, c.Target)
	combo = append(combo, c.Arguments...)
	return fmt.Sprintf("%v", combo)
}
