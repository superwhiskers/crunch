/*

crunch - utilities for taking bytes out of things
copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>

this source code form is subject to the terms of the mozilla public
license, v. 2.0. if a copy of the mpl was not distributed with this
file, you can obtain one at http://mozilla.org/MPL/2.0/.

*/

package v1

import "testing"

/*

tests

*/

func TestErrorError(t *testing.T) {

	expected := "crunch: buffer: read exceeds buffer capacity"

	out := BufferOverreadError.Error()
	if out != expected {

		t.Fatalf("expected string does not match the one gotten (got \"%s\", expected \"%s\")", out, expected)

	}

}
