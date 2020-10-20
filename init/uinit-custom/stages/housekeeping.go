package stages

import (
	"bufio"
	"os"
	"uinit-custom/config"
)

//Networking implements IStage
type Housekeeping struct {
	finals []string
}

//String ..
func (n Housekeeping) String() string {
	return "House Keeping"
}

//Finalise ..
func (n Housekeeping) Finalise() []string {
	return n.finals
}

//Run ..
func (n Housekeeping) Run(c config.Config) error {

	err := writeLines("1")
	if err != nil {
		return err
	}

	return nil
}

func writeLines(line string) error {

	// overwrite file if it exists
	file, err := os.OpenFile("/sys/fs/cgroup/memory/memory.use_hierarchy", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logf("Error setting memory use_hierarchy:  " + err.Error())
		return err
	}

	defer file.Close()

	// new writer w/ default 4096 buffer size
	w := bufio.NewWriter(file)

	_, err = w.WriteString(line)
	if err != nil {
		logf("Error setting memory use_hierarchy:  " + err.Error())
		return err
	}

	// flush outstanding data
	return w.Flush()
}
