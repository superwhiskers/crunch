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

// UnsafeBuffer implements a faster and even-lower-memory buffer type in go that handles bytes easily. it is not safe
// for concurrent usage out of the box, you are required to handle that yourself by calling the Lock and Unlock methods on it
type UnsafeBuffer struct {
	buf  []byte
	beg  uintptr
	off  uintptr
	cap  int64

	sync.Mutex
}

// NewUnsafeBuffer initilaizes a new UnsafeBuffer with the provided byte slice(s) stored inside in the order provided
func NewUnsafeBuffer(out **UnsafeBuffer, slices ...[]byte) {

	*out = &UnsafeBuffer{
		buf:  []byte{},
	}

	switch len(slices) {

	case 0:
		break

	case 1:
		(*out).buf = slices[0]
		break

	default:
		for _, s := range slices {

			(*out).buf = append((*out).buf, s...)

		}

	}

	(*out).refresh()
	return

}

/* internal use methods */

/* byte buffer methods */

// write writes a slice of bytes to the buffer at the specified offset
func (b *UnsafeBuffer) write(off uintptr, data []byte) {

	for i := range data {

		*(*byte)(unsafe.Pointer(off + uintptr(i))) = data[i]

	}

	return
}

// read reads n bytes from the buffer at the specified offset
// TODO: implement this faster somehow
func (b *UnsafeBuffer) read(out *[]byte, off uintptr, n int64) {

	*out = b.buf[off-b.beg : (off-b.beg)+uintptr(n)]
	return

}

// seek seeks to position off of the byte buffer relative to the current position or exact
func (b *UnsafeBuffer) seek(off int64, relative bool) {

	if relative == true {

		b.off = b.off + uintptr(off)

	} else {

		b.off = b.beg + uintptr(off)

	}

	return

}

// after returns the amount of bytes located after the current position or the specified one
func (b *UnsafeBuffer) after(out *int64, off ...int64) {

	if len(off) == 0 {

		*out = b.cap - int64(b.off)

	}
	*out = b.cap - off[0]

	return

}

/* generic methods */

// grow grows the buffer by n bytes
func (b *UnsafeBuffer) grow(n int64) {

	b.buf = append(b.buf, make([]byte, n)...)
	b.refresh()

	return

}

// refresh updates the internal statistics of the byte buffer forcefully
func (b *UnsafeBuffer) refresh() {

	b.cap = int64(len(b.buf))
	b.off = uintptr(unsafe.Pointer(&b.buf[0])) + (b.off - b.beg)
	b.beg = uintptr(unsafe.Pointer(&b.buf[0]))

	return

}

/* public methods */

// Bytes stores the internal byte slice of the buffer in out
func (b *UnsafeBuffer) Bytes(out *[]byte) {

	*out = b.buf
	return

}

// Capacity stores the capacity of the buffer in out
func (b *UnsafeBuffer) Capacity(out *int64) {

	*out = b.cap
	return

}

// Offset stores the current offset of the buffer in out
func (b *UnsafeBuffer) Offset(out *int64) {

	*out = int64(b.off)
	return

}

// Refresh updates the cached internal statistics of the buffer forcefully
func (b *UnsafeBuffer) Refresh() {

	b.refresh()
	return

}

// Grow makes the buffer's capacity bigger by n bytes
func (b *UnsafeBuffer) Grow(n int64) {

	b.grow(n)
	return

}

// Seek seeks to position off of the buffer relative to the current position or exact
func (b *UnsafeBuffer) Seek(off int64, relative bool) {

	b.seek(off, relative)
	return

}

// After stores the amount of bytes located after the current position or the specified one in out
func (b *UnsafeBuffer) After(out *int64, off ...int64) {

	b.after(out, off...)
	return

}

// ReadBytes stores the next n bytes from the specified offset without modifying the internal offset value in out
func (b *UnsafeBuffer) ReadBytes(out *[]byte, off, n int64) {

	b.read(out, b.off + uintptr(off), n)
	return

}

// ReadBytesNext stores the next n bytes from the current offset and moves the offset forward the amount of bytes read in out
func (b *UnsafeBuffer) ReadBytesNext(out *[]byte, n int64) {

	b.read(out, b.off, n)
	b.seek(n, true)
	return

}

// WriteBytes writes bytes to the buffer at the specified offset without modifying the internal offset value
func (b *UnsafeBuffer) WriteBytes(off int64, data []byte) {

	b.write(b.off + uintptr(off), data)
	return

}

// WriteBytesNext writes bytes to the buffer at the current offset and moves the offset forward the amount of bytes written
func (b *UnsafeBuffer) WriteBytesNext(data []byte) {

	b.write(b.off, data)
	b.seek(int64(len(data)), true)
	return

}

