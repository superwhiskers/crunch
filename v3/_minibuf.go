/*

crunch - utilities for taking bytes out of things
copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>

this source code form is subject to the terms of the mozilla public
license, v. 2.0. if a copy of the mpl was not distributed with this
file, you can obtain one at http://mozilla.org/MPL/2.0/.

*/

package v3

import "unsafe"

// MiniBuffer implements a fast and low-memory buffer type in go that
// handles multiple types of data easily. it lacks the overwrite/read
// and underwrite/read checks that Buffer has
type MiniBuffer struct {
	buf  []byte
	off  int64
	cap  int64
	boff int64
	bcap int64

	// temp?
	obuf unsafe.Pointer
}

// NewMiniBuffer initilaizes a new MiniBuffer with the provided byte
// slice(s) stored inside in the order provided
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
		var (
			i = int64(0)
			n = int64(len(slices))
		)
		{
		append_loop:
			(*out).buf = append((*out).buf, slices[i]...)
			i++
			if i < n {

				goto append_loop

			}
		}

	}

	(*out).Refresh()

}

/* bitfield methods */

// ReadBit stores the bit located at the specified offset without
// modifying the internal offset value in out
func (b *MiniBuffer) ReadBit(out *byte, off int64) {

	*out = (b.buf[off/8] >> (7 - uint64(off%8))) & 1

}

// ReadBitNext stores the next bit from the current offset and moves
// the offset forward a bit in out
func (b *MiniBuffer) ReadBitNext(out *byte) {

	b.ReadBit(out, b.boff)
	b.SeekBit(1, true)

}

// ReadBits stores the next n bits from the specified offset without
// modifying the internal offset value in out
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

// ReadBitsNext stores the next n bits from the current offset and
// moves the offset forward the amount of bits read in out
func (b *MiniBuffer) ReadBitsNext(out *uint64, n int64) {

	b.ReadBits(out, b.boff, n)
	b.SeekBit(n, true)

}

// SetBit sets the bit located at the specified offset without
// modifying the internal offset value
func (b *MiniBuffer) SetBit(off int64) {

	b.buf[off/8] |= (1 << uint(7-(off%8)))

}

// SetBitNext sets the next bit from the current offset and moves the
// offset forward a bit
func (b *MiniBuffer) SetBitNext() {

	b.SetBit(b.boff)
	b.SeekBit(1, true)

}

// ClearBit clears the bit located at the specified offset without
// modifying the internal offset value
func (b *MiniBuffer) ClearBit(off int64) {

	b.buf[off/8] &= ^(1 << uint(7-(off%8)))

}

// ClearBitNext clears the next bit from the current offset and moves
// the offset forward a bit
func (b *MiniBuffer) ClearBitNext() {

	b.ClearBit(b.boff)
	b.SeekBit(1, true)

}

// SetBits sets the next n bits from the specified offset without
// modifying the internal offset value
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

// SetBitsNext sets the next n bits from the current offset and moves
// the offset forward the amount of bits set
func (b *MiniBuffer) SetBitsNext(data uint64, n int64) {

	b.SetBits(b.boff, data, n)
	b.SeekBit(n, true)

}

// FlipBit flips the bit located at the specified offset without
// modifying the internal offset value
func (b *MiniBuffer) FlipBit(off int64) {

	b.buf[off/8] ^= (1 << uint(7-(off%8)))

}

// FlipBitNext flips the next bit from the current offset and moves
// the offset forward a bit
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

// SeekBit seeks to bit position off of the the buffer relative to
// the current position or exact
func (b *MiniBuffer) SeekBit(off int64, relative bool) {

	if relative {

		b.boff += off

	} else {

		b.boff = off

	}

}

// AfterBit stores the amount of bits located after the current bit
// position or the specified one in out
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

// WriteBytes writes bytes to the buffer at the specified offset
// without modifying the internal offset value
func (b *MiniBuffer) WriteBytes(off int64, data []byte) {

	/*
	   i'm just leaving this here incase this new
	   method proves to be slower in some edge cases
	*/
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
		p = unsafe.Pointer(uintptr(b.obuf) + uintptr(off))
		i = int64(0)
		n = int64(len(data))
	)
	{
	write_loop:
		*(*byte)(unsafe.Pointer(uintptr(p) + uintptr(i))) = data[i]
		i++
		if i < n {

			goto write_loop

		}
	}

}

