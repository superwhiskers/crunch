/*

crunch - utilities for taking bytes out of things
copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>

this source code form is subject to the terms of the mozilla public
license, v. 2.0. if a copy of the mpl was not distributed with this
file, you can obtain one at http://mozilla.org/MPL/2.0/.

*/

package v3

import "unsafe"

// Buffer implements a buffer type in go that handles multiple types
// of data easily. it has overwrite/read checks for extra safety
type Buffer struct {
	buf  []byte
	off  int64
	cap  int64
	boff int64
	bcap int64
}

// NewBuffer initilaizes a new Buffer with the provided byte slice(s)
// stored inside in the order provided
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
		var (
			i = int64(0)
			n = int64(len(slices))
		)
		{
		append_loop:
			buf.buf = append(buf.buf, slices[i]...)
			i++
			if i < n {

				goto append_loop

			}
		}

	}

	buf.Refresh()
	return

}

/* internal use methods */

/* bitfield methods */

// ReadBit returns the bit located at the specified offset without
// modifying the internal offset value
func (b *Buffer) ReadBit(off int64) byte {

	if off > (b.bcap - 1) {

		panic(BufferOverreadError)

	}

	if off < 0x00 {

		panic(BufferUnderreadError)

	}

	return (b.buf[off/8] >> (7 - uint64(off%8))) & 1

}

// ReadBitNext returns the next bit from the current offset and moves
// the offset forward a bit
func (b *Buffer) ReadBitNext() (out byte) {

	out = b.ReadBit(b.boff)
	b.SeekBit(1, true)
	return

}

// ReadBits returns the next n bits from the specified offset without
// modifying the internal offset value
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

// ReadBitsNext returns the next n bits from the current offset and
// moves the offset forward the amount of bits read
func (b *Buffer) ReadBitsNext(n int64) (out uint64) {

	out = b.ReadBits(b.boff, n)
	b.SeekBit(n, true)
	return

}

