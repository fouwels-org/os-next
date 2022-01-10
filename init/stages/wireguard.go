// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package stages

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"time"

	"github.com/fouwels/os-next/init/config"
	"github.com/fouwels/os-next/init/external/qrterminal"
	"github.com/fouwels/os-next/init/filesystem"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

//Wireguard implements IStage
type Wireguard struct {
	finals []string
}

//String ..
func (n *Wireguard) String() string {
	return "wireguard"
}

//Policy ..
func (n *Wireguard) Policy() Policy {
	return PolicyHard
}

//Finalise ..
func (n *Wireguard) Finalise() []string {
	return n.finals
}

//Run ..
func (n *Wireguard) Run(c config.Config) error {

	const _keyroot = "/var/config"

	for _, wloop := range c.Secondary.Wireguard {
		wg := wloop //Prevent loop ref capture

		keypath := fmt.Sprintf("%v/%v", _keyroot, wg.Device)

		wgc, err := wgctrl.New()
		if err != nil {
			return err
		}

		wgkey := wgtypes.Key{}

		skey, err := ioutil.ReadFile(filepath.Clean(keypath + ".private"))
		if err != nil {

			n.finals = append(n.finals, fmt.Sprintf("private key generated for %v", wg.Device))

			wgkey, err = wgtypes.GeneratePrivateKey()
			if err != nil {
				return fmt.Errorf("failed to generate private key: %v", err)
			}

			err = filesystem.WriteSync(filepath.Clean(keypath+".private"), []byte(wgkey.String()))
			if err != nil {
				return fmt.Errorf("failed to save wg key: %v", err)
			}

		} else {
			wgkey, err = wgtypes.ParseKey(string(skey))
			if err != nil {
				return fmt.Errorf("failed to parse loaded private key: %w", err)
			}

			n.finals = append(n.finals, "private key Loaded")
		}

		err = filesystem.WriteSync(filepath.Clean(keypath+".pub.qr"), []byte(n.writeQR(wgkey.PublicKey())))
		if err != nil {
			return fmt.Errorf("failed to write public key QR: %v", err)
		}
		n.finals = append(n.finals, fmt.Sprintf("public key written to %v", filepath.Clean(keypath+".qr")))

		wgpeers := []wgtypes.PeerConfig{}

		for _, v := range wg.Peers {

			vkey, err := wgtypes.ParseKey(v.PublicKey)
			if err != nil {
				return fmt.Errorf("failed to parse key for %v: %w", v.Endpoint, err)
			}

			vudp, err := net.ResolveUDPAddr("udp", v.Endpoint)
			if err != nil {
				return fmt.Errorf("failed to resolve endpoint for %v: %w", v.Endpoint, err)
			}

			keepalive := time.Duration(v.PersistentKeepalive) * time.Second

			ipnets := []net.IPNet{}

			for _, ap := range v.AllowedIps {
				_, ipnet, err := net.ParseCIDR(ap)
				if err != nil {
					return fmt.Errorf("failed to parse allowedIP for %v: %w", ap, err)
				}

				ipnets = append(ipnets, *ipnet)
			}

			wgp := wgtypes.PeerConfig{
				PublicKey:                   vkey,
				Endpoint:                    vudp,
				PersistentKeepaliveInterval: &keepalive,
				AllowedIPs:                  ipnets,
			}

			wgpeers = append(wgpeers, wgp)
		}

		cfg := wgtypes.Config{
			PrivateKey: &wgkey,
			ListenPort: &wg.ListenPort,
			Peers:      wgpeers,
		}

		err = wgc.ConfigureDevice(wg.Device, cfg)
		if err != nil {
			return fmt.Errorf("failed to configure wireguard device: %v", err)
		}
	}

	return nil
}
func (n *Wireguard) writeQR(publicKey wgtypes.Key) string {

	var buf bytes.Buffer

	config := qrterminal.Config{
		Level:      qrterminal.L,
		HalfBlocks: false,
		Writer:     &buf,
		BlackChar:  qrterminal.BLACK,
		WhiteChar:  qrterminal.WHITE,
		QuietZone:  1,
	}

	err := qrterminal.GenerateWithConfig(publicKey.String(), config)
	if err != nil {
		n.finals = append(n.finals, fmt.Sprintf("failed to generate QR: %v", err))
	}

	return buf.String()
}
