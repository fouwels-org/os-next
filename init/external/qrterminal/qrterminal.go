// SPDX-FileCopyrightText: Copyright 2019 Mark Percival <m@mdp.im>
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package qrterminal

import (
	"io"
	"strings"

	"rsc.io/qr"
)

const WHITE = "\033[47m  \033[0m"
const BLACK = "\033[40m  \033[0m"

// Use ascii blocks to form the QR Code
const BLACK_WHITE = "▄"
const BLACK_BLACK = " "
const WHITE_BLACK = "▀"
const WHITE_WHITE = "█"

// Level - the QR Code's redundancy level
const H = qr.H
const M = qr.M
const L = qr.L

// default is 4-pixel-wide white quiet zone
const QUIET_ZONE = 4

//Config for generating a barcode
type Config struct {
	Level          qr.Level
	Writer         io.Writer
	HalfBlocks     bool
	BlackChar      string
	BlackWhiteChar string
	WhiteChar      string
	WhiteBlackChar string
	QuietZone      int
}

func (c *Config) writeFullBlocks(w io.Writer, code *qr.Code) error {
	white := c.WhiteChar
	black := c.BlackChar

	// Frame the barcode in a 1 pixel border
	_, err := w.Write([]byte(stringRepeat(stringRepeat(white, code.Size+c.QuietZone*2)+"\n", c.QuietZone))) // top border
	if err != nil {
		return err
	}
	for i := 0; i <= code.Size; i++ {
		_, err := w.Write([]byte(stringRepeat(white, c.QuietZone))) // left border
		if err != nil {
			return err
		}
		for j := 0; j <= code.Size; j++ {
			if code.Black(j, i) {
				_, err := w.Write([]byte(black))
				if err != nil {
					return err
				}
			} else {
				_, err := w.Write([]byte(white))
				if err != nil {
					return err
				}
			}
		}
		_, err = w.Write([]byte(stringRepeat(white, c.QuietZone-1) + "\n")) // right border
		if err != nil {
			return err
		}
	}
	_, err = w.Write([]byte(stringRepeat(stringRepeat(white, code.Size+c.QuietZone*2)+"\n", c.QuietZone-1))) // bottom border
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) writeHalfBlocks(w io.Writer, code *qr.Code) error {
	ww := c.WhiteChar
	bb := c.BlackChar
	wb := c.WhiteBlackChar
	bw := c.BlackWhiteChar
	// Frame the barcode in a 4 pixel border
	// top border
	if c.QuietZone%2 != 0 {
		_, err := w.Write([]byte(stringRepeat(bw, code.Size+c.QuietZone*2) + "\n"))
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(stringRepeat(stringRepeat(ww, code.Size+c.QuietZone*2)+"\n", c.QuietZone/2)))
		if err != nil {
			return err
		}
	} else {
		_, err := w.Write([]byte(stringRepeat(stringRepeat(ww, code.Size+c.QuietZone*2)+"\n", c.QuietZone/2)))
		if err != nil {
			return err
		}
	}
	for i := 0; i <= code.Size; i += 2 {
		_, err := w.Write([]byte(stringRepeat(ww, c.QuietZone))) // left border
		if err != nil {
			return err
		}
		for j := 0; j <= code.Size; j++ {
			next_black := false
			if i+1 < code.Size {
				next_black = code.Black(j, i+1)
			}
			curr_black := code.Black(j, i)
			if curr_black && next_black {
				_, err := w.Write([]byte(bb))
				if err != nil {
					return err
				}
			} else if curr_black && !next_black {
				_, err := w.Write([]byte(bw))
				if err != nil {
					return err
				}
			} else if !curr_black && !next_black {
				_, err := w.Write([]byte(ww))
				if err != nil {
					return err
				}
			} else {
				_, err := w.Write([]byte(wb))
				if err != nil {
					return err
				}
			}
		}
		_, err = w.Write([]byte(stringRepeat(ww, c.QuietZone-1) + "\n")) // right border
		if err != nil {
			return err
		}
	}
	// bottom border
	if c.QuietZone%2 == 0 {
		_, err := w.Write([]byte(stringRepeat(stringRepeat(ww, code.Size+c.QuietZone*2)+"\n", c.QuietZone/2-1)))
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(stringRepeat(wb, code.Size+c.QuietZone*2) + "\n"))
		if err != nil {
			return err
		}
	} else {
		_, err := w.Write([]byte(stringRepeat(stringRepeat(ww, code.Size+c.QuietZone*2)+"\n", c.QuietZone/2)))
		if err != nil {
			return err
		}
	}

	return nil
}

func stringRepeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(s, count)
}

// GenerateWithConfig expects a string to encode and a config
func GenerateWithConfig(text string, config Config) error {
	if config.QuietZone < 1 {
		config.QuietZone = 1 // at least 1-pixel-wide white quiet zone
	}
	w := config.Writer
	code, _ := qr.Encode(text, config.Level)
	if config.HalfBlocks {
		return config.writeHalfBlocks(w, code)
	} else {
		return config.writeFullBlocks(w, code)
	}
}

// Generate a QR Code and write it out to io.Writer
func Generate(text string, l qr.Level, w io.Writer) error {
	config := Config{
		Level:     l,
		Writer:    w,
		BlackChar: BLACK,
		WhiteChar: WHITE,
		QuietZone: QUIET_ZONE,
	}
	return GenerateWithConfig(text, config)
}

// Generate a QR Code with half blocks and write it out to io.Writer
func GenerateHalfBlock(text string, l qr.Level, w io.Writer) error {
	config := Config{
		Level:          l,
		Writer:         w,
		HalfBlocks:     true,
		BlackChar:      BLACK_BLACK,
		WhiteBlackChar: WHITE_BLACK,
		WhiteChar:      WHITE_WHITE,
		BlackWhiteChar: BLACK_WHITE,
		QuietZone:      QUIET_ZONE,
	}
	return GenerateWithConfig(text, config)
}
