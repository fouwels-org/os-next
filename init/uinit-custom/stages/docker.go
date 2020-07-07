package stages

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"uinit-custom/config"
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

	cmd := exec.Command("dockerd")
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr

	stdoutp, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	cmd.Env = append(cmd.Env, "DOCKER_RAMDISK")
	cmd.Start()

	scanner := bufio.NewScanner(stdoutp)
	line := ""

	for scanner.Scan() {
		line = scanner.Text()
		fmt.Printf("%v", line)
		if strings.Contains(line, "API listen on /var/run/docker.sock") {

			resp, err := executeOne("docker version", "")
			if err != nil {
				return fmt.Errorf("docker version: %v", err)
			}
			d.finals = append(d.finals, fmt.Sprintf("Docker version: %v", resp))

			return nil
		}
	}
	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("Error while scanning: %v", err)
	}

	return fmt.Errorf("Docker exit before running status detected")

}
