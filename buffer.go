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

// Buffer implements a concurrent-safe buffer type in go that handles multiple types of data
type Buffer struct {
	buf  []byte
	off  int64
	cap  int64
	boff int64
	bcap int64

	sync.Mutex
}

// NewBuffer initilaizes a new Buffer with the provided byte slice(s) stored inside in the order provided
func NewBuffer(slices ...[]byte) (buf *Buffer) {

	buf = &Buffer{
		buf:  []byte{},
		off:  0x00,
		boff: 0x00,
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

/* bitfield methods */

// readbit reads a bit from the bitfield at the specified offset
func (b *Buffer) readbit(off int64) byte {

	if off > (b.bcap - 1) {

		panic(BitfieldOverreadError)

	}

	b.Lock()
	defer b.Unlock()

	i, o := (off / 8), uint(7-(off%8))
	return atob((b.buf[i] & (1 << o)) != 0)

}

// readbits reads n bits from the bitfield at the specified offset
func (b *Buffer) readbits(off, n int64) (v uint64) {

	if (off + n) > b.bcap {

		panic(BitfieldOverreadError)

	}

	for i := int64(0); i < n; i++ {

		v = (v << uint64(1)) | uint64(b.readbit(off+i))

	}

	return

}

// setbit sets a bit in the bitfield to the specified value
func (b *Buffer) setbit(off, data int64) {

	if off > (b.bcap - 1) {

		panic(BitfieldOverwriteError)

	}

	if data != 1 && data != 0 {

		panic(BitfieldInvalidBitError)

	}

	b.Lock()
	defer b.Unlock()

	i, o := (off / 8), uint(7-(off%8))
	switch data {

	case 0:
		b.buf[i] &= ^(1 << o)

	case 1:
		b.buf[i] |= (1 << o)

	}

}

// setbits sets n bits in the bitfield to the specified value at the specified offset
func (b *Buffer) setbits(off, data, n int64) {

	if off+n > (b.bcap - 1) {

		panic(BitfieldOverwriteError)

	}

	for i := int64(0); i < n; i++ {

		b.setbit(off+i, (data>>uint64(n-i-1))&1)

	}

}

// flipbit flips a bit in the bitfield
func (b *Buffer) flipbit(off int64) {

	if off > (b.bcap - 1) {

		panic(BitfieldOverwriteError)

	}

	b.Lock()
	defer b.Unlock()

	i, o := (off / 8), uint(7-(off%8))
	b.buf[i] ^= (1 << o)

}

// clearallbits sets all of the buffer's bits to 0
func (b *Buffer) clearallbits() {

	b.Lock()
	defer b.Unlock()

	for i := range b.buf {

		b.buf[i] = 0

	}

}

// setallbits sets all of the buffer's bits to 1
func (b *Buffer) setallbits() {

	b.Lock()
	defer b.Unlock()

	for i := range b.buf {

		b.buf[i] = 0xFF

	}

}

// flipallbits flips all of the buffer's bits
func (b *Buffer) flipallbits() {

	b.Lock()
	defer b.Unlock()

	for i := range b.buf {

		b.buf[i] = ^b.buf[i]

	}

}

// seekbit seeks to position off of the bitfield relative to the current position or exact
func (b *Buffer) seekbit(off int64, relative bool) {

	b.Lock()
	defer b.Unlock()

	if relative == true {

		b.boff = b.boff + off

	} else {

		b.boff = off

	}

	return

}

// afterbit returns the amount of bits located after the current position or the specified one
func (b *Buffer) afterbit(off ...int64) int64 {

	if len(off) == 0 {

		return b.bcap - b.boff

	}
	return b.bcap - off[0]

}

/* byte buffer methods */

// write writes a slice of bytes to the buffer at the specified offset
func (b *Buffer) write(off int64, data []byte) {

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
func (b *Buffer) writeComplex(off int64, idata interface{}, size IntegerSize, endianness Endianness) {

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
func (b *Buffer) read(off, n int64) []byte {

	if (off + n) > b.cap {

		panic(ByteBufferOverreadError)

	}

	b.Lock()
	defer b.Unlock()

	return b.buf[off : off+n]

}

// readComplex reads a slice of bytes from the buffer at the specified offset with the specified endianness and integer type
func (b *Buffer) readComplex(off, n int64, size IntegerSize, endianness Endianness) interface{} {

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

// seek seeks to position off of the byte buffer relative to the current position or exact
func (b *Buffer) seek(off int64, relative bool) {

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
func (b *Buffer) after(off ...int64) int64 {

	if len(off) == 0 {

		return b.cap - b.off

	}
	return b.cap - off[0]

}

/* generic methods */

// grow grows the buffer by n bytes
func (b *Buffer) grow(n int64) {

	b.Lock()

	b.buf = append(b.buf, make([]byte, n)...)

	b.Unlock()

	b.refresh()

	return

}

// refresh updates the internal statistics of the byte buffer forcefully
func (b *Buffer) refresh() {

	b.Lock()
	defer b.Unlock()

	b.cap = int64(len(b.buf))
	b.bcap = b.cap * 8

	return

}

// alignbit aligns the bit offset to the byte offset
func (b *Buffer) alignbit() {

	b.Lock()
	defer b.Unlock()

	b.boff = b.off * 8

}

// alignbyte aligns the byte offset to the bit offset
func (b *Buffer) alignbyte() {

	b.Lock()
	defer b.Unlock()

	b.off = b.boff / 8

}

/* public methods */

// Bytes returns the internal byte slice of the buffer
func (b *Buffer) Bytes() []byte {

	return b.buf

}

// Capacity returns the capacity of the buffer
func (b *Buffer) Capacity() int64 {

	return b.cap

}

// BitCapacity returns the bit capacity of the buffer
func (b *Buffer) BitCapacity() int64 {

	return b.bcap

}

// Offset returns the current offset
func (b *Buffer) Offset() int64 {

	return b.off

}

// BitOffset returns the current bit offset
func (b *Buffer) BitOffset() int64 {

	return b.boff

}

// Refresh updates the cached internal statistics of the buffer forcefully
func (b *Buffer) Refresh() {

	b.refresh()
	return

}

// Grow makes the buffer's capacity bigger by n bytes
func (b *Buffer) Grow(n int64) {

	b.grow(n)
	return

}

// Seek seeks to position off of the buffer relative to the current position or exact
func (b *Buffer) Seek(off int64, relative bool) {

	b.seek(off, relative)
	return

}

// SeekBit seeks to bit position off of the the buffer relative to the current position or exact
func (b *Buffer) SeekBit(off int64, relative bool) {

	b.seekbit(off, relative)
	return

}

// AlignBit aligns the bit offset to the byte offset
func (b *Buffer) AlignBit() {

	b.alignbit()
	return

}

// AlignByte aligns the byte offset to the bit offset
func (b *Buffer) AlignByte() {

	b.alignbyte()
	return

}

// After returns the amount of bytes located after the current position or the specified one
func (b *Buffer) After(off ...int64) int64 {

	return b.after(off...)

}

// AfterBit returns the amount of bits located after the current bit position or the specified one
func (b *Buffer) AfterBit(off ...int64) int64 {

	return b.afterbit(off...)

}

// ReadBytes returns the next n bytes from the specified offset without modifying the internal offset value
func (b *Buffer) ReadBytes(off, n int64) []byte {

	return b.read(off, n)

}

// ReadComplex returns the next n uint8/uint16/uint32/uint64-s from the specified offset without modifying the internal offset value
func (b *Buffer) ReadComplex(off, n int64, size IntegerSize, endianness Endianness) interface{} {

	return b.readComplex(off, n, size, endianness)

}

// ReadBytesNext returns the next n bytes from the current offset and moves the offset forward the amount of bytes read
func (b *Buffer) ReadBytesNext(n int64) (out []byte) {

	out = b.read(b.off, n)
	b.seek(n, true)
	return

}

// ReadComplexNext returns the next n uint8/uint16/uint32/uint64-s from the current offset and moves the offset forward the amount of bytes read
func (b *Buffer) ReadComplexNext(n int64, size IntegerSize, endianness Endianness) (out interface{}) {

	out = b.readComplex(b.off, n, size, endianness)
	b.seek(n*int64(size), true)
	return

}

// WriteByte writes a byte to the buffer at the specified offset without modifying the internal offset value
func (b *Buffer) WriteByte(off int64, data byte) {

	b.write(off, []byte{data})
	return

}

// WriteBytes writes bytes to the buffer at the specified offset without modifying the internal offset value
func (b *Buffer) WriteBytes(off int64, data []byte) {

	b.write(off, data)
	return

}

// WriteComplex writes a uint8/uint16/uint32/uint64 to the buffer at the specified offset without modifying the internal offset value
func (b *Buffer) WriteComplex(off int64, data interface{}, size IntegerSize, endianness Endianness) {

	b.writeComplex(off, data, size, endianness)
	return

}

// WriteByteNext writes a byte to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *Buffer) WriteByteNext(data byte) {

	b.write(b.off, []byte{data})
	b.seek(1, true)
	return

}

// WriteBytesNext writes bytes to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *Buffer) WriteBytesNext(data []byte) {

	b.write(b.off, data)
	b.seek(int64(len(data)), true)
	return

}

// WriteComplexNext writes a uint8/uint16/uint32/uint64 to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *Buffer) WriteComplexNext(data interface{}, size IntegerSize, endianness Endianness) {

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

// ReadBit returns the bit located at the specified offset without modifying the internal offset value
func (b *Buffer) ReadBit(off int64) byte {

	return b.readbit(off)

}

// ReadBits returns the next n bits from the specified offset without modifying the internal offset value
func (b *Buffer) ReadBits(off, n int64) uint64 {

	return b.readbits(off, n)

}

// ReadBitNext returns the next bit from the current offset and moves the offset forward a bit
func (b *Buffer) ReadBitNext() (out byte) {

	out = b.readbit(b.boff)
	b.seekbit(1, true)
	return

}

// ReadBitsNext returns the next n bits from the current offset and moves the offset forward the amount of bits read
func (b *Buffer) ReadBitsNext(n int64) (out uint64) {

	out = b.readbits(b.boff, n)
	b.seekbit(n, true)
	return

}

// SetBit sets the bit located at the specified offset without modifying the internal offset value
func (b *Buffer) SetBit(off, data int64) {

	b.setbit(off, data)
	return

}

// SetBits sets the next n bits from the specified offset without modifying the internal offset value
func (b *Buffer) SetBits(off, data, n int64) {

	b.setbits(off, data, n)
	return

}

// SetBitNext sets the next bit from the current offset and moves the offset forward a bit
func (b *Buffer) SetBitNext(data int64) {

	b.setbit(b.boff, data)
	b.seekbit(1, true)
	return

}

// SetBitsNext sets the next n bits from the current offset and moves the offset forward the amount of bits set
func (b *Buffer) SetBitsNext(data, n int64) {

	b.setbits(b.boff, data, n)
	b.seekbit(n, true)
	return

}

// FlipBit flips the bit located at the specified offset without modifying the internal offset value
func (b *Buffer) FlipBit(off int64) {

	b.flipbit(off)
	return

}

// FlipBitNext flips the next bit from the current offset and moves the offset forward a bit
func (b *Buffer) FlipBitNext() {

	b.flipbit(b.boff)
	b.seekbit(1, true)
	return

}

// ClearAllBits sets all of the buffer's bits to 0
func (b *Buffer) ClearAllBits() {

	b.clearallbits()
	return

}

// SetAllBits sets all of the buffer's bits to 1
func (b *Buffer) SetAllBits() {

	b.setallbits()
	return

}

// FlipAllBits flips all of the buffer's bits
func (b *Buffer) FlipAllBits() {

	b.flipallbits()
	return

}
