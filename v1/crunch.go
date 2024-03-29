/*

crunch - utilities for taking bytes out of things
Copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.

*/

// Package v1 provides various utilities for manipulating bytes and bits easily
package v1

import "encoding/binary"

var (
	// LittleEndian represents the little-endian byte order
	LittleEndian = binary.LittleEndian

	// BigEndian represents the big-endian byte order
	BigEndian = binary.BigEndian
)

// IntegerSize represents the size of the integer read or written (in bytes)
type IntegerSize int

const (
	// Unsigned8 represents the 8-bit unsigned integer size
	Unsigned8 = 1

	// Unsigned16 represents the 16-bit unsigned integer size
	Unsigned16 = 2

	// Unsigned32 represents the 32-bit unsigned integer size
	Unsigned32 = 4

	// Unsigned64 represents the 64-bit unsigned integer size
	Unsigned64 = 8
)
