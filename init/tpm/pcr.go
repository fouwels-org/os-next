// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: LicenseRef-MS-RSL

package tpm

type PCR struct {
	ID    PCRID
	Value []byte
}

type PCRID int

const (
	PCR0UEFI PCRID = iota
	PCR1UEFIConfiguration
	PCR2OROM
	PCR3OROMConfiguration
	PCR4MBR
	PCR5MBRConfiguration
	PCR6PowerState
	PCR7SecureBoot
	PCR8KernelHash
	PCR16Debug = 16
)
