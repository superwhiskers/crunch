/*

crunch - utilities for taking bytes out of things
copyright (c) 2018 superwhiskers <whiskerdev@protonmail.com>

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
type Error struct{
	scope string
	error string
}

// Error formats the error held in a Error as a string
func (e Error) Error() string {

	return fmt.Sprintf("crunch: %s: %s", e.scope, e.error)

}

var (
	ByteBufferOverreadError = Error{
		scope: "bytebuffer",
		error: "read exceeds buffer capacity",
	}
	ByteBufferOverwriteError = Error{
		scope: "bytebuffer",
		error: "write exceeds buffer capacity",
	}
	ByteBufferInvalidIntegerSizeError = Error{
		scope: "bytebuffer",
		error: "invalid integer size specified",
	}
	ByteBufferInvalidEndiannessError = Error{
		scope: "bytebuffer",
		error: "invalid endianness specified",
	}
	ByteBufferInvalidByteCountError = Error{
		scope: "bytebuffer",
		error: "invalid byte count requested",
	}
	BitfieldInvalidBitError = Error{
		scope: "bitfield",
		error: "invalid bit value specified",
	}
	BitfieldOverreadError = Error{
		scope: "bitfield",
		error: "read exceeds bitfield capacity",
	}
	BitfieldOverwriteError = Error{
		scope: "bitfield",
		error: "write exceeds bitfield capacity",
	}
)
