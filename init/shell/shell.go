// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/fouwels/os-next/init/journal"
)

//Shell static instance of ShellUtil{}
var Executor ShellExecutor = ShellExecutor{}

//Command ..
type Command struct {
	Executable IExecutable
	Arguments  []string
	Env        []string
}

//String ..
func (c Command) String() string {

	combo := []string{}
	combo = append(combo, c.Executable.String())
	combo = append(combo, c.Arguments...)
	return fmt.Sprintf("%v", combo)
}

//ShellExecutor ..
type ShellExecutor struct {
}

func (s *ShellExecutor) ExecuteOne(c Command) (string, error) {

	cmd := exec.Command(c.Executable.Target(), c.Arguments...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

//Execute ..
func (s *ShellExecutor) Execute(commands []Command) error {
	for _, c := range commands {
		out, err := s.ExecuteOne(c)
		if err != nil {
			return fmt.Errorf("%v: %v (%v)", c, string(out), err)
		}
	}

	return nil
}

//ExecuteDaemon execute in daemon mode
func (s *ShellExecutor) ExecuteDaemon(c Command, writer io.Writer) error {

	cmd := exec.Command(c.Executable.Target(), c.Arguments...)

	cmd.Env = append(cmd.Env, c.Env...)

	cmd.Stderr = writer
	cmd.Stdout = writer

	go func(cmd *exec.Cmd, writer io.Writer) {
		for {
			err := cmd.Run()
			_, err = writer.Write([]byte(fmt.Sprintf("**daemon exit with err - restarting**: %v\n", err)))
			if err != nil {
				journal.Logfln("failed to write to writer for daemon %v: %v", c, err)
			}
			time.Sleep(5 * time.Second)
		}
	}(cmd, writer)

	return nil
}

//ExecuteInteractive Execute with attached TTY
func (s *ShellExecutor) ExecuteInteractive(commands []Command) error {
	for _, c := range commands {
		err := s.ExecuteInteractiveOne(c)
		if err != nil {
			return fmt.Errorf("%v failed: %w", c, err)
		}
	}

	return nil
}

func (s *ShellExecutor) ExecuteInteractiveOne(c Command) error {

	cmd := exec.Command(c.Executable.Target(), c.Arguments...)

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
