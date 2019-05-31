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
	"sync"
	"unsafe"
)

// Buffer implements a concurrent-safe buffer type in go that handles multiple types of data
type Buffer struct {
	buf  []byte
	off  int64
	cap  int64
	boff int64
	bcap int64

	// temp?
	obuf uintptr

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

	default:
		for _, s := range slices {

			buf.buf = append(buf.buf, s...)

		}

	}

	buf.Refresh()
	return

}

/* internal use methods */

/* bitfield methods */

// ReadBit returns the bit located at the specified offset without modifying the internal offset value
func (b *Buffer) ReadBit(off int64) byte {

	if off > (b.bcap - 1) {

		panic(BufferOverreadError)

	}

	if off < 0x00 {

		panic(BufferUnderreadError)

	}

	b.Lock()
	defer b.Unlock()

	return (b.buf[off/8] >> (7 - uint64(off%8))) & 1

}

// ReadBitNext returns the next bit from the current offset and moves the offset forward a bit
func (b *Buffer) ReadBitNext() (out byte) {

	out = b.ReadBit(b.boff)
	b.SeekBit(1, true)
	return

}

// ReadBits returns the next n bits from the specified offset without modifying the internal offset value
func (b *Buffer) ReadBits(off, n int64) (out uint64) {

	i := int64(0)

	{
	read_loop:
		out = (out << uint64(1)) | uint64(b.ReadBit(off+i))
		i++
		if i < n {

			goto read_loop

		}
	}

	return

}

// ReadBitsNext returns the next n bits from the current offset and moves the offset forward the amount of bits read
func (b *Buffer) ReadBitsNext(n int64) (out uint64) {

	out = b.ReadBits(b.boff, n)
	b.SeekBit(n, true)
	return

}

