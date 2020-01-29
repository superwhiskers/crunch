/*

crunch - utilities for taking bytes out of things
copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>

this source code form is subject to the terms of the mozilla public
license, v. 2.0. if a copy of the mpl was not distributed with this
file, you can obtain one at http://mozilla.org/MPL/2.0/.

*/

package v2

import "sync"

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

	var (
		bout byte

		i = int64(0)
	)
	{
	read_loop:
		b.readbit(&bout, off+i)
		*out = (*out << uint64(1)) | uint64(bout)
		i++
		if i < n {

			goto read_loop

		}
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

	i := int64(0)

	{
	write_loop:
		b.setbit(off+i, byte((data>>uint64(n-i-1))&1))
		i++
		if i < n {

			goto write_loop

		}
	}

}

// flipbit flips a bit in the bitfield
func (b *MiniBuffer) flipbit(off int64) {

	b.buf[off/8] ^= (1 << uint(7-(off%8)))

}

// clearallbits sets all of the buffer's bits to 0
func (b *MiniBuffer) clearallbits() {

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

// setallbits sets all of the buffer's bits to 1
func (b *MiniBuffer) setallbits() {

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

// flipallbits flips all of the buffer's bits
func (b *MiniBuffer) flipallbits() {

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

	var (
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
	}

}

// writeU16LE writes a slice of uint16s to the buffer at the specified offset in little-endian
func (b *MiniBuffer) writeU16LE(off int64, data []uint16) {

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

// writeU16BE writes a slice of uint16s to the buffer at the specified offset in big-endian
func (b *MiniBuffer) writeU16BE(off int64, data []uint16) {

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

// writeU32LE writes a slice of uint32s to the buffer at the specified offset in little-endian
func (b *MiniBuffer) writeU32LE(off int64, data []uint32) {

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

// writeU32BE writes a slice of uint32s to the buffer at the specified offset in big-endian
func (b *MiniBuffer) writeU32BE(off int64, data []uint32) {

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

// writeU64LE writes a slice of uint64s to the buffer at the specified offset in little-endian
func (b *MiniBuffer) writeU64LE(off int64, data []uint64) {

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

// writeU64BE writes a slice of uint64s to the buffer at the specified offset in big-endian
func (b *MiniBuffer) writeU64BE(off int64, data []uint64) {

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

// read reads n bytes from the buffer at the specified offset
func (b *MiniBuffer) read(out *[]byte, off, n int64) {

	*out = b.buf[off : off+n]

}

// readU16LE reads a slice of uint16s from the buffer at the specified offset in little-endian
func (b *MiniBuffer) readU16LE(out *[]uint16, off, n int64) {

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

// readU16BE reads a slice of uint16s from the buffer at the specified offset in big-endian
func (b *MiniBuffer) readU16BE(out *[]uint16, off, n int64) {

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

// readU32LE reads a slice of uint32s from the buffer at the specified offset in little-endian
func (b *MiniBuffer) readU32LE(out *[]uint32, off, n int64) {

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

// readU32BE reads a slice of uint32s from the buffer at the specified offset in big-endian
func (b *MiniBuffer) readU32BE(out *[]uint32, off, n int64) {

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

// readU64LE reads a slice of uint64s from the buffer at the specified offset in little-endian
func (b *MiniBuffer) readU64LE(out *[]uint64, off, n int64) {

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

// readU64BE reads a slice of uint64s from the buffer at the specified offset in big-endian
func (b *MiniBuffer) readU64BE(out *[]uint64, off, n int64) {

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

// reset resets the buffer
func (b *MiniBuffer) reset() {

	b.buf = []byte{}
	b.off = 0x00
	b.boff = 0x00
	b.cap = 0
	b.bcap = 0

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

// Reset resets the entire buffer
func (b *MiniBuffer) Reset() {

	b.reset()

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

// ReadU16LE reads a slice of uint16s from the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *MiniBuffer) ReadU16LE(out *[]uint16, off, n int64) {

	b.readU16LE(out, off, n)

}

// ReadU16BE reads a slice of uint16s from the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *MiniBuffer) ReadU16BE(out *[]uint16, off, n int64) {

	b.readU16BE(out, off, n)

}

// ReadU32LE reads a slice of uint32s from the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *MiniBuffer) ReadU32LE(out *[]uint32, off, n int64) {

	b.readU32LE(out, off, n)

}

// ReadU32BE reads a slice of uint32s from the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *MiniBuffer) ReadU32BE(out *[]uint32, off, n int64) {

	b.readU32BE(out, off, n)

}

// ReadU64LE reads a slice of uint64s from the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *MiniBuffer) ReadU64LE(out *[]uint64, off, n int64) {

	b.readU64LE(out, off, n)

}

// ReadU64BE reads a slice of uint64s from the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *MiniBuffer) ReadU64BE(out *[]uint64, off, n int64) {

	b.readU64BE(out, off, n)

}

// ReadBytesNext stores the next n bytes from the current offset and moves the offset forward the amount of bytes read in out
func (b *MiniBuffer) ReadBytesNext(out *[]byte, n int64) {

	b.read(out, b.off, n)
	b.seek(n, true)

}

// ReadU16LENext reads a slice of uint16s from the buffer at the current offset in little-endian and moves the offset forward the amount of bytes read
func (b *MiniBuffer) ReadU16LENext(out *[]uint16, n int64) {

	b.readU16LE(out, b.off, n)
	b.seek(n*2, true)

}

// ReadU16BENext reads a slice of uint16s from the buffer at the current offset in big-endian and moves the offset forward the amount of bytes read
func (b *MiniBuffer) ReadU16BENext(out *[]uint16, n int64) {

	b.readU16BE(out, b.off, n)
	b.seek(n*2, true)

}

// ReadU32LENext reads a slice of uint32s from the buffer at the current offset in little-endian and moves the offset forward the amount of bytes read
func (b *MiniBuffer) ReadU32LENext(out *[]uint32, n int64) {

	b.readU32LE(out, b.off, n)
	b.seek(n*4, true)

}

// ReadU32BENext reads a slice of uint32s from the buffer at the current offset in big-endian and moves the offset forward the amount of bytes read
func (b *MiniBuffer) ReadU32BENext(out *[]uint32, n int64) {

	b.readU32BE(out, b.off, n)
	b.seek(n*4, true)

}

// ReadU64LENext reads a slice of uint64s from the buffer at the current offset in little-endian and moves the offset forward the amount of bytes read
func (b *MiniBuffer) ReadU64LENext(out *[]uint64, n int64) {

	b.readU64LE(out, b.off, n)
	b.seek(n*8, true)

}

// ReadU64BENext reads a slice of uint64s from the buffer at the current offset in big-endian and moves the offset forward the amount of bytes read
func (b *MiniBuffer) ReadU64BENext(out *[]uint64, n int64) {

	b.readU64BE(out, b.off, n)
	b.seek(n*8, true)

}

// WriteBytes writes bytes to the buffer at the specified offset without modifying the internal offset value
func (b *MiniBuffer) WriteBytes(off int64, data []byte) {

	b.write(off, data)

}

// WriteU16LE writes a slice of uint16s to the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *MiniBuffer) WriteU16LE(off int64, data []uint16) {

	b.writeU16LE(off, data)

}

// WriteU16BE writes a slice of uint16s to the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *MiniBuffer) WriteU16BE(off int64, data []uint16) {

	b.writeU16BE(off, data)

}

// WriteU32LE writes a slice of uint32s to the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *MiniBuffer) WriteU32LE(off int64, data []uint32) {

	b.writeU32LE(off, data)

}

// WriteU32BE writes a slice of uint32s to the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *MiniBuffer) WriteU32BE(off int64, data []uint32) {

	b.writeU32BE(off, data)

}

// WriteU64LE writes a slice of uint64s to the buffer at the specfied offset in little-endian without modifying the internal offset value
func (b *MiniBuffer) WriteU64LE(off int64, data []uint64) {

	b.writeU64LE(off, data)

}

// WriteU64BE writes a slice of uint64s to the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *MiniBuffer) WriteU64BE(off int64, data []uint64) {

	b.writeU64BE(off, data)

}

// WriteBytesNext writes bytes to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteBytesNext(data []byte) {

	b.write(b.off, data)
	b.seek(int64(len(data)), true)

}

// WriteU16LENext writes a slice of uint16s to the buffer at the current offset in little-endian and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteU16LENext(data []uint16) {

	b.writeU16LE(b.off, data)
	b.seek(int64(len(data))*2, true)

}

// WriteU16BENext writes a slice of uint16s to the buffer at the current offset in big-endian and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteU16BENext(data []uint16) {

	b.writeU16BE(b.off, data)
	b.seek(int64(len(data))*2, true)

}

// WriteU32LENext writes a slice of uint32s to the buffer at the current offset in little-endian and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteU32LENext(data []uint32) {

	b.writeU32LE(b.off, data)
	b.seek(int64(len(data))*4, true)

}

// WriteU32BENext writes a slice of uint32s to the buffer at the current offset in big-endian and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteU32BENext(data []uint32) {

	b.writeU32BE(b.off, data)
	b.seek(int64(len(data))*4, true)

}

// WriteU64LENext writes a slice of uint64s to the buffer at the current offset in little-endian and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteU64LENext(data []uint64) {

	b.writeU64LE(b.off, data)
	b.seek(int64(len(data))*8, true)

}

// WriteU64BENext writes a slice of uint64s to the buffer at the current offset in big-endian and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteU64BENext(data []uint64) {

	b.writeU64BE(b.off, data)
	b.seek(int64(len(data))*8, true)

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
