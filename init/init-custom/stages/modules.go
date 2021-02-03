package stages

import (
	"fmt"
	"init-custom/config"
	"init-custom/util"
)

//Modules implementes IStage
type Modules struct {
	finals []string
}

//String ..
func (m *Modules) String() string {
	return "modules"
}

//Finalise ..
func (m *Modules) Finalise() []string {
	return m.finals
}

//Run ..
func (m *Modules) Run(c config.Config) error {

	lok := 0
	errs := []error{}
	for _, v := range c.Primary.Modules {

		com := []util.Command{{Target: "/sbin/modprobe", Arguments: []string{v}}}
		err := util.Shell.Execute(com)

		if err != nil {
			errs = append(errs, err)
		} else {
			lok++
		}
	}

	m.finals = append(m.finals, fmt.Sprintf("loaded %v/%v modules ok", lok, len(c.Primary.Modules)))
	m.finals = append(m.finals, fmt.Sprintf("Errors: %v", errs))

	return nil
}
