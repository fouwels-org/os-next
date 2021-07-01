// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: LicenseRef-MS-RSL

package tpm

import (
	"fmt"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

const _device = "/dev/tpmrm0"
const _hash = tpm2.AlgSHA256

func ReadPCRs() ([]PCR, error) {

	tpm, err := tpm2.OpenTPM(_device)
	if err != nil {
		return nil, fmt.Errorf("couldn't open TPM %s: %s", _device, err)
	}

	defer tpm.Close()

	PCRs := []PCRID{
		PCR0UEFI,
		PCR1UEFIConfiguration,
		PCR2OROM,
		PCR3OROMConfiguration,
		PCR4MBR,
		PCR5MBRConfiguration,
		PCR6PowerState,
		PCR7SecureBoot,
		PCR8KernelHash,
		PCR16Debug,
	}

	plist := []PCR{}
	for _, v := range PCRs {

		p, err := tpm2.ReadPCR(tpm, int(v), _hash)
		if err != nil {
			return nil, fmt.Errorf("couldn't read PCR %v: %s", v, err)
		}
		plist = append(plist, PCR{ID: v, Value: p})
	}

	return plist, nil
}

func PCRExtend(pcr PCR) error {

	tpm, err := tpm2.OpenTPM(_device)
	if err != nil {
		return fmt.Errorf("couldn't open TPM %s: %s", _device, err)
	}
	defer tpm.Close()

	err = tpm2.PCRExtend(tpm, tpmutil.Handle(pcr.ID), _hash, pcr.Value, "")
	if err != nil {
		return fmt.Errorf("failed to extend: %w", err)
	}

	return nil
}
