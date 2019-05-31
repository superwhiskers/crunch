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

// TODO: mirror over modified bitfield methods to Buffer for performance

package crunch

import (
	"sync"
	"unsafe"
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

	// temp?
	obuf uintptr

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

	(*out).Refresh()

}

/* bitfield methods */

// ReadBit stores the bit located at the specified offset without modifying the internal offset value in out
func (b *MiniBuffer) ReadBit(out *byte, off int64) {

	*out = (b.buf[off/8] >> (7 - uint64(off%8))) & 1

}

// ReadBitNext stores the next bit from the current offset and moves the offset forward a bit in out
func (b *MiniBuffer) ReadBitNext(out *byte) {

	b.ReadBit(out, b.boff)
	b.SeekBit(1, true)

}

// ReadBits stores the next n bits from the specified offset without modifying the internal offset value in out
func (b *MiniBuffer) ReadBits(out *uint64, off, n int64) {

	var (
		bout byte

		i = int64(0)
	)
	{
	read_loop:
		b.ReadBit(&bout, off+i)
		*out = (*out << uint64(1)) | uint64(bout)
		i++
		if i < n {

			goto read_loop

		}
	}

}

// ReadBitsNext stores the next n bits from the current offset and moves the offset forward the amount of bits read in out
func (b *MiniBuffer) ReadBitsNext(out *uint64, n int64) {

	b.ReadBits(out, b.boff, n)
	b.SeekBit(n, true)

}

// SetBit sets the bit located at the specified offset without modifying the internal offset value
func (b *MiniBuffer) SetBit(off int64) {

	b.buf[off/8] |= (1 << uint(7-(off%8)))

}

// SetBitNext sets the next bit from the current offset and moves the offset forward a bit
func (b *MiniBuffer) SetBitNext() {

	b.SetBit(b.boff)
	b.SeekBit(1, true)

}

// ClearBit clears the bit located at the specified offset without modifying the internal offset value
func (b *MiniBuffer) ClearBit(off int64) {

	b.buf[off/8] &= ^(1 << uint(7-(off%8)))

}

// ClearBitNext clears the next bit from the current offset and moves the offset forward a bit
func (b *MiniBuffer) ClearBitNext() {

	b.ClearBit(b.boff)
	b.SeekBit(1, true)

}

// SetBits sets the next n bits from the specified offset without modifying the internal offset value
func (b *MiniBuffer) SetBits(off int64, data uint64, n int64) {

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
func (b *MiniBuffer) SetBitsNext(data uint64, n int64) {

	b.SetBits(b.boff, data, n)
	b.SeekBit(n, true)

}

// FlipBit flips the bit located at the specified offset without modifying the internal offset value
func (b *MiniBuffer) FlipBit(off int64) {

	b.buf[off/8] ^= (1 << uint(7-(off%8)))

}

// FlipBitNext flips the next bit from the current offset and moves the offset forward a bit
func (b *MiniBuffer) FlipBitNext() {

	b.FlipBit(b.boff)
	b.SeekBit(1, true)

}

// ClearAllBits sets all of the buffer's bits to 0
func (b *MiniBuffer) ClearAllBits() {

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
func (b *MiniBuffer) SetAllBits() {

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
func (b *MiniBuffer) FlipAllBits() {

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
func (b *MiniBuffer) SeekBit(off int64, relative bool) {

	if relative {

		b.boff += off

	} else {

		b.boff = off

	}

}

// AfterBit stores the amount of bits located after the current bit position or the specified one in out
func (b *MiniBuffer) AfterBit(out *int64, off ...int64) {

	if len(off) == 0 {

		*out = b.bcap - b.boff - 1
		return

	}
	*out = b.bcap - off[0] - 1

}

// AlignBit aligns the bit offset to the byte offset
func (b *MiniBuffer) AlignBit() {

	b.boff = b.off * 8

}

/* byte buffer methods */

// WriteBytes writes bytes to the buffer at the specified offset without modifying the internal offset value
func (b *MiniBuffer) WriteBytes(off int64, data []byte) {

	/* i'm just leaving this here incase this new method proves to be slower in some edge cases */
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

}

// WriteBytesNext writes bytes to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteBytesNext(data []byte) {

	b.WriteBytes(b.off, data)
	b.SeekByte(int64(len(data)), true)

}

// WriteU16LE writes a slice of uint16s to the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *MiniBuffer) WriteU16LE(off int64, data []uint16) {

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

}

// WriteU16LENext writes a slice of uint16s to the buffer at the current offset in little-endian and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteU16LENext(data []uint16) {

	b.WriteU16LE(b.off, data)
	b.SeekByte(int64(len(data))*2, true)

}

// WriteU16BE writes a slice of uint16s to the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *MiniBuffer) WriteU16BE(off int64, data []uint16) {

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

}

