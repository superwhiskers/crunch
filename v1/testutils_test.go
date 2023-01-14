/*

crunch - utilities for taking bytes out of things
Copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.

*/

package v1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func panicChecker(t *testing.T, errs ...Error) {

	if r := recover(); r != nil {

		for i, err := range errs {

			if cmp.Equal(r, err, cmp.Comparer(func(x, y Error) bool {

				return x.scope == y.scope &&
					x.error == y.error

			})) {

				break

			}

			if i == len(errs)-1 {

				t.Fatalf("none of the expected panic return value(s) do not match the one gotten (got \"%s\", expected %v)", r, errs)

			}

		}

	} else {

		t.Fatalf("none of the expected panics were triggered")

	}

}
