package stages

import (
	"fmt"
	"init-custom/config"
	"os"
	"os/exec"
	"time"
)

//Docker implementes IStage
type Docker struct {
	finals []string
}

//String ..
func (d Docker) String() string {
	return "Docker"
}

//Finalise ..
func (d Docker) Finalise() []string {
	return d.finals
}

//Run ..
func (d Docker) Run(c config.Config) error {

	// Start Docker

	cmd := exec.Command("/usr/bin/dockerd")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "DOCKER_RAMDISK=true")
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Failed to start dockerd: %w", err)
	}

	response := ""
	for i := 0; i < 5; i++ {

		resp, err := executeOne(command{command: "/usr/bin/docker", arguments: []string{"version"}}, "")
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		response = resp
	}

	if response == "" {
		return fmt.Errorf("Failed to get docker version, docker did not start correctly")
	}

	d.finals = append(d.finals, fmt.Sprintf("Docker version: %v", response))

	return nil

}
