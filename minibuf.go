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

import "sync"

// MiniBuffer implements a concurrent-safe buffer type in go that handles multiple types of data with very low memory allocations
// complex reading and writing functions are removed as they are harder to optimize and aren't hard to do without
type MiniBuffer struct {
	buf  []byte
	off  int64
	cap  int64
	boff int64
	bcap int64

	sync.Mutex
}

// NewMiniBuffer initilaizes a new MiniBuffer with the provided byte slice(s) stored inside in the order provided
func NewMiniBuffer(out *MiniBuffer, slices ...[]byte) {

	out = &MiniBuffer{
		buf:  []byte{},
		off:  0x00,
		boff: 0x00,
	}

	switch len(slices) {

	case 0:
		break

	case 1:
		out.buf = slices[0]
		break

	default:
		for _, s := range slices {

			out.buf = append(out.buf, s...)

		}

	}

	out.refresh()
	return

}

/* internal use methods */

/* bitfield methods */

// readbit reads a bit from the bitfield at the specified offset
func (b *MiniBuffer) readbit(out *byte, off int64) {

	if off > (b.bcap - 1) {

		panic(BitfieldOverreadError)

	}

	b.Lock()
	defer b.Unlock()

	*out = atob((b.buf[off/8] & (1 << uint(7-(off%8)))) != 0)
	return

}

// readbits reads n bits from the bitfield at the specified offset
func (b *MiniBuffer) readbits(out *uint64, off, n int64) {

	if (off + n) > b.bcap {

		panic(BitfieldOverreadError)

	}

	var bout byte
	for i := int64(0); i < n; i++ {

		b.readbit(&bout, off+i)
		*out = (*out << uint64(1)) | uint64(bout)

	}

	return

}

// setbit sets a bit in the bitfield to the specified value
func (b *MiniBuffer) setbit(off int64, data byte) {

	if off > (b.bcap - 1) {

		panic(BitfieldOverwriteError)

	}

	if data != 1 && data != 0 {

		panic(BitfieldInvalidBitError)

	}

	b.Lock()
	defer b.Unlock()

	switch data {

	case 0:
		b.buf[off/8] &= ^(1 << uint(7-(off%8)))

	case 1:
		b.buf[off/8] |= (1 << uint(7-(off%8)))

	}

	return

}

// setbits sets n bits in the bitfield to the specified value at the specified offset
func (b *MiniBuffer) setbits(off int64, data uint64, n int64) {

	if off+n > (b.bcap - 1) {

		panic(BitfieldOverwriteError)

	}

	for i := int64(0); i < n; i++ {

		b.setbit(off+i, byte((data>>uint64(n-i-1))&1))

	}

	return

}

// flipbit flips a bit in the bitfield
func (b *MiniBuffer) flipbit(off int64) {

	if off > (b.bcap - 1) {

		panic(BitfieldOverwriteError)

	}

	b.Lock()
	defer b.Unlock()

	b.buf[off/8] ^= (1 << uint(7-(off%8)))
	return

}

// clearallbits sets all of the buffer's bits to 0
func (b *MiniBuffer) clearallbits() {

	b.Lock()
	defer b.Unlock()

	for i := range b.buf {

		b.buf[i] = 0

	}

	return

}

// setallbits sets all of the buffer's bits to 1
func (b *MiniBuffer) setallbits() {

	b.Lock()
	defer b.Unlock()

	for i := range b.buf {

		b.buf[i] = 0xFF

	}

	return

}

// flipallbits flips all of the buffer's bits
func (b *MiniBuffer) flipallbits() {

	b.Lock()
	defer b.Unlock()

	for i := range b.buf {

		b.buf[i] = ^b.buf[i]

	}

	return

}

