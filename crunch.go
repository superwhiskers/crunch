/*

crunch - utilities for taking bytes out of things
copyright (c) 2019 superwhiskers <whiskerdev@protonmail.com>

this program is free software: you can redistribute it and/or modify
it under the terms of the gnu lesser general public license as published by
the free software foundation, either version 3 of the license, or
(at your option) any later version.

this program is distributed in the hope that it will be useful,
but without any warranty; without even the implied warranty of
merchantability or fitness for a particular purpose.  see the
gnu lesser general public license for more details.

you should have received a copy of the gnu lesser general public license
along with this program.  if not, see <https://www.gnu.org/licenses/>.

*/

// Package crunch provides various utilities for manipulating bytes and bits easily
package crunch

// Endianness represents the endianness of the value read or written
type Endianness int

const (
	// LittleEndian represents the little-endian byte order
	LittleEndian Endianness = iota

	// BigEndian represents the big-endian byte order
	BigEndian
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