// SetBit sets the bit located at the specified offset without modifying the internal offset value
func (b *Buffer) SetBit(off int64) {

	if off > (b.bcap - 1) {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	b.Lock()
	defer b.Unlock()

	b.buf[off/8] |= (1 << uint(7-(off%8)))

}

// SetBitNext sets the next bit from the current offset and moves the offset forward a bit
func (b *Buffer) SetBitNext() {

	b.SetBit(b.boff)
	b.SeekBit(1, true)

}

// ClearBit clears the bit located at the specified offset without modifying the internal offset value
func (b *Buffer) ClearBit(off int64) {

	if off > (b.bcap - 1) {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	b.Lock()
	defer b.Unlock()

	b.buf[off/8] &= ^(1 << uint(7-(off%8)))

}

// ClearBitNext clears the next bit from the current offset and moves the offset forward a bit
func (b *Buffer) ClearBitNext() {

	b.ClearBit(b.boff)
	b.SeekBit(1, true)

}

// SetBits sets the next n bits from the specified offset without modifying the internal offset value
func (b *Buffer) SetBits(off int64, data uint64, n int64) {

	i := int64(0)

	{
	write_loop:
		if byte((data>>uint64(n-i-1))&1) == 0 {

			b.ClearBit(off + i)

		} else {

			b.SetBit(off + i)

		}
		i++
		if i < n {

			goto write_loop

		}
	}

}

// SetBitsNext sets the next n bits from the current offset and moves the offset forward the amount of bits set
func (b *Buffer) SetBitsNext(data uint64, n int64) {

	b.SetBits(b.boff, data, n)
	b.SeekBit(n, true)

}

// FlipBit flips the bit located at the specified offset without modifying the internal offset value
func (b *Buffer) FlipBit(off int64) {

	if off > (b.bcap - 1) {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	b.Lock()
	defer b.Unlock()

	b.buf[off/8] ^= (1 << uint(7-(off%8)))

}

// FlipBitNext flips the next bit from the current offset and moves the offset forward a bit
func (b *Buffer) FlipBitNext() {

	b.FlipBit(b.boff)
	b.SeekBit(1, true)

}

// ClearAllBits sets all of the buffer's bits to 0
func (b *Buffer) ClearAllBits() {

	b.Lock()
	defer b.Unlock()

	var (
		i = int64(0)
		n = int64(len(b.buf))
	)
	{
	clear_loop:
		b.buf[i] = 0
		i++
		if i < n {

			goto clear_loop

		}
	}

}

// SetAllBits sets all of the buffer's bits to 1
func (b *Buffer) SetAllBits() {

	b.Lock()
	defer b.Unlock()

	var (
		i = int64(0)
		n = int64(len(b.buf))
	)
	{
	set_loop:
		b.buf[i] = 0xFF
		i++
		if i < n {

			goto set_loop

		}
	}

}

// FlipAllBits flips all of the buffer's bits
func (b *Buffer) FlipAllBits() {

	b.Lock()
	defer b.Unlock()

	var (
		i = int64(0)
		n = int64(len(b.buf))
	)
	{
	flip_loop:
		b.buf[i] = ^b.buf[i]
		i++
		if i < n {

			goto flip_loop

		}

	}

}

// SeekBit seeks to bit position off of the the buffer relative to the current position or exact
func (b *Buffer) SeekBit(off int64, relative bool) {

	b.Lock()
	defer b.Unlock()

	if relative {

		b.boff += off

	} else {

		b.boff = off

	}

}

// AfterBit returns the amount of bits located after the current bit position or the specified one
func (b *Buffer) AfterBit(off ...int64) int64 {

	if len(off) == 0 {

		return b.bcap - b.boff - 1

	}
	return b.bcap - off[0] - 1

}

// AlignBit aligns the bit offset to the byte offset
func (b *Buffer) AlignBit() {

	b.Lock()
	defer b.Unlock()

	b.boff = b.off * 8

}

/* byte buffer methods */

// WriteBytes writes bytes to the buffer at the specified offset without modifying the internal offset value
func (b *Buffer) WriteBytes(off int64, data []byte) {

	if (off + int64(len(data))) > b.cap {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	b.Lock()

	/* same as in minibuffer, leaving here incase this new method proves to be slower in some edge cases */
	/*var (
		i = int64(0)
		n = int64(len(data))
	)
	{
	write_loop:
		b.buf[off+i] = data[i]
		i++
		if i < n {

			goto write_loop

		}
	}*/

	var (
		p = uintptr(off) + b.obuf
		i = int64(0)
		n = int64(len(data))
	)
	{
	write_loop:
		*(*byte)(unsafe.Pointer(p)) = data[i]
		i++
		p++
		if i < n {

			goto write_loop

		}
	}

	b.Unlock()

}

// WriteBytesNext writes bytes to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *Buffer) WriteBytesNext(data []byte) {

	b.WriteBytes(b.off, data)
	b.SeekByte(int64(len(data)), true)

}

// WriteByte writes a byte to the buffer at the specified offset without modifying the internal offset value
func (b *Buffer) WriteByte(off int64, data byte) {

	b.WriteBytes(off, []byte{data})

}

// WriteByteNext writes a byte to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *Buffer) WriteByteNext(data byte) {

	b.WriteBytes(b.off, []byte{data})
	b.SeekByte(1, true)

}

// WriteU16LE writes a slice of uint16s to the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *Buffer) WriteU16LE(off int64, data []uint16) {

	if (off + int64(len(data))*2) > b.cap {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	b.Lock()

	var (
		i = 0
		n = len(data)
	)
	{
	write_loop:
		b.buf[off+int64(i*2)] = byte(data[i])
		b.buf[off+int64(1+(i*2))] = byte(data[i] >> 8)

		i++
		if i < n {

			goto write_loop

		}
	}

	b.Unlock()

}

// WriteU16LENext writes a slice of uint16s to the buffer at the current offset in little-endian and moves the offset forward the amount of bytes written
func (b *Buffer) WriteU16LENext(data []uint16) {

	b.WriteU16LE(b.off, data)
	b.SeekByte(int64(len(data))*2, true)

}

// WriteU16BE writes a slice of uint16s to the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *Buffer) WriteU16BE(off int64, data []uint16) {

	if (off + int64(len(data))*2) > b.cap {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	b.Lock()

	var (
		i = 0
		n = len(data)
	)
	{
	write_loop:
		b.buf[off+int64(i*2)] = byte(data[i] >> 8)
		b.buf[off+int64(1+(i*2))] = byte(data[i])

		i++
		if i < n {

			goto write_loop

		}
	}

	b.Unlock()

}

// WriteU16BENext writes a slice of uint16s to the buffer at the current offset in big-endian and moves the offset forward the amount of bytes written
func (b *Buffer) WriteU16BENext(data []uint16) {

	b.WriteU16BE(b.off, data)
	b.SeekByte(int64(len(data))*2, true)

}

// WriteU32LE writes a slice of uint32s to the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *Buffer) WriteU32LE(off int64, data []uint32) {

	if (off + int64(len(data))*4) > b.cap {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	b.Lock()

	var (
		i = 0
		n = len(data)
	)
	{
	write_loop:
		b.buf[off+int64(i*4)] = byte(data[i])
		b.buf[off+int64(1+(i*4))] = byte(data[i] >> 8)
		b.buf[off+int64(2+(i*4))] = byte(data[i] >> 16)
		b.buf[off+int64(3+(i*4))] = byte(data[i] >> 24)

		i++
		if i < n {

			goto write_loop

		}
	}

	b.Unlock()

}

// WriteU32LENext writes a slice of uint32s to the buffer at the current offset in little-endian and moves the offset forward the amount of bytes written
func (b *Buffer) WriteU32LENext(data []uint32) {

	b.WriteU32LE(b.off, data)
	b.SeekByte(int64(len(data))*4, true)

}

// WriteU32BE writes a slice of uint32s to the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *Buffer) WriteU32BE(off int64, data []uint32) {

	if (off + int64(len(data))*4) > b.cap {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	b.Lock()

	var (
		i = 0
		n = len(data)
	)
	{
	write_loop:
		b.buf[off+int64(i*4)] = byte(data[i] >> 24)
		b.buf[off+int64(1+(i*4))] = byte(data[i] >> 16)
		b.buf[off+int64(2+(i*4))] = byte(data[i] >> 8)
		b.buf[off+int64(3+(i*4))] = byte(data[i])

		i++
		if i < n {

			goto write_loop

		}
	}

	b.Unlock()

}

// WriteU32BENext writes a slice of uint32s to the buffer at the current offset in big-endian and moves the offset forward the amount of bytes written
func (b *Buffer) WriteU32BENext(data []uint32) {

	b.WriteU32BE(b.off, data)
	b.SeekByte(int64(len(data))*4, true)

}

// WriteU64LE writes a slice of uint64s to the buffer at the specfied offset in little-endian without modifying the internal offset value
func (b *Buffer) WriteU64LE(off int64, data []uint64) {

	if (off + int64(len(data))*8) > b.cap {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	b.Lock()

	var (
		i = 0
		n = len(data)
	)
	{
	write_loop:
		b.buf[off+int64(i*8)] = byte(data[i])
		b.buf[off+int64(1+(i*8))] = byte(data[i] >> 8)
		b.buf[off+int64(2+(i*8))] = byte(data[i] >> 16)
		b.buf[off+int64(3+(i*8))] = byte(data[i] >> 24)
		b.buf[off+int64(4+(i*8))] = byte(data[i] >> 32)
		b.buf[off+int64(5+(i*8))] = byte(data[i] >> 40)
		b.buf[off+int64(6+(i*8))] = byte(data[i] >> 48)
		b.buf[off+int64(7+(i*8))] = byte(data[i] >> 56)

		i++
		if i < n {

			goto write_loop

		}
	}

	b.Unlock()

}

// WriteU64LENext writes a slice of uint64s to the buffer at the current offset in little-endian and moves the offset forward the amount of bytes written
func (b *Buffer) WriteU64LENext(data []uint64) {

	b.WriteU64LE(b.off, data)
	b.SeekByte(int64(len(data))*8, true)

}

// WriteU64BE writes a slice of uint64s to the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *Buffer) WriteU64BE(off int64, data []uint64) {

	if (off + int64(len(data))*8) > b.cap {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	b.Lock()

	var (
		i = 0
		n = len(data)
	)
	{
	write_loop:
		b.buf[off+int64(i*8)] = byte(data[i] >> 56)
		b.buf[off+int64(1+(i*8))] = byte(data[i] >> 48)
		b.buf[off+int64(2+(i*8))] = byte(data[i] >> 40)
		b.buf[off+int64(3+(i*8))] = byte(data[i] >> 32)
		b.buf[off+int64(4+(i*8))] = byte(data[i] >> 24)
		b.buf[off+int64(5+(i*8))] = byte(data[i] >> 16)
		b.buf[off+int64(6+(i*8))] = byte(data[i] >> 8)
		b.buf[off+int64(7+(i*8))] = byte(data[i])

		i++
		if i < n {

			goto write_loop

		}
	}

	b.Unlock()

}

// WriteU64BENext writes a slice of uint64s to the buffer at the current offset in big-endian and moves the offset forward the amount of bytes written
func (b *Buffer) WriteU64BENext(data []uint64) {

	b.WriteU64BE(b.off, data)
	b.SeekByte(int64(len(data))*8, true)

}

// ReadBytes returns the next n bytes from the specified offset without modifying the internal offset value
func (b *Buffer) ReadBytes(off, n int64) []byte {

	if (off + n) > b.cap {

		panic(BufferOverreadError)

	}

	if off < 0x00 {

		panic(BufferUnderreadError)

	}

	b.Lock()
	defer b.Unlock()

	return b.buf[off : off+n]

}

// ReadBytesNext returns the next n bytes from the current offset and moves the offset forward the amount of bytes read
func (b *Buffer) ReadBytesNext(n int64) (out []byte) {

	out = b.ReadBytes(b.off, n)
	b.SeekByte(n, true)
	return

}

// ReadByte returns the next byte from the specified offset without modifying the internal offset value
func (b *Buffer) ReadByte(off int64) byte {

	return b.ReadBytes(off, 1)[0]

}

// ReadByteNext returns the next byte from the current offset and moves the offset forward a byte
func (b *Buffer) ReadByteNext() (out byte) {

	out = b.ReadBytes(b.off, 1)[0]
	b.SeekByte(1, true)
	return

}

// ReadU16LE reads a slice of uint16s from the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *Buffer) ReadU16LE(off, n int64) (out []uint16) {

	if (off + n*2) > b.cap {

		panic(BufferOverreadError)

	}

	if off < 0x00 {

		panic(BufferUnderreadError)

	}

	out = make([]uint16, n)

	b.Lock()
	defer b.Unlock()

	i := int64(0)
	{
	read_loop:
		out[i] = uint16(b.buf[off+(i*2)]) |
			uint16(b.buf[off+(1+(i*2))])<<8

		i++
		if i < n {

			goto read_loop

		}
	}

	return

}

// ReadU16LENext reads a slice of uint16s from the buffer at the current offset in little-endian and moves the offset forward the amount of bytes read
func (b *Buffer) ReadU16LENext(n int64) (out []uint16) {

	out = b.ReadU16LE(b.off, n)
	b.SeekByte(n*2, true)
	return

}

// ReadU16BE reads a slice of uint16s from the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *Buffer) ReadU16BE(off, n int64) (out []uint16) {

	if (off + n*2) > b.cap {

		panic(BufferOverreadError)

	}

	if off < 0x00 {

		panic(BufferUnderreadError)

	}

	out = make([]uint16, n)

	b.Lock()
	defer b.Unlock()

	i := int64(0)
	{
	read_loop:
		out[i] = uint16(b.buf[off+(1+(i*2))]) |
			uint16(b.buf[off+(i*2)])<<8

		i++
		if i < n {

			goto read_loop

		}
	}

	return

}

// ReadU16BENext reads a slice of uint16s from the buffer at the current offset in big-endian and moves the offset forward the amount of bytes read
func (b *Buffer) ReadU16BENext(n int64) (out []uint16) {

	out = b.ReadU16BE(b.off, n)
	b.SeekByte(n*2, true)
	return

}

// ReadU32LE reads a slice of uint32s from the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *Buffer) ReadU32LE(off, n int64) (out []uint32) {

	if (off + n*4) > b.cap {

		panic(BufferOverreadError)

	}

	if off < 0x00 {

		panic(BufferUnderreadError)

	}

	out = make([]uint32, n)

	b.Lock()
	defer b.Unlock()

	i := int64(0)
	{
	read_loop:
		out[i] = uint32(b.buf[off+(i*4)]) |
			uint32(b.buf[off+(1+(i*4))])<<8 |
			uint32(b.buf[off+(2+(i*4))])<<16 |
			uint32(b.buf[off+(3+(i*4))])<<24

		i++
		if i < n {

			goto read_loop

		}
	}

	return

}

