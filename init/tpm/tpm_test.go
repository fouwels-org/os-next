// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package tpm

import (
	"crypto/sha256"
	"log"
	"testing"
)

func TestReadPCR(t *testing.T) {

	values, err := ReadPCRs()
	if err != nil {
		t.Fatalf("%v", err)
	}

	for _, v := range values {
		log.Printf("%v 0x%X (%X)", v.ID, v.Value, v.Value[len(v.Value)-3:])
	}

}

func TestPCRExtend(t *testing.T) {
	_id := PCR8KernelHash

	value := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

	sha := sha256.Sum256(value)

	log.Printf("hash: %v", sha[:])

	err := PCRExtend(PCR{
		ID:    _id,
		Value: sha[:],
	})
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestPCRChainExtend(t *testing.T) {
	_id := PCR8KernelHash

	value := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	sha := sha256.Sum256(value)

	values, err := ReadPCRs()
	if err != nil {
		t.Fatalf("%v", err)
	}

	for _, v := range values {
		log.Printf("%v 0x%X (%X)", v.ID, v.Value, v.Value[len(v.Value)-3:])
	}

	err = PCRExtend(PCR{
		ID:    _id,
		Value: sha[:],
	})
	if err != nil {
		t.Fatalf("%v", err)
	}
}
