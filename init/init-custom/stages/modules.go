package stages

import (
	"init-custom/config"
)

//Modules implementes IStage
type Modules struct {
	finals []string
}

//String ..
func (m Modules) String() string {
	return "Modules"
}

//Finalise ..
func (m Modules) Finalise() []string {
	return m.finals
}

//Run ..
func (m Modules) Run(c config.Config) error {

	for _, v := range c.Primary.Modules {
		com := command{command: "/sbin/modprobe", arguments: []string{v}}
		_, err := executeOne(com, "")

		logf("Executing Command: %v", com)
		if err != nil {
			logf("Command failed: %v ", err)
			continue
		}
	}
	return nil
}
