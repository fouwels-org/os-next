package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/go-playground/validator.v9"
)

//Header ..
type Header struct {
	Site    string `validate:"required"`
	Comment string
}

//Config ..
type Config struct {
	Header     Header     `validate:"required"`
	Modules    []string   `validate:"required"`
	Networking Networking `validate:"required"`
}

//Networking ..
type Networking struct {
	Networks []NetworkingNetwork `validate:"required,dive"`
}

//NetworkingNetwork ..
type NetworkingNetwork struct {
	Device         string `validate:"required"`
	DHCP           bool
	IPV6           bool
	Address        string `validate:"required_without=DHCP"`
	DefaultGateway string `validate:"required_with=Address"`
	SubnetMask     string `validate:"required_with=Address"`
}

//LoadConfig ..
func LoadConfig(path string) (c Config, e error) {

	c = Config{}
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return Config{}, err
	}

	defer func() {
		ferr := f.Close()
		if ferr != nil {
			e = fmt.Errorf("Failed to close file: %v", ferr)
		}
	}()

	jd := json.NewDecoder(f)
	err = jd.Decode(&c)
	if err != nil {
		return Config{}, err
	}

	validate := validator.New()

	err = validate.Struct(c)

	if err != nil {
		return Config{}, err
	}

	return c, nil
}
