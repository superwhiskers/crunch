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
BenchmarkBufferWriteBytes-4       	50000000	        36.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferReadBytes-4        	20000000	       100 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferWriteU32LE-4       	30000000	        47.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferReadU32LE-4        	10000000	       170 ns/op	       8 B/op	       1 allocs/op
BenchmarkBufferReadBit-4          	20000000	       113 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferReadBits-4         	 5000000	       253 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferSetBit-4           	20000000	       101 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferClearBit-4         	20000000	        99.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferWriteBytes-4   	200000000	         6.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferReadBytes-4    	2000000000	         1.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferWriteU32LE-4   	100000000	        21.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferReadU32LE-4    	100000000	        11.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferReadBit-4      	2000000000	         1.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferReadBits-4     	100000000	        14.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferSetBit-4       	500000000	         3.05 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferClearBit-4     	500000000	         3.16 ns/op	       0 B/op	       0 allocs/op
BenchmarkStdByteBufferWrite-4     	50000000	        23.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkStdByteBufferRead-4      	200000000	         7.18 ns/op	       0 B/op	       0 allocs/op
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
