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

//Finalise ..
func (m *Modules) Finalise() []string {
	return m.finals
}

//Run ..
func (m *Modules) Run(c config.Config) error {

	lok := 0
	errs := []error{}
	for _, v := range c.Primary.Modules {

		com := []shell.Command{{Executable: shell.Modprobe, Arguments: []string{v}}}
		err := shell.Executor.Execute(com)

		if err != nil {
			errs = append(errs, err)
		} else {
			lok++
		}
	}

	m.finals = append(m.finals, fmt.Sprintf("loaded %v/%v modules ok", lok, len(c.Primary.Modules)))
	m.finals = append(m.finals, fmt.Sprintf("errors: %v", errs))

	return nil
}
