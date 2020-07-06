package stages_test

import (
	"testing"
	"uinit-custom/config"
	"uinit-custom/stages"
)

func TestWireguard(t *testing.T) {
	wg := stages.Wireguard{}

	c, err := config.LoadConfig("../config.json")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	s, err := config.LoadSecrets("../secrets.json")
	if err != nil {
		t.Fatalf("Failed to load secrets: %v", err)
	}

	err = wg.Run(c, s)
	if err != nil {
		t.Fatalf("%v", err)
	}
}
