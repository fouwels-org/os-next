// SPDX-FileCopyrightText: Copyright (c) 2014 Paul Querna
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package otp

import (
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"net/url"
	"strconv"
	"strings"
)

// Error when attempting to convert the secret from base32 to raw bytes.
var ErrValidateSecretInvalidBase32 = errors.New("decoding of secret as base32 failed")

// The user provided passcode length was not expected.
var ErrValidateInputInvalidLength = errors.New("input length unexpected")

// When generating a Key, the Issuer must be set.
var ErrGenerateMissingIssuer = errors.New("issuer must be set")

// When generating a Key, the Account Name must be set.
var ErrGenerateMissingAccountName = errors.New("accountName must be set")

// Key represents an TOTP or HTOP key.
type Key struct {
	orig string
	url  *url.URL
}

// NewKeyFromURL creates a new Key from an TOTP or HOTP url.
//
// The URL format is documented here:
//   https://github.com/google/google-authenticator/wiki/Key-Uri-Format
//
func NewKeyFromURL(orig string) (*Key, error) {
	s := strings.TrimSpace(orig)

	u, err := url.Parse(s)

	if err != nil {
		return nil, err
	}

	return &Key{
		orig: s,
		url:  u,
	}, nil
}

func (k *Key) String() string {
	return k.orig
}

// Type returns "hotp" or "totp".
func (k *Key) Type() string {
	return k.url.Host
}

// Issuer returns the name of the issuing organization.
func (k *Key) Issuer() string {
	q := k.url.Query()

	issuer := q.Get("issuer")

	if issuer != "" {
		return issuer
	}

	p := strings.TrimPrefix(k.url.Path, "/")
	i := strings.Index(p, ":")

	if i == -1 {
		return ""
	}

	return p[:i]
}

// AccountName returns the name of the user's account.
func (k *Key) AccountName() string {
	p := strings.TrimPrefix(k.url.Path, "/")
	i := strings.Index(p, ":")

	if i == -1 {
		return p
	}

	return p[i+1:]
}

// Secret returns the opaque secret for this Key.
func (k *Key) Secret() string {
	q := k.url.Query()

	return q.Get("secret")
}

// Period returns a tiny int representing the rotation time in seconds.
func (k *Key) Period() uint64 {
	q := k.url.Query()

	if u, err := strconv.ParseUint(q.Get("period"), 10, 64); err == nil {
		return u
	}

	// If no period is defined 30 seconds is the default per (rfc6238)
	return 30
}

// URL returns the OTP URL as a string
func (k *Key) URL() string {
	return k.url.String()
}

// Algorithm represents the hashing function to use in the HMAC
// operation needed for OTPs.
type Algorithm int

const (
	AlgorithUnknown Algorithm = iota
	AlgorithmSHA256
	AlgorithmSHA512
)

func (a Algorithm) String() string {
	switch a {
	case AlgorithmSHA256:
		return "SHA256"
	case AlgorithmSHA512:
		return "SHA512"
	}

	return "UNKNOWN"
}

func (a Algorithm) Hash() hash.Hash {
	switch a {
	case AlgorithmSHA256:
		return sha256.New()
	case AlgorithmSHA512:
		return sha512.New()
	}

	return sha512.New()
}

// Digits represents the number of digits present in the
// user's OTP passcode. Six and Eight are the most common values.
type Digits int

const (
	DigitsSix   Digits = 6
	DigitsEight Digits = 8
)

// Format converts an integer into the zero-filled size for this Digits.
func (d Digits) Format(in int32) string {
	f := fmt.Sprintf("%%0%dd", d)
	return fmt.Sprintf(f, in)
}

// Length returns the number of characters for this Digits.
func (d Digits) Length() int {
	return int(d)
}

func (d Digits) String() string {
	return fmt.Sprintf("%d", d)
}
