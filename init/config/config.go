// SPDX-FileCopyrightText: 2020 Lagoni Engineering
// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"encoding/json"
	"fmt"
	"log"
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
			e = fmt.Errorf("failed to close file: %v", ferr)
		}
	}()

	jd := json.NewDecoder(f)
	err = jd.Decode(&config)
	if err != nil {
		return err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("could not seek: %w", err)
	}

	jd.DisallowUnknownFields() // Force errors
	err = jd.Decode(&config)
	if err != nil {
		log.Printf("Warning: %v", err)
	}

	return nil
}
