package util

import (
	"fmt"
	"os"
	"path/filepath"
)

//File static instance of FileUtil{}
var File FileUtil = FileUtil{}

//FileUtil ..
type FileUtil struct {
}

//SetFile ..
func (c *FileUtil) SetFile(path string, value string, filemode os.FileMode) (e error) {

	// #nosec G304 (CWE-22).
	// N/A: filemode parameter is incorrectly being flagged.
	// This is intended to be variable, and does not represent a path traversal.
	f, err := os.OpenFile(filepath.Clean(path), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filemode)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}

	// #nosec G307. Double defer is safe for file.Writer
	defer f.Close()

	_, err = fmt.Fprintf(f, "%v", value)
	if err != nil {
		return fmt.Errorf("failed to write: %v", err)
	}

	err = f.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync on %v: %v", f.Name(), err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("failed to close on %v: %v", f.Name(), err)
	}

	return nil
}
