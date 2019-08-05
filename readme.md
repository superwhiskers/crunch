# crunch

[![godoc](https://godoc.org/github.com/superwhiskers/crunch?status.svg)](https://godoc.org/github.com/superwhiskers/crunch)&nbsp;[![travis](https://travis-ci.org/superwhiskers/crunch.svg?branch=master)](https://travis-ci.org/superwhiskers/crunch#)&nbsp;[![codecov](https://codecov.io/gh/superwhiskers/crunch/branch/master/graph/badge.svg)](https://codecov.io/gh/superwhiskers/crunch)&nbsp;[![go report card](https://goreportcard.com/badge/github.com/superwhiskers/crunch)](https://goreportcard.com/report/github.com/superwhiskers/crunch)&nbsp;[![edit on repl.it](https://img.shields.io/badge/try%20it%20on-repl.it-%2359646A.svg)](https://repl.it/github/https://github.com/superwhiskers/crunch?ref=button)


manipulate bytes and bits in golang with ease

## install

```
$ go get github.com/superwhiskers/crunch
```

## benchmarks

both `Buffer` and `MiniBuffer` perform on average more than twice as fast as `bytes.Buffer` in both writing and reading
```
BenchmarkBufferWriteBytes-4             2000000000               1.45 ns/op            0 B/op          0 allocs/op
BenchmarkBufferReadBytes-4              2000000000               0.84 ns/op            0 B/op          0 allocs/op
BenchmarkBufferWriteU32LE-4             200000000                9.28 ns/op            0 B/op          0 allocs/op
BenchmarkBufferReadU32LE-4              50000000                26.6 ns/op             8 B/op          1 allocs/op
BenchmarkBufferReadBit-4                2000000000               0.84 ns/op            0 B/op          0 allocs/op
BenchmarkBufferReadBits-4               1000000000               2.25 ns/op            0 B/op          0 allocs/op
BenchmarkBufferSetBit-4                 1000000000               2.07 ns/op            0 B/op          0 allocs/op
BenchmarkBufferClearBit-4               1000000000               2.05 ns/op            0 B/op          0 allocs/op
BenchmarkMiniBufferWriteBytes-4         2000000000               1.42 ns/op            0 B/op          0 allocs/op
BenchmarkMiniBufferReadBytes-4          2000000000               0.58 ns/op            0 B/op          0 allocs/op
BenchmarkMiniBufferWriteU32LE-4         200000000                8.79 ns/op            0 B/op          0 allocs/op
BenchmarkMiniBufferReadU32LE-4          500000000                3.91 ns/op            0 B/op          0 allocs/op
BenchmarkMiniBufferReadBit-4            2000000000               0.64 ns/op            0 B/op          0 allocs/op
BenchmarkMiniBufferReadBits-4           300000000                5.64 ns/op            0 B/op          0 allocs/op
BenchmarkMiniBufferSetBit-4             1000000000               2.03 ns/op            0 B/op          0 allocs/op
BenchmarkMiniBufferClearBit-4           1000000000               2.10 ns/op            0 B/op          0 allocs/op
BenchmarkStdByteBufferWrite-4           200000000                9.42 ns/op            0 B/op          0 allocs/op
BenchmarkStdByteBufferRead-4            500000000                3.36 ns/op            0 B/op          0 allocs/op
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
	buf.SeekByte(-1, true)
	
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
