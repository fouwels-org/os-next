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
