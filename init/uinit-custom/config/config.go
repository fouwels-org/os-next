package config

import (
	"encoding/json"
	"os"

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
	Networks  []NetworkingNetwork   `validate:"required,dive"`
	Wireguard []NetworkingWireguard `validate:"required,dive"`
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

//NetworkingWireguard ..
type NetworkingWireguard struct {
	Device  string                    `validate:"required"`
	Address string                    `validate:"required"`
	Peers   []NetworkingWireguardPeer `validate:"required"`
}

//NetworkingWireguardPeer ..
type NetworkingWireguardPeer struct {
	Endpoint   string   `validate:"required"`
	PublicKey  string   `validate:"required"`
	AllowedIPs []string `validate:"required"`
}

//Secrets ..
type Secrets struct {
	Header Header `validate:"required"`
}

//LoadConfig ..
func LoadConfig(path string) (Config, error) {

	c := Config{}

	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

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

//LoadSecrets ..
func LoadSecrets(path string) (Secrets, error) {

	s := Secrets{}

	f, err := os.Open(path)
	if err != nil {
		return Secrets{}, err
	}
	defer f.Close()

	jd := json.NewDecoder(f)
	err = jd.Decode(&s)
	if err != nil {
		return Secrets{}, err
	}

	validate := validator.New()

	err = validate.Struct(s)
	if err != nil {
		return Secrets{}, err
	}
	return s, nil
}
