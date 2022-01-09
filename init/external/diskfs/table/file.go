// SPDX-FileCopyrightText: 2017 Avi Deitcher
// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package table

import "io"

// File interface that can be read from and written to.
// Normally implemented as actual os.File, but useful as a separate interface so can easily
// use alternate implementations.
type File interface {
	io.ReaderAt
	io.WriterAt
	io.Seeker
}
