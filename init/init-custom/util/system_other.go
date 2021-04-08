// +build !linux

package util

import (
	"fmt"
)

//StringInfo ..
func (s *SystemUtil) StringInfo() (string, error) {

	return "", fmt.Errorf("stringInfo is not supported on this platform")
}

//SetConsoleLogLevel sets the console level with syslog(2).
func (s *SystemUtil) SetConsoleLogLevel(level KLogLevel) error {
	return fmt.Errorf("setting console log level not supported on this platform")
}
