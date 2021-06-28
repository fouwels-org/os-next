// SPDX-FileCopyrightText: 2020 Lagoni Engineering
// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
//
// SPDX-License-Identifier: Apache-2.0

package shell

const Login Executable = "/bin/login"
const Ntpd Executable = "/sbin/ntpd"
const Modprobe Executable = "/sbin/modprobe"
const Hwclock Executable = "/sbin/hwclock"
const IP Executable = "/sbin/ip"
const Udhcp Executable = "/sbin/udhcpc"
const Dockerd Executable = "/usr/bin/dockerd"
const Docker Executable = "/usr/bin/docker"
const Mkdir Executable = "/bin/mkdir"
const Mount Executable = "/bin/mount"
const Ash Executable = "/bin/ash"
const Blkid Executable = "/sbin/blkid"
const Mkdosfs Executable = "/sbin/mkdosfs"
const Lsblk Executable = "/bin/lsblk"
const Mke2fs Executable = "/bin/mke2fs"

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