// WriteU16BENext writes a slice of uint16s to the buffer at the current offset in big-endian and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteU16BENext(data []uint16) {

	b.WriteU16BE(b.off, data)
	b.SeekByte(int64(len(data))*2, true)

}

// WriteU32LE writes a slice of uint32s to the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *MiniBuffer) WriteU32LE(off int64, data []uint32) {

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

}

// WriteU32LENext writes a slice of uint32s to the buffer at the current offset in little-endian and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteU32LENext(data []uint32) {

	b.WriteU32LE(b.off, data)
	b.SeekByte(int64(len(data))*4, true)

}

// WriteU32BE writes a slice of uint32s to the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *MiniBuffer) WriteU32BE(off int64, data []uint32) {

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

}

// WriteU32BENext writes a slice of uint32s to the buffer at the current offset in big-endian and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteU32BENext(data []uint32) {

	b.WriteU32BE(b.off, data)
	b.SeekByte(int64(len(data))*4, true)

}

// WriteU64LE writes a slice of uint64s to the buffer at the specfied offset in little-endian without modifying the internal offset value
func (b *MiniBuffer) WriteU64LE(off int64, data []uint64) {

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

}

// WriteU64LENext writes a slice of uint64s to the buffer at the current offset in little-endian and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteU64LENext(data []uint64) {

	b.WriteU64LE(b.off, data)
	b.SeekByte(int64(len(data))*8, true)

}

// WriteU64BE writes a slice of uint64s to the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *MiniBuffer) WriteU64BE(off int64, data []uint64) {

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

}

// WriteU64BENext writes a slice of uint64s to the buffer at the current offset in big-endian and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteU64BENext(data []uint64) {

	b.WriteU64BE(b.off, data)
	b.SeekByte(int64(len(data))*8, true)

}

// ReadBytes stores the next n bytes from the specified offset without modifying the internal offset value in out
func (b *MiniBuffer) ReadBytes(out *[]byte, off, n int64) {

	*out = b.buf[off : off+n]

}

// ReadBytesNext stores the next n bytes from the current offset and moves the offset forward the amount of bytes read in out
func (b *MiniBuffer) ReadBytesNext(out *[]byte, n int64) {

	b.ReadBytes(out, b.off, n)
	b.SeekByte(n, true)

}

// ReadU16LE reads a slice of uint16s from the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *MiniBuffer) ReadU16LE(out *[]uint16, off, n int64) {

	i := int64(0)
	{
	read_loop:
		(*out)[i] = uint16(b.buf[off+(i*2)]) |
			uint16(b.buf[off+(1+(i*2))])<<8

		i++
		if i < n {

			goto read_loop

		}
	}

}

// ReadU16LENext reads a slice of uint16s from the buffer at the current offset in little-endian and moves the offset forward the amount of bytes read
func (b *MiniBuffer) ReadU16LENext(out *[]uint16, n int64) {

	b.ReadU16LE(out, b.off, n)
	b.SeekByte(n*2, true)

}

// ReadU16BE reads a slice of uint16s from the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *MiniBuffer) ReadU16BE(out *[]uint16, off, n int64) {

	i := int64(0)
	{
	read_loop:
		(*out)[i] = uint16(b.buf[off+(1+(i*2))]) |
			uint16(b.buf[off+(i*2)])<<8

		i++
		if i < n {

			goto read_loop

		}
	}

}

// ReadU16BENext reads a slice of uint16s from the buffer at the current offset in big-endian and moves the offset forward the amount of bytes read
func (b *MiniBuffer) ReadU16BENext(out *[]uint16, n int64) {

	b.ReadU16BE(out, b.off, n)
	b.SeekByte(n*2, true)

}

// ReadU32LE reads a slice of uint32s from the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *MiniBuffer) ReadU32LE(out *[]uint32, off, n int64) {

	i := int64(0)
	{
	read_loop:
		(*out)[i] = uint32(b.buf[off+(i*4)]) |
			uint32(b.buf[off+(1+(i*4))])<<8 |
			uint32(b.buf[off+(2+(i*4))])<<16 |
			uint32(b.buf[off+(3+(i*4))])<<24

		i++
		if i < n {

			goto read_loop

		}
	}

}

