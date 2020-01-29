/*

crunch - utilities for taking bytes out of things
copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>

this source code form is subject to the terms of the mozilla public
license, v. 2.0. if a copy of the mpl was not distributed with this
file, you can obtain one at http://mozilla.org/MPL/2.0/.

*/

package v1

import (
	"encoding/binary"
	"sync"
)

// MiniBuffer implements a fast and low-memory buffer type in go that handles multiple types of data easily. it is not safe
// for concurrent usage out of the box, you are required to handle that yourself by calling the Lock and Unlock methods on it.
// it also lacks the overwrite/read and underwrite/read checks that Buffer has
type MiniBuffer struct {
	buf  []byte
	off  int64
	cap  int64
	boff int64
	bcap int64

	sync.Mutex
}

// NewMiniBuffer initilaizes a new MiniBuffer with the provided byte slice(s) stored inside in the order provided
func NewMiniBuffer(out **MiniBuffer, slices ...[]byte) {

	*out = &MiniBuffer{
		buf:  []byte{},
		off:  0x00,
		boff: 0x00,
	}

	switch len(slices) {

	case 0:
		break

	case 1:
		(*out).buf = slices[0]

	default:
		for _, s := range slices {

			(*out).buf = append((*out).buf, s...)

		}

	}

	(*out).refresh()

}

/* internal use methods */

/* bitfield methods */

// readbit reads a bit from the bitfield at the specified offset
func (b *MiniBuffer) readbit(out *byte, off int64) {

	*out = atob((b.buf[off/8] & (1 << uint(7-(off%8)))) != 0)

}

// readbits reads n bits from the bitfield at the specified offset
func (b *MiniBuffer) readbits(out *uint64, off, n int64) {

	var bout byte
	for i := int64(0); i < n; i++ {

		b.readbit(&bout, off+i)
		*out = (*out << uint64(1)) | uint64(bout)

	}

}

// setbit sets a bit in the bitfield to the specified value
func (b *MiniBuffer) setbit(off int64, data byte) {

	switch data {

	case 0:
		b.buf[off/8] &= ^(1 << uint(7-(off%8)))

	case 1:
		b.buf[off/8] |= (1 << uint(7-(off%8)))

	default:
		panic(BufferInvalidBitError)

	}

}

// setbits sets n bits in the bitfield to the specified value at the specified offset
func (b *MiniBuffer) setbits(off int64, data uint64, n int64) {

	for i := int64(0); i < n; i++ {

		b.setbit(off+i, byte((data>>uint64(n-i-1))&1))

	}

}

// flipbit flips a bit in the bitfield
func (b *MiniBuffer) flipbit(off int64) {

	b.buf[off/8] ^= (1 << uint(7-(off%8)))

}

// clearallbits sets all of the buffer's bits to 0
func (b *MiniBuffer) clearallbits() {

	for i := range b.buf {

		b.buf[i] = 0

	}

}

// setallbits sets all of the buffer's bits to 1
func (b *MiniBuffer) setallbits() {

	for i := range b.buf {

		b.buf[i] = 0xFF

	}

}

// flipallbits flips all of the buffer's bits
func (b *MiniBuffer) flipallbits() {

	for i := range b.buf {

		b.buf[i] = ^b.buf[i]

	}

}

// seekbit seeks to position off of the bitfield relative to the current position or exact
func (b *MiniBuffer) seekbit(off int64, relative bool) {

	if relative {

		b.boff += off

	} else {

		b.boff = off

	}

}

// afterbit returns the amount of bits located after the current position or the specified one
func (b *MiniBuffer) afterbit(out *int64, off ...int64) {

	if len(off) == 0 {

		*out = b.bcap - b.boff - 1
		return

	}
	*out = b.bcap - off[0] - 1

}

/* byte buffer methods */

// write writes a slice of bytes to the buffer at the specified offset
func (b *MiniBuffer) write(off int64, data []byte) {

	i := int64(0)
	n := int64(len(data))

	{
	write_loop:
		b.buf[off+i] = data[i]
		i++
		if i < n {

			goto write_loop

		}
	}

}

