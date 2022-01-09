// SPDX-FileCopyrightText: 2017 Avi Deitcher
// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package diskfs

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"os-next/init/external/diskfs/table"

	"golang.org/x/sys/unix"
)

// Disk is a reference to a single disk block device or image that has been Create() or Open()
type Disk struct {
	File              *os.File
	Info              os.FileInfo
	Type              Type
	Size              int64
	LogicalBlocksize  int64
	PhysicalBlocksize int64
	DefaultBlocks     bool
}

// Open a Disk from a path to a device in read-write exclusive mode
// Should pass a path to a block device e.g. /dev/sda or a path to a file /tmp/foo.img
// The provided device must exist at the time you call Open()
func Open(device string) (Disk, error) {

	_, err := os.Stat(device)
	if os.IsNotExist(err) {
		return Disk{}, fmt.Errorf("provided device %v does not exist", device)
	} else if err != nil {
		return Disk{}, fmt.Errorf("unable to stat device %v: %w", device, err)
	}

	f, err := os.OpenFile(filepath.Clean(device), os.O_RDWR|os.O_EXCL, 0600)
	if err != nil {
		return Disk{}, fmt.Errorf("could not open device %s exclusively for writing", device)
	}

	// get device information
	devInfo, err := f.Stat()
	if err != nil {
		return Disk{}, fmt.Errorf("could not get info for device %s: %v", f.Name(), err)
	}
	mode := devInfo.Mode()

	var diskType Type
	var size int64

	lblksize := int64(_defaultBlocksize)
	pblksize := int64(_defaultBlocksize)
	defaultBlocks := true

	switch {
	case mode.IsRegular():

		diskType = File
		size = devInfo.Size()

		if size <= 0 {
			return Disk{}, fmt.Errorf("could not get file size for device %s", f.Name())
		}

	case mode&os.ModeDevice != 0:

		diskType = Device

		file, err := os.Open(f.Name())
		if err != nil {
			return Disk{}, fmt.Errorf("error opening block device %s: %s", f.Name(), err)
		}

		size, err = file.Seek(0, io.SeekEnd)
		if err != nil {
			return Disk{}, fmt.Errorf("error seeking to end of block device %s: %s", f.Name(), err)
		}

		lblksize, pblksize, err = getSectorSizes(f)

		defaultBlocks = false
		if err != nil {
			return Disk{}, fmt.Errorf("unable to get block sizes for device %s: %v", f.Name(), err)
		}

	default:

		return Disk{}, fmt.Errorf("device %s is neither a block device nor a regular file", f.Name())

	}

	ret := Disk{
		File:              f,
		Info:              devInfo,
		Type:              diskType,
		Size:              size,
		LogicalBlocksize:  lblksize,
		PhysicalBlocksize: pblksize,
		DefaultBlocks:     defaultBlocks,
	}

	return ret, nil
}

func (d *Disk) Close() error {
	err := d.File.Close()
	if err != nil {
		return fmt.Errorf("failed to close: %v", err)
	}

	return nil
}

// ReReadPartitionTable forces the kernel to re-read the partition table
// on the disk.
//
// It is done via an ioctl call with request as BLKRRPART.
func (d *Disk) ReReadPartitionTable() error {

	fd := d.File.Fd()
	_, err := unix.IoctlGetInt(int(fd), 0x125f) //BLKRRPART
	if err != nil {
		return fmt.Errorf("unable to re-read partition table: %v", err)
	}
	return nil
}

// GetTable retrieves a PartitionTable for a Disk
//
// returns an error if the Disk is invalid or does not exist, or the partition table is unknown
func (d *Disk) GetTable() (table.Table, error) {

	t, err := table.Read(d.File, int(d.LogicalBlocksize), int(d.PhysicalBlocksize))
	if err != nil {
		return table.Table{}, fmt.Errorf("GPT: %w", err)
	}

	return t, nil
}

// CheckMagicFlag checks if a specific flag is held in the MBR bootstrap code area.
// Magic flags are used check  disk state if the partition table cannot be read.
func (d *Disk) GetFlag() (uint32, error) {
	const _flagAddress = 0x00E0
	b := make([]byte, 16)

	_, err := d.File.ReadAt(b, _flagAddress)
	if err != nil {
		return 0x00, fmt.Errorf("failed to read address: %w", err)
	}
	return binary.LittleEndian.Uint32(b), nil
}

func (d *Disk) SetFlag(flag uint32) error {
	const _flagAddress = 0x00E0

	b := make([]byte, 16)
	binary.LittleEndian.PutUint32(b, flag)

	n, err := d.File.WriteAt(b, _flagAddress)
	if err != nil {
		return fmt.Errorf("failed to write address: %w", err)
	}
	if n != 16 {
		return fmt.Errorf("wrote %v bytes, expected %v", n, 16)
	}
	return nil
}

// SetTable applies a partition.Table implementation to a Disk
//
// The Table can have zero, one or more Partitions, each of which is unique to its
// implementation. E.g. MBR partitions in mbr.Table look different from GPT partitions in gpt.Table
//
// Actual writing of the table is delegated to the individual implementation
func (d *Disk) SetTable(table table.Table) error {

	err := table.Write(d.File, d.Size)
	if err != nil {
		return fmt.Errorf("failed to write partition table: %v", err)
	}

	// the partition table needs to be re-read only if
	// the disk file is an actual block device
	if d.Type == Device {
		err = d.ReReadPartitionTable()
		if err != nil {
			return fmt.Errorf("unable to re-read the partition table. Kernel still uses old partition table: %v", err)
		}
	}
	return nil
}

// getSectorSizes get the logical and physical sector sizes for a block device
func getSectorSizes(f *os.File) (int64, int64, error) {
	/*
		ioctl(fd, BLKPBSZGET, &physicalsectsize);

	*/
	fd := f.Fd()
	logicalSectorSize, err := unix.IoctlGetInt(int(fd), _blksszGet)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to get device logical sector size: %v", err)
	}
	physicalSectorSize, err := unix.IoctlGetInt(int(fd), _blkpbszGet)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to get device physical sector size: %v", err)
	}
	return int64(logicalSectorSize), int64(physicalSectorSize), nil
}
