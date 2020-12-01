package stages

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func executeOne(c command, stdin string) (string, error) {

	buffer := bytes.Buffer{}
	buffer.Write([]byte(stdin))

	// #nosec G204 (CWE-78).
	// N/A: subprocesses are executed dynamically by design, based on fixed
	// configuration within calling functions.
	cmd := exec.Command(c.command, c.arguments...)
	cmd.Stdin = &buffer
	out, err := cmd.CombinedOutput()

	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

func execute(commands []command) error {
	for _, c := range commands {
		out, err := executeOne(c, "")
		if err != nil {
			return fmt.Errorf("%v failed: %v %w", c, string(out), err)
		}
	}

	return nil
}

func logf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Printf("[uinit] %v", message)
}

type command struct {
	command   string
	arguments []string
}

func (c command) String() string {

	combo := []string{}
	combo = append(combo, c.command)
	combo = append(combo, c.arguments...)
	return fmt.Sprintf("%v", combo)
}

func setFile(path string, value string, filemode os.FileMode) (e error) {

	// #nosec G304 (CWE-22).
	// N/A: filemode parameter is incorrectly being flagged.
	// This is intended to be variable, and does not represent a path traversal.
	f, err := os.OpenFile(filepath.Clean(path), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filemode)
	if err != nil {
		return fmt.Errorf("Failed to open: %w", err)
	}

	// #nosec G307. Double defer is safe for file.Writer
	defer f.Close()

	_, err = fmt.Fprintf(f, "%v", value)
	if err != nil {
		return fmt.Errorf("Failed to write: %v", err)
	}

	err = f.Sync()
	if err != nil {
		return fmt.Errorf("Failed to sync on %v: %v", f.Name(), err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("Failed to close on %v: %v", f.Name(), err)
	}

	return nil
}
