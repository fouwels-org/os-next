package disks

import (
	"fmt"
	"init/shell"
	"log"
	"os"
	"os/exec"
	"strings"
)

const TYPE_EFI = "ef"
const TYPE_LINUX_FS = "83"

//AddPartition attempts to create a partition of a given size, on a given slot of a device
// Specify size = -1 to use max available.
func AddPartition(label string, device string, size int, index int) error {

	// check the device exists
	_, err := os.Stat(device)

	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("device %v does not exist", device)
	} else if err != nil {
		return fmt.Errorf("failed to check device identifier %v: %w", device, err)
	}

	// attempt to fdisk partition
	cmd := exec.Command(string(shell.Fdisk), device)
	cmd.Stderr = log.Writer()

	in, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to open stdin pipe: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("fdisk: %w", err)
	}

	// establish partition type
	partType := TYPE_LINUX_FS
	if label == "BOOT" {
		// force to EFI if label is BOOT..
		partType = TYPE_EFI
	}

	fdisksize := ""
	if size != -1 {
		fdisksize = fmt.Sprintf("+%vM", size)
	}

	commands := []string{
		"n",                        //new partition
		"p",                        //primary
		fmt.Sprintf("%v", index+1), //partition number
		"",                         //default start sector
		fdisksize,                  //of size MB
		"t",                        //set type
		partType,                   //set type initially. number is not prompted for if <2 partitions..
		fmt.Sprintf("%v", index+1), //partition number
		partType,                   //set type
		"w",                        //write partition
	}

	for _, c := range commands {
		_, err = in.Write([]byte(c))
		if err != nil {
			return fmt.Errorf("failed to write in: %w", err)
		}
		_, err = in.Write([]byte("\n"))
		if err != nil {
			return fmt.Errorf("failed to write in: %w", err)
		}
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("recieved err from wait: %w", err)
	}

	return nil
}

//FormatPartition ..
func FormatPartition(label string, device string, format string, index int) error {

	var fqname string

	if strings.Contains(device, "nvme") {
		fqname = fmt.Sprintf("%vp%v", device, index+1) //if nvme, build up with the p (nvme0n1p1)
	} else {
		fqname = fmt.Sprintf("%v%v", device, index+1) //otherwise, naively assume is sequental (sda1)
	}

	var err error
	var result string
	switch format {
	case "vfat":
		result, err = shell.Executor.ExecuteOne(shell.Command{
			Executable: shell.Mkdosfs,
			Arguments: []string{
				"-n", label,
				fqname,
			},
		})
	case "ext4":
		result, err = shell.Executor.ExecuteOne(shell.Command{
			Executable: shell.Mke2fs,
			Arguments: []string{
				"-t", format,
				"-L", label,
				fqname,
			},
		})
	default:
		return fmt.Errorf("file system is not supported: %v", format)
	}

	if err != nil {
		return fmt.Errorf("failed: %v %v", result, err)
	}

	return nil
}
