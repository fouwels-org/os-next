// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package config

type PrimaryFile struct {
	Primary Primary
	Header  Header
}

type SecondaryFile struct {
	Secondary Secondary
	Header    Header
}
