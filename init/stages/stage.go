package stages

import "init/config"

type Policy int

const (
	PolicyHard Policy = iota // Hard fail, failure of stage aborts bootup
	PolicySoft               // Soft fail, failure of stage continued boot
)

//IStage ..
type IStage interface {
	//Called during the sequental stage init
	Run(config.Config) error
	//Friendly stage name
	String() string
	//Final strings, called after all stages have been initialized. (eg. to render the stage acquired DHCP IP, or generated SSH public key)
	Finalise() []string
	//Policy
	Policy() Policy
}
