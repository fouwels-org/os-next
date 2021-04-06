package stages

import (
	"fmt"
	"init-custom/config"
	"init-custom/util"
)

//Systeminfo implements IStage
type Systeminfo struct {
	finals []string
}

//String ..
func (n *Systeminfo) String() string {
	return "system info"
}

//Finalise ..
func (n *Systeminfo) Finalise() []string {
	return n.finals
}

//Run ..
func (n *Systeminfo) Run(c config.Config) error {

	data, err := util.System.StringInfo()
	if err != nil {
		return fmt.Errorf("Failed to get system info: %w", err)
	}

	err = util.File.SetFile("/tmp/info.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}
