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

	// temp?
	obuf unsafe.Pointer
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

	/*
	   same as in minibuffer, leaving here in case this new
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

// WriteU16LE writes a slice of uint16s to the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) WriteU16LE(off int64, data []uint16) {
	if (off + int64(len(data))*2) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
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

// WriteU16LENextwrites a slice of uint16s to the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteU16LENext(data []uint16) {
	b.WriteU16LE(b.off, data)
	b.SeekByte(int64(len(data))*2, true)
}

// WriteU16BE writes a slice of uint16s to the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) WriteU16BE(off int64, data []uint16) {
	if (off + int64(len(data))*2) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
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

// WriteU16BENextwrites a slice of uint16s to the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteU16BENext(data []uint16) {
	b.WriteU16BE(b.off, data)
	b.SeekByte(int64(len(data))*2, true)
}

// WriteU32LE writes a slice of uint32s to the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) WriteU32LE(off int64, data []uint32) {
	if (off + int64(len(data))*4) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
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

// WriteU32LENextwrites a slice of uint32s to the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteU32LENext(data []uint32) {
	b.WriteU32LE(b.off, data)
	b.SeekByte(int64(len(data))*4, true)
}

// WriteU32BE writes a slice of uint32s to the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) WriteU32BE(off int64, data []uint32) {
	if (off + int64(len(data))*4) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
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

// WriteU32BENextwrites a slice of uint32s to the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteU32BENext(data []uint32) {
	b.WriteU32BE(b.off, data)
	b.SeekByte(int64(len(data))*4, true)
}

// WriteU64LE writes a slice of uint64s to the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) WriteU64LE(off int64, data []uint64) {
	if (off + int64(len(data))*8) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
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

// WriteU64LENextwrites a slice of uint64s to the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteU64LENext(data []uint64) {
	b.WriteU64LE(b.off, data)
	b.SeekByte(int64(len(data))*8, true)
}

// WriteU64BE writes a slice of uint64s to the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) WriteU64BE(off int64, data []uint64) {
	if (off + int64(len(data))*8) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
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

// WriteU64BENextwrites a slice of uint64s to the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteU64BENext(data []uint64) {
	b.WriteU64BE(b.off, data)
	b.SeekByte(int64(len(data))*8, true)
}


// WriteI16LE writes a slice of int16s to the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) WriteI16LE(off int64, data []int16) {
	if (off + int64(len(data))*2) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
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

// WriteI16LENextwrites a slice of int16s to the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteI16LENext(data []int16) {
	b.WriteI16LE(b.off, data)
	b.SeekByte(int64(len(data))*2, true)
}

// WriteI16BE writes a slice of int16s to the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) WriteI16BE(off int64, data []int16) {
	if (off + int64(len(data))*2) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
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

// WriteI16BENextwrites a slice of int16s to the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteI16BENext(data []int16) {
	b.WriteI16BE(b.off, data)
	b.SeekByte(int64(len(data))*2, true)
}

// WriteI32LE writes a slice of int32s to the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) WriteI32LE(off int64, data []int32) {
	if (off + int64(len(data))*4) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
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

// WriteI32LENextwrites a slice of int32s to the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteI32LENext(data []int32) {
	b.WriteI32LE(b.off, data)
	b.SeekByte(int64(len(data))*4, true)
}

// WriteI32BE writes a slice of int32s to the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) WriteI32BE(off int64, data []int32) {
	if (off + int64(len(data))*4) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
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

// WriteI32BENextwrites a slice of int32s to the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteI32BENext(data []int32) {
	b.WriteI32BE(b.off, data)
	b.SeekByte(int64(len(data))*4, true)
}

// WriteI64LE writes a slice of int64s to the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) WriteI64LE(off int64, data []int64) {
	if (off + int64(len(data))*8) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
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

// WriteI64LENextwrites a slice of int64s to the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteI64LENext(data []int64) {
	b.WriteI64LE(b.off, data)
	b.SeekByte(int64(len(data))*8, true)
}

// WriteI64BE writes a slice of int64s to the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) WriteI64BE(off int64, data []int64) {
	if (off + int64(len(data))*8) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
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

// WriteI64BENextwrites a slice of int64s to the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteI64BENext(data []int64) {
	b.WriteI64BE(b.off, data)
	b.SeekByte(int64(len(data))*8, true)
}

//TODO(superwhiskers): add tests for the following generation directives

// WriteF32LE writes a slice of float32s to the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) WriteF32LE(off int64, data []float32) {
	if (off + int64(len(data))*4) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
	{
	write_loop:
		b.buf[off+int64(i*4)] = byte(*(*uint32)(unsafe.Pointer(&data[i])))
		b.buf[off+int64(1+(i*4))] = byte(*(*uint32)(unsafe.Pointer(&data[i])) >> 8)
		b.buf[off+int64(2+(i*4))] = byte(*(*uint32)(unsafe.Pointer(&data[i])) >> 16)
		b.buf[off+int64(3+(i*4))] = byte(*(*uint32)(unsafe.Pointer(&data[i])) >> 24)
		i++
		if i < n {
			goto write_loop
		}
	}
}

// WriteF32LENextwrites a slice of float32s to the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteF32LENext(data []float32) {
	b.WriteF32LE(b.off, data)
	b.SeekByte(int64(len(data))*4, true)
}

// WriteF32BE writes a slice of float32s to the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) WriteF32BE(off int64, data []float32) {
	if (off + int64(len(data))*4) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
	{
	write_loop:
		b.buf[off+int64(i*4)] = byte(*(*uint32)(unsafe.Pointer(&data[i])) >> 24)
		b.buf[off+int64(1+(i*4))] = byte(*(*uint32)(unsafe.Pointer(&data[i])) >> 16)
		b.buf[off+int64(2+(i*4))] = byte(*(*uint32)(unsafe.Pointer(&data[i])) >> 8)
		b.buf[off+int64(3+(i*4))] = byte(*(*uint32)(unsafe.Pointer(&data[i])))
		i++
		if i < n {
			goto write_loop
		}
	}
}

// WriteF32BENextwrites a slice of float32s to the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteF32BENext(data []float32) {
	b.WriteF32BE(b.off, data)
	b.SeekByte(int64(len(data))*4, true)
}

// WriteF64LE writes a slice of float64s to the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) WriteF64LE(off int64, data []float64) {
	if (off + int64(len(data))*8) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
	{
	write_loop:
		b.buf[off+int64(i*8)] = byte(*(*uint64)(unsafe.Pointer(&data[i])))
		b.buf[off+int64(1+(i*8))] = byte(*(*uint64)(unsafe.Pointer(&data[i])) >> 8)
		b.buf[off+int64(2+(i*8))] = byte(*(*uint64)(unsafe.Pointer(&data[i])) >> 16)
		b.buf[off+int64(3+(i*8))] = byte(*(*uint64)(unsafe.Pointer(&data[i])) >> 24)
		b.buf[off+int64(4+(i*8))] = byte(*(*uint64)(unsafe.Pointer(&data[i])) >> 32)
		b.buf[off+int64(5+(i*8))] = byte(*(*uint64)(unsafe.Pointer(&data[i])) >> 40)
		b.buf[off+int64(6+(i*8))] = byte(*(*uint64)(unsafe.Pointer(&data[i])) >> 48)
		b.buf[off+int64(7+(i*8))] = byte(*(*uint64)(unsafe.Pointer(&data[i])) >> 56)
		i++
		if i < n {
			goto write_loop
		}
	}
}

// WriteF64LENextwrites a slice of float64s to the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteF64LENext(data []float64) {
	b.WriteF64LE(b.off, data)
	b.SeekByte(int64(len(data))*8, true)
}

// WriteF64BE writes a slice of float64s to the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) WriteF64BE(off int64, data []float64) {
	if (off + int64(len(data))*8) > b.cap {
		panic(BufferOverwriteError)
	}
	if off < 0 {
		panic(BufferUnderwriteError)
	}
	i := 0
	n := len(data)
	{
	write_loop:
		b.buf[off+int64(i*8)] = byte(*(*uint64)(unsafe.Pointer(&data[i])) >> 56)
		b.buf[off+int64(1+(i*8))] = byte(*(*uint64)(unsafe.Pointer(&data[i])) >> 48)
		b.buf[off+int64(2+(i*8))] = byte(*(*uint64)(unsafe.Pointer(&data[i])) >> 40)
		b.buf[off+int64(3+(i*8))] = byte(*(*uint64)(unsafe.Pointer(&data[i])) >> 32)
		b.buf[off+int64(4+(i*8))] = byte(*(*uint64)(unsafe.Pointer(&data[i])) >> 24)
		b.buf[off+int64(5+(i*8))] = byte(*(*uint64)(unsafe.Pointer(&data[i])) >> 16)
		b.buf[off+int64(6+(i*8))] = byte(*(*uint64)(unsafe.Pointer(&data[i])) >> 8)
		b.buf[off+int64(7+(i*8))] = byte(*(*uint64)(unsafe.Pointer(&data[i])))
		i++
		if i < n {
			goto write_loop
		}
	}
}

// WriteF64BENextwrites a slice of float64s to the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) WriteF64BENext(data []float64) {
	b.WriteF64BE(b.off, data)
	b.SeekByte(int64(len(data))*8, true)
}

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

// ReadU16LE reads a slice of uint16s from the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) ReadU16LE(off, n int64) (out []uint16) {
	if (off + n*2) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]uint16, n)
	i := int64(0)
	{
	read_loop:
		out[i] = uint16(b.buf[off+(i*2)]) | uint16(b.buf[off+(1+(i*2))])<<8
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadU16LENextreads a slice of uint16s from the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadU16LENext(n int64) (out []uint16) {
	out = b.ReadU16LE(b.off, n)
	b.SeekByte(n*2, true)
	return
}

// ReadU16BE reads a slice of uint16s from the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) ReadU16BE(off, n int64) (out []uint16) {
	if (off + n*2) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]uint16, n)
	i := int64(0)
	{
	read_loop:
		out[i] = uint16(b.buf[off+(1+(i*2))]) | uint16(b.buf[off+(i*2)])<<8
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadU16BENextreads a slice of uint16s from the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadU16BENext(n int64) (out []uint16) {
	out = b.ReadU16BE(b.off, n)
	b.SeekByte(n*2, true)
	return
}

// ReadU32LE reads a slice of uint32s from the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) ReadU32LE(off, n int64) (out []uint32) {
	if (off + n*4) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]uint32, n)
	i := int64(0)
	{
	read_loop:
		out[i] = uint32(b.buf[off+(i*4)]) | uint32(b.buf[off+(1+(i*4))])<<8 | uint32(b.buf[off+(2+(i*4))])<<16 | uint32(b.buf[off+(3+(i*4))])<<24
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadU32LENextreads a slice of uint32s from the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadU32LENext(n int64) (out []uint32) {
	out = b.ReadU32LE(b.off, n)
	b.SeekByte(n*4, true)
	return
}

// ReadU32BE reads a slice of uint32s from the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) ReadU32BE(off, n int64) (out []uint32) {
	if (off + n*4) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]uint32, n)
	i := int64(0)
	{
	read_loop:
		out[i] = uint32(b.buf[off+(3+(i*4))]) | uint32(b.buf[off+(2+(i*4))])<<8 | uint32(b.buf[off+(1+(i*4))])<<16 | uint32(b.buf[off+(i*4)])<<24
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadU32BENextreads a slice of uint32s from the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadU32BENext(n int64) (out []uint32) {
	out = b.ReadU32BE(b.off, n)
	b.SeekByte(n*4, true)
	return
}

// ReadU64LE reads a slice of uint64s from the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) ReadU64LE(off, n int64) (out []uint64) {
	if (off + n*8) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]uint64, n)
	i := int64(0)
	{
	read_loop:
		out[i] = uint64(b.buf[off+(i*8)]) | uint64(b.buf[off+(1+(i*8))])<<8 | uint64(b.buf[off+(2+(i*8))])<<16 | uint64(b.buf[off+(3+(i*8))])<<24 | uint64(b.buf[off+(4+(i*8))])<<32 | uint64(b.buf[off+(5+(i*8))])<<40 | uint64(b.buf[off+(6+(i*8))])<<48 | uint64(b.buf[off+(7+(i*8))])<<56
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadU64LENextreads a slice of uint64s from the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadU64LENext(n int64) (out []uint64) {
	out = b.ReadU64LE(b.off, n)
	b.SeekByte(n*8, true)
	return
}

// ReadU64BE reads a slice of uint64s from the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) ReadU64BE(off, n int64) (out []uint64) {
	if (off + n*8) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]uint64, n)
	i := int64(0)
	{
	read_loop:
		out[i] = uint64(b.buf[off+(7+(i*8))]) | uint64(b.buf[off+(6+(i*8))])<<8 | uint64(b.buf[off+(5+(i*8))])<<16 | uint64(b.buf[off+(4+(i*8))])<<24 | uint64(b.buf[off+(3+(i*8))])<<32 | uint64(b.buf[off+(2+(i*8))])<<40 | uint64(b.buf[off+(1+(i*8))])<<48 | uint64(b.buf[off+(i*8)])<<56
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadU64BENextreads a slice of uint64s from the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadU64BENext(n int64) (out []uint64) {
	out = b.ReadU64BE(b.off, n)
	b.SeekByte(n*8, true)
	return
}


// ReadI16LE reads a slice of int16s from the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) ReadI16LE(off, n int64) (out []int16) {
	if (off + n*2) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]int16, n)
	i := int64(0)
	{
	read_loop:
		out[i] = int16(b.buf[off+(i*2)]) | int16(b.buf[off+(1+(i*2))])<<8
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadI16LENextreads a slice of int16s from the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadI16LENext(n int64) (out []int16) {
	out = b.ReadI16LE(b.off, n)
	b.SeekByte(n*2, true)
	return
}

// ReadI16BE reads a slice of int16s from the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) ReadI16BE(off, n int64) (out []int16) {
	if (off + n*2) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]int16, n)
	i := int64(0)
	{
	read_loop:
		out[i] = int16(b.buf[off+(1+(i*2))]) | int16(b.buf[off+(i*2)])<<8
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadI16BENextreads a slice of int16s from the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadI16BENext(n int64) (out []int16) {
	out = b.ReadI16BE(b.off, n)
	b.SeekByte(n*2, true)
	return
}

// ReadI32LE reads a slice of int32s from the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) ReadI32LE(off, n int64) (out []int32) {
	if (off + n*4) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]int32, n)
	i := int64(0)
	{
	read_loop:
		out[i] = int32(b.buf[off+(i*4)]) | int32(b.buf[off+(1+(i*4))])<<8 | int32(b.buf[off+(2+(i*4))])<<16 | int32(b.buf[off+(3+(i*4))])<<24
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadI32LENextreads a slice of int32s from the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadI32LENext(n int64) (out []int32) {
	out = b.ReadI32LE(b.off, n)
	b.SeekByte(n*4, true)
	return
}

// ReadI32BE reads a slice of int32s from the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) ReadI32BE(off, n int64) (out []int32) {
	if (off + n*4) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]int32, n)
	i := int64(0)
	{
	read_loop:
		out[i] = int32(b.buf[off+(3+(i*4))]) | int32(b.buf[off+(2+(i*4))])<<8 | int32(b.buf[off+(1+(i*4))])<<16 | int32(b.buf[off+(i*4)])<<24
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadI32BENextreads a slice of int32s from the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadI32BENext(n int64) (out []int32) {
	out = b.ReadI32BE(b.off, n)
	b.SeekByte(n*4, true)
	return
}

// ReadI64LE reads a slice of int64s from the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) ReadI64LE(off, n int64) (out []int64) {
	if (off + n*8) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]int64, n)
	i := int64(0)
	{
	read_loop:
		out[i] = int64(b.buf[off+(i*8)]) | int64(b.buf[off+(1+(i*8))])<<8 | int64(b.buf[off+(2+(i*8))])<<16 | int64(b.buf[off+(3+(i*8))])<<24 | int64(b.buf[off+(4+(i*8))])<<32 | int64(b.buf[off+(5+(i*8))])<<40 | int64(b.buf[off+(6+(i*8))])<<48 | int64(b.buf[off+(7+(i*8))])<<56
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadI64LENextreads a slice of int64s from the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadI64LENext(n int64) (out []int64) {
	out = b.ReadI64LE(b.off, n)
	b.SeekByte(n*8, true)
	return
}

// ReadI64BE reads a slice of int64s from the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) ReadI64BE(off, n int64) (out []int64) {
	if (off + n*8) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]int64, n)
	i := int64(0)
	{
	read_loop:
		out[i] = int64(b.buf[off+(7+(i*8))]) | int64(b.buf[off+(6+(i*8))])<<8 | int64(b.buf[off+(5+(i*8))])<<16 | int64(b.buf[off+(4+(i*8))])<<24 | int64(b.buf[off+(3+(i*8))])<<32 | int64(b.buf[off+(2+(i*8))])<<40 | int64(b.buf[off+(1+(i*8))])<<48 | int64(b.buf[off+(i*8)])<<56
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadI64BENextreads a slice of int64s from the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadI64BENext(n int64) (out []int64) {
	out = b.ReadI64BE(b.off, n)
	b.SeekByte(n*8, true)
	return
}

//TODO(superwhiskers): add tests for the following generation directives

// ReadF32LE reads a slice of float32s from the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) ReadF32LE(off, n int64) (out []float32) {
	if (off + n*4) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]float32, n)
	i := int64(0)
	var u uint32
	{
	read_loop:
		u = (uint32(b.buf[off+(i*4)]) | uint32(b.buf[off+(1+(i*4))])<<8 | uint32(b.buf[off+(2+(i*4))])<<16 | uint32(b.buf[off+(3+(i*4))])<<24)
		out[i] = *(*float32)(unsafe.Pointer(&u))
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadF32LENextreads a slice of float32s from the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadF32LENext(n int64) (out []float32) {
	out = b.ReadF32LE(b.off, n)
	b.SeekByte(n*4, true)
	return
}

// ReadF32BE reads a slice of float32s from the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) ReadF32BE(off, n int64) (out []float32) {
	if (off + n*4) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]float32, n)
	i := int64(0)
	var u uint32
	{
	read_loop:
		u = (uint32(b.buf[off+(3+(i*4))]) | uint32(b.buf[off+(2+(i*4))])<<8 | uint32(b.buf[off+(1+(i*4))])<<16 | uint32(b.buf[off+(i*4)])<<24)
		out[i] = *(*float32)(unsafe.Pointer(&u))
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadF32BENextreads a slice of float32s from the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadF32BENext(n int64) (out []float32) {
	out = b.ReadF32BE(b.off, n)
	b.SeekByte(n*4, true)
	return
}

// ReadF64LE reads a slice of float64s from the buffer at the
// specified offset in little-endian without modifying the internal
// offset value
func (b *Buffer) ReadF64LE(off, n int64) (out []float64) {
	if (off + n*8) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]float64, n)
	i := int64(0)
	var u uint64
	{
	read_loop:
		u = (uint64(b.buf[off+(i*8)]) | uint64(b.buf[off+(1+(i*8))])<<8 | uint64(b.buf[off+(2+(i*8))])<<16 | uint64(b.buf[off+(3+(i*8))])<<24 | uint64(b.buf[off+(4+(i*8))])<<32 | uint64(b.buf[off+(5+(i*8))])<<40 | uint64(b.buf[off+(6+(i*8))])<<48 | uint64(b.buf[off+(7+(i*8))])<<56)
		out[i] = *(*float64)(unsafe.Pointer(&u))
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadF64LENextreads a slice of float64s from the buffer at the
// current offset in little-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadF64LENext(n int64) (out []float64) {
	out = b.ReadF64LE(b.off, n)
	b.SeekByte(n*8, true)
	return
}

// ReadF64BE reads a slice of float64s from the buffer at the
// specified offset in big-endian without modifying the internal
// offset value
func (b *Buffer) ReadF64BE(off, n int64) (out []float64) {
	if (off + n*8) > b.cap {
		panic(BufferOverreadError)
	}
	if off < 0 {
		panic(BufferUnderreadError)
	}
	out = make([]float64, n)
	i := int64(0)
	var u uint64
	{
	read_loop:
		u = (uint64(b.buf[off+(7+(i*8))]) | uint64(b.buf[off+(6+(i*8))])<<8 | uint64(b.buf[off+(5+(i*8))])<<16 | uint64(b.buf[off+(4+(i*8))])<<24 | uint64(b.buf[off+(3+(i*8))])<<32 | uint64(b.buf[off+(2+(i*8))])<<40 | uint64(b.buf[off+(1+(i*8))])<<48 | uint64(b.buf[off+(i*8)])<<56)
		out[i] = *(*float64)(unsafe.Pointer(&u))
		i++
		if i < n {
			goto read_loop
		}
	}
	return
}

// ReadF64BENextreads a slice of float64s from the buffer at the
// current offset in big-endian and moves the offset forward the
// amount of bytes written
func (b *Buffer) ReadF64BENext(n int64) (out []float64) {
	out = b.ReadF64BE(b.off, n)
	b.SeekByte(n*8, true)
	return
}

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

	if len(b.buf) > 0 {

		b.obuf = unsafe.Pointer(&b.buf[0])

	}

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
