// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 K. Fouwels <k@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package config

//Config ..
type Config struct {
	Primary   Primary
	Secondary Secondary
}

//Primary ..
type Primary struct {
	Modules    []string
	Filesystem Filesystem
}

//Secondary ..
type Secondary struct {
	Modules    []string
	Networking Networking
	Wireguard  []Wireguard
	Time       Time
}

//Header ..
type Header struct {
	Site    string
	Comment string
}

//Networking ..
type Networking struct {
	Networks    []NetworkingNetwork
	Routes      []Route
	Nameservers []string
}

//NetworkingNetwork ..
type NetworkingNetwork struct {
	Device         string
	DHCP           bool
	IPV6           bool
	Type           string
	Addresses      []string
	DefaultGateway string `yaml:"default-gateway"`
}

//Filesystem ..
type Filesystem struct {
	Devices []FilesystemDevice
}

//FilesystemDevice ..
type FilesystemDevice struct {
	Label      string
	MountPoint string
	FileSystem string
}

//Wireguard ..
type Wireguard struct {
	Device     string
	ListenPort int `yaml:"listen-port"`
	Peers      []WireguardPeer
}

//WireguardPeer ..
type WireguardPeer struct {
	PublicKey           string `yaml:"public-key"`
	Endpoint            string
	AllowedIps          []string `yaml:"allowed-ips"`
	PersistentKeepalive int      `yaml:"persistent-keepalive"`
}

//Route ..
type Route struct {
	Address string
	Device  string
}

//Time ..
type Time struct {
	NTP     bool
	HWClock bool
	Servers []string
}
