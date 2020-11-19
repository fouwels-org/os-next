package config_test

import (
	"init-custom/config"
	"testing"
)

func TestLoads(t *testing.T) {

	_, err := config.LoadConfig("../config.json")
	if err != nil {
		t.Fatalf("%v", err)
	}
}
