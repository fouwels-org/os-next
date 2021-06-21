package diskfs_test

import (
	"init/contrib/diskfs"
	"init/contrib/diskfs/table"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
)

const _megabyte = 1024 * 1024
const _gigabyte = 1024 * _megabyte

func dummyDisk(t *testing.T, size int) string {
	const device = "/tmp/device.bin"

	// create dummy disk
	fd, err := os.Create(device)
	if err != nil {
		log.Fatal("Failed to create device")
	}
	defer fd.Close()

	_, err = fd.Seek(int64(size)-1, 0)
	if err != nil {
		log.Fatal("Failed to seek")
	}
	_, err = fd.Write([]byte{0})
	if err != nil {
		log.Fatal("Write failed")
	}

	err = fd.Close()
	if err != nil {
		log.Fatal("Failed to close file")
	}

	return device
}

func TestOpenDisk(t *testing.T) {
	device := dummyDisk(t, 1000*_megabyte)

	pm, err := diskfs.Open(device)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer pm.Close()
	log.Printf("%+v", pm)
}

func TestGetInvalidTable(t *testing.T) {
	device := dummyDisk(t, 1000*_megabyte)

	pm, err := diskfs.Open(device)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer pm.Close()

	_, err = pm.GetTable()
	if err == nil {
		t.Fatalf("error was not raised")
	}
}

func TestGetValidTable(t *testing.T) {

	device := dummyDisk(t, 1000*_megabyte)

	pm, err := diskfs.Open(device)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer pm.Close()

	tab := table.NewTable(1*_gigabyte, uuid.New())

	err = pm.SetTable(tab)
	if err != nil {
		t.Fatalf("failed to set: %v", err)
	}

	table, err := pm.GetTable()
	if err != nil {
		t.Fatalf("error getting table: %v", err)
	}

	log.Printf("%+v", table)
}

func TestWritePartitionTable(t *testing.T) {

	const _size = 10000 * _megabyte
	device := dummyDisk(t, _size)

	pm, err := diskfs.Open(device)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer pm.Close()

	// new table
	tab := table.NewTable(4*_size, uuid.New())

	err = pm.SetTable(tab)
	if err != nil {
		t.Fatalf("failed to set: %v", err)
	}

	tab.Partitions[0] = table.Partition{
		Start:      0 * _megabyte,
		Size:       500 * _megabyte,
		Type:       table.GPTEFISystemPartition,
		Name:       "EFI",
		GUID:       uuid.New(),
		Attributes: 0x00,
	}
	tab.Partitions[1] = table.Partition{
		Start:      500 * _megabyte,
		Size:       500 * _megabyte,
		Type:       table.GPTLinuxFilesystem,
		Name:       "SYS",
		GUID:       uuid.New(),
		Attributes: 0x00,
	}
	tab.Partitions[2] = table.Partition{
		Start:      1000 * _megabyte,
		Size:       0 * _megabyte,
		Type:       table.GPTLinuxFilesystem,
		Name:       "MAX",
		GUID:       uuid.New(),
		Attributes: 0x00,
	}

	// write table
	err = pm.SetTable(tab)
	if err != nil {
		t.Fatalf("failed to set table: %v", err)
	}

	// read back table
	_, err = pm.GetTable()
	if err != nil {
		t.Fatalf("error getting table: %v", err)
	}
}

func TestWriteReadTableBadSectorAllocation(t *testing.T) {

	device := dummyDisk(t, 1*_gigabyte)

	pm, err := diskfs.Open(device)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer pm.Close()

	uid := uuid.New()
	if err != nil {
		t.Fatalf("failed to create uid: %v", err)
	}

	// new table
	tab := table.NewTable(1024*1024, uid)
	pm.SetTable(tab)

	// read table
	tab, err = pm.GetTable()
	if err != nil {
		t.Fatalf("error getting table: %v", err)
	}

	// set part
	puid := uuid.New()
	if err != nil {
		t.Fatalf("failed to create puid: %v", err)
	}

	in := table.Partition{
		Start:      0,
		Size:       1024*1024*512 + 1,
		Type:       table.GPTEFISystemPartition,
		Name:       "YOLO",
		GUID:       puid,
		Attributes: 0x00,
	}

	tab.Partitions[0] = in

	// write table
	err = pm.SetTable(tab)
	log.Printf("%v", err)
	if err == nil {
		t.Fatalf("bad size not caught")
	}

	// read table
	tab, err = pm.GetTable()
	if err != nil {
		t.Fatalf("error getting table: %v", err)
	}

	// set part
	puid = uuid.New()
	if err != nil {
		t.Fatalf("failed to create puid: %v", err)
	}

	in = table.Partition{
		Start:      1,
		Size:       1024 * 1024 * 512,
		Type:       table.GPTEFISystemPartition,
		Name:       "YOLO",
		GUID:       puid,
		Attributes: 0x00,
	}

	tab.Partitions[0] = in

	// write table
	err = pm.SetTable(tab)
	log.Printf("%v", err)
	if err == nil {
		t.Fatalf("bad start not caught")
	}
}

func TestFlags(t *testing.T) {

	device := dummyDisk(t, 10*_megabyte)

	pm, err := diskfs.Open(device)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer pm.Close()

	f, err := pm.GetFlag()
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if f != 0x00 {
		t.Fatalf("flag was not 0x00")
	}

	err = pm.SetFlag(0xFF0000AA)
	if err != nil {
		t.Fatalf("set: %v", err)
	}

	f, err = pm.GetFlag()
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if f != 0xFF0000AA {
		t.Fatalf("flag was not 0xFF0000AA")
	}
}
