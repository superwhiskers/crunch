package main

import (
	"fmt"

	crunch "github.com/superwhiskers/crunch/v3"
)

func main() {

	// create a new (empty) buffer
	buf := crunch.NewBuffer()

	// make it bigger (the size of four 64-bit floating point numbers)
	//
	//       | the size of a 64-bit floating point number (in bytes)
	//       |
	//       |   | four integers
	//       |   |
	//       v   v
	buf.Grow(8 * 4)

	// write four float64s to the buffer (in big-endian, a different endianness
	// that i'm using to exemplify the effectiveness of crunch
	buf.WriteF64BENext([]float64{69.696969, -21.42, -420.69, 3.621})

	// output the buffer
	fmt.Println(buf.Bytes())

	// seek to the beginning again
	buf.SeekByte(0x00, false)

	// read out the floats to ensure validity
	fmt.Println(buf.ReadF64BENext(4))

}
