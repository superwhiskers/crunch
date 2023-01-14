/*

crunch - utilities for taking bytes out of things
Copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.

*/

package v2

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
