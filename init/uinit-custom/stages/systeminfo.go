package stages

import (
	"encoding/json"
	"log"
	"os/user"
	"uinit-custom/config"
	"github.com/zcalusic/sysinfo"
)

//Networking implements IStage
type Systeminfo struct {
	finals []string
}

//String ..
func (n Systeminfo) String() string {
	return "House Keeping"
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

	err = writeLines("/tmp/info.json",string(data))
	if err != nil {
		return err
	}

	return nil
}