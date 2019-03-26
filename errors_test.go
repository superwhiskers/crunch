package crunch

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
