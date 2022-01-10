// SPDX-FileCopyrightText: Copyright (c) 2014 Paul Querna
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package totp

import (
	"fmt"
	"hash"

	"github.com/fouwels/os-next/init/external/otp/common"
	"github.com/fouwels/os-next/init/external/otp/hotp"

	"math"
	"time"
)

// ValidateOpts provides options for ValidateCustom().
type ValidateOpts struct {
	// Number of seconds a TOTP hash is valid for.
	Period uint
	// Digits as part of the input.
	Digits common.Digits
	// Algorithm to use for HMAC
	Algorithm func() hash.Hash
}

// GenerateCodeCustom takes a timepoint and produces a passcode using a
// secret and the provided opts. (Under the hood, this is making an adapted
// call to hotp.GenerateCodeCustom)
func GenerateCode(secret string, t time.Time, opts ValidateOpts) (passcode string, err error) {

	if opts.Period == 0 {
		return "", fmt.Errorf("missing period")
	}

	if opts.Digits == 0 {
		return "", fmt.Errorf("missing digits")
	}

	if opts.Algorithm == nil {
		return "", fmt.Errorf("missing algorithm")
	}

	counter := uint64(math.Floor(float64(t.Unix()) / float64(opts.Period)))
	passcode, err = hotp.GenerateCode(secret, counter, hotp.ValidateOpts{
		Digits:    opts.Digits,
		Algorithm: opts.Algorithm,
	})
	if err != nil {
		return "", err
	}
	return passcode, nil
}

// ValidateCustom validates a TOTP given a user specified time and custom options.
// Most users should use Validate() to provide an interpolatable TOTP experience.
func Validate(passcode string, secret string, t time.Time, opts ValidateOpts) (bool, error) {

	if opts.Period == 0 {
		return false, fmt.Errorf("missing period")
	}

	if opts.Digits == 0 {
		return false, fmt.Errorf("missing digits")
	}

	if opts.Algorithm == nil {
		return false, fmt.Errorf("missing algorithm")
	}

	counters := []uint64{}
	counter := int64(math.Floor(float64(t.Unix()) / float64(opts.Period)))

	counters = append(counters, uint64(counter))

	for _, counter := range counters {
		rv, err := hotp.Validate(passcode, counter, secret, hotp.ValidateOpts{
			Digits:    opts.Digits,
			Algorithm: opts.Algorithm,
		})

		if err != nil {
			return false, err
		}

		if rv {
			return true, nil
		}
	}

	return false, nil
}