// writeComplex writes a slice of bytes to the buffer at the specified offset with the specified endianness and integer type
func (b *MiniBuffer) writeComplex(off int64, idata interface{}, size IntegerSize, endian binary.ByteOrder) {

	var (
		data  []byte
		tdata []byte
	)

	switch size {

	case Unsigned8:
		data = idata.([]byte)

	case Unsigned16:
		adata := idata.([]uint16)
		data = make([]byte, 2*len(adata))
		for i := 0; i < len(adata); i++ {

			tdata = []byte{0x00, 0x00}
			endian.PutUint16(tdata, adata[i])

			data[0+(i*2)] = tdata[0]
			data[1+(i*2)] = tdata[1]

		}

	case Unsigned32:
		adata := idata.([]uint32)
		data = make([]byte, 4*len(adata))
		for i := 0; i < len(adata); i++ {

			tdata = []byte{0x00, 0x00, 0x00, 0x00}
			endian.PutUint32(tdata, adata[i])

			data[0+(i*4)] = tdata[0]
			data[1+(i*4)] = tdata[1]
			data[2+(i*4)] = tdata[2]
			data[3+(i*4)] = tdata[3]

		}

	case Unsigned64:
		adata := idata.([]uint64)
		data = make([]byte, 8*len(adata))
		for i := 0; i < len(adata); i++ {

			tdata = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
			endian.PutUint64(tdata, adata[i])

			data[0+(i*8)] = tdata[0]
			data[1+(i*8)] = tdata[1]
			data[2+(i*8)] = tdata[2]
			data[3+(i*8)] = tdata[3]
			data[4+(i*8)] = tdata[4]
			data[5+(i*8)] = tdata[5]
			data[6+(i*8)] = tdata[6]
			data[7+(i*8)] = tdata[7]

		}

	default:
		panic(BufferInvalidIntegerSizeError)

	}

	b.write(off, data)

}

// read reads n bytes from the buffer at the specified offset
func (b *MiniBuffer) read(out *[]byte, off, n int64) {

	*out = b.buf[off : off+n]

}

// readComplex reads a slice of bytes from the buffer at the specified offset with the specified endianness and integer type
func (b *MiniBuffer) readComplex(out *interface{}, off, n int64, size IntegerSize, endian binary.ByteOrder) {

	var data []byte
	b.read(&data, off, n)

	switch size {

	case Unsigned8:
		*out = data

	case Unsigned16:
		*out = make([]uint16, n)

		for i := int64(0); i < n; i++ {

			(*out).([]uint16)[i] = endian.Uint16(data[i*2 : (i+1)*2])

		}

	case Unsigned32:
		*out = make([]uint32, n)

		for i := int64(0); i < n; i++ {

			(*out).([]uint32)[i] = endian.Uint32(data[i*4 : (i+1)*4])

		}

	case Unsigned64:
		*out = make([]uint64, n)

		for i := int64(0); i < n; i++ {

			(*out).([]uint64)[i] = endian.Uint64(data[i*8 : (i+1)*8])

		}

	default:
		panic(BufferInvalidIntegerSizeError)

	}

}

// seek seeks to position off of the byte buffer relative to the current position or exact
func (b *MiniBuffer) seek(off int64, relative bool) {

	if relative {

		b.off += off

	} else {

		b.off = off

	}

}

// after returns the amount of bytes located after the current position or the specified one
func (b *MiniBuffer) after(out *int64, off ...int64) {

	if len(off) == 0 {

		*out = b.cap - b.off - 1
		return

	}
	*out = b.cap - off[0] - 1

}

/* generic methods */

// grow grows the buffer by n bytes
func (b *MiniBuffer) grow(n int64) {

	b.buf = append(b.buf, make([]byte, n)...)
	b.refresh()

}

// refresh updates the internal statistics of the byte buffer forcefully
func (b *MiniBuffer) refresh() {

	b.cap = int64(len(b.buf))
	b.bcap = b.cap * 8

}

// alignbit aligns the bit offset to the byte offset
func (b *MiniBuffer) alignbit() {

	b.boff = b.off * 8

}

// alignbyte aligns the byte offset to the bit offset
func (b *MiniBuffer) alignbyte() {

	b.off = b.boff / 8

}

/* public methods */

// Bytes stores the internal byte slice of the buffer in out
func (b *MiniBuffer) Bytes(out *[]byte) {

	*out = b.buf

}

// Capacity stores the capacity of the buffer in out
func (b *MiniBuffer) Capacity(out *int64) {

	*out = b.cap

}

// BitCapacity stores the bit capacity of the buffer in out
func (b *MiniBuffer) BitCapacity(out *int64) {

	*out = b.bcap

}

// Offset stores the current offset of the buffer in out
func (b *MiniBuffer) Offset(out *int64) {

	*out = b.off

}

// BitOffset stores the current bit offset of the buffer in out
func (b *MiniBuffer) BitOffset(out *int64) {

	*out = b.boff

}

// Refresh updates the cached internal statistics of the buffer forcefully
func (b *MiniBuffer) Refresh() {

	b.refresh()

}

// Grow makes the buffer's capacity bigger by n bytes
func (b *MiniBuffer) Grow(n int64) {

	b.grow(n)

}

// Seek seeks to position off of the buffer relative to the current position or exact
func (b *MiniBuffer) Seek(off int64, relative bool) {

	b.seek(off, relative)

}

// SeekBit seeks to bit position off of the the buffer relative to the current position or exact
func (b *MiniBuffer) SeekBit(off int64, relative bool) {

	b.seekbit(off, relative)

}

// AlignBit aligns the bit offset to the byte offset
func (b *MiniBuffer) AlignBit() {

	b.alignbit()

}

