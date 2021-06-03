// +build !linux

package util

import (
	"fmt"
)

//SetConsoleLogLevel sets the console level with syslog(2).
func (s *SystemUtil) SetConsoleLogLevel(level KLogLevel) error {
	return fmt.Errorf("setting console log level not supported on this platform")
}
