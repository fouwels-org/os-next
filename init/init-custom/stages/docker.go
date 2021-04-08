package stages

import (
	"fmt"
	"init-custom/config"
	"init-custom/util"
	"os"
	"time"
)

//Docker implementes IStage
type Docker struct {
	finals []string
}

//String ..
func (d *Docker) String() string {
	return "docker"
}

//Finalise ..
func (d *Docker) Finalise() []string {
	return d.finals
}

//Run ..
func (d *Docker) Run(c config.Config) error {

	const _logpath string = "/var/lib/docker/docker.log"

	// Start Docker
	// Set path to allow docker to find containerd
	command := util.Command{
		Target:    "/usr/bin/dockerd",
		Arguments: []string{},
		Env:       []string{"DOCKER_RAMDISK=true", "PATH=/sbin:/usr/sbin:/bin:/usr/bin"},
	}

	b, err := os.Create(_logpath)
	if err != nil {
		return fmt.Errorf("failed to create docker log at %v: %w", _logpath, err)
	}

	err = util.Shell.ExecuteDaemon(command, b)
	if err != nil {
		return fmt.Errorf("failed to start dockerd: %w", err)
	}

	started := false
	for i := 0; i < 5; i++ {

		commands := []util.Command{}
		commands = append(commands, util.Command{Target: "/usr/bin/docker", Arguments: []string{"version"}})

		err := util.Shell.Execute(commands)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		started = true
	}

	if !started {
		return fmt.Errorf("failed to get docker version, docker did not start correctly")
	}

	d.finals = append(d.finals, fmt.Sprintf("logging to %v", _logpath))
	return nil
}
