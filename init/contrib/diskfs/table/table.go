// SPDX-FileCopyrightText: 2017 Avi Deitcher
// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
//
// SPDX-License-Identifier: MIT

package table

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"

	"github.com/google/uuid"
)

// Table represents a partition table to be applied to a disk or read from a disk
type Table struct {
	Partitions         []Partition // slice of Partition
	LogicalSectorSize  int         // logical size of a sector
	PhysicalSectorSize int         // physical size of the sector
	GUID               uuid.UUID   // disk GUID
	partitionArraySize int         // how many entries are in the partition array size
	partitionEntrySize uint32      // size of the partition entry in the table, usually 128 bytes
	primaryHeader      uint64      // LBA of primary header, always 1
	secondaryHeader    uint64      // LBA of secondary header, always last sectors on disk
	firstDataSector    uint64      // LBA of first data sector
	lastDataSector     uint64      // LBA of last data sector
	diskSize           uint64      // The size of the host disk
}

func NewTable(diskSize uint64, GUID uuid.UUID) Table {

	// how many sectors on the disk?
	diskSectors := uint64(diskSize) / uint64(_logicalSectorSize)

	// how many sectors used for partition entries?
	partSectors := uint64(_partitionArraySize) * uint64(_partitionEntrySize) / uint64(_logicalSectorSize)

	t := Table{
		LogicalSectorSize:  _logicalSectorSize,
		PhysicalSectorSize: _physicalSectorSize,
		primaryHeader:      1,
		GUID:               GUID,
		partitionArraySize: _partitionArraySize,
		partitionEntrySize: _partitionEntrySize,
		firstDataSector:    2 + partSectors,
		secondaryHeader:    diskSectors - 1,
		lastDataSector:     diskSectors - 1 - partSectors,
		diskSize:           diskSize,
		Partitions:         make([]Partition, _partitionArraySize),
	}

	for _, v := range t.Partitions {
		v.GUID = GPTUnused
		v.Type = GPTUnused
	}

	return t
}

// Write writes a GPT to disk
// Must be passed the util.File to which to write and the size of the disk
func (t *Table) Write(f File, size int64) error {

	// write the protective MBR
	// write the primary GPT header
	// write the primary partition array
	// write the secondary partition array
	// write the secondary GPT header
	var written int
	var err error

	// check table
	for k, v := range t.Partitions {
		if v.Start%_logicalSectorSize != 0 {
			return fmt.Errorf("partition %v start (%v) is not a multiple of logical sector size (%v)", k, v.Start, _logicalSectorSize)
		}
		if v.Size%_logicalSectorSize != 0 {
			return fmt.Errorf("partition %v size (%v) is not a multiple of logical sector size (%v)", k, v.Size, _logicalSectorSize)
		}
	}

	protectiveMBR := t.generateProtectiveMBR()
	written, err = f.WriteAt(protectiveMBR, 0)
	if err != nil {
		return fmt.Errorf("error writing protective MBR to disk: %v", err)
	}
	if written != len(protectiveMBR) {
		return fmt.Errorf("wrote %d bytes of protective MBR instead of %d", written, len(protectiveMBR))
	}

	primaryHeader, err := t.encodeHeaderBytes(true)
	if err != nil {
		return fmt.Errorf("encoding primary GPT header: %v", err)
	}
	written, err = f.WriteAt(primaryHeader, int64(t.LogicalSectorSize))
	if err != nil {
		return fmt.Errorf("writing primary GPT: %v", err)
	}
	if written != len(primaryHeader) {
		return fmt.Errorf("%d bytes of primary GPT header instead of %d", written, len(primaryHeader))
	}

	partitionArray, err := t.encodePartitionBytes()
	if err != nil {
		return fmt.Errorf("encoding primary GPT partitions: %v", err)
	}
	written, err = f.WriteAt(partitionArray, int64(t.LogicalSectorSize*int(t.partitionArraySector(true))))
	if err != nil {
		return fmt.Errorf("writing primary partitions to disk: %v", err)
	}
	if written != len(partitionArray) {
		return fmt.Errorf("wrote %d bytes of primary partition array instead of %d", written, len(primaryHeader))
	}

	written, err = f.WriteAt(partitionArray, int64(t.LogicalSectorSize*int(t.partitionArraySector(false))))
	if err != nil {
		return fmt.Errorf("writing secondary partitions to disk: %v", err)
	}
	if written != len(partitionArray) {
		return fmt.Errorf("wrote %d bytes of secondary partition array instead of %d", written, len(primaryHeader))
	}

	secondaryHeader, err := t.encodeHeaderBytes(false)
	if err != nil {
		return fmt.Errorf("encoding secondary GPT header: %v", err)
	}
	written, err = f.WriteAt(secondaryHeader, int64(t.secondaryHeader)*int64(t.LogicalSectorSize))
	if err != nil {
		return fmt.Errorf("writing secondary GPT to disk: %v", err)
	}
	if written != len(secondaryHeader) {
		return fmt.Errorf("wrote %d bytes of secondary GPT header instead of %d", written, len(secondaryHeader))
	}

	return nil
}

