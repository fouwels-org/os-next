package util

// DeviceType ..
type DeviceType struct {
	DEV      string
	LABEL    string
	UUID     string
	FSTYPE   string
	PARTUUID string
}

//Partitions ..
type Partitions struct {
	BootPartition   string
	DataPartition   string
	ConfigPartition string

	BootFile string

	DefaultDevConfig string
	DefaultDevData   string
}

//Disk static instance of DiskUtil{}
var Disk DiskUtil = DiskUtil{}
