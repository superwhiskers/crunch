/*

crunch - utilities for taking bytes out of things
Copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.

*/

package v2

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

/*

utilities

*/

var BufferComparer = cmp.Comparer(func(x, y *Buffer) bool {

	return cmp.Equal(x.buf, y.buf) &&
		x.off == y.off &&
		x.cap == y.cap &&
		x.boff == y.boff &&
		x.bcap == y.bcap

})

/*

tests

*/

func TestNewBuffer(t *testing.T) {

	var (
		expected1 = &Buffer{
			buf:  []byte{0x00, 0x00, 0x00, 0x00},
			off:  0x00,
			cap:  4,
			boff: 0x00,
			bcap: 32,
		}
		expected2 = &Buffer{
			buf:  []byte{},
			off:  0x00,
			cap:  0,
			boff: 0x00,
			bcap: 0,
		}
		expected3 = &Buffer{
			buf:  []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			off:  0x00,
			cap:  8,
			boff: 0x00,
			bcap: 64,
		}
	)

	out := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})
	if !cmp.Equal(expected1, out, BufferComparer) {

		t.Fatalf("expected buffer does not match the one gotten (got %#v, expected %#v)", out, expected1)

	}

	out = NewBuffer()
	if !cmp.Equal(expected2, out, BufferComparer) {

		t.Fatalf("expected buffer does not match the one gotten (got %#v, expected %#v)", out, expected2)

	}

	out = NewBuffer([]byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00, 0x00, 0x00})
	if !cmp.Equal(expected3, out, BufferComparer) {

		t.Fatalf("expected buffer does not match the one gotten (got %#v, expected %#v)", out, expected3)

	}

}

