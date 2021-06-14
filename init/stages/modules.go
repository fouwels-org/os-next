package stages

import (
	"fmt"
	"init/config"
	"init/shell"
)

//Modules implementes IStage
type Modules struct {
	finals []string
}

//String ..
func (m *Modules) String() string {
	return "modules"
}

//Policy ..
func (m *Modules) Policy() Policy {
	return PolicySoft
}

//Finalise ..
func (m *Modules) Finalise() []string {
	return m.finals
}

//Run ..
func (m *Modules) Run(c config.Config) error {

	lok := 0
	errs := []error{}

	//Append secondary modules if the secondary config has been loaded
	modules := append(c.Primary.Modules, c.Secondary.Modules...)

	for _, v := range modules {

		com := []shell.Command{{Executable: shell.Modprobe, Arguments: []string{v}}}
		err := shell.Executor.Execute(com)

		if err != nil {
			errs = append(errs, err)
		} else {
			lok++
		}
	}

	m.finals = append(m.finals, fmt.Sprintf("loaded %v/%v modules ok", lok, len(c.Primary.Modules)))

	return fmt.Errorf("%v", errs)
}
