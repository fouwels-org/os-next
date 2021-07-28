// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package console

import (
	"log"
	"testing"
)

func TestGenerateAuthenticator(t *testing.T) {

	hash, err := generateAuthenticator("super-secure")
	if err != nil {
		t.Fatalf("%v", err)
	}

	log.Printf("%v", hash)
}

func TestCheckAuthenticator(t *testing.T) {

	err := checkAuthenticator("JDJhJDEwJFpDamNuZzVGMGFOZ0NwYlRYUVdSWnVaUkh5WmdSUXIvOXFtbzYySGJ2dEFTbnIzTm9DazhT", "super-secure")
	if err != true {
		t.Fatalf("ok check failed")
	}

	err = checkAuthenticator("JDJhJDEwJFpDamNuZzVGMGFOZ0NwYlRYUVdSWnVaUkh5WmdSUXIvOXFtbzYySGJ2dEFTbnIzTm9DazhT", "supr-secure")
	if err != false {
		t.Fatalf("not ok check failed")
	}

	err = checkAuthenticator("JDJhJDEwJFpDamNuZzVGMGFOZ0NwYlRYUVdSWnVaUkh5WmdSUXIvOXFtbzYySGJ2dEFTbnIzTm9DazhT", "")
	if err != false {
		t.Fatalf("none check failed")
	}

	err = checkAuthenticator("", "super-secure")
	if err != false {
		t.Fatalf("none2 check failed")
	}

	err = checkAuthenticator("", "")
	if err != false {
		t.Fatalf("none3 check failed")
	}

}
