package util

import "testing"

func GetBLKDevices(t *testing.T) {
	d := DiskUtil{}

	err, devs := d.GetBLKDevices()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%v", devs)
}