// Read a partition table from a disk
// must be passed the util.File from which to read, and the logical and physical block sizes
//
// if successful, returns a gpt.Table struct
// returns errors if fails at any stage reading the disk or processing the bytes on disk as a GPT
func Read(f File, logicalBlockSize, physicalBlockSize int) (Table, error) {

	// read the data off of the disk
	b := make([]byte, _gptSize+logicalBlockSize*2)

	read, err := f.ReadAt(b, 0)

	if err != nil {
		return Table{}, fmt.Errorf("error reading device: %v", err)
	}
	if read != len(b) {
		return Table{}, fmt.Errorf("read only %d bytes of GPT from device instead of expected %d", read, len(b))
	}
	return decodeTableBytes(b, logicalBlockSize, physicalBlockSize)
}

func (t *Table) generateProtectiveMBR() []byte {

	b := make([]byte, 512)

	// we don't do anything to the first 446 bytes

	// Add MBR signature
	copy(b[510:], []byte{0x55, 0xaa})

	// create the single all disk partition
	parts := b[446 : 446+16]

	parts[0] = 0x00 // non-bootable
	parts[1] = 0x00 // ignore CHS entirely
	parts[2] = 0x00 // ignore CHS entirely
	parts[3] = 0x00 // ignore CHS entirely
	parts[4] = 0xee // partition type 0xee
	parts[5] = 0x00 // ignore CHS entirely
	parts[6] = 0x00 // ignore CHS entirely
	parts[7] = 0x00 // ignore CHS entirely

	// start LBA 1
	binary.LittleEndian.PutUint32(parts[8:12], 1)

	// end LBA
	endSector := t.secondaryHeader
	binary.LittleEndian.PutUint32(parts[12:16], uint32(endSector))

	return b
}

// partitionArraySector get the sector that holds the primary or secondary partition array
func (t *Table) partitionArraySector(primary bool) uint64 {
	if primary {
		return t.primaryHeader + 1
	}
	return t.secondaryHeader - uint64(t.partitionArraySize)*uint64(t.partitionEntrySize)/uint64(t.LogicalSectorSize)
}

// generatePartitionBytes write the bytes for the partition array
func (t *Table) encodePartitionBytes() ([]byte, error) {

	// generate the partition bytes
	partSize := t.partitionEntrySize * uint32(t.partitionArraySize)
	bpart := make([]byte, partSize)
	for i, p := range t.Partitions {

		// write the primary partition entry
		b2, err := p.encodeBytes(t.diskSize)
		if err != nil {
			return nil, fmt.Errorf("encoding partition %d: %w", i, err)
		}

		slotStart := i * int(t.partitionEntrySize)
		slotEnd := slotStart + int(t.partitionEntrySize)

		copy(bpart[slotStart:slotEnd], b2)
	}
	return bpart, nil
}

// encodeHeaderBytes write just the gpt header to bytes
func (t *Table) encodeHeaderBytes(primary bool) ([]byte, error) {
	b := make([]byte, t.LogicalSectorSize)

	// 8 bytes "EFI PART" signature - endianness on this?
	copy(b[0:8], _efiSignature)
	// 4 bytes revision 1.0
	copy(b[8:12], _efiRevision)
	// 4 bytes header size
	copy(b[12:16], _efiHeaderSize)
	// 4 bytes CRC32/zlib of header with this field zeroed out - must calculate then come back
	copy(b[16:20], []byte{0x00, 0x00, 0x00, 0x00})
	// 4 bytes zeroes reserved
	copy(b[20:24], _efiZeroes)

	// which LBA are we?
	if primary {
		binary.LittleEndian.PutUint64(b[24:32], t.primaryHeader)
		binary.LittleEndian.PutUint64(b[32:40], t.secondaryHeader)
	} else {
		binary.LittleEndian.PutUint64(b[24:32], t.secondaryHeader)
		binary.LittleEndian.PutUint64(b[32:40], t.primaryHeader)
	}

	// usable LBAs for partitions
	binary.LittleEndian.PutUint64(b[40:48], t.firstDataSector)
	binary.LittleEndian.PutUint64(b[48:56], t.lastDataSector)

	// 16 bytes disk GUID
	guid := t.GUID
	if guid == GPTUnused {
		guid = uuid.New()
	}
	copy(b[56:72], bytesToUUIDBytes(guid[0:16]))

	// starting LBA of array of partition entries
	binary.LittleEndian.PutUint64(b[72:80], t.partitionArraySector(primary))

	// how many entries?
	binary.LittleEndian.PutUint32(b[80:84], uint32(t.partitionArraySize))
	// how big is a single entry?
	binary.LittleEndian.PutUint32(b[84:88], 0x80)

	// we need a CRC/zlib of the partition entries, so we do those first, then append the bytes
	bpart, err := t.encodePartitionBytes()
	if err != nil {
		return nil, fmt.Errorf("encoding partitions: %w", err)
	}
	checksum := crc32.ChecksumIEEE(bpart)
	binary.LittleEndian.PutUint32(b[88:92], checksum)

	// calculate checksum of entire header and place 4 bytes of offset 16 = 0x10
	checksum = crc32.ChecksumIEEE(b[0:92])
	binary.LittleEndian.PutUint32(b[16:20], checksum)

	// zeroes to the end of the sector
	for i := 92; i < t.LogicalSectorSize; i++ {
		b[i] = 0x00
	}

	return b, nil
}

