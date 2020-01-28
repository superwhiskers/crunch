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
