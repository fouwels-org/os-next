// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fouwels/os-next/init/journal"

	"gopkg.in/yaml.v2"
)

//LoadConfig ..
func LoadConfig(path string, config interface{}) (e error) {

	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	y := yaml.NewDecoder(f)
	err = y.Decode(config)
	if err != nil {
		return err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("could not seek: %w", err)
	}

	y = yaml.NewDecoder(f)
	y.SetStrict(true)
	err = y.Decode(config)
	if err != nil {
		journal.Logfln("Warning: %v", err)
	}

	return nil
}
