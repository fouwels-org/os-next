package config_test

import (
	"testing"
	"uinit-custom/config"
)

func TestLoads(t *testing.T) {

	_, err := config.LoadConfig("../config.json")
	if err != nil {
		t.Fatalf("%v", err)
	}

	_, err = config.LoadSecrets("../secrets.json")
	if err != nil {
		t.Fatalf("%v", err)
	}
}
