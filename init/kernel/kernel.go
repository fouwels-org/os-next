// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
//
// SPDX-License-Identifier: Apache-2.0

package kernel

import (
	"fmt"

	"golang.org/x/sys/unix"
)

//KLogLevel are the log levels used by printk.
type KLogLevel uintptr

//These are the log levels used by printk as described in syslog(2).
const (
	KLogEmergency KLogLevel = 0
	KLogAlert     KLogLevel = 1
	KLogCritical  KLogLevel = 2
	KLogError     KLogLevel = 3
	KLogWarning   KLogLevel = 4
	KLogNotice    KLogLevel = 5
	KLogInfo      KLogLevel = 6
	KLogDebug     KLogLevel = 7
)

//SetLogLevel sets the kernel level with
func SetLogLevel(level KLogLevel) error {
	if _, _, err := unix.Syscall(unix.SYS_SYSLOG, unix.SYSLOG_ACTION_CONSOLE_LEVEL, 0, uintptr(level)); err != 0 {
		return fmt.Errorf("failed to set kernel kog level to %d: %v", level, err)
	}
	return nil
}
