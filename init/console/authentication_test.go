package console

import (
	"os-next/init/journal"
	"testing"
)

func TestGenerateAuthenticator(t *testing.T) {

	hash, err := generateAuthenticator("super-secure")
	if err != nil {
		t.Fatalf("%v", err)
	}

	journal.Logfln("%v", hash)
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

func TestGenerateTotp(t *testing.T) {

	secret, url, err := generateTotp()
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Logf("%+v", secret)
	t.Logf("%+v", url)
}

func TestCheckTotp(t *testing.T) {

	code := "296835"

	secret := "Q2GYK34AKCXJKLY47T4ZIZYJZWVIH2QLGIGZVMGX2XXXWTNYU7HQ"

	res, err := checkTotp(secret, code)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if res != true {
		t.Fatalf("check failed")
	}
}
