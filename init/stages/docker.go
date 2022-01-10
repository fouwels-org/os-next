// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package stages

import (
	"fmt"
	"os"
	"time"

	"github.com/fouwels/os-next/init/config"
	"github.com/fouwels/os-next/init/shell"
)

//Docker implementes IStage
type Docker struct {
	finals []string
}

//String ..
func (d *Docker) String() string {
	return "docker"
}

//Policy ..
func (d *Docker) Policy() Policy {
	return PolicyHard
}

//Finalise ..
func (d *Docker) Finalise() []string {
	return d.finals
}

//Run ..
func (d *Docker) Run(c config.Config) error {

	const _logpath string = "/var/log/docker/docker.log"

	// Start Docker
	// Set path to allow docker to find containerd
	command := shell.Command{
		Executable: shell.Dockerd,
		Arguments:  []string{},
		Env:        []string{"DOCKER_RAMDISK=true", "PATH=/sbin:/usr/sbin:/bin:/usr/bin"},
	}

	b, err := os.Create(_logpath)
	if err != nil {
		return fmt.Errorf("failed to create docker log at %v: %w", _logpath, err)
	}

	err = shell.Executor.ExecuteDaemon(command, b)
	if err != nil {
		return fmt.Errorf("failed to start dockerd: %w", err)
	}

	started := false
	for i := 0; i < 5; i++ {

		commands := []shell.Command{}
		commands = append(commands, shell.Command{Executable: shell.Docker, Arguments: []string{"version"}})

		err := shell.Executor.Execute(commands)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		started = true
	}

	if !started {
		return fmt.Errorf("failed to get docker version, docker did not start correctly")
	}
	return nil
}
