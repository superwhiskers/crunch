/*

crunch - utilities for taking bytes out of things
copyright (c) 2019 superwhiskers <whiskerdev@protonmail.com>

this program is free software: you can redistribute it and/or modify
it under the terms of the gnu lesser general public license as published by
the free software foundation, either version 3 of the license, or
(at your option) any later version.

this program is distributed in the hope that it will be useful,
but without any warranty; without even the implied warranty of
merchantability or fitness for a particular purpose.  see the
gnu lesser general public license for more details.

you should have received a copy of the gnu lesser general public license
along with this program.  if not, see <https://www.gnu.org/licenses/>.

*/

package crunch

import (
	"encoding/binary"
	"sync"
)

// ByteBuffer implements a concurrent-safe byte buffer implementation in go
type ByteBuffer struct {
	buf []byte
	off int64
	cap int64

	sync.Mutex
}

// NewByteBuffer initilaizes a new ByteBuffer with the provided byte slice(s) stored inside in the order provided
func NewByteBuffer(slices ...[]byte) (buf *ByteBuffer) {

	buf = &ByteBuffer{
		buf: []byte{},
		off: 0x00,
	}

	switch len(slices) {

	case 0:
		break

	case 1:
		buf.buf = slices[0]
		break

	default:
		for _, s := range slices {

			buf.buf = append(buf.buf, s...)

		}

	}

	buf.refresh()

	return

}

/* internal use methods */

// write writes a slice of bytes to the buffer at the specified offset
func (b *ByteBuffer) write(off int64, data []byte) {

	if (off + int64(len(data))) > b.cap {

		panic(ByteBufferOverwriteError)

	}

	b.Lock()

	for i, byt := range data {

		b.buf[off+int64(i)] = byt

	}

	b.Unlock()

	return

}

// writeComplex writes a slice of bytes to the buffer at the specified offset with the specified endianness and integer type
func (b *ByteBuffer) writeComplex(off int64, idata interface{}, size IntegerSize, endianness Endianness) {

	var data []byte
	switch size {

	case Unsigned8:
		// literally just a byte array
		// if you did this, you should probably be using the regular write methods bc those are more efficient than this one
		data = idata.([]byte)
		break

	case Unsigned16:
		var tdata []byte
		adata := idata.([]uint16)
		data = make([]byte, 2*len(adata))

		switch endianness {

		case LittleEndian:
			for i := 0; i < len(adata); i++ {

				tdata = []byte{0, 0}
				binary.LittleEndian.PutUint16(tdata, adata[i])

				data[0+(i*2)] = tdata[0]
				data[1+(i*2)] = tdata[1]

			}
			break

		case BigEndian:
			for i := 0; i < len(adata); i++ {

				tdata = []byte{0, 0}
				binary.BigEndian.PutUint16(tdata, adata[i])

				data[0+(i*2)] = tdata[0]
				data[1+(i*2)] = tdata[1]

			}
			break

		default:
			panic(ByteBufferInvalidEndiannessError)

		}
		break

	case Unsigned32:
		var tdata []byte
		adata := idata.([]uint32)
		data = make([]byte, 4*len(adata))

		switch endianness {

		case LittleEndian:
			for i := 0; i < len(adata); i++ {

				tdata = []byte{0, 0, 0, 0}
				binary.LittleEndian.PutUint32(tdata, adata[i])

				data[0+(i*4)] = tdata[0]
				data[1+(i*4)] = tdata[1]
				data[2+(i*4)] = tdata[2]
				data[3+(i*4)] = tdata[3]

			}
			break

		case BigEndian:
			for i := 0; i < len(adata); i++ {

				tdata = []byte{0, 0, 0, 0}
				binary.BigEndian.PutUint32(tdata, adata[i])

				data[0+(i*4)] = tdata[0]
				data[1+(i*4)] = tdata[1]
				data[2+(i*4)] = tdata[2]
				data[3+(i*4)] = tdata[3]

			}
			break

		default:
			panic(ByteBufferInvalidEndiannessError)

		}
		break

	case Unsigned64:
		var tdata []byte
		adata := idata.([]uint64)
		data = make([]byte, 8*len(adata))

		switch endianness {

		case LittleEndian:
			for i := 0; i < len(adata); i++ {

				tdata = []byte{0, 0, 0, 0, 0, 0, 0, 0}
				binary.LittleEndian.PutUint64(tdata, adata[i])

				data[0+(i*8)] = tdata[0]
				data[1+(i*8)] = tdata[1]
				data[2+(i*8)] = tdata[2]
				data[3+(i*8)] = tdata[3]
				data[4+(i*8)] = tdata[4]
				data[5+(i*8)] = tdata[5]
				data[6+(i*8)] = tdata[6]
				data[7+(i*8)] = tdata[7]

			}
			break

		case BigEndian:
			for i := 0; i < len(adata); i++ {

				tdata = []byte{0, 0, 0, 0, 0, 0, 0, 0}
				binary.BigEndian.PutUint64(tdata, adata[i])

				data[0+(i*8)] = tdata[0]
				data[1+(i*8)] = tdata[1]
				data[2+(i*8)] = tdata[2]
				data[3+(i*8)] = tdata[3]
				data[4+(i*8)] = tdata[4]
				data[5+(i*8)] = tdata[5]
				data[6+(i*8)] = tdata[6]
				data[7+(i*8)] = tdata[7]

			}
			break

		default:
			panic(ByteBufferInvalidEndiannessError)

		}
		break

	default:
		panic(ByteBufferInvalidIntegerSizeError)

	}

	b.write(off, data)

}

