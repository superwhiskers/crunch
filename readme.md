<h1 align="center"><img height="340" src="https://github.com/superwhiskers/crunch/raw/canon/.github/cookie.png" alt="cookie with a bite taken out of it"/><br />crunch</h1>

<p align="center">
	<b>a library for easily manipulating bits and bytes in golang</b>
</p>

<p align="center">
	<a href="https://github.com/superwhiskers/crunch/blob/canon/LICENSE.txt">
		<img src="https://img.shields.io/badge/license-MPL--2.0-brightgreen" alt="license" />
	</a>
  <a href="https://pkg.go.dev/github.com/superwhiskers/crunch/v3">
    <img src="https://pkg.go.dev/badge/github.com/superwhiskers/crunch/v3.svg" alt="documentation" />
  </a>
	<a href="https://github.com/superwhiskers/crunch/actions/workflows/go.yml">
		<img src="https://github.com/superwhiskers/crunch/actions/workflows/go.yml/badge.svg" alt="go" />
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
goos: linux
goarch: amd64
pkg: github.com/superwhiskers/crunch/v3
cpu: Intel(R) Core(TM) i5-4300M CPU @ 2.60GHz
BenchmarkBufferWriteBytes-4       	612593820	         1.948 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferReadBytes-4        	1000000000	         0.5476 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferWriteU32LE-4       	125229171	         9.528 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferReadU32LE-4        	44677784	        24.73 ns/op	       8 B/op	       1 allocs/op
BenchmarkBufferReadBit-4          	1000000000	         0.5709 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferReadBits-4         	620840577	         1.869 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferSetBit-4           	600202990	         1.929 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferClearBit-4         	625814206	         1.993 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferGrow-4             	252735192	         6.638 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferWriteBytes-4   	577940112	         2.112 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferReadBytes-4    	1000000000	         0.5531 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferWriteU32LE-4   	116178949	        12.90 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferReadU32LE-4    	189681555	         6.296 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferWriteF32LE-4   	121033429	        13.06 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferReadF32LE-4    	170977377	         7.244 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferReadBit-4      	1000000000	         0.4730 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferReadBits-4     	249350655	         4.968 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferSetBit-4       	566985802	         2.173 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferClearBit-4     	531134203	         2.195 ns/op	       0 B/op	       0 allocs/op
BenchmarkMiniBufferGrow-4         	271458589	         4.434 ns/op	       0 B/op	       0 allocs/op
BenchmarkStdByteBufferWrite-4     	127322588	         9.249 ns/op	       0 B/op	       0 allocs/op
BenchmarkStdByteBufferRead-4      	319744824	         3.742 ns/op	       0 B/op	       0 allocs/op
BenchmarkStdByteBufferGrow-4      	274847331	         4.367 ns/op	       0 B/op	       0 allocs/op
```

## examples

examples can be found in the [examples](https://github.com/superwhiskers/crunch/blob/canon/examples) directory

## acknowledgements

icon (cookie logo) made by [freepik](https://www.freepik.com/) from [flaticon.com](https://www.flaticon.com)