func TestBufferBytes(t *testing.T) {

	var expected = []byte{0x01, 0x02, 0x03, 0x04}

	buf := NewBuffer([]byte{0x01, 0x02, 0x03, 0x04})

	out := buf.Bytes()
	if !cmp.Equal(expected, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

}

func TestBufferCapacity(t *testing.T) {

	var expected int64 = 4

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	out := buf.Capacity()
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferBitCapacity(t *testing.T) {

	var expected int64 = 32

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	out := buf.BitCapacity()
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferOffset(t *testing.T) {

	var expected int64 = 0x02

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.Seek(0x02, false)

	out := buf.Offset()
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferBitOffset(t *testing.T) {

	var expected int64 = 0x0f

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.SeekBit(0x0f, false)

	out := buf.BitOffset()
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferRefresh(t *testing.T) {

	var (
		expected1 int64 = 4
		expected2 int64 = 32
	)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.Refresh()

	if expected1 != buf.cap || expected2 != buf.bcap {

		t.Fatalf("expected int64(s) do not match the ones gotten (got %d and %d, expected %d and %d)", buf.cap, buf.bcap, expected1, expected2)

	}

}

func TestBufferReset(t *testing.T) {

	var expected = []byte{}

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.Reset()

	if !cmp.Equal(expected, buf.buf) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", buf.buf, expected)

	}

}

func TestBufferGrow(t *testing.T) {

	var expected int64 = 4

	buf := NewBuffer([]byte{0x00, 0x00})

	buf.Grow(2)

	out := buf.Capacity()
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferSeek(t *testing.T) {

	var expected int64 = 0x04

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.Seek(2, true)
	buf.Seek(2, true)

	out := buf.Offset()
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferSeekBit(t *testing.T) {

	var expected int64 = 0x20

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.SeekBit(16, true)
	buf.SeekBit(16, true)

	out := buf.BitOffset()
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferAlignBit(t *testing.T) {

	var expected int64 = 0x20

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.Seek(0x04, false)
	buf.AlignBit()

	out := buf.BitOffset()
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferAlignByte(t *testing.T) {

	var expected int64 = 0x04

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.SeekBit(0x20, false)
	buf.AlignByte()

	out := buf.Offset()
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferAfter(t *testing.T) {

	var expected int64 = 2

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.Seek(0x01, false)
	out := buf.After()

	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

	out = buf.After(0x01)

	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferAfterBit(t *testing.T) {

	var expected int64 = 16

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.SeekBit(0x0f, false)
	out := buf.AfterBit()

	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

	out = buf.AfterBit(0x0f)

	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferReadByte(t *testing.T) {

	var expected byte = 0x03

	buf := NewBuffer([]byte{0x01, 0x02, 0x03, 0x04})

	out := buf.ReadByte(0x02)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferReadByteNext(t *testing.T) {

	var expected byte = 0x01

	buf := NewBuffer([]byte{0x01, 0x02, 0x03, 0x04})

	out := buf.ReadByteNext()
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferReadBytes(t *testing.T) {

	var expected = []byte{0x03, 0x04}

	buf := NewBuffer([]byte{0x01, 0x02, 0x03, 0x04})

	out := buf.ReadBytes(0x02, 2)
	if !cmp.Equal(expected, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

}

func TestBufferReadBytesNext(t *testing.T) {

	var expected = []byte{0x01, 0x02}

	buf := NewBuffer([]byte{0x01, 0x02, 0x03, 0x04})

	out := buf.ReadBytesNext(2)
	if !cmp.Equal(expected, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

}

func TestBufferReadUNEN(t *testing.T) {

	var (
		expected1 = []uint16{0x01}
		expected2 = []uint16{0x100}
		expected3 = []uint32{0x01}
		expected4 = []uint32{0x1000000}
		expected5 = []uint64{0x01}
		expected6 = []uint64{0x100000000000000}
	)

	buf := NewBuffer([]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	out1 := buf.ReadU16LE(0x00, 1)
	if !cmp.Equal(out1, expected1) {

		t.Fatalf("expected uint16 array does not match the one gotten (got %#v, expected %#v)", out1, expected1)

	}

	out1 = buf.ReadU16BE(0x00, 1)
	if !cmp.Equal(out1, expected2) {

		t.Fatalf("expected uint16 array does not match the one gotten (got %#v, expected %#v)", out1, expected2)

	}

	out2 := buf.ReadU32LE(0x00, 1)
	if !cmp.Equal(out2, expected3) {

		t.Fatalf("expected uint32 array does not match the one gotten (got %#v, expected %#v)", out2, expected3)

	}

	out2 = buf.ReadU32BE(0x00, 1)
	if !cmp.Equal(out2, expected4) {

		t.Fatalf("expected uint32 array does not match the one gotten (got %#v, expected %#v)", out2, expected4)

	}

	out3 := buf.ReadU64LE(0x00, 1)
	if !cmp.Equal(out3, expected5) {

		t.Fatalf("expected uint64 array does not match the one gotten (got %#v, expected %#v)", out3, expected5)

	}

	out3 = buf.ReadU64BE(0x00, 1)
	if !cmp.Equal(out3, expected6) {

		t.Fatalf("expected uint64 array does not match the one gotten (got %#v, expected %#v)", out3, expected6)

	}

}

func TestBufferReadUNENNext(t *testing.T) {

	var (
		expected1 = []uint16{0x01}
		expected2 = []uint16{0x100}
		expected3 = []uint32{0x01}
		expected4 = []uint32{0x1000000}
		expected5 = []uint64{0x01}
		expected6 = []uint64{0x100000000000000}
	)

	buf := NewBuffer([]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	out1 := buf.ReadU16LENext(1)
	if !cmp.Equal(out1, expected1) {

		t.Fatalf("expected uint16 array does not match the one gotten (got %#v, expected %#v)", out1, expected1)

	}

	off := buf.Offset()
	if off != 2 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.Seek(0x00, false)

	out1 = buf.ReadU16BENext(1)
	if !cmp.Equal(out1, expected2) {

		t.Fatalf("expected uint16 array does not match the one gotten (got %#v, expected %#v)", out1, expected2)

	}

	off = buf.Offset()
	if off != 2 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.Seek(0x00, false)

	out2 := buf.ReadU32LENext(1)
	if !cmp.Equal(out2, expected3) {

		t.Fatalf("expected uint32 array does not match the one gotten (got %#v, expected %#v)", out2, expected3)

	}

	off = buf.Offset()
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.Seek(0x00, false)

	out2 = buf.ReadU32BENext(1)
	if !cmp.Equal(out2, expected4) {

		t.Fatalf("expected uint32 array does not match the one gotten (got %#v, expected %#v)", out2, expected4)

	}

	off = buf.Offset()
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.Seek(0x00, false)

	out3 := buf.ReadU64LENext(1)
	if !cmp.Equal(out3, expected5) {

		t.Fatalf("expected uint64 array does not match the one gotten (got %#v, expected %#v)", out3, expected5)

	}

	off = buf.Offset()
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.Seek(0x00, false)

	out3 = buf.ReadU64BENext(1)
	if !cmp.Equal(out3, expected6) {

		t.Fatalf("expected uint64 array does not match the one gotten (got %#v, expected %#v)", out3, expected6)

	}

	off = buf.Offset()
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.Seek(0x00, false)

}

func TestBufferWriteByte(t *testing.T) {

	var expected byte = 0x04

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteByte(0x03, 0x04)

	out := buf.ReadByte(0x03)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferWriteByteNext(t *testing.T) {

	var expected byte = 0x01

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteByteNext(0x01)

	out := buf.ReadByte(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferWriteBytes(t *testing.T) {

	var expected = []byte{0x01, 0x02}

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteBytes(0x02, []byte{0x01, 0x02})

	out := buf.ReadBytes(0x02, 2)
	if !cmp.Equal(expected, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

}

func TestBufferWriteBytesNext(t *testing.T) {

	var expected = []byte{0x01, 0x02}

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteBytesNext([]byte{0x01, 0x02})

	out := buf.ReadBytes(0x00, 2)
	if !cmp.Equal(expected, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

}

func TestBufferWriteUNEN(t *testing.T) {

	var expected byte = 0x01

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	buf.WriteU16LE(0x00, []uint16{0x01, 0x01})

	out := buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	out = buf.ReadBytes(0x02, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.WriteU16BE(0x00, []uint16{0x100, 0x100})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	out = buf.ReadBytes(0x02, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.WriteU32LE(0x00, []uint32{0x01})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.WriteU32BE(0x00, []uint32{0x1000000})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.WriteU64LE(0x00, []uint64{0x01})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.WriteU64BE(0x00, []uint64{0x100000000000000})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

}

func TestBufferWriteUNENNext(t *testing.T) {

	var expected byte = 0x01

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	buf.WriteU16LENext([]uint16{0x01})

	out := buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off := buf.Offset()
	if off != 2 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.Seek(0x00, false)

	buf.WriteU16BENext([]uint16{0x100})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.Offset()
	if off != 2 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.Seek(0x00, false)

	buf.WriteU32LENext([]uint32{0x01})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.Offset()
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.Seek(0x00, false)

	buf.WriteU32BENext([]uint32{0x1000000})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.Offset()
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.Seek(0x00, false)

	buf.WriteU64LENext([]uint64{0x01})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.Offset()
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.Seek(0x00, false)

	buf.WriteU64BENext([]uint64{0x100000000000000})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.Offset()
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.Seek(0x00, false)

}

func TestBufferReadBit(t *testing.T) {

	var expected byte = 1

	buf := NewBuffer([]byte{0x01, 0x00, 0x00, 0x00})

	out := buf.ReadBit(0x07)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferReadBitNext(t *testing.T) {

	var expected byte = 1

	buf := NewBuffer([]byte{0xff, 0x00, 0x00, 0x00})

	out := buf.ReadBitNext()
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferReadBits(t *testing.T) {

	var expected uint64 = 5

	buf := NewBuffer([]byte{0x0d, 0x00, 0x00, 0x00})

	out := buf.ReadBits(0x05, 3)
	if expected != out {

		t.Fatalf("expected uint64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferReadBitsNext(t *testing.T) {

	var expected uint64 = 13

	buf := NewBuffer([]byte{0x0d, 0x00, 0x00, 0x00})

	out := buf.ReadBitsNext(8)
	if expected != out {

		t.Fatalf("expected uint64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferSetBit(t *testing.T) {

	var expected byte = 1

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.SetBit(0x00, 1)

	out := buf.ReadBit(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferSetBitNext(t *testing.T) {

	var expected byte = 1

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.SetBitNext(1)

	out := buf.ReadBit(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferSetBits(t *testing.T) {

	var expected uint64 = 5

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.SetBits(0x06, 5, 3)

	out := buf.ReadBits(0x06, 3)
	if expected != out {

		t.Fatalf("expected uint64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferSetBitsNext(t *testing.T) {

	var expected uint64 = 13

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.SetBitsNext(13, 8)

	out := buf.ReadBits(0x00, 8)
	if expected != out {

		t.Fatalf("expected uint64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferFlipBit(t *testing.T) {

	var expected byte

	buf := NewBuffer([]byte{0xff, 0xff, 0xff, 0xff})

	buf.FlipBit(0x01)

	out := buf.ReadBit(0x01)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferFlipBitNext(t *testing.T) {

	var expected byte = 1

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.FlipBitNext()

	out := buf.ReadBit(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferClearAllBits(t *testing.T) {

	var expected = []byte{0x00, 0x00, 0x00, 0x00}

	buf := NewBuffer([]byte{0xff, 0x00, 0xff, 0x00})

	buf.ClearAllBits()

	if !cmp.Equal(expected, buf.buf) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", buf.buf, expected)

	}

}

func TestBufferSetAllBits(t *testing.T) {

	var expected = []byte{0xff, 0xff, 0xff, 0xff}

	buf := NewBuffer([]byte{0x00, 0xff, 0x00, 0xff})

	buf.SetAllBits()

	if !cmp.Equal(expected, buf.buf) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", buf.buf, expected)

	}

}

func TestBufferFlipAllBits(t *testing.T) {

	var expected = []byte{0xff, 0x00, 0xff, 0x00}

	buf := NewBuffer([]byte{0x00, 0xff, 0x00, 0xff})

	buf.FlipAllBits()

	if !cmp.Equal(expected, buf.buf) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", buf.buf, expected)

	}

}

func TestBufferReadbitPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.readbit(0x20)

}

func TestBufferReadbitPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.readbit(-0x01)

}

func TestBufferSetbitPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.setbit(0x20, 1)

}

func TestBufferSetbitPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.setbit(-0x01, 1)

}

func TestBufferSetbitPanic3(t *testing.T) {

	defer panicChecker(t, BufferInvalidBitError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.setbit(0x00, 2)

}

func TestBufferFlipbitPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.flipbit(0x20)

}

func TestBufferFlipbitPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.flipbit(-0x01)

}

func TestBufferWritePanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.write(0x04, []byte{0x01, 0x01})

}

func TestBufferWritePanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.write(-0x01, []byte{0x01, 0x01})

}

func TestBufferReadPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.read(0x20, 1)

}

func TestBufferReadPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.read(-0x01, 1)

}

func TestBufferWriteU16LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.writeU16LE(0x00, []uint16{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteU16LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.writeU16LE(-0x01, []uint16{0x01, 0x02})

}

func TestBufferWriteU16BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.writeU16BE(0x00, []uint16{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteU16BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.writeU16BE(-0x01, []uint16{0x01, 0x02})

}

func TestBufferWriteU32LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.writeU32LE(0x00, []uint32{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteU32LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.writeU32LE(-0x01, []uint32{0x01})

}

func TestBufferWriteU32BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.writeU32BE(0x00, []uint32{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteU32BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.writeU32BE(-0x01, []uint32{0x01})

}

func TestBufferWriteU64LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.writeU64LE(0x00, []uint64{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteU64LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	buf.writeU64LE(-0x01, []uint64{0x01})

}

func TestBufferWriteU64BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.writeU64BE(0x00, []uint64{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteU64BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	buf.writeU64BE(-0x01, []uint64{0x01})

}

func TestBufferReadU16LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.readU16LE(0x04, 1)

}

func TestBufferReadU16LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.readU16LE(-0x01, 1)

}

func TestBufferReadU16BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.readU16BE(0x04, 1)

}

func TestBufferReadU16BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.readU16BE(-0x01, 1)

}

func TestBufferReadU32LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.readU32LE(0x04, 1)

}

func TestBufferReadU32LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.readU32LE(-0x01, 1)

}

func TestBufferReadU32BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.readU32BE(0x04, 1)

}

func TestBufferReadU32BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.readU32BE(-0x01, 1)

}

func TestBufferReadU64LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.readU64LE(0x04, 1)

}

func TestBufferReadU64LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	_ = buf.readU64LE(-0x01, 1)

}

func TestBufferReadU64BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.readU64BE(0x04, 1)

}

func TestBufferReadU64BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	_ = buf.readU64BE(-0x01, 1)

}

/*

benchmarks

*/

func BenchmarkBufferWrite(b *testing.B) {

	b.ReportAllocs()

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {

		buf.WriteBytes(0x00, []byte{0x01, 0x02})

	}

}

func BenchmarkBufferRead(b *testing.B) {

	b.ReportAllocs()

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	var out []byte
	for n := 0; n < b.N; n++ {

		out = buf.ReadBytes(0x00, 2)

	}

	_ = out

}

func BenchmarkBufferWriteU32LE(b *testing.B) {

	b.ReportAllocs()

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {

		buf.WriteU32LE(0x00, []uint32{0x01, 0x02})

	}

}

func BenchmarkBufferReadU32LE(b *testing.B) {

	b.ReportAllocs()

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	var out []uint32
	for n := 0; n < b.N; n++ {

		out = buf.ReadU32LE(0x00, 2)

	}

	_ = out

}
