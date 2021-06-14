package disks

import (
	"log"
	"testing"
)

func TestGetBlkid(t *testing.T) {
	b, err := GetBlkid("")

	log.Printf("%+v", b[0])

	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetBlkidMap(t *testing.T) {
	b, err := GetBlkid("")

	bmap := b.LabelMap()

	log.Printf("%+v", bmap)

	if err != nil {
		t.Fatalf("%v", err)
	}
}
