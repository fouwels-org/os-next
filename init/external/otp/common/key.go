// SPDX-FileCopyrightText: Copyright (c) 2014 Paul Querna
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"net/url"
	"strconv"
	"strings"
)

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
