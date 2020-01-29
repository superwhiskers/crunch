/*

crunch - utilities for taking bytes out of things
copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>

this source code form is subject to the terms of the mozilla public
license, v. 2.0. if a copy of the mpl was not distributed with this
file, you can obtain one at http://mozilla.org/MPL/2.0/.

*/

package v2

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

	// BufferInvalidBitError represents an instance in which an invalid bit was
	// passed to one of the buffer's methods
	BufferInvalidBitError = Error{
		scope: "buffer",
		error: "invalid bit value specified",
	}
)
