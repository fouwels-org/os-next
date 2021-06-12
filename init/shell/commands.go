package shell

const Login Executable = "/bin/login"
const Ntpd Executable = "/bin/ntpd"
const Modprobe Executable = "/bin/modprobe"
const Hwclock Executable = "/bin/hwclock"
const IP Executable = "/bin/ip"
const Udhcp Executable = "/bin/udhcpc"
const Dockerd Executable = "/bin/dockerd"
const Docker Executable = "/bin/docker"
const Mkdir Executable = "/bin/mkdir"
const Mount Executable = "/bin/mount"
const Ash Executable = "/bin/ash"

//IExecutable exists to force use of defined Excutable const, disable naked strings being acceptable as arguments to shell.Executor
type IExecutable interface {
	String() string
	Target() string
}

//Executable ..
type Executable string

//String ..
func (e Executable) String() string {
	return string(e)
}

//Target ..
func (e Executable) Target() string {
	return string(e)
}
