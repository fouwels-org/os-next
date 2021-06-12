package shell

import (
	"fmt"
	"os"
)

//SelfTest tests defined executables to ensure they exist
func SelfTest() error {
	es := []Executable{
		Login,
		Ntpd,
		Modprobe,
		Hwclock,
		IP,
		Udhcp,
		Dockerd,
		Docker,
		Mkdir,
		Mount,
	}

	for _, e := range es {
		err := testOne(e)
		if err != nil {
			return err
		}
	}
	return nil
}

func testOne(e Executable) error {
	_, err := os.Stat(e.Target())
	if err != nil {
		return fmt.Errorf("failed to stat %v: %w", e, err)
	}

	return nil
}
