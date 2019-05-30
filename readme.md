# crunch

[![godoc](https://godoc.org/github.com/superwhiskers/crunch?status.svg)](https://godoc.org/github.com/superwhiskers/crunch)&nbsp;[![travis](https://travis-ci.org/superwhiskers/crunch.svg?branch=master)](https://travis-ci.org/superwhiskers/crunch#)&nbsp;[![codecov](https://codecov.io/gh/superwhiskers/crunch/branch/master/graph/badge.svg)](https://codecov.io/gh/superwhiskers/crunch)&nbsp;[![go report card](https://goreportcard.com/badge/github.com/superwhiskers/crunch)](https://goreportcard.com/report/github.com/superwhiskers/crunch)&nbsp;[![edit on repl.it](https://img.shields.io/badge/try%20it%20on-repl.it-%2359646A.svg)](https://repl.it/github/https://github.com/superwhiskers/crunch?ref=button)


manipulate bytes and bits in golang with ease

## install

```
$ go get github.com/superwhiskers/crunch
```

## benchmarks (temporary, will be updated once changes are mirrored between the buffers)

`MiniBuffer` performs on average more than twice as fast as `bytes.Buffer` in both writing and reading
```
BenchmarkBufferWrite-4            	50000000	         37.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferRead-4             	20000000	         106 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferWriteU32LE-4       	30000000	         46.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferReadU32LE-4        	10000000	         156 ns/op	       8 B/op	       1 allocs/op
BenchmarkMiniBufferWriteBytes-4   	200000000	         6.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferReadBytes-4    	2000000000	         1.47 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferWriteU32LE-4   	100000000	         21.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferReadU32LE-4    	100000000	         12.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferReadBit-4      	2000000000	         1.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferReadBits-4     	100000000	         14.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferSetBit-4       	500000000	         3.11 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferClearBit-4     	500000000	         3.13 ns/op	       0 B/op	       0 allocs/op
BenchmarkStdByteBufferWrite-4     	100000000	         23.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkStdByteBufferRead-4      	200000000	         7.09 ns/op	       0 B/op	       0 allocs/op
```

## example

```golang
package main

import (
	"fmt"
	"github.com/superwhiskers/crunch"
)

func main() {

	// creates a new buffer with four zeroes
	buf := crunch.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})
	
	// write the byte `0x01` to the first offset, and move the offset forward one
	buf.WriteByteNext(0x01)
	
	// write the byte `0x01` to the second offset, and move the offset forward one
	buf.WriteByteNext(0x01)
	
	// seek the offset back one
	buf.Seek(-1, true)
	
	// write the bytes `0x02` and `0x03` to the second and third offsets, respectively
	buf.WriteBytesNext([]byte{0x02, 0x03})
	
	// write the byte `0x04` to offset `0x03`
	buf.WriteByte(0x03, 0x04)
	
	// output the buffer's contents to the console
	fmt.Printf("%v\n", buf.Bytes())
	
}
```

## license

[lgplv3](https://www.gnu.org/licenses/lgpl-3.0.en.html)
