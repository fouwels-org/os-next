package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

//LoadConfig ..
func LoadConfig(path string, config interface{}) (e error) {

	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return err
	}

	defer func() {
		ferr := f.Close()
		if ferr != nil {
			e = fmt.Errorf("Failed to close file: %v", ferr)
		}
	}()

	jd := json.NewDecoder(f)
	err = jd.Decode(&config)
	if err != nil {
		return err
	}

	return nil
}
