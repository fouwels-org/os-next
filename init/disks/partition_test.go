package disks

import (
	"log"
	"os"
	"testing"
)

func TestAddPartition(t *testing.T) {

	const device = "/tmp/device.bin"
	const gb = 1

	// create dummy "disk"
	size := int64(gb * 1024 * 1024 * 1024)
	fd, err := os.Create(device)
	if err != nil {
		log.Fatal("Failed to create device")
	}
	_, err = fd.Seek(size-1, 0)
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

	// test disk
	err = AddPartition("BOOT", device, 300, 0)
	if err != nil {
		t.Fatalf("Failed to add BOOT: %v", err)
	}

	err = FormatPartition("BOOT", device, "vfat", 0)
	if err != nil {
		t.Fatalf("Failed to format BOOT: %v", err)
	}

	err = AddPartition("CONFIG", device, 500, 1)
	if err != nil {
		t.Fatalf("Failed to add CONFIG: %v", err)
	}

	err = FormatPartition("CONFIG", device, "ext4", 0)
	if err != nil {
		t.Fatalf("Failed to format CONFIG: %v", err)
	}

	err = AddPartition("DATA", device, 0, 2)
	if err != nil {
		t.Fatalf("Failed to add DATA: %v", err)
	}

	err = FormatPartition("DATA", device, "ext4", 0)
	if err != nil {
		t.Fatalf("Failed to format DATA: %v", err)
	}
}
