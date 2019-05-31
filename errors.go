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
	// BufferOverreadError represents an instance in which a read attempted to
	// read past the buffer itself
	BufferOverreadError = Error{
		scope: "buffer",
		error: "read exceeds buffer capacity",
	}

	// BufferUnderreadError represents an instance in which a read attempted to
	// read before the buffer itself
	BufferUnderreadError = Error{
		scope: "buffer",
		error: "read offset is less than zero",
	}

	// BufferOverwriteError represents an instance in which a write attempted to
	// write past the buffer itself
	BufferOverwriteError = Error{
		scope: "buffer",
		error: "write offset exceeds buffer capacity",
	}

	// BufferUnderwriteError represents an instance in which a write attempted to
	// write before the buffer itself
	BufferUnderwriteError = Error{
		scope: "buffer",
		error: "write offset is less than zero",
	}

	// BufferInvalidByteCountError represents an instance in which an invalid byte
	// count was passed to one of the buffer's methods
	BufferInvalidByteCountError = Error{
		scope: "buffer",
		error: "invalid byte count requested",
	}
)