// read reads n bytes from the buffer at the specified offset
func (b *ByteBuffer) read(off, n int64) []byte {

	if (off + n) > b.cap {

		panic(ByteBufferOverreadError)

	}

	b.Lock()
	defer b.Unlock()

	return b.buf[off : off+n]

}

// readComplex reads a slice of bytes from the buffer at the specified offset with the specified endianness and integer type
func (b *ByteBuffer) readComplex(off, n int64, size IntegerSize, endianness Endianness) interface{} {

	data := b.read(off, n)

	switch size {

	case Unsigned8:
		return data

	case Unsigned16:
		idata := make([]uint16, n)

		switch endianness {

		case LittleEndian:
			for i := int64(0); i < n; i++ {

				idata[i] = binary.LittleEndian.Uint16(data[i*2 : (i+1)*2])

			}
			break

		case BigEndian:
			for i := int64(0); i < n; i++ {

				idata[i] = binary.BigEndian.Uint16(data[i*2 : (i+1)*2])

			}
			break

		default:
			panic(ByteBufferInvalidEndiannessError)

		}

		return idata

	case Unsigned32:
		idata := make([]uint32, n)

		switch endianness {

		case LittleEndian:
			for i := int64(0); i < n; i++ {

				idata[i] = binary.LittleEndian.Uint32(data[i*4 : (i+1)*4])

			}
			break

		case BigEndian:
			for i := int64(0); i < n; i++ {

				idata[i] = binary.BigEndian.Uint32(data[i*4 : (i+1)*4])

			}
			break

		default:
			panic(ByteBufferInvalidEndiannessError)

		}

		return idata

	case Unsigned64:
		idata := make([]uint64, n)

		switch endianness {

		case LittleEndian:
			for i := int64(0); i < n; i++ {

				idata[i] = binary.LittleEndian.Uint64(data[i*8 : (i+1)*8])

			}
			break

		case BigEndian:
			for i := int64(0); i < n; i++ {

				idata[i] = binary.BigEndian.Uint64(data[(i * 8) : (i+1)*8])

			}
			break

		default:
			panic(ByteBufferInvalidEndiannessError)

		}

		return idata

	default:
		panic(ByteBufferInvalidIntegerSizeError)

	}

}

// grow grows the buffer by n bytes
func (b *ByteBuffer) grow(n int64) {

	b.Lock()

	b.buf = append(b.buf, make([]byte, n)...)

	b.Unlock()

	b.refresh()

	return

}

// refresh updates the internal statistics of the byte buffer forcefully
func (b *ByteBuffer) refresh() {

	b.Lock()
	defer b.Unlock()

	b.cap = int64(len(b.buf))

	return

}

// seek seeks to position off of the byte buffer relative to the current position or exact
func (b *ByteBuffer) seek(off int64, relative bool) {

	b.Lock()
	defer b.Unlock()

	if relative == true {

		b.off = b.off + off

	} else {

		b.off = off

	}

	return

}

// after returns the amount of bytes located after the current position or the specified one
func (b *ByteBuffer) after(off ...int64) int64 {

	if len(off) == 0 {

		return b.cap - b.off

	}
	return b.cap - off[0]

}

