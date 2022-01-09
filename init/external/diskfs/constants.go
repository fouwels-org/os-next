// SPDX-FileCopyrightText: 2017 Avi Deitcher
// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package diskfs

// when we use a disk image with a GPT, we cannot get the logical sector size from the disk via the kernel
//    so we use the default sector size of 512, per Rod Smith
const (
	_defaultBlocksize = 512
	_firstblock       = 2048
	_blksszGet        = 0x1268
	_blkpbszGet       = 0x127b
)

// Type represents the type of disk this is
type Type int

const (
	// File is a file-based disk image
	File Type = iota
	// Device is an OS-managed block device
	Device
)
