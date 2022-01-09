// SPDX-FileCopyrightText: 2017 Avi Deitcher
// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package table

import "github.com/google/uuid"

const (
	_gptSize            = 128 * 128
	_physicalSectorSize = 512
	_logicalSectorSize  = 512
	_partitionArraySize = 128
	_partitionEntrySize = 128
)

var (
	_efiSignature  = []byte{0x45, 0x46, 0x49, 0x20, 0x50, 0x41, 0x52, 0x54}
	_efiRevision   = []byte{0x00, 0x00, 0x01, 0x00}
	_efiHeaderSize = []byte{0x5c, 0x00, 0x00, 0x00}
	_efiZeroes     = []byte{0x00, 0x00, 0x00, 0x00}
)

// List of GUID partition types
var (
	GPTUnused                   uuid.UUID = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	GPTMbrBoot                  uuid.UUID = uuid.MustParse("024DEE41-33E7-11D3-9D69-0008C781F39F")
	GPTEFISystemPartition       uuid.UUID = uuid.MustParse("C12A7328-F81F-11D2-BA4B-00A0C93EC93B")
	GPTBiosBoot                 uuid.UUID = uuid.MustParse("21686148-6449-6E6F-744E-656564454649")
	GPTMicrosoftReserved        uuid.UUID = uuid.MustParse("E3C9E316-0B5C-4DB8-817D-F92DF00215AE")
	GPTMicrosoftBasicData       uuid.UUID = uuid.MustParse("EBD0A0A2-B9E5-4433-87C0-68B6B72699C7")
	GPTMicrosoftLDMMetadata     uuid.UUID = uuid.MustParse("5808C8AA-7E8F-42E0-85D2-E1E90434CFB3")
	GPTMicrosoftLDMData         uuid.UUID = uuid.MustParse("AF9B60A0-1431-4F62-BC68-3311714A69AD")
	GPTMicrosoftWindowsRecovery uuid.UUID = uuid.MustParse("DE94BBA4-06D1-4D40-A16A-BFD50179D6AC")
	GPTLinuxFilesystem          uuid.UUID = uuid.MustParse("0FC63DAF-8483-4772-8E79-3D69D8477DE4")
	GPTLinuxRaid                uuid.UUID = uuid.MustParse("A19D880F-05FC-4D3B-A006-743F0F84911E")
	GPTLinuxRootX86             uuid.UUID = uuid.MustParse("44479540-F297-41B2-9AF7-D131D5F0458A")
	GPTLinuxRootX86_64          uuid.UUID = uuid.MustParse("4F68BCE3-E8CD-4DB1-96E7-FBCAF984B709")
	GPTLinuxRootArm32           uuid.UUID = uuid.MustParse("69DAD710-2CE4-4E3C-B16C-21A1D49ABED3")
	GPTLinuxRootArm64           uuid.UUID = uuid.MustParse("B921B045-1DF0-41C3-AF44-4C6F280D3FAE")
	GPTLinuxSwap                uuid.UUID = uuid.MustParse("0657FD6D-A4AB-43C4-84E5-0933C84B4F4F")
	GPTLinuxLVM                 uuid.UUID = uuid.MustParse("E6D6D379-F507-44C2-A23C-238F2A3DF928")
	GPTLinuxDMCrypt             uuid.UUID = uuid.MustParse("7FFEC5C9-2D00-49B7-8941-3EA10A5586B7")
	GPTLinuxLUKS                uuid.UUID = uuid.MustParse("CA7D7CCB-63ED-4C53-861C-1742536059CC")
	GPTVMWareFilesystem         uuid.UUID = uuid.MustParse("AA31E02A-400F-11DB-9590-000C2911D1B8")
)