// ReadU32LENext reads a slice of uint32s from the buffer at the current offset in little-endian and moves the offset forward the amount of bytes read
func (b *MiniBuffer) ReadU32LENext(out *[]uint32, n int64) {

	b.ReadU32LE(out, b.off, n)
	b.SeekByte(n*4, true)

}

// ReadU32BE reads a slice of uint32s from the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *MiniBuffer) ReadU32BE(out *[]uint32, off, n int64) {

	i := int64(0)
	{
	read_loop:
		(*out)[i] = uint32(b.buf[off+(3+(i*4))]) |
			uint32(b.buf[off+(2+(i*4))])<<8 |
			uint32(b.buf[off+(1+(i*4))])<<16 |
			uint32(b.buf[off+(i*4)])<<24

		i++
		if i < n {

			goto read_loop

		}
	}

}

// ReadU32BENext reads a slice of uint32s from the buffer at the current offset in big-endian and moves the offset forward the amount of bytes read
func (b *MiniBuffer) ReadU32BENext(out *[]uint32, n int64) {

	b.ReadU32BE(out, b.off, n)
	b.SeekByte(n*4, true)

}

// ReadU64LE reads a slice of uint64s from the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *MiniBuffer) ReadU64LE(out *[]uint64, off, n int64) {

	i := int64(0)
	{
	read_loop:
		(*out)[i] = uint64(b.buf[off+(i*8)]) |
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

}

// ReadU64LENext reads a slice of uint64s from the buffer at the current offset in little-endian and moves the offset forward the amount of bytes read
func (b *MiniBuffer) ReadU64LENext(out *[]uint64, n int64) {

	b.ReadU64LE(out, b.off, n)
	b.SeekByte(n*8, true)

}

// ReadU64BE reads a slice of uint64s from the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *MiniBuffer) ReadU64BE(out *[]uint64, off, n int64) {

	i := int64(0)
	{
	read_loop:
		(*out)[i] = uint64(b.buf[off+(7+(i*8))]) |
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

}

// ReadU64BENext reads a slice of uint64s from the buffer at the current offset in big-endian and moves the offset forward the amount of bytes read
func (b *MiniBuffer) ReadU64BENext(out *[]uint64, n int64) {

	b.ReadU64BE(out, b.off, n)
	b.SeekByte(n*8, true)

}

// SeekByte seeks to position off of the buffer relative to the current position or exact
func (b *MiniBuffer) SeekByte(off int64, relative bool) {

	if relative {

		b.off += off

	} else {

		b.off = off

	}

}

// AfterByte stores the amount of bytes located after the current position or the specified one in out
func (b *MiniBuffer) AfterByte(out *int64, off ...int64) {

	if len(off) == 0 {

		*out = b.cap - b.off - 1
		return

	}
	*out = b.cap - off[0] - 1

}

// AlignByte aligns the byte offset to the bit offset
func (b *MiniBuffer) AlignByte() {

	b.off = b.boff / 8

}

/* generic methods */

// Grow makes the buffer's capacity bigger by n bytes
func (b *MiniBuffer) Grow(n int64) {

	b.buf = append(b.buf, make([]byte, n)...)
	b.Refresh()

}

// Refresh updates the cached internal statistics of the buffer forcefully
func (b *MiniBuffer) Refresh() {

	b.cap = int64(len(b.buf))
	b.bcap = b.cap * 8

	if len(b.buf) > 0 {

		b.obuf = (uintptr)(unsafe.Pointer(&b.buf[0]))

	}

}

// Reset resets the entire buffer
func (b *MiniBuffer) Reset() {

	b.buf = []byte{}
	b.off = 0x00
	b.boff = 0x00
	b.cap = 0
	b.bcap = 0

}

/* value retrieval */

// Bytes stores the internal byte slice of the buffer in out
func (b *MiniBuffer) Bytes(out *[]byte) {

	*out = b.buf

}

// ByteCapacity stores the capacity of the buffer in out
func (b *MiniBuffer) ByteCapacity(out *int64) {

	*out = b.cap

}

// BitCapacity stores the bit capacity of the buffer in out
func (b *MiniBuffer) BitCapacity(out *int64) {

	*out = b.bcap

}

// ByteOffset stores the current offset of the buffer in out
func (b *MiniBuffer) ByteOffset(out *int64) {

	*out = b.off

}

// BitOffset stores the current bit offset of the buffer in out
func (b *MiniBuffer) BitOffset(out *int64) {

	*out = b.boff

}
