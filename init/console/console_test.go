// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package console

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {

	ok := checkPasswordHash("42be3e081457a3ff83372d810c6b84de70aeb57336e24f8715b7903e9ab8f1a2", "super-secure")
	if !ok {
		t.Fatalf("ok check failed")
	}

	nok := checkPasswordHash("42be3e081457a3ff83372d810c6b84de70aeb57336e24f8715b7903e9ab8f1a2", "supr-secure")
	if nok {
		t.Fatalf("not ok check failed")
	}

	none := checkPasswordHash("42be3e081457a3ff83372d810c6b84de70aeb57336e24f8715b7903e9ab8f1a2", "")
	if none {
		t.Fatalf("none check failed")
	}

	none2 := checkPasswordHash("", "super-secure")
	if none2 {
		t.Fatalf("none2 check failed")
	}

	none3 := checkPasswordHash("", "")
	if none3 {
		t.Fatalf("none3 check failed")
	}

}
