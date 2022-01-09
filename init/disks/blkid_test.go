// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package disks

import (
	"testing"

	"os-next/init/journal"
)

func TestGetBlkid(t *testing.T) {
	b, err := GetBlkid("")

	journal.Logfln("%+v", b[0])

	if err != nil {
		t.Fatalf("%v", err)
	}
}
