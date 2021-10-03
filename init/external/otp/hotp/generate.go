// SPDX-FileCopyrightText: Copyright (c) 2014 Paul Querna
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package hotp

import (
	"hash"
	"os-next/init/external/otp/common"

	"crypto/rand"
	"encoding/base32"
	"fmt"
	"net/url"
)

// GenerateOpts provides options for .Generate()
type GenerateOpts struct {
	// Name of the issuing Organization/Company.
	Issuer string
	// Name of the User's Account (eg, email address)
	AccountName string
	// Size in size of the generated Secret.
	SecretSize uint
	// Digits to request.
	Digits common.Digits
	// Algorithm to use for HMAC.
	Algorithm func() hash.Hash
	// Algorithm ID to pass in URL parameters
	AlgorithmID string
}

// Generate creates a new HOTP Key.
func Generate(opts GenerateOpts) (*common.Key, error) {
	// url encode the Issuer/AccountName
	if opts.Issuer == "" {
		return nil, fmt.Errorf("missing issuer")
	}

	if opts.AccountName == "" {
		return nil, fmt.Errorf("missing account name")
	}

	if opts.SecretSize == 0 {
		return nil, fmt.Errorf("missing secret size")
	}

	if opts.Digits == 0 {
		return nil, fmt.Errorf("missing digits")
	}

	if opts.AlgorithmID == "" {
		return nil, fmt.Errorf("missing algorithm ID")
	}

	// otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example

	v := url.Values{}

	secret := make([]byte, opts.SecretSize)
	r := rand.Reader
	_, err := r.Read(secret)
	if err != nil {
		return nil, err
	}

	base := base32.StdEncoding.WithPadding(base32.NoPadding)

	v.Set("secret", base.EncodeToString(secret))
	v.Set("issuer", opts.Issuer)
	v.Set("algorithm", opts.AlgorithmID)
	v.Set("digits", opts.Digits.String())

	u := url.URL{
		Scheme:   "otpauth",
		Host:     "hotp",
		Path:     "/" + opts.Issuer + ":" + opts.AccountName,
		RawQuery: v.Encode(),
	}

	return common.NewKeyFromURL(u.String())
}
