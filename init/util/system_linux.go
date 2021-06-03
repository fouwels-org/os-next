// +build linux

package util

import (
	"fmt"

	"golang.org/x/sys/unix"
)

//SetConsoleLogLevel sets the console level with syslog(2).
func (s *SystemUtil) SetConsoleLogLevel(level KLogLevel) error {
	if _, _, err := unix.Syscall(unix.SYS_SYSLOG, unix.SYSLOG_ACTION_CONSOLE_LEVEL, 0, uintptr(level)); err != 0 {
		return fmt.Errorf("failed to set kernel kog level to %d: %v", level, err)
	}
	return nil
}
