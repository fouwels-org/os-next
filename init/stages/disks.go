package stages

import (
	"fmt"
	"init/config"
	"init/disks"
	"log"
)

//Disks implementes IStage
type Disks struct {
	finals []string
}

//String ..
func (m *Disks) String() string {
	return "disks"
}

//Policy ..
func (m *Disks) Policy() Policy {
	return PolicyHard
}

//Finalise ..
func (m *Disks) Finalise() []string {
	return m.finals
}

//Run ..
func (m *Disks) Run(c config.Config) error {

	// Provisioning will not be started if BLKID Label=BOOT exists.
	b, err := disks.GetBlkid("")
	if err != nil {
		return fmt.Errorf("failed to get blkid: %v", err)
	}
	bmap := b.LabelMap()
	_, ok := bmap["BOOT"]
	if ok {
		m.finals = append(m.finals, "disk were not re-initialised, BLKID BOOT exists")
		return nil
	}

	// Provision disks
	for _, v := range c.Primary.Filesystem.Devices {
		log.Printf("creating %+v", v)
		err := disks.AddPartition(v.Label, v.Device, v.Size, v.Index)
		if err != nil {
			return fmt.Errorf("failed to create partition for %v: %w", v.Label, err)
		}

		log.Printf("formatting %+v", v)
		err = disks.FormatPartition(v.Label, v.Device, v.FileSystem, v.Index)
		if err != nil {
			return fmt.Errorf("failed to format partition for %v: %w", v.Label, err)
		}
	}

	return nil
}
