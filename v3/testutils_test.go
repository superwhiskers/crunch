/*

crunch - utilities for taking bytes out of things
copyright (c) 2019 superwhiskers <whiskerdev@protonmail.com>

this source code form is subject to the terms of the mozilla public
license, v. 2.0. if a copy of the mpl was not distributed with this
file, you can obtain one at http://mozilla.org/MPL/2.0/.

*/

package crunch

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