// AlignByte aligns the byte offset to the bit offset
func (b *MiniBuffer) AlignByte() {

	b.alignbyte()

}

// After stores the amount of bytes located after the current position or the specified one in out
func (b *MiniBuffer) After(out *int64, off ...int64) {

	b.after(out, off...)

}

// AfterBit stores the amount of bits located after the current bit position or the specified one in out
func (b *MiniBuffer) AfterBit(out *int64, off ...int64) {

	b.afterbit(out, off...)

}

// ReadBytes stores the next n bytes from the specified offset without modifying the internal offset value in out
func (b *MiniBuffer) ReadBytes(out *[]byte, off, n int64) {

	b.read(out, off, n)

}

// ReadComplex stores the next n uint8/uint16/uint32/uint64-s from the specified offset without modifying the internal offset value
func (b *MiniBuffer) ReadComplex(out *interface{}, off, n int64, size IntegerSize, endian binary.ByteOrder) {

	b.readComplex(out, off, n, size, endian)

}

// ReadBytesNext stores the next n bytes from the current offset and moves the offset forward the amount of bytes read in out
func (b *MiniBuffer) ReadBytesNext(out *[]byte, n int64) {

	b.read(out, b.off, n)
	b.seek(n, true)

}

// ReadComplexNext stores the next n uint8/uint16/uint32/uint64-s from the current offset and moves the offset forward the amount of bytes read
func (b *MiniBuffer) ReadComplexNext(out *interface{}, n int64, size IntegerSize, endian binary.ByteOrder) {

	b.readComplex(out, b.off, n, size, endian)
	b.seek(n*int64(size), true)

}

// WriteBytes writes bytes to the buffer at the specified offset without modifying the internal offset value
func (b *MiniBuffer) WriteBytes(off int64, data []byte) {

	b.write(off, data)

}

// WriteComplex writes a uint8/uint16/uint32/uint64 to the buffer at the specified offset without modifying the internal offset value
func (b *MiniBuffer) WriteComplex(off int64, data interface{}, size IntegerSize, endian binary.ByteOrder) {

	b.writeComplex(off, data, size, endian)

}

// WriteBytesNext writes bytes to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteBytesNext(data []byte) {

	b.write(b.off, data)
	b.seek(int64(len(data)), true)

}

// WriteComplexNext writes a uint8/uint16/uint32/uint64 to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteComplexNext(data interface{}, size IntegerSize, endian binary.ByteOrder) {

	b.writeComplex(b.off, data, size, endian)

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

}

// ReadBit stores the bit located at the specified offset without modifying the internal offset value in out
func (b *MiniBuffer) ReadBit(out *byte, off int64) {

	b.readbit(out, off)

}

// ReadBits stores the next n bits from the specified offset without modifying the internal offset value in out
func (b *MiniBuffer) ReadBits(out *uint64, off, n int64) {

	b.readbits(out, off, n)

}

// ReadBitNext stores the next bit from the current offset and moves the offset forward a bit in out
func (b *MiniBuffer) ReadBitNext(out *byte) {

	b.readbit(out, b.boff)
	b.seekbit(1, true)

}

// ReadBitsNext stores the next n bits from the current offset and moves the offset forward the amount of bits read in out
func (b *MiniBuffer) ReadBitsNext(out *uint64, n int64) {

	b.readbits(out, b.boff, n)
	b.seekbit(n, true)

}

// SetBit sets the bit located at the specified offset without modifying the internal offset value
func (b *MiniBuffer) SetBit(off int64, data byte) {

	b.setbit(off, data)

}

// SetBits sets the next n bits from the specified offset without modifying the internal offset value
func (b *MiniBuffer) SetBits(off int64, data uint64, n int64) {

	b.setbits(off, data, n)

}

// SetBitNext sets the next bit from the current offset and moves the offset forward a bit
func (b *MiniBuffer) SetBitNext(data byte) {

	b.setbit(b.boff, data)
	b.seekbit(1, true)

}

// SetBitsNext sets the next n bits from the current offset and moves the offset forward the amount of bits set
func (b *MiniBuffer) SetBitsNext(data uint64, n int64) {

	b.setbits(b.boff, data, n)
	b.seekbit(n, true)

}

// FlipBit flips the bit located at the specified offset without modifying the internal offset value
func (b *MiniBuffer) FlipBit(off int64) {

	b.flipbit(off)

}

// FlipBitNext flips the next bit from the current offset and moves the offset forward a bit
func (b *MiniBuffer) FlipBitNext() {

	b.flipbit(b.boff)
	b.seekbit(1, true)

}

// ClearAllBits sets all of the buffer's bits to 0
func (b *MiniBuffer) ClearAllBits() {

	b.clearallbits()

}

// SetAllBits sets all of the buffer's bits to 1
func (b *MiniBuffer) SetAllBits() {

	b.setallbits()

}

// FlipAllBits flips all of the buffer's bits
func (b *MiniBuffer) FlipAllBits() {

	b.flipallbits()

}
