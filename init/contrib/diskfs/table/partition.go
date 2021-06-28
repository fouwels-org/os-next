// SPDX-FileCopyrightText: 2017 Avi Deitcher
// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
//
// SPDX-License-Identifier: MIT

package table

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"unicode/utf16"

	"github.com/google/uuid"
)

// PartitionEntrySize fixed size of a GPT partition entry
const PartitionEntrySize = 128

// Partition represents the structure of a single partition on the disk
type Partition struct {
	Start      uint64    // start sector for the partition
	Size       uint64    // size of the partition in bytes
	Type       uuid.UUID // parttype for the partition
	Name       string    // name for the partition
	GUID       uuid.UUID // partition GUID
	Attributes uint64    // Attributes flags
}

func reverseSlice(s interface{}) {
	size := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// toBytes return the 128 bytes for this partition
func (p *Partition) encodeBytes(diskSize uint64) ([]byte, error) {
	b := make([]byte, PartitionEntrySize)

	// if the Type is Unused, just return all zeroes
	if p.Type == GPTUnused {
		return b, nil
	}

	copy(b[0:16], bytesToUUIDBytes(p.Type[0:16]))

	// partition identifier GUID is next 16 bytes
	copy(b[16:32], bytesToUUIDBytes(p.GUID[0:16]))

	startSector := p.Start / _logicalSectorSize

	var endSector uint64
	if p.Size == 0 {
		endSector = (diskSize / _logicalSectorSize) - 1
	} else {
		endSector = (p.Start + p.Size) / _logicalSectorSize
	}

	binary.LittleEndian.PutUint64(b[32:40], startSector)
	binary.LittleEndian.PutUint64(b[40:48], endSector)
	binary.LittleEndian.PutUint64(b[48:56], p.Attributes)

	// now the partition name - it is UTF16LE encoded, max 36 code units for 72 bytes
	r := make([]rune, 0, len(p.Name))
	// first convert to runes
	for _, s := range p.Name {
		r = append(r, rune(s))
	}
	if len(r) > 36 {
		return nil, fmt.Errorf("cannot use %s as partition name, has %d Unicode code units, maximum size is 36", p.Name, len(r))
	}
	// next convert the runes to uint16
	nameb := utf16.Encode(r)
	// and then convert to little-endian bytes
	for i, u := range nameb {
		pos := 56 + i*2
		binary.LittleEndian.PutUint16(b[pos:pos+2], u)
	}

	return b, nil
}

// decodePartitionBytes create a partition entry from bytes
func decodePartitionBytes(b []byte, logicalSectorSize int, physicalSectorSize int) (Partition, error) {
	if len(b) != PartitionEntrySize {
		return Partition{}, fmt.Errorf("data for partition was %d bytes instead of expected %d", len(b), PartitionEntrySize)
	}
	// is it all zeroes?
	typeGUID, err := uuid.FromBytes(bytesToUUIDBytes(b[0:16]))
	if err != nil {
		return Partition{}, fmt.Errorf("decoding type GUID: %w", err)
	}
	uuid, err := uuid.FromBytes(bytesToUUIDBytes(b[16:32]))
	if err != nil {
		return Partition{}, fmt.Errorf("decoding partition GUID: %w", err)
	}
	firstLBA := binary.LittleEndian.Uint64(b[32:40])
	lastLBA := binary.LittleEndian.Uint64(b[40:48])
	attribs := binary.LittleEndian.Uint64(b[48:56])

	// get the partition name
	nameb := b[56:]
	u := make([]uint16, 0, 72)
	for i := 0; i < len(nameb); i += 2 {
		// strip any 0s off of the end
		entry := binary.LittleEndian.Uint16(nameb[i : i+2])
		if entry == 0 {
			break
		}
		u = append(u, entry)
	}
	r := utf16.Decode(u)
	name := string(r)

	return Partition{
		Start:      firstLBA,
		Size:       (lastLBA - firstLBA) * uint64(logicalSectorSize),
		Name:       name,
		GUID:       uuid,
		Attributes: attribs,
		Type:       typeGUID,
	}, nil
}

// WriteContents fills the partition with the contents provided
// reads from beginning of reader to exactly size of partition in bytes
func (p *Partition) WriteContents(f File, contents io.Reader) (uint64, error) {

	total := uint64(0)

	b := make([]byte, _physicalSectorSize)

	for {

		read, err := contents.Read(b)
		if err != nil && err != io.EOF {
			return total, fmt.Errorf("could not read contents to pass to partition: %v", err)
		}

		tmpTotal := uint64(read) + total
		if uint64(tmpTotal) > p.Size {
			return total, fmt.Errorf("requested to write at least %d bytes to partition but maximum size is %d", tmpTotal, p.Size)
		}

		if read > 0 {

			var written int
			written, err = f.WriteAt(b[:read], int64(p.Start+total))
			if err != nil {
				return total, fmt.Errorf("error writing to file: %v", err)
			}

			total = total + uint64(written)
		}
		// increment our total
		// is this the end of the data?
		if err == io.EOF {
			break
		}
	}
	// did the total written equal the size of the partition?
	if uint64(total) != p.Size {
		return total, fmt.Errorf("write %d bytes to partition but actual size is %d", total, p.Size)
	}
	return total, nil
}

// ReadContents reads the contents of the partition into a writer
// streams the entire partition to the writer
func (p *Partition) ReadContents(f File, out io.Writer) (int64, error) {

	// chunks of physical sector size for efficient writing
	b := make([]byte, _physicalSectorSize)

	total := int64(0)

	// loop in physical sector sizes
	for {
		read, err := f.ReadAt(b, int64(p.Start)+total)

		if err != nil && err != io.EOF {
			return total, fmt.Errorf("reading from file: %w", err)
		}
		if read > 0 {
			n, err := out.Write(b[:read])
			if err != nil {
				return int64(n), fmt.Errorf("failed to write: %w", err)
			}
		}
		// increment our total
		total += int64(read)
		// is this the end of the data?
		if err == io.EOF || total >= int64(p.Size) {
			break
		}
	}
	return total, nil
}