// ReadU32LENext reads a slice of uint32s from the buffer at the current offset in little-endian and moves the offset forward the amount of bytes read
func (b *Buffer) ReadU32LENext(n int64) (out []uint32) {

	out = b.ReadU32LE(b.off, n)
	b.SeekByte(n*4, true)
	return

}

// ReadU32BE reads a slice of uint32s from the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *Buffer) ReadU32BE(off, n int64) (out []uint32) {

	if (off + n*4) > b.cap {

		panic(BufferOverreadError)

	}

	if off < 0x00 {

		panic(BufferUnderreadError)

	}

	out = make([]uint32, n)

	b.Lock()
	defer b.Unlock()

	i := int64(0)
	{
	read_loop:
		out[i] = uint32(b.buf[off+(3+(i*4))]) |
			uint32(b.buf[off+(2+(i*4))])<<8 |
			uint32(b.buf[off+(1+(i*4))])<<16 |
			uint32(b.buf[off+(i*4)])<<24

		i++
		if i < n {

			goto read_loop

		}
	}

	return

}

// ReadU32BENext reads a slice of uint32s from the buffer at the current offset in big-endian and moves the offset forward the amount of bytes read
func (b *Buffer) ReadU32BENext(n int64) (out []uint32) {

	out = b.ReadU32BE(b.off, n)
	b.SeekByte(n*4, true)
	return

}

