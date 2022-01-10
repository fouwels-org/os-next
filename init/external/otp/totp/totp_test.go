// SPDX-FileCopyrightText: Copyright (c) 2014 Paul Querna
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package totp

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"hash"
	"testing"
	"time"

	"github.com/fouwels/os-next/init/external/otp/common"
)

type tc struct {
	TS     int64
	TOTP   string
	Mode   func() hash.Hash
	Secret string
}

var (
	secSha256 = base32.StdEncoding.EncodeToString([]byte("12345678901234567890123456789012"))
	secSha512 = base32.StdEncoding.EncodeToString([]byte("1234567890123456789012345678901234567890123456789012345678901234"))

	rfcMatrixTCs = []tc{
		{59, "46119246", sha256.New, secSha256},
		{59, "90693936", sha512.New, secSha512},
		{1111111109, "68084774", sha256.New, secSha256},
		{1111111109, "25091201", sha512.New, secSha512},
		{1111111111, "67062674", sha256.New, secSha256},
		{1111111111, "99943326", sha512.New, secSha512},
		{1234567890, "91819424", sha256.New, secSha256},
		{1234567890, "93441116", sha512.New, secSha512},
		{2000000000, "90698825", sha256.New, secSha256},
		{2000000000, "38618901", sha512.New, secSha512},
		{20000000000, "77737706", sha256.New, secSha256},
		{20000000000, "47863826", sha512.New, secSha512},
	}
)

//
// Test vectors from http://tools.ietf.org/html/rfc6238#appendix-B
// NOTE -- the test vectors are documented as having the SAME
// secret -- this is WRONG -- they have a variable secret
// depending upon the hmac algorithm:
// 		http://www.rfc-editor.org/errata_search.php?rfc=6238
// this only took a few hours of head/desk interaction to figure out.
//
func TestValidateRFCMatrix(t *testing.T) {
	for _, tx := range rfcMatrixTCs {

		opts := ValidateOpts{
			Digits:    common.DigitsEight,
			Algorithm: tx.Mode,
			Period:    30,
		}

		valid, err := Validate(tx.TOTP, tx.Secret, time.Unix(tx.TS, 0).UTC(), opts)
		if err != nil {
			t.Fatalf("unexpected error totp=%s ts=%v: %v", tx.TOTP, tx.TS, err)
		}

		if !valid {
			t.Fatalf("unexpected totp failure totp=%s ts=%v", tx.TOTP, tx.TS)
		}
	}
}

func TestGenerateRFCTCs(t *testing.T) {
	for _, tx := range rfcMatrixTCs {
		passcode, err := GenerateCode(tx.Secret, time.Unix(tx.TS, 0).UTC(),
			ValidateOpts{
				Digits:    common.DigitsEight,
				Algorithm: tx.Mode,
				Period:    30,
			})

		if err != nil {
			t.Fatalf("%v", err)
		}

		if passcode != tx.TOTP {
			t.Fatalf("%v != %v", passcode, tx.TOTP)
		}
	}
}

func TestGenerate(t *testing.T) {

	opts := GenerateOpts{
		Issuer:      "SnakeOil",
		AccountName: "alice@example.com",
		SecretSize:  20,
		Period:      30,
		Digits:      common.DigitsSix,
		Algorithm:   sha256.New,
		AlgorithmID: "SHA256",
	}

	k, err := Generate(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if k.Issuer() != "SnakeOil" {
		t.Fatalf("unexpected issuer: %v", k.Issuer())
	}
	if k.AccountName() != "alice@example.com" {
		t.Fatalf("unexpected account name: %v", k.AccountName())
	}
	if len(k.Secret()) != 32 {
		t.Fatalf("unexpected secret length: %v", len(k.Secret()))
	}

	opts.SecretSize = 64

	k, err = Generate(opts)
	if err != nil {
		t.Fatalf("unexpected error for size 64: %v", err)
	}
	if len(k.Secret()) == 32 {
		t.Fatalf("unexpected secret length for size 64: %v", len(k.Secret()))
	}

	opts.SecretSize = 13

	k, err = Generate(opts)
	if err != nil {
		t.Fatalf("unexpected error for size 13: %v", err)
	}
	if len(k.Secret()) == 32 {
		t.Fatalf("unexpected secret length for size 13: %v", len(k.Secret()))
	}
}

func TestGenerateValidate(t *testing.T) {

	opts := GenerateOpts{
		Issuer:      "SnakeOil",
		AccountName: "alice@example.com",
		SecretSize:  32,
		Period:      30,
		Digits:      common.DigitsEight,
		Algorithm:   sha1.New,
		AlgorithmID: "SHA512",
	}

	vopts := ValidateOpts{
		Period:    30,
		Digits:    common.DigitsEight,
		Algorithm: sha512.New,
	}

	k, err := Generate(opts)
	if err != nil {
		t.Fatalf("could not generate: %v", err)
	}

	pass, err := GenerateCode(k.Secret(), time.Now(), vopts)
	if err != nil {
		t.Fatalf("could not validate: %v", err)
	}

	ok, err := Validate(pass, k.Secret(), time.Now(), vopts)
	if err != nil {
		t.Fatalf("could not validate")
	}
	if !ok {
		t.Fatalf("could not valiate as ok")
	}

	ok, err = Validate(pass, k.Secret(), time.Now().Add(100*time.Second), vopts)
	if err != nil {
		t.Fatalf("could not validate")
	}
	if ok {
		t.Fatalf("could not valiate as not ok out of time +")
	}

	ok, err = Validate(pass, k.Secret(), time.Now().Add(-100*time.Second), vopts)
	if err != nil {
		t.Fatalf("could not validate")
	}
	if ok {
		t.Fatalf("could not valiate as not ok out of time -")
	}

}
