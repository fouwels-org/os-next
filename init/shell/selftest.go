// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package shell

import (
	"fmt"
	"os"
)

//SelfTest tests defined executables to ensure they exist
func SelfTest() error {
	es := []Executable{
		Sntp,
		Modprobe,
		Hwclock,
		DHCP,
		Dockerd,
		Docker,
		Mount,
		Blkid,
	}

	for _, e := range es {
		err := testOne(e)
		if err != nil {
			return err
		}
	}
	return nil
}

func testOne(e Executable) error {
	_, err := os.Stat(e.Target())
	if err != nil {
		return fmt.Errorf("failed to stat %v: %w", e, err)
	}

	return nil
}
