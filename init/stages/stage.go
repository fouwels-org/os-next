package stages

import "init/config"

//IStage ..
type IStage interface {
	//Called during the sequental stage init
	Run(config.Config) error
	//Friendly stage name
	String() string
	//Final strings, called after all stages have been initialized. (eg. to render the stage acquired DHCP IP, or generated SSH public key)
	Finalise() []string
}
