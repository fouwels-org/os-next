package stages

import (
	"fmt"
	"net"
	"time"
	"uinit-custom/config"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"golang.zx2c4.com/wireguard/wgctrl"
)

//Wireguard implements IStage
type Wireguard struct {
	finals []string
}

//String ..
func (w Wireguard) String() string {
	return "Wireguard"
}

//Finalise ..
func (w Wireguard) Finalise() []string {

	return w.finals
}

//Run ..
func (w *Wireguard) Run(c config.Config, s config.Secrets) error {

	commands := []string{}

	// Create WG
	wg, err := wgctrl.New()
	if err != nil {
		return err
	}

	for _, wr := range c.Networking.Wireguard {
		// Create Network Devices
		commands = append(commands, fmt.Sprintf("ip link add dev %v type wireguard", wr.Device))
		commands = append(commands, fmt.Sprintf("ip address add dev %v %v", wr.Device, wr.Address))

		err = execute(commands)
		if err != nil {
			return err
		}

		// Create Keys
		key, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return err
		}

		wgPeers := []wgtypes.PeerConfig{}

		for _, p := range wr.Peers {

			// parse endpoint
			udpnet, err := net.ResolveUDPAddr("udp", p.Endpoint)
			if err != nil {
				return err
			}

			// parse allowedIPS
			netips := []net.IPNet{}
			for _, n := range p.AllowedIPs {
				_, ipnet, err := net.ParseCIDR(n)
				if err != nil {
					return err
				}
				netips = append(netips, *ipnet)
			}

			// parse keepalive
			keepAlive := 5 * time.Second

			// parse public key
			pubkey, err := wgtypes.ParseKey(p.PublicKey)
			if err != nil {
				return err
			}

			// generate peer
			wgPeers = append(wgPeers, wgtypes.PeerConfig{
				PublicKey:                   pubkey,
				Endpoint:                    udpnet,
				PersistentKeepaliveInterval: &keepAlive,
				ReplaceAllowedIPs:           true,
				AllowedIPs:                  netips,
			})
		}

		// configure device
		wg.ConfigureDevice(wr.Device, wgtypes.Config{
			PrivateKey:   &key,
			ReplacePeers: true,
			Peers:        wgPeers,
		})

		w.finals = append(w.finals, fmt.Sprintf("Interface %v added with public key %v", wr.Device, key.PublicKey()))

	}

	return nil
}
