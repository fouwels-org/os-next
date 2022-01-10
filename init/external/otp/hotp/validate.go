// SPDX-FileCopyrightText: Copyright (c) 2014 Paul Querna
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package hotp

import (
	"hash"

	"github.com/fouwels/os-next/init/external/otp/common"

	"crypto/hmac"
	"crypto/subtle"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

// ValidateOpts provides options for ValidateCustom().
type ValidateOpts struct {
	// Digits as part of the input.
	Digits common.Digits
	// Algorithm to use for HMAC.
	Algorithm func() hash.Hash
}

// GenerateCodeCustom uses a counter and secret value and options struct to
// create a passcode.
func GenerateCode(secret string, counter uint64, opts ValidateOpts) (passcode string, err error) {
	// As noted in issue #10 and #17 this adds support for TOTP secrets that are
	// missing their padding.
	secret = strings.TrimSpace(secret)
	if n := len(secret) % 8; n != 0 {
		secret = secret + strings.Repeat("=", 8-n)
	}

	// As noted in issue #24 Google has started producing base32 in lower case,
	// but the StdEncoding (and the RFC), expect a dictionary of only upper case letters.
	secret = strings.ToUpper(secret)

	secretBytes, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to decode secret as base32")
	}

	buf := make([]byte, 8)
	mac := hmac.New(opts.Algorithm, secretBytes)
	binary.BigEndian.PutUint64(buf, counter)

	_, err = mac.Write(buf)
	if err != nil {
		return "", err
	}
	sum := mac.Sum(nil)

	// "Dynamic truncation" in RFC 4226
	// http://tools.ietf.org/html/rfc4226#section-5.4
	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	l := opts.Digits.Length()
	mod := int32(value % int64(math.Pow10(l)))

	return opts.Digits.Format(mod), nil
}

func Validate(passcode string, counter uint64, secret string, opts ValidateOpts) (bool, error) {
	passcode = strings.TrimSpace(passcode)

	if len(passcode) != opts.Digits.Length() {
		return false, fmt.Errorf("digits is incorrect length")
	}

	otpstr, err := GenerateCode(secret, counter, opts)
	if err != nil {
		return false, err
	}

	if subtle.ConstantTimeCompare([]byte(otpstr), []byte(passcode)) == 1 {
		return true, nil
	}

	return false, nil
}
