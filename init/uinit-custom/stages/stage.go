package stages

import "uinit-custom/config"

//IStage ..
type IStage interface {
	Run(config.Config, config.Secrets) error
	String() string
}