// ReadU64LE reads a slice of uint64s from the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *Buffer) ReadU64LE(off, n int64) (out []uint64) {

	if (off + n*8) > b.cap {

		panic(BufferOverreadError)

	}

	if off < 0x00 {

		panic(BufferUnderreadError)

	}

	out = make([]uint64, n)

	b.Lock()
	defer b.Unlock()

	i := int64(0)
	{
	read_loop:
		out[i] = uint64(b.buf[off+(i*8)]) |
			uint64(b.buf[off+(1+(i*8))])<<8 |
			uint64(b.buf[off+(2+(i*8))])<<16 |
			uint64(b.buf[off+(3+(i*8))])<<24 |
			uint64(b.buf[off+(4+(i*8))])<<32 |
			uint64(b.buf[off+(5+(i*8))])<<40 |
			uint64(b.buf[off+(6+(i*8))])<<48 |
			uint64(b.buf[off+(7+(i*8))])<<56

		i++
		if i < n {

			goto read_loop

		}
	}

	return

}

// ReadU64LENext reads a slice of uint64s from the buffer at the current offset in little-endian and moves the offset forward the amount of bytes read
func (b *Buffer) ReadU64LENext(n int64) (out []uint64) {

	out = b.ReadU64LE(b.off, n)
	b.SeekByte(n*8, true)
	return

}