// WriteBytesNext writes bytes to the buffer at the current offset
// and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteBytesNext(data []byte) {

	b.WriteBytes(b.off, data)
	b.SeekByte(int64(len(data)), true)

}

//generator:complex MiniBuffer Write U 16 LE

//generator:complex MiniBuffer Write U 16 BE

//generator:complex MiniBuffer Write U 32 LE

//generator:complex MiniBuffer Write U 32 BE

//generator:complex MiniBuffer Write U 64 LE

//generator:complex MiniBuffer Write U 64 BE


//generator:complex MiniBuffer Write I 16 LE

//generator:complex MiniBuffer Write I 16 BE

//generator:complex MiniBuffer Write I 32 LE

//generator:complex MiniBuffer Write I 32 BE

//generator:complex MiniBuffer Write I 64 LE

//generator:complex MiniBuffer Write I 64 BE

//TODO(superwhiskers): add tests for the following generation directives

//generator:complex MiniBuffer Write F 32 LE

//generator:complex MiniBuffer Write F 32 BE

//generator:complex MiniBuffer Write F 64 LE

//generator:complex MiniBuffer Write F 64 BE

// ReadBytes stores the next n bytes from the specified offset
// without modifying the internal offset value in out
func (b *MiniBuffer) ReadBytes(out *[]byte, off, n int64) {

	*out = b.buf[off : off+n]

}

// ReadBytesNext stores the next n bytes from the current offset and
// moves the offset forward the amount of bytes read in out
func (b *MiniBuffer) ReadBytesNext(out *[]byte, n int64) {

	b.ReadBytes(out, b.off, n)
	b.SeekByte(n, true)

}

//generator:complex MiniBuffer Read U 16 LE

//generator:complex MiniBuffer Read U 16 BE

//generator:complex MiniBuffer Read U 32 LE

//generator:complex MiniBuffer Read U 32 BE

//generator:complex MiniBuffer Read U 64 LE

//generator:complex MiniBuffer Read U 64 BE


//generator:complex MiniBuffer Read I 16 LE

//generator:complex MiniBuffer Read I 16 BE

//generator:complex MiniBuffer Read I 32 LE

//generator:complex MiniBuffer Read I 32 BE

//generator:complex MiniBuffer Read I 64 LE

//generator:complex MiniBuffer Read I 64 BE

//TODO(superwhiskers): add tests for the following generation directives

//generator:complex MiniBuffer Read F 32 LE

//generator:complex MiniBuffer Read F 32 BE

//generator:complex MiniBuffer Read F 64 LE

//generator:complex MiniBuffer Read F 64 BE

// SeekByte seeks to position off of the buffer relative to the
// current position or exact
func (b *MiniBuffer) SeekByte(off int64, relative bool) {

	if relative {

		b.off += off

	} else {

		b.off = off

	}

}

// AfterByte stores the amount of bytes located after the current
// position or the specified one in out
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

// TruncateLeft truncates the buffer on the left side
func (b *MiniBuffer) TruncateLeft(n int64) {

	b.buf = b.buf[n:b.cap]
	b.Refresh()

}

// TruncateRight truncates the buffer on the right side
func (b *MiniBuffer) TruncateRight(n int64) {

	b.buf = b.buf[0x00 : b.cap-n]
	b.Refresh()

}

// Grow makes the buffer's capacity bigger by n bytes
func (b *MiniBuffer) Grow(n int64) {

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
func (b *MiniBuffer) Refresh() {

	b.cap = int64(len(b.buf))
	b.bcap = b.cap * 8

	if len(b.buf) > 0 {

		b.obuf = unsafe.Pointer(&b.buf[0])

	}

}

// Reset resets the entire buffer
func (b *MiniBuffer) Reset() {

	b.buf = b.buf[0:0]
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
