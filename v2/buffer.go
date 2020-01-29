/*

crunch - utilities for taking bytes out of things
copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>

this source code form is subject to the terms of the mozilla public
license, v. 2.0. if a copy of the mpl was not distributed with this
file, you can obtain one at http://mozilla.org/MPL/2.0/.

*/

package v2

import "sync"

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

		panic(BufferOverreadError)

	}

	if off < 0x00 {

		panic(BufferUnderreadError)

	}

	b.Lock()
	defer b.Unlock()

	return atob((b.buf[off/8] & (1 << uint(7-(off%8)))) != 0)

}

// readbits reads n bits from the bitfield at the specified offset
func (b *Buffer) readbits(off, n int64) (out uint64) {

	i := int64(0)

	{
	read_loop:
		out = (out << uint64(1)) | uint64(b.readbit(off+i))
		i++
		if i < n {

			goto read_loop

		}
	}

	return

}

// setbit sets a bit in the bitfield to the specified value
func (b *Buffer) setbit(off int64, data byte) {

	if off > (b.bcap - 1) {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	b.Lock()
	defer b.Unlock()

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
func (b *Buffer) setbits(off int64, data uint64, n int64) {

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
func (b *Buffer) flipbit(off int64) {

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

// clearallbits sets all of the buffer's bits to 0
func (b *Buffer) clearallbits() {

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

// setallbits sets all of the buffer's bits to 1
func (b *Buffer) setallbits() {

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

// flipallbits flips all of the buffer's bits
func (b *Buffer) flipallbits() {

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

// seekbit seeks to position off of the bitfield relative to the current position or exact
func (b *Buffer) seekbit(off int64, relative bool) {

	b.Lock()
	defer b.Unlock()

	if relative {

		b.boff += off

	} else {

		b.boff = off

	}

}

// afterbit returns the amount of bits located after the current position or the specified one
func (b *Buffer) afterbit(off ...int64) int64 {

	if len(off) == 0 {

		return b.bcap - b.boff - 1

	}
	return b.bcap - off[0] - 1

}

/* byte buffer methods */

// write writes a slice of bytes to the buffer at the specified offset
func (b *Buffer) write(off int64, data []byte) {

	if (off + int64(len(data))) > b.cap {

		panic(BufferOverwriteError)

	}

	if off < 0x00 {

		panic(BufferUnderwriteError)

	}

	b.Lock()

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

	b.Unlock()

}

// writeU16LE writes a slice of uint16s to the buffer at the specified offset in little-endian
func (b *Buffer) writeU16LE(off int64, data []uint16) {

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

// writeU16BE writes a slice of uint16s to the buffer at the specified offset in big-endian
func (b *Buffer) writeU16BE(off int64, data []uint16) {

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

// writeU32LE writes a slice of uint32s to the buffer at the specified offset in little-endian
func (b *Buffer) writeU32LE(off int64, data []uint32) {

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

// writeU32BE writes a slice of uint32s to the buffer at the specified offset in big-endian
func (b *Buffer) writeU32BE(off int64, data []uint32) {

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

// writeU64LE writes a slice of uint64s to the buffer at the specified offset in little-endian
func (b *Buffer) writeU64LE(off int64, data []uint64) {

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

// writeU64BE writes a slice of uint64s to the buffer at the specified offset in big-endian
func (b *Buffer) writeU64BE(off int64, data []uint64) {

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

// read reads n bytes from the buffer at the specified offset
func (b *Buffer) read(off, n int64) []byte {

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

// readU16LE reads a slice of uint16s from the buffer at the specified offset in little-endian
func (b *Buffer) readU16LE(off, n int64) (out []uint16) {

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

// readU16BE reads a slice of uint16s from the buffer at the specified offset in big-endian
func (b *Buffer) readU16BE(off, n int64) (out []uint16) {

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

// readU32LE reads a slice of uint32s from the buffer at the specified offset in little-endian
func (b *Buffer) readU32LE(off, n int64) (out []uint32) {

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

// readU32BE reads a slice of uint32s from the buffer at the specified offset in big-endian
func (b *Buffer) readU32BE(off, n int64) (out []uint32) {

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

// readU64LE reads a slice of uint64s from the buffer at the specified offset in little-endian
func (b *Buffer) readU64LE(off, n int64) (out []uint64) {

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

// readU64BE reads a slice of uint64s from the buffer at the specified offset in big-endian
func (b *Buffer) readU64BE(off, n int64) (out []uint64) {

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

// seek seeks to position off of the byte buffer relative to the current position or exact
func (b *Buffer) seek(off int64, relative bool) {

	b.Lock()
	defer b.Unlock()

	if relative {

		b.off += off

	} else {

		b.off = off

	}

}

// after returns the amount of bytes located after the current position or the specified one
func (b *Buffer) after(off ...int64) int64 {

	if len(off) == 0 {

		return b.cap - b.off - 1

	}
	return b.cap - off[0] - 1

}

/* generic methods */

// grow grows the buffer by n bytes
func (b *Buffer) grow(n int64) {

	b.Lock()

	b.buf = append(b.buf, make([]byte, n)...)

	b.Unlock()

	b.refresh()

}

// refresh updates the internal statistics of the byte buffer forcefully
func (b *Buffer) refresh() {

	b.Lock()
	defer b.Unlock()

	b.cap = int64(len(b.buf))
	b.bcap = b.cap * 8

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

// reset resets the buffer
func (b *Buffer) reset() {

	b.Lock()
	defer b.Unlock()

	b.buf = []byte{}
	b.off = 0x00
	b.boff = 0x00
	b.cap = 0
	b.bcap = 0

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

// Offset returns the current offset of the buffer
func (b *Buffer) Offset() int64 {

	return b.off

}

// BitOffset returns the current bit offset of the buffer
func (b *Buffer) BitOffset() int64 {

	return b.boff

}

// Refresh updates the cached internal statistics of the buffer forcefully
func (b *Buffer) Refresh() {

	b.refresh()

}

// Reset resets the entire buffer
func (b *Buffer) Reset() {

	b.reset()

}

// Grow makes the buffer's capacity bigger by n bytes
func (b *Buffer) Grow(n int64) {

	b.grow(n)

}

// Seek seeks to position off of the buffer relative to the current position or exact
func (b *Buffer) Seek(off int64, relative bool) {

	b.seek(off, relative)

}

// SeekBit seeks to bit position off of the the buffer relative to the current position or exact
func (b *Buffer) SeekBit(off int64, relative bool) {

	b.seekbit(off, relative)

}

// AlignBit aligns the bit offset to the byte offset
func (b *Buffer) AlignBit() {

	b.alignbit()

}

// AlignByte aligns the byte offset to the bit offset
func (b *Buffer) AlignByte() {

	b.alignbyte()

}

// After returns the amount of bytes located after the current position or the specified one
func (b *Buffer) After(off ...int64) int64 {

	return b.after(off...)

}

// AfterBit returns the amount of bits located after the current bit position or the specified one
func (b *Buffer) AfterBit(off ...int64) int64 {

	return b.afterbit(off...)

}

// ReadByte returns the next byte from the specified offset without modifying the internal offset value
func (b *Buffer) ReadByte(off int64) byte {

	return b.read(off, 1)[0]

}

// ReadBytes returns the next n bytes from the specified offset without modifying the internal offset value
func (b *Buffer) ReadBytes(off, n int64) []byte {

	return b.read(off, n)

}

// ReadU16LE reads a slice of uint16s from the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *Buffer) ReadU16LE(off, n int64) []uint16 {

	return b.readU16LE(off, n)

}

// ReadU16BE reads a slice of uint16s from the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *Buffer) ReadU16BE(off, n int64) []uint16 {

	return b.readU16BE(off, n)

}

// ReadU32LE reads a slice of uint32s from the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *Buffer) ReadU32LE(off, n int64) []uint32 {

	return b.readU32LE(off, n)

}

// ReadU32BE reads a slice of uint32s from the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *Buffer) ReadU32BE(off, n int64) []uint32 {

	return b.readU32BE(off, n)

}

// ReadU64LE reads a slice of uint64s from the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *Buffer) ReadU64LE(off, n int64) []uint64 {

	return b.readU64LE(off, n)

}

// ReadU64BE reads a slice of uint64s from the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *Buffer) ReadU64BE(off, n int64) []uint64 {

	return b.readU64BE(off, n)

}

// ReadByteNext returns the next byte from the current offset and moves the offset forward a byte
func (b *Buffer) ReadByteNext() (out byte) {

	out = b.read(b.off, 1)[0]
	b.seek(1, true)
	return

}

// ReadBytesNext returns the next n bytes from the current offset and moves the offset forward the amount of bytes read
func (b *Buffer) ReadBytesNext(n int64) (out []byte) {

	out = b.read(b.off, n)
	b.seek(n, true)
	return

}

// ReadU16LENext reads a slice of uint16s from the buffer at the current offset in little-endian and moves the offset forward the amount of bytes read
func (b *Buffer) ReadU16LENext(n int64) (out []uint16) {

	out = b.readU16LE(b.off, n)
	b.seek(n*2, true)
	return

}

// ReadU16BENext reads a slice of uint16s from the buffer at the current offset in big-endian and moves the offset forward the amount of bytes read
func (b *Buffer) ReadU16BENext(n int64) (out []uint16) {

	out = b.readU16BE(b.off, n)
	b.seek(n*2, true)
	return

}

// ReadU32LENext reads a slice of uint32s from the buffer at the current offset in little-endian and moves the offset forward the amount of bytes read
func (b *Buffer) ReadU32LENext(n int64) (out []uint32) {

	out = b.readU32LE(b.off, n)
	b.seek(n*4, true)
	return

}

// ReadU32BENext reads a slice of uint32s from the buffer at the current offset in big-endian and moves the offset forward the amount of bytes read
func (b *Buffer) ReadU32BENext(n int64) (out []uint32) {

	out = b.readU32BE(b.off, n)
	b.seek(n*4, true)
	return

}

// ReadU64LENext reads a slice of uint64s from the buffer at the current offset in little-endian and moves the offset forward the amount of bytes read
func (b *Buffer) ReadU64LENext(n int64) (out []uint64) {

	out = b.readU64LE(b.off, n)
	b.seek(n*8, true)
	return

}

// ReadU64BENext reads a slice of uint64s from the buffer at the current offset in big-endian and moves the offset forward the amount of bytes read
func (b *Buffer) ReadU64BENext(n int64) (out []uint64) {

	out = b.readU64BE(b.off, n)
	b.seek(n*8, true)
	return

}

// WriteByte writes a byte to the buffer at the specified offset without modifying the internal offset value
func (b *Buffer) WriteByte(off int64, data byte) {

	b.write(off, []byte{data})

}

// WriteBytes writes bytes to the buffer at the specified offset without modifying the internal offset value
func (b *Buffer) WriteBytes(off int64, data []byte) {

	b.write(off, data)

}

// WriteU16LE writes a slice of uint16s to the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *Buffer) WriteU16LE(off int64, data []uint16) {

	b.writeU16LE(off, data)

}

// WriteU16BE writes a slice of uint16s to the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *Buffer) WriteU16BE(off int64, data []uint16) {

	b.writeU16BE(off, data)

}

// WriteU32LE writes a slice of uint32s to the buffer at the specified offset in little-endian without modifying the internal offset value
func (b *Buffer) WriteU32LE(off int64, data []uint32) {

	b.writeU32LE(off, data)

}

// WriteU32BE writes a slice of uint32s to the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *Buffer) WriteU32BE(off int64, data []uint32) {

	b.writeU32BE(off, data)

}

// WriteU64LE writes a slice of uint64s to the buffer at the specfied offset in little-endian without modifying the internal offset value
func (b *Buffer) WriteU64LE(off int64, data []uint64) {

	b.writeU64LE(off, data)

}

// WriteU64BE writes a slice of uint64s to the buffer at the specified offset in big-endian without modifying the internal offset value
func (b *Buffer) WriteU64BE(off int64, data []uint64) {

	b.writeU64BE(off, data)

}

// WriteByteNext writes a byte to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *Buffer) WriteByteNext(data byte) {

	b.write(b.off, []byte{data})
	b.seek(1, true)

}

// WriteBytesNext writes bytes to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *Buffer) WriteBytesNext(data []byte) {

	b.write(b.off, data)
	b.seek(int64(len(data)), true)

}

// WriteU16LENext writes a slice of uint16s to the buffer at the current offset in little-endian and moves the offset forward the amount of bytes written
func (b *Buffer) WriteU16LENext(data []uint16) {

	b.writeU16LE(b.off, data)
	b.seek(int64(len(data))*2, true)

}

// WriteU16BENext writes a slice of uint16s to the buffer at the current offset in big-endian and moves the offset forward the amount of bytes written
func (b *Buffer) WriteU16BENext(data []uint16) {

	b.writeU16BE(b.off, data)
	b.seek(int64(len(data))*2, true)

}

// WriteU32LENext writes a slice of uint32s to the buffer at the current offset in little-endian and moves the offset forward the amount of bytes written
func (b *Buffer) WriteU32LENext(data []uint32) {

	b.writeU32LE(b.off, data)
	b.seek(int64(len(data))*4, true)

}

// WriteU32BENext writes a slice of uint32s to the buffer at the current offset in big-endian and moves the offset forward the amount of bytes written
func (b *Buffer) WriteU32BENext(data []uint32) {

	b.writeU32BE(b.off, data)
	b.seek(int64(len(data))*4, true)

}

// WriteU64LENext writes a slice of uint64s to the buffer at the current offset in little-endian and moves the offset forward the amount of bytes written
func (b *Buffer) WriteU64LENext(data []uint64) {

	b.writeU64LE(b.off, data)
	b.seek(int64(len(data))*8, true)

}

// WriteU64BENext writes a slice of uint64s to the buffer at the current offset in big-endian and moves the offset forward the amount of bytes written
func (b *Buffer) WriteU64BENext(data []uint64) {

	b.writeU64BE(b.off, data)
	b.seek(int64(len(data))*8, true)

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
func (b *Buffer) SetBit(off int64, data byte) {

	b.setbit(off, data)

}

// SetBits sets the next n bits from the specified offset without modifying the internal offset value
func (b *Buffer) SetBits(off int64, data uint64, n int64) {

	b.setbits(off, data, n)

}

// SetBitNext sets the next bit from the current offset and moves the offset forward a bit
func (b *Buffer) SetBitNext(data byte) {

	b.setbit(b.boff, data)
	b.seekbit(1, true)

}

// SetBitsNext sets the next n bits from the current offset and moves the offset forward the amount of bits set
func (b *Buffer) SetBitsNext(data uint64, n int64) {

	b.setbits(b.boff, data, n)
	b.seekbit(n, true)

}

// FlipBit flips the bit located at the specified offset without modifying the internal offset value
func (b *Buffer) FlipBit(off int64) {

	b.flipbit(off)

}

// FlipBitNext flips the next bit from the current offset and moves the offset forward a bit
func (b *Buffer) FlipBitNext() {

	b.flipbit(b.boff)
	b.seekbit(1, true)

}

// ClearAllBits sets all of the buffer's bits to 0
func (b *Buffer) ClearAllBits() {

	b.clearallbits()

}

// SetAllBits sets all of the buffer's bits to 1
func (b *Buffer) SetAllBits() {

	b.setallbits()

}

// FlipAllBits flips all of the buffer's bits
func (b *Buffer) FlipAllBits() {

	b.flipallbits()

}
