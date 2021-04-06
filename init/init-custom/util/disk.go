package util

// BLKDevice ..
type BLKDevice struct {
	Device   string
	Label    string
	UUID     string
	FsType   string
	PartUUID string
}

//Disk static instance of DiskUtil{}
var Disk DiskUtil = DiskUtil{}
