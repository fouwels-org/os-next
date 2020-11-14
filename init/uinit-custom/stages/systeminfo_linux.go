// +build linux

package stages

import (
	"encoding/json"
	"log"
	"os/user"
	"uinit-custom/config"

	"github.com/zcalusic/sysinfo"
)

//Systeminfo implements IStage
type Systeminfo struct {
	finals []string
}

//String ..
func (n Systeminfo) String() string {
	return "System Info"
}

//Finalise ..
func (n Systeminfo) Finalise() []string {
	return n.finals
}

//Run ..
func (n Systeminfo) Run(c config.Config) error {
	current, err := user.Current()
	if err != nil {
		return err
	}

	if current.Uid != "0" {
		log.Fatal("requires superuser privilege")
	}

	var si sysinfo.SysInfo

	si.GetSysInfo()

	data, err := json.MarshalIndent(&si, "", "  ")
	if err != nil {
		return err
	}

	err = setFile("/tmp/info.json", string(data), 0644)
	if err != nil {
		return err
	}

	return nil
}
