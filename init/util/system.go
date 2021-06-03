package util

//SystemUtil ..
type SystemUtil struct {
}

//System static instance of SystemUtil
var System = SystemUtil{}

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