/* public methods */

// Bytes returns the internal byte slice of the buffer
func (b *ByteBuffer) Bytes() []byte {

	return b.buf

}

// Capacity returns the capacity of the buffer
func (b *ByteBuffer) Capacity() int64 {

	return b.cap

}

// Offset returns the current offset
func (b *ByteBuffer) Offset() int64 {

	return b.off

}

// Refresh updates the cached internal statistics of the byte buffer forcefully
func (b *ByteBuffer) Refresh() {

	b.refresh()
	return

}

// Grow makes the buffer's capacity bigger by n bytes
func (b *ByteBuffer) Grow(n int64) {

	b.grow(n)
	return

}

// Seek seeks to position off of the byte buffer relative to the current position or exact
func (b *ByteBuffer) Seek(off int64, relative bool) {

	b.seek(off, relative)
	return

}

// After returns the amount of bytes located after the current position or the specified one
func (b *ByteBuffer) After(off ...int64) int64 {

	return b.after(off...)

}

// ReadByte returns the next byte from the specified offset without modifying the internal offset value
func (b *ByteBuffer) ReadByte(off int64) byte {

	return b.read(off, 1)[0]

}

// ReadBytes returns the next n bytes from the specified offset without modifying the internal offset value
func (b *ByteBuffer) ReadBytes(off, n int64) []byte {

	return b.read(off, n)

}

// ReadComplex returns the next n uint8/uint16/uint32/uint64-s from the specified offset without modifying the internal offset value
func (b *ByteBuffer) ReadComplex(off, n int64, size IntegerSize, endianness Endianness) interface{} {

	return b.readComplex(off, n, size, endianness)

}

// ReadByteNext returns the next byte from the current offset and moves the offset forward a byte
func (b *ByteBuffer) ReadByteNext() (out byte) {

	out = b.read(b.off, 1)[0]
	b.seek(1, true)
	return

}

// ReadBytesNext returns the next n bytes from the current offset and moves the offset forward the amount of bytes read
func (b *ByteBuffer) ReadBytesNext(n int64) (out []byte) {

	out = b.read(b.off, n)
	b.seek(n, true)
	return

}

// ReadComplexNext returns the next n uint8/uint16/uint32/uint64-s from the current offset and moves the offset forward the amount of bytes read
func (b *ByteBuffer) ReadComplexNext(n int64, size IntegerSize, endianness Endianness) (out interface{}) {

	out = b.readComplex(b.off, n, size, endianness)
	b.seek(n*int64(size), true)
	return

}

// WriteByte writes a byte to the buffer at the specified offset without modifying the internal offset value
func (b *ByteBuffer) WriteByte(off int64, data byte) {

	b.write(off, []byte{data})
	return

}

// WriteBytes writes bytes to the buffer at the specified offset without modifying the internal offset value
func (b *ByteBuffer) WriteBytes(off int64, data []byte) {

	b.write(off, data)
	return

}

// WriteComplex writes a uint8/uint16/uint32/uint64 to the buffer at the specified offset without modifying the internal offset value
func (b *ByteBuffer) WriteComplex(off int64, data interface{}, size IntegerSize, endianness Endianness) {

	b.writeComplex(off, data, size, endianness)
	return

}

// WriteByteNext writes a byte to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *ByteBuffer) WriteByteNext(data byte) {

	b.write(b.off, []byte{data})
	b.seek(1, true)
	return

}

// WriteBytesNext writes bytes to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *ByteBuffer) WriteBytesNext(data []byte) {

	b.write(b.off, data)
	b.seek(int64(len(data)), true)
	return

}

// WriteComplexNext writes a uint8/uint16/uint32/uint64 to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *ByteBuffer) WriteComplexNext(data interface{}, size IntegerSize, endianness Endianness) {

	b.writeComplex(b.off, data, size, endianness)

	switch size {

	case Unsigned8:
		b.seek(int64(len(data.([]uint8))*int(size)), true)

	case Unsigned16:
		b.seek(int64(len(data.([]uint16))*int(size)), true)

	case Unsigned32:
		b.seek(int64(len(data.([]uint32))*int(size)), true)

	case Unsigned64:
		b.seek(int64(len(data.([]uint64))*int(size)), true)

	}

	return

}
