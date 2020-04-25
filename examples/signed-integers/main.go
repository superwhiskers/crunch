package main

import (
	"fmt"

	crunch "github.com/superwhiskers/crunch/v3"
)

func main() {

	// create a new (empty) buffer
	buf := crunch.NewBuffer()

	// make it bigger (the size of two 32-bit signed integers)
	//
	//       | the size of a 32-bit signed integer (in bytes)
	//       |
	//       |   | two integers
	//       |   | 
	//       v   v
	buf.Grow(4 * 2)

	// write two int32s to the buffer (in little-endian, the most common endianness)
	buf.WriteI32LENext([]int32{-42, 69})

	// output the buffer
	fmt.Println(buf.Bytes())

	// seek to the beginning again
	buf.SeekByte(0x00, false)

	// read out the integers to ensure validity
	fmt.Println(buf.ReadI32LENext(2))

}