// ReadU64BE reads a slice of uint64s from the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *Buffer) ReadU64BE(off, n int64) (out []uint64) {

	if (off + n*8) > b.cap {

		panic(BufferOverreadError)

	}

	if off < 0x00 {

		panic(BufferUnderreadError)

	}

	out = make([]uint64, n)

	b.Lock()
	defer b.Unlock()
	i := int64(0)
	{
	read_loop:
		out[i] = uint64(b.buf[off+(7+(i*8))]) |
			uint64(b.buf[off+(6+(i*8))])<<8 |
			uint64(b.buf[off+(5+(i*8))])<<16 |
			uint64(b.buf[off+(4+(i*8))])<<24 |
			uint64(b.buf[off+(3+(i*8))])<<32 |
			uint64(b.buf[off+(2+(i*8))])<<40 |
			uint64(b.buf[off+(1+(i*8))])<<48 |
			uint64(b.buf[off+(i*8)])<<56

		i++
		if i < n {

			goto read_loop

		}
	}

	return

}

// ReadU64BENext reads a slice of uint64s from the buffer at the current offset in big-endian and moves the offset forward the amount of bytes read
func (b *Buffer) ReadU64BENext(n int64) (out []uint64) {

	out = b.ReadU64BE(b.off, n)
	b.SeekByte(n*8, true)
	return

}

