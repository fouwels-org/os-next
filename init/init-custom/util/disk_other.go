// +build !linux

package util

import "fmt"

//DiskUtil ..
type DiskUtil struct {
}

//OpenLUKSvolumes Decrypts the volume or returns an error if decryption fails
func (d *DiskUtil) OpenLUKSvolumes() error {
	return fmt.Errorf("Not supported on this platform")
}

//FormatBlkid returns the available blockId as a deviceType from the underlying OS
func (d *DiskUtil) FormatBlkid() ([]DeviceType, error) {
	return []DeviceType{}, fmt.Errorf("Not supported on this platform")
}

// FindLabelledDevices returns the config and data partitions based on the formatted partitions.
// string: Config partition e.g /dev/sda2
// string: Data partition e.g /dev/sda3
// error: nil if partitons are found otherwise !nil
func (d *DiskUtil) FindLabelledDevices(partitions Partitions) (string, string, error) {
	return "", "", fmt.Errorf("Not supported on this platform")
}
