// +build !linux

package stages

import (
	"init-custom/config"
)

//Systeminfo implements IStage
type Systeminfo struct {
	finals []string
}

//String ..
func (n *Systeminfo) String() string {
	return "System Info"
}

//Finalise ..
func (n *Systeminfo) Finalise() []string {
	return n.finals
}

//Run ..
func (n *Systeminfo) Run(c config.Config) error {
	strerr := "zcalusic/sysinfo is not supported on this platform"
	n.finals = append(n.finals, strerr)
	return nil
}