// decodeTableBytes read a partition table from a byte slice
func decodeTableBytes(b []byte, logicalBlockSize, physicalBlockSize int) (Table, error) {

	// minimum size - gpt entries + header + LBA0 for (protective) MBR
	minSize := _gptSize + logicalBlockSize*2
	if len(b) < minSize {
		return Table{}, fmt.Errorf("data for partition was %d bytes instead of expected minimum %d", len(b), minSize)
	}

	// GPT starts at LBA1
	gpt := b[logicalBlockSize:]
	// start with fixed headers
	efiSignature := gpt[0:8]
	efiRevision := gpt[8:12]
	efiHeaderSize := gpt[12:16]
	efiHeaderCrcBytes := append(make([]byte, 0, 4), gpt[16:20]...)
	efiHeaderCrc := binary.LittleEndian.Uint32(efiHeaderCrcBytes)
	efiZeroes := gpt[20:24]
	primaryHeader := binary.LittleEndian.Uint64(gpt[24:32])
	secondaryHeader := binary.LittleEndian.Uint64(gpt[32:40])
	firstDataSector := binary.LittleEndian.Uint64(gpt[40:48])
	lastDataSector := binary.LittleEndian.Uint64(gpt[48:56])
	diskGUID, err := uuid.FromBytes(bytesToUUIDBytes(gpt[56:72]))
	if err != nil {
		return Table{}, fmt.Errorf("unable to read guid from disk: %v", err)
	}
	partitionEntryFirstLBA := binary.LittleEndian.Uint64(gpt[72:80])
	partitionEntryCount := binary.LittleEndian.Uint32(gpt[80:84])
	partitionEntrySize := binary.LittleEndian.Uint32(gpt[84:88])
	partitionEntryChecksum := binary.LittleEndian.Uint32(gpt[88:92])

	// once we have the header CRC, zero it out
	copy(gpt[16:20], []byte{0x00, 0x00, 0x00, 0x00})
	if !bytes.Equal(efiSignature, _efiSignature) {
		return Table{}, fmt.Errorf("invalid EFI Signature %v", efiSignature)
	}
	if !bytes.Equal(efiRevision, _efiRevision) {
		return Table{}, fmt.Errorf("invalid EFI Revision %v", efiRevision)
	}
	if !bytes.Equal(efiHeaderSize, _efiHeaderSize) {
		return Table{}, fmt.Errorf("invalid EFI Header size %v", efiHeaderSize)
	}
	if !bytes.Equal(efiZeroes, _efiZeroes) {
		return Table{}, fmt.Errorf("invalid EFI Header, expected zeroes, got %v", efiZeroes)
	}

	// get the checksum
	checksum := crc32.ChecksumIEEE(gpt[0:92])
	if efiHeaderCrc != checksum {
		return Table{}, fmt.Errorf("invalid EFI Header Checksum, expected %v, got %v", checksum, efiHeaderCrc)
	}

	// now for partitions
	partArrayStart := partitionEntryFirstLBA * uint64(logicalBlockSize)
	partArrayEnd := partArrayStart + uint64(partitionEntryCount*partitionEntrySize)
	bpart := b[partArrayStart:partArrayEnd]

	// we need a CRC/zlib of the partition entries, so we do those first, then append the bytes
	checksum = crc32.ChecksumIEEE(bpart)
	if partitionEntryChecksum != checksum {
		return Table{}, fmt.Errorf("invalid EFI Partition Entry Checksum, expected %v, got %v", checksum, partitionEntryChecksum)
	}

	parts := make([]Partition, 0, partitionEntryCount)
	count := int(partitionEntryCount)

	for i := 0; i < count; i++ {

		// write the primary partition entry
		start := i * int(partitionEntrySize)
		end := start + int(partitionEntrySize)

		p, err := decodePartitionBytes(bpart[start:end], logicalBlockSize, physicalBlockSize)
		if err != nil {
			return Table{}, fmt.Errorf("decoding partition %d: %v", i, err)
		}

		parts = append(parts, p)
	}

	table := Table{
		LogicalSectorSize:  logicalBlockSize,
		PhysicalSectorSize: physicalBlockSize,
		partitionEntrySize: partitionEntrySize,
		primaryHeader:      primaryHeader,
		secondaryHeader:    secondaryHeader,
		firstDataSector:    firstDataSector,
		lastDataSector:     lastDataSector,
		partitionArraySize: int(partitionEntryCount),
		GUID:               diskGUID,
		Partitions:         parts,
	}

	return table, nil
}