// seekbit seeks to position off of the bitfield relative to the current position or exact
func (b *MiniBuffer) seekbit(off int64, relative bool) {

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
func (b *MiniBuffer) afterbit(out *int64, off ...int64) {

	if len(off) == 0 {

		*out = b.bcap - b.boff

	}
	*out = b.bcap - off[0]
	return

}

/* byte buffer methods */

// write writes a slice of bytes to the buffer at the specified offset
func (b *MiniBuffer) write(off int64, data []byte) {

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

// read reads n bytes from the buffer at the specified offset
func (b *MiniBuffer) read(out *[]byte, off, n int64) {

	if (off + n) > b.cap {

		panic(ByteBufferOverreadError)

	}

	b.Lock()
	defer b.Unlock()

	*out = b.buf[off : off+n]
	return

}

// seek seeks to position off of the byte buffer relative to the current position or exact
func (b *MiniBuffer) seek(off int64, relative bool) {

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
func (b *MiniBuffer) after(out *int64, off ...int64) {

	if len(off) == 0 {

		*out = b.cap - b.off

	}
	*out = b.cap - off[0]
	return

}

/* generic methods */

// grow grows the buffer by n bytes
func (b *MiniBuffer) grow(n int64) {

	b.Lock()

	b.buf = append(b.buf, make([]byte, n)...)

	b.Unlock()

	b.refresh()

	return

}

// refresh updates the internal statistics of the byte buffer forcefully
func (b *MiniBuffer) refresh() {

	b.Lock()
	defer b.Unlock()

	b.cap = int64(len(b.buf))
	b.bcap = b.cap * 8

	return

}

// alignbit aligns the bit offset to the byte offset
func (b *MiniBuffer) alignbit() {

	b.Lock()
	defer b.Unlock()

	b.boff = b.off * 8

}

// alignbyte aligns the byte offset to the bit offset
func (b *MiniBuffer) alignbyte() {

	b.Lock()
	defer b.Unlock()

	b.off = b.boff / 8

}

/* public methods */

// Bytes stores the internal byte slice of the buffer in out
func (b *MiniBuffer) Bytes(out *[]byte) {

	*out = b.buf
	return

}

// Capacity stores the capacity of the buffer in out
func (b *MiniBuffer) Capacity(out *int64) {

	*out = b.cap
	return

}

// BitCapacity stores the bit capacity of the buffer in out
func (b *MiniBuffer) BitCapacity(out *int64) {

	*out = b.bcap
	return

}

// Offset stores the current offset of the buffer in out
func (b *MiniBuffer) Offset(out *int64) {

	*out = b.off
	return

}

// BitOffset stores the current bit offset of the buffer in out
func (b *MiniBuffer) BitOffset(out *int64) {

	*out = b.boff
	return

}

// Refresh updates the cached internal statistics of the buffer forcefully
func (b *MiniBuffer) Refresh() {

	b.refresh()
	return

}

// Grow makes the buffer's capacity bigger by n bytes
func (b *MiniBuffer) Grow(n int64) {

	b.grow(n)
	return

}

// Seek seeks to position off of the buffer relative to the current position or exact
func (b *MiniBuffer) Seek(off int64, relative bool) {

	b.seek(off, relative)
	return

}

// SeekBit seeks to bit position off of the the buffer relative to the current position or exact
func (b *MiniBuffer) SeekBit(off int64, relative bool) {

	b.seekbit(off, relative)
	return

}

// AlignBit aligns the bit offset to the byte offset
func (b *MiniBuffer) AlignBit() {

	b.alignbit()
	return

}

// AlignByte aligns the byte offset to the bit offset
func (b *MiniBuffer) AlignByte() {

	b.alignbyte()
	return

}

// After stores the amount of bytes located after the current position or the specified one in out
func (b *MiniBuffer) After(out *int64, off ...int64) {

	b.after(out, off...)
	return

}

// AfterBit stores the amount of bits located after the current bit position or the specified one in out
func (b *MiniBuffer) AfterBit(out *int64, off ...int64) {

	b.afterbit(out, off...)
	return

}

// ReadBytes stores the next n bytes from the specified offset without modifying the internal offset value in out
func (b *MiniBuffer) ReadBytes(out *[]byte, off, n int64) {

	b.read(out, off, n)
	return

}

// ReadBytesNext stores the next n bytes from the current offset and moves the offset forward the amount of bytes read in out
func (b *MiniBuffer) ReadBytesNext(out *[]byte, n int64) {

	b.read(out, b.off, n)
	b.seek(n, true)
	return

}

// WriteBytes writes bytes to the buffer at the specified offset without modifying the internal offset value
func (b *MiniBuffer) WriteBytes(off int64, data []byte) {

	b.write(off, data)
	return

}

// WriteBytesNext writes bytes to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *MiniBuffer) WriteBytesNext(data []byte) {

	b.write(b.off, data)
	b.seek(int64(len(data)), true)
	return

}

// ReadBit stores the bit located at the specified offset without modifying the internal offset value in out
func (b *MiniBuffer) ReadBit(out *byte, off int64) {

	b.readbit(out, off)
	return

}

// ReadBits stores the next n bits from the specified offset without modifying the internal offset value in out
func (b *MiniBuffer) ReadBits(out *uint64, off, n int64) {

	b.readbits(out, off, n)
	return

}

// ReadBitNext stores the next bit from the current offset and moves the offset forward a bit in out
func (b *MiniBuffer) ReadBitNext(out *byte) {

	b.readbit(out, b.boff)
	b.seekbit(1, true)
	return

}

// ReadBitsNext stores the next n bits from the current offset and moves the offset forward the amount of bits read in out
func (b *MiniBuffer) ReadBitsNext(out *uint64, n int64) {

	b.readbits(out, b.boff, n)
	b.seekbit(n, true)
	return

}

// SetBit sets the bit located at the specified offset without modifying the internal offset value
func (b *MiniBuffer) SetBit(off int64, data byte) {

	b.setbit(off, data)
	return

}

// SetBits sets the next n bits from the specified offset without modifying the internal offset value
func (b *MiniBuffer) SetBits(off int64, data uint64, n int64) {

	b.setbits(off, data, n)
	return

}

// SetBitNext sets the next bit from the current offset and moves the offset forward a bit
func (b *MiniBuffer) SetBitNext(data byte) {

	b.setbit(b.boff, data)
	b.seekbit(1, true)
	return

}

// SetBitsNext sets the next n bits from the current offset and moves the offset forward the amount of bits set
func (b *MiniBuffer) SetBitsNext(data uint64, n int64) {

	b.setbits(b.boff, data, n)
	b.seekbit(n, true)
	return

}

// FlipBit flips the bit located at the specified offset without modifying the internal offset value
func (b *MiniBuffer) FlipBit(off int64) {

	b.flipbit(off)
	return

}

// FlipBitNext flips the next bit from the current offset and moves the offset forward a bit
func (b *MiniBuffer) FlipBitNext() {

	b.flipbit(b.boff)
	b.seekbit(1, true)
	return

}

// ClearAllBits sets all of the buffer's bits to 0
func (b *MiniBuffer) ClearAllBits() {

	b.clearallbits()
	return

}

// SetAllBits sets all of the buffer's bits to 1
func (b *MiniBuffer) SetAllBits() {

	b.setallbits()
	return

}

// FlipAllBits flips all of the buffer's bits
func (b *MiniBuffer) FlipAllBits() {

	b.flipallbits()
	return

}
