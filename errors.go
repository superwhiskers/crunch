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

package crunch

import "fmt"

// Error implements a custom error type used in crunch
type Error struct {
	scope string
	error string
}

// Error formats the error held in a Error as a string
func (e Error) Error() string {

	return fmt.Sprintf("crunch: %s: %s", e.scope, e.error)

}

var (
	// ByteBufferOverreadError represents an instance in which an overread of a ByteBuffer has been attempted
	ByteBufferOverreadError = Error{
		scope: "bytebuffer",
		error: "read exceeds buffer capacity",
	}

	// ByteBufferOverwriteError represents an instance in which an overwrite of a ByteBuffer has been attempted
	ByteBufferOverwriteError = Error{
		scope: "bytebuffer",
		error: "write exceeds buffer capacity",
	}

	// ByteBufferInvalidIntegerSizeError represents an instance in which an invalid integer size was provided to the read/write function of a ByteBuffer
	ByteBufferInvalidIntegerSizeError = Error{
		scope: "bytebuffer",
		error: "invalid integer size specified",
	}

	// ByteBufferInvalidEndiannessError represents an instance in which an invalid endianness was provided to the read/write function of a ByteBuffer
	ByteBufferInvalidEndiannessError = Error{
		scope: "bytebuffer",
		error: "invalid endianness specified",
	}

	// ByteBufferInvalidByteCountError represents an instance in which an invalid byte count was provided to a function
	ByteBufferInvalidByteCountError = Error{
		scope: "bytebuffer",
		error: "invalid byte count requested",
	}

	// BitfieldInvalidBitError represents an instance in which an invalid bit was provided to a function
	BitfieldInvalidBitError = Error{
		scope: "bitfield",
		error: "invalid bit value specified",
	}

	// BitfieldOverreadError represents an instance in which an overread of a Bitfield has been attempted
	BitfieldOverreadError = Error{
		scope: "bitfield",
		error: "read exceeds bitfield capacity",
	}

	// BitfieldOverwriteError represents an instance in which an overwrite of a Bitfield has been attempted
	BitfieldOverwriteError = Error{
		scope: "bitfield",
		error: "write exceeds bitfield capacity",
	}
)
