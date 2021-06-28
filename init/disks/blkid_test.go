// SPDX-FileCopyrightText: 2020 Lagoni Engineering
// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
//
// SPDX-License-Identifier: Apache-2.0

package disks

import (
	"log"
	"testing"
)

func TestGetBlkid(t *testing.T) {
	b, err := GetBlkid("")

	log.Printf("%+v", b[0])

	if err != nil {
		t.Fatalf("%v", err)
	}
}
