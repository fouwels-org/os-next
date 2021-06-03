// +build !linux

// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libinit

import "fmt"

type mount struct {
	Source string
	Target string
	FSType string
	Flags  uintptr
	Opts   string
}

func (m mount) create() error {
	return fmt.Errorf("creating mounts is not supported on this platform")
}

func (m mount) String() string {
	return fmt.Sprintf("mount -t %q -o %s %q %q flags %#x", m.FSType, m.Opts, m.Source, m.Target, m.Flags)
}