// SeekByte seeks to position off of the buffer relative to the current position or exact
func (b *Buffer) SeekByte(off int64, relative bool) {

	b.Lock()
	defer b.Unlock()

	if relative {

		b.off += off

	} else {

		b.off = off

	}

}

// AfterByte returns the amount of bytes located after the current position or the specified one
func (b *Buffer) AfterByte(off ...int64) int64 {

	if len(off) == 0 {

		return b.cap - b.off - 1

	}
	return b.cap - off[0] - 1

}

// AlignByte aligns the byte offset to the bit offset
func (b *Buffer) AlignByte() {

	b.Lock()
	defer b.Unlock()

	b.off = b.boff / 8

}

/* generic methods */

// Grow makes the buffer's capacity bigger by n bytes
func (b *Buffer) Grow(n int64) {

	b.Lock()

	b.buf = append(b.buf, make([]byte, n)...)

	b.Unlock()

	b.Refresh()

}

// Refresh updates the cached internal statistics of the buffer forcefully
func (b *Buffer) Refresh() {

	b.Lock()
	defer b.Unlock()

	b.cap = int64(len(b.buf))
	b.bcap = b.cap * 8

	if len(b.buf) > 0 {

		b.obuf = (uintptr)(unsafe.Pointer(&b.buf[0]))

	}

}

// Reset resets the entire buffer
func (b *Buffer) Reset() {

	b.Lock()
	defer b.Unlock()

	b.buf = []byte{}
	b.off = 0x00
	b.boff = 0x00
	b.cap = 0
	b.bcap = 0

}

/* value retrieval */

// Bytes returns the internal byte slice of the buffer
func (b *Buffer) Bytes() []byte {

	return b.buf

}

// ByteCapacity returns the capacity of the buffer
func (b *Buffer) ByteCapacity() int64 {

	return b.cap

}

// BitCapacity returns the bit capacity of the buffer
func (b *Buffer) BitCapacity() int64 {

	return b.bcap

}

// ByteOffset returns the current offset of the buffer
func (b *Buffer) ByteOffset() int64 {

	return b.off

}

// BitOffset returns the current bit offset of the buffer
func (b *Buffer) BitOffset() int64 {

	return b.boff

}
