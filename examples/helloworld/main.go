package main

import (
	"fmt"

	"github.com/superwhiskers/crunch"
)

func main() {

	// create a new empty buffer
	buf := crunch.NewBuffer()

	// expand it to have 12 null bytes
	buf.Grow(12)

	// write "hello, world" to it
	buf.WriteBytesNext([]byte("hello, world"))

	// output the buffer to the console
	fmt.Println(string(buf.Bytes()))

}
