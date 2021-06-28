// +build linux

// Copyright 2014-2019 the u-root Authors. All rights reserved
// SPDX-FileCopyrightText: 2014 the u-root Authors
//
// SPDX-License-Identifier: BSD-3-Clause

package libinit

import (
	"fmt"
	"syscall"
)

type mount struct {
	Source string
	Target string
	FSType string
	Flags  uintptr
	Opts   string
}

func (m mount) create() error {
	return syscall.Mount(m.Source, m.Target, m.FSType, m.Flags, m.Opts)
}

func (m mount) String() string {
	return fmt.Sprintf("mount -t %q -o %s %q %q flags %#x", m.FSType, m.Opts, m.Source, m.Target, m.Flags)
}
