<h1 align="center"><img height="340" src="https://github.com/superwhiskers/crunch/raw/master/.github/cookie.png" alt="cookie with a bite taken out of it"/><br />crunch</h1>

<p align="center">
	<b>a library for easily manipulating bits and bytes in golang</b>
</p>

<p align="center">
	<a href="https://github.com/superwhiskers/crunch/blob/master/LICENSE.txt">
		<img src="https://img.shields.io/badge/license-MPL--2.0-brightgreen" alt="license" />
	</a>
  <a href="https://pkg.go.dev/pkg.go.dev/github.com/superwhiskers/crunch/v3">
    <img src="https://pkg.go.dev/badge/pkg.go.dev/github.com/superwhiskers/crunch/v3.svg" alt="documentation" />
  </a>
	<a href="https://travis-ci.org/superwhiskers/crunch">
		<img src="https://travis-ci.org/superwhiskers/crunch.svg?branch=master" alt="travis" />
	</a>
	<a href="https://codecov.io/gh/superwhiskers/crunch">
		<img src="https://codecov.io/gh/superwhiskers/crunch/branch/master/graph/badge.svg" alt="codecov" />
	</a>
	<a href="https://goreportcard.com/report/github.com/superwhiskers/crunch">
		<img src="https://goreportcard.com/badge/github.com/superwhiskers/crunch" alt="go report card" />
	</a>
	<a href="https://repl.it/github/https://github.com/superwhiskers/crunch?ref=button">
		<img src="https://img.shields.io/badge/try%20it%20on-repl.it-%2359646A.svg" alt="try it on repl.it" />
	</a>
</p>

<p align="center">
	<a href="#features">features</a> | <a href="#installation">installation</a> | <a href="#benchmarks">benchmarks</a> | <a href="#examples">examples</a>
</p>

## features

- **feature-rich**: supports reading and writing integers of varying sizes in both little and big endian
- **performant**: performs more than twice as fast as the standard library's `bytes.Buffer`
- **simple and familiar**: has a consistent and easy-to-use api
- **licensed under the mpl-2.0**: use it anywhere you wish, just don't change it privately

## installation

### install with the `go` tool

```bash
$ go get github.com/superwhiskers/crunch/v3
```

then, just import it in your project like this. easy!

```go
package "yourpackage"

import crunch "github.com/superwhiskers/crunch/v3"
```

### install using git submodules (not recommended)

```bash
# this assumes that you are in a git repository
$ git submodule add https://github.com/superwhiskers/crunch path/to/where/you/want/crunch
```

then, you can import it like this

```go
package "yourpackage"

import crunch "github.com/your-username/project/path/to/crunch/v3"
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

## examples

examples can be found in the [examples](https://github.com/superwhiskers/crunch/blob/master/examples) directory

## acknowledgements

icon (cookie logo) made by [freepik](https://www.freepik.com/) from [flaticon.com](https://www.flaticon.com)