// SetBit sets the bit located at the specified offset without
// modifying the internal offset value
func (b *Buffer) SetBit(off int64) {

	if off > (b.bcap - 1) {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	b.buf[off/8] |= (1 << uint(7-(off%8)))

}

// SetBitNext sets the next bit from the current offset and moves the
// offset forward a bit
func (b *Buffer) SetBitNext() {

	b.SetBit(b.boff)
	b.SeekBit(1, true)

}

// ClearBit clears the bit located at the specified offset without
// modifying the internal offset value
func (b *Buffer) ClearBit(off int64) {

	if off > (b.bcap - 1) {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	b.buf[off/8] &= ^(1 << uint(7-(off%8)))

}

// ClearBitNext clears the next bit from the current offset and moves
// the offset forward a bit
func (b *Buffer) ClearBitNext() {

	b.ClearBit(b.boff)
	b.SeekBit(1, true)

}

// SetBits sets the next n bits from the specified offset without
// modifying the internal offset value
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

// SetBitsNext sets the next n bits from the current offset and moves
// the offset forward the amount of bits set
func (b *Buffer) SetBitsNext(data uint64, n int64) {

	b.SetBits(b.boff, data, n)
	b.SeekBit(n, true)

}

// FlipBit flips the bit located at the specified offset without
// modifying the internal offset value
func (b *Buffer) FlipBit(off int64) {

	if off > (b.bcap - 1) {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	b.buf[off/8] ^= (1 << uint(7-(off%8)))

}

// FlipBitNext flips the next bit from the current offset and moves
// the offset forward a bit
func (b *Buffer) FlipBitNext() {

	b.FlipBit(b.boff)
	b.SeekBit(1, true)

}

// ClearAllBits sets all of the buffer's bits to 0
func (b *Buffer) ClearAllBits() {

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

// SeekBit seeks to bit position off of the the buffer relative to
// the current position or exact
func (b *Buffer) SeekBit(off int64, relative bool) {

	if relative {

		b.boff += off

	} else {

		b.boff = off

	}

}

// AfterBit returns the amount of bits located after the current bit
// position or the specified one
func (b *Buffer) AfterBit(off ...int64) int64 {

	if len(off) == 0 {

		return b.bcap - b.boff - 1

	}
	return b.bcap - off[0] - 1

}

// AlignBit aligns the bit offset to the byte offset
func (b *Buffer) AlignBit() {

	b.boff = b.off * 8

}

/* byte buffer methods */

// WriteBytes writes bytes to the buffer at the specified offset
// without modifying the internal offset value
func (b *Buffer) WriteBytes(off int64, data []byte) {

	if (off + int64(len(data))) > b.cap {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	copy(b.buf[off:], data)

}

// WriteBytesNext writes bytes to the buffer at the current offset
// and moves the offset forward the amount of bytes written
func (b *Buffer) WriteBytesNext(data []byte) {

	b.WriteBytes(b.off, data)
	b.SeekByte(int64(len(data)), true)

}

// WriteByte writes a byte to the buffer at the specified offset
// without modifying the internal offset value
func (b *Buffer) WriteByte(off int64, data byte) {

	b.WriteBytes(off, []byte{data})

}

// WriteByteNext writes a byte to the buffer at the current
// offset and moves the offset forward the amount of bytes written
func (b *Buffer) WriteByteNext(data byte) {

	b.WriteBytes(b.off, []byte{data})
	b.SeekByte(1, true)

}

//generator:complex Buffer Write U 16 LE

//generator:complex Buffer Write U 16 BE

//generator:complex Buffer Write U 32 LE

//generator:complex Buffer Write U 32 BE

//generator:complex Buffer Write U 64 LE

//generator:complex Buffer Write U 64 BE

//generator:complex Buffer Write I 16 LE

//generator:complex Buffer Write I 16 BE

//generator:complex Buffer Write I 32 LE

//generator:complex Buffer Write I 32 BE

//generator:complex Buffer Write I 64 LE

//generator:complex Buffer Write I 64 BE

//generator:complex Buffer Write F 32 LE

//generator:complex Buffer Write F 32 BE

//generator:complex Buffer Write F 64 LE

//generator:complex Buffer Write F 64 BE

// ReadBytes returns the next n bytes from the specified offset
// without modifying the internal offset value
func (b *Buffer) ReadBytes(off, n int64) []byte {

	if (off + n) > b.cap {

		panic(BufferOverreadError)

	}

	if off < 0x00 {

		panic(BufferUnderreadError)

	}

	return b.buf[off : off+n]

}

// ReadBytesNext returns the next n bytes from the current offset
// and moves the offset forward the amount of bytes read
func (b *Buffer) ReadBytesNext(n int64) (out []byte) {

	out = b.ReadBytes(b.off, n)
	b.SeekByte(n, true)
	return

}

// ReadByte returns the next byte from the specified offset without
// modifying the internal offset value
func (b *Buffer) ReadByte(off int64) byte {

	return b.ReadBytes(off, 1)[0]

}

// ReadByteNext returns the next byte from the current offset and
// moves the offset forward a byte
func (b *Buffer) ReadByteNext() (out byte) {

	out = b.ReadBytes(b.off, 1)[0]
	b.SeekByte(1, true)
	return

}

//generator:complex Buffer Read U 16 LE

//generator:complex Buffer Read U 16 BE

//generator:complex Buffer Read U 32 LE

//generator:complex Buffer Read U 32 BE

//generator:complex Buffer Read U 64 LE

//generator:complex Buffer Read U 64 BE

//generator:complex Buffer Read I 16 LE

//generator:complex Buffer Read I 16 BE

//generator:complex Buffer Read I 32 LE

//generator:complex Buffer Read I 32 BE

//generator:complex Buffer Read I 64 LE

//generator:complex Buffer Read I 64 BE

//generator:complex Buffer Read F 32 LE

//generator:complex Buffer Read F 32 BE

//generator:complex Buffer Read F 64 LE

//generator:complex Buffer Read F 64 BE

// SeekByte seeks to position off of the buffer relative to the
// current position or exact
func (b *Buffer) SeekByte(off int64, relative bool) {

	if relative {

		b.off += off

	} else {

		b.off = off

	}

}

// AfterByte returns the amount of bytes located after the current
// position or the specified one
func (b *Buffer) AfterByte(off ...int64) int64 {

	if len(off) == 0 {

		return b.cap - b.off - 1

	}
	return b.cap - off[0] - 1

}

// AlignByte aligns the byte offset to the bit offset
func (b *Buffer) AlignByte() {

	b.off = b.boff / 8

}

/* generic methods */

// TruncateLeft truncates the buffer on the left side
func (b *Buffer) TruncateLeft(n int64) {

	if n < 0 {

		panic(BufferInvalidByteCountError)

	}

	b.buf = b.buf[n:b.cap]
	b.Refresh()

}

// TruncateRight truncates the buffer on the right side
func (b *Buffer) TruncateRight(n int64) {

	if n < 0 {

		panic(BufferInvalidByteCountError)

	}

	b.buf = b.buf[0x00 : b.cap-n]
	b.Refresh()

}

// Grow makes the buffer's capacity bigger by n bytes
func (b *Buffer) Grow(n int64) {

	if n < 0 {

		panic(BufferInvalidByteCountError)

	}

	if n <= int64(cap(b.buf))-b.cap {

		b.buf = b.buf[0 : b.cap+n]
		b.Refresh()
		return

	}
	tmp := make([]byte, b.cap+n, (int64(cap(b.buf))+n)*2)
	copy(tmp, b.buf)
	b.buf = tmp
	b.Refresh()

}

// Refresh updates the cached internal statistics of the buffer forcefully
func (b *Buffer) Refresh() {

	b.cap = int64(len(b.buf))
	b.bcap = b.cap * 8

}

// Reset resets the entire buffer
func (b *Buffer) Reset() {

	b.buf = b.buf[0:0]
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
