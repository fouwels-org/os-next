// SPDX-FileCopyrightText: Copyright (c) 2014 Paul Querna
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package totp

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"hash"

	"github.com/fouwels/os-next/init/external/otp/common"

	"net/url"
	"strconv"
)

// GenerateOpts provides options for Generate().  The default values
// are compatible with Google-Authenticator.
type GenerateOpts struct {
	// Name of the issuing Organization/Company.
	Issuer string
	// Name of the User's Account (eg, email address)
	AccountName string
	// Number of seconds a TOTP hash is valid for. Defaults to 30 seconds.
	Period uint
	// Size in size of the generated Secret. Defaults to 20 bytes.
	SecretSize uint
	// Digits to request. Defaults to 6.
	Digits common.Digits
	// Algorithm to use for HMAC. Defaults to SHA1.
	Algorithm func() hash.Hash
	// Algorithm ID to pass in URL parameters
	AlgorithmID string
}

// Generate a new TOTP Key.
func Generate(opts GenerateOpts) (*common.Key, error) {
	// url encode the Issuer/AccountName
	if opts.Issuer == "" {
		return nil, fmt.Errorf("missing issuer")
	}

	if opts.AccountName == "" {
		return nil, fmt.Errorf("missing account name")
	}

	if opts.Period == 0 {
		return nil, fmt.Errorf("missing period")
	}

	if opts.SecretSize == 0 {
		return nil, fmt.Errorf("missing secret size")
	}

	if opts.Digits == 0 {
		return nil, fmt.Errorf("missing digits")
	}

	if opts.Algorithm == nil {
		return nil, fmt.Errorf("missing algorithm")
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
	v.Set("period", strconv.FormatUint(uint64(opts.Period), 10))
	v.Set("algorithm", opts.AlgorithmID)
	v.Set("digits", opts.Digits.String())

	u := url.URL{
		Scheme:   "otpauth",
		Host:     "totp",
		Path:     "/" + opts.Issuer + ":" + opts.AccountName,
		RawQuery: v.Encode(),
	}

	return common.NewKeyFromURL(u.String())
}
