// +build !linux

package util

import (
	"fmt"
)

//StringInfo ..
func (s *SystemUtil) StringInfo() (string, error) {

	return "", fmt.Errorf("StringInfo is not supported on this platform")
}


//SetConsoleLogLevel sets the console level with syslog(2).
func (s *SystemUtil) SetConsoleLogLevel(level KLogLevel) error {
	return fmt.Errorf("Setting console log level not supported on this platform")
}