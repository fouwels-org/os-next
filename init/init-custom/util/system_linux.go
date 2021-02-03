// +build linux

package util

import (
	"encoding/json"
	"fmt"
	"os/user"
	"golang.org/x/sys/unix"
	"github.com/zcalusic/sysinfo"
)

//StringInfo ..
func (s *SystemUtil) StringInfo() (string, error) {

	current, err := user.Current()
	if err != nil {
		return "", err
	}

	if current.Uid != "0" {
		return "", fmt.Errorf("Requires superuser privilege")
	}

	var si sysinfo.SysInfo

	si.GetSysInfo()

	data, err := json.MarshalIndent(&si, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

//SetConsoleLogLevel sets the console level with syslog(2).
func (s *SystemUtil) SetConsoleLogLevel(level KLogLevel) error {
	if _, _, err := unix.Syscall(unix.SYS_SYSLOG, unix.SYSLOG_ACTION_CONSOLE_LEVEL, 0, uintptr(level)); err != 0 {
		return fmt.Errorf("failed to set kernel kog level to %d: %v", level, err)
	}
	return nil
}
