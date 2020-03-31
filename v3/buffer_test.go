/*

crunch - utilities for taking bytes out of things
copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>

this source code form is subject to the terms of the mozilla public
license, v. 2.0. if a copy of the mpl was not distributed with this
file, you can obtain one at http://mozilla.org/MPL/2.0/.

*/

package v3

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

func TestBufferByteCapacity(t *testing.T) {

	var expected int64 = 4

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	out := buf.ByteCapacity()
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

func TestBufferByteOffset(t *testing.T) {

	var expected int64 = 0x02

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.SeekByte(0x02, false)

	out := buf.ByteOffset()
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

func TestBufferTruncateLeft(t *testing.T) {

	var (
		expected1 byte  = 0x02
		expected2 int64 = 1
	)

	buf := NewBuffer([]byte{0x01, 0x02})

	buf.TruncateLeft(1)

	t.Log(buf.buf)

	out1 := buf.ReadByte(0x00)
	if expected1 != out1 {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out1, expected1)

	}

	out2 := buf.ByteCapacity()
	if expected2 != out2 {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out2, expected2)

	}

}

func TestBufferTruncateRight(t *testing.T) {

	var (
		expected1 byte  = 0x01
		expected2 int64 = 1
	)

	buf := NewBuffer([]byte{0x01, 0x02})

	buf.TruncateRight(1)

	t.Log(buf.buf)

	out1 := buf.ReadByte(0x00)
	if expected1 != out1 {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out1, expected1)

	}

	out2 := buf.ByteCapacity()
	if expected2 != out2 {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out2, expected2)

	}

}

func TestBufferGrow(t *testing.T) {

	var expected int64 = 4

	buf := NewBuffer([]byte{0x00, 0x00})

	buf.Grow(2)

	out := buf.ByteCapacity()
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferGrow2(t *testing.T) {

	var expected int64 = 4

	buf := NewBuffer(make([]byte, 2, 4))

	buf.Grow(2)

	var out = buf.ByteCapacity()
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferSeekByte(t *testing.T) {

	var expected int64 = 0x04

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.SeekByte(2, true)
	buf.SeekByte(2, true)

	out := buf.ByteOffset()
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

	buf.SeekByte(0x04, false)
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

	out := buf.ByteOffset()
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferAfterByte(t *testing.T) {

	var expected int64 = 2

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.SeekByte(0x01, false)
	out := buf.AfterByte()

	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

	out = buf.AfterByte(0x01)

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

func TestBufferReadSNEN(t *testing.T) {

	var (
		expected_readu16le = []uint16{0x01, 0x00}
		expected_readu16be = []uint16{0x100, 0x00}
		expected_readu32le = []uint32{0x01, 0x00}
		expected_readu32be = []uint32{0x1000000, 0x00}
		expected_readu64le = []uint64{0x01, 0x01}
		expected_readu64be = []uint64{0x100000000000000, 0x100000000000000}

		expected_readi16le = []int16{0x01, 0x00}
		expected_readi16be = []int16{0x100, 0x00}
		expected_readi32le = []int32{0x01, 0x00}
		expected_readi32be = []int32{0x1000000, 0x00}
		expected_readi64le = []int64{0x01, 0x01}
		expected_readi64be = []int64{0x100000000000000, 0x100000000000000}
	)

	buf := NewBuffer([]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	//
	// u16le
	//

	out1 := buf.ReadU16LENext(2)
	if !cmp.Equal(out1, expected_readu16le) {

		t.Fatalf("expected uint16 array does not match the one gotten (got %#v, expected %#v)", out1, expected_readu16le)

	}

	off := buf.ByteOffset()
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u16be
	//

	out1 = buf.ReadU16BENext(2)
	if !cmp.Equal(out1, expected_readu16be) {

		t.Fatalf("expected uint16 array does not match the one gotten (got %#v, expected %#v)", out1, expected_readu16be)

	}

	off = buf.ByteOffset()
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u32le
	//

	out2 := buf.ReadU32LENext(2)
	if !cmp.Equal(out2, expected_readu32le) {

		t.Fatalf("expected uint32 array does not match the one gotten (got %#v, expected %#v)", out2, expected_readu32le)

	}

	off = buf.ByteOffset()
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u32be
	//

	out2 = buf.ReadU32BENext(2)
	if !cmp.Equal(out2, expected_readu32be) {

		t.Fatalf("expected uint32 array does not match the one gotten (got %#v, expected %#v)", out2, expected_readu32be)

	}

	off = buf.ByteOffset()
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u64le
	//

	out3 := buf.ReadU64LENext(2)
	if !cmp.Equal(out3, expected_readu64le) {

		t.Fatalf("expected uint64 array does not match the one gotten (got %#v, expected %#v)", out3, expected_readu64le)

	}

	off = buf.ByteOffset()
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u64be
	//

	out3 = buf.ReadU64BENext(2)
	if !cmp.Equal(out3, expected_readu64be) {

		t.Fatalf("expected uint64 array does not match the one gotten (got %#v, expected %#v)", out3, expected_readu64be)

	}

	off = buf.ByteOffset()
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i16le
	//

	out4 := buf.ReadI16LENext(2)
	if !cmp.Equal(out4, expected_readi16le) {

		t.Fatalf("expected int16 array does not match the one gotten (got %#v, expected %#v)", out4, expected_readi16le)

	}
	off = buf.ByteOffset()
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i16be
	//

	out4 = buf.ReadI16BENext(2)
	if !cmp.Equal(out4, expected_readi16be) {

		t.Fatalf("expected int16 array does not match the one gotten (got %#v, expected %#v)", out4, expected_readi16be)

	}
	off = buf.ByteOffset()
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i32le
	//

	out5 := buf.ReadI32LENext(2)
	if !cmp.Equal(out5, expected_readi32le) {

		t.Fatalf("expected int32 array does not match the one gotten (got %#v, expected %#v)", out5, expected_readi32le)

	}
	off = buf.ByteOffset()
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i32be
	//

	out5 = buf.ReadI32BENext(2)
	if !cmp.Equal(out5, expected_readi32be) {

		t.Fatalf("expected int32 array does not match the one gotten (got %#v, expected %#v)", out5, expected_readi32be)

	}
	off = buf.ByteOffset()
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i64le
	//

	out6 := buf.ReadI64LENext(2)
	if !cmp.Equal(out6, expected_readi64le) {

		t.Fatalf("expected int64 array does not match the one gotten (got %#v, expected %#v)", out6, expected_readi64le)

	}
	off = buf.ByteOffset()
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i64be
	//

	out6 = buf.ReadI64BENext(2)
	if !cmp.Equal(out6, expected_readi64be) {

		t.Fatalf("expected int64 array does not match the one gotten (got %#v, expected %#v)", out6, expected_readi64be)

	}
	off = buf.ByteOffset()
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

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

func TestBufferWriteSNEN(t *testing.T) {

	var expected byte = 0x01

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	//
	// u16le
	//

	buf.WriteU16LENext([]uint16{0x01, 0x01})

	out := buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	out = buf.ReadBytes(0x02, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off := buf.ByteOffset()
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u16be
	//

	buf.WriteU16BENext([]uint16{0x100, 0x100})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	out = buf.ReadBytes(0x02, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.ByteOffset()
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u32le
	//

	buf.WriteU32LENext([]uint32{0x01, 0x01})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	out = buf.ReadBytes(0x04, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.ByteOffset()
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u32be
	//

	buf.WriteU32BENext([]uint32{0x1000000, 0x1000000})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	out = buf.ReadBytes(0x04, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.ByteOffset()
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u64le
	//

	buf.WriteU64LENext([]uint64{0x01, 0x01})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	out = buf.ReadBytes(0x08, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.ByteOffset()
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u64be
	//

	buf.WriteU64BENext([]uint64{0x100000000000000, 0x100000000000000})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	out = buf.ReadBytes(0x08, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.ByteOffset()
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i16le
	//

	buf.WriteI16LENext([]int16{0x01, 0x01})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	out = buf.ReadBytes(0x02, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.ByteOffset()
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i16be
	//

	buf.WriteI16BENext([]int16{0x100, 0x100})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	out = buf.ReadBytes(0x02, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.ByteOffset()
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i32le
	//

	buf.WriteI32LENext([]int32{0x01, 0x01})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	out = buf.ReadBytes(0x04, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.ByteOffset()
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i32be
	//

	buf.WriteI32BENext([]int32{0x1000000, 0x1000000})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	out = buf.ReadBytes(0x04, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.ByteOffset()
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i64le
	//

	buf.WriteI64LENext([]int64{0x01, 0x01})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	out = buf.ReadBytes(0x08, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.ByteOffset()
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i64be
	//

	buf.WriteI64BENext([]int64{0x100000000000000, 0x100000000000000})

	out = buf.ReadBytes(0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	out = buf.ReadBytes(0x08, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	off = buf.ByteOffset()
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

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

	buf.SetBit(0x00)

	out := buf.ReadBit(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferSetBitNext(t *testing.T) {

	var expected byte = 1

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.SetBitNext()

	out := buf.ReadBit(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferClearBit(t *testing.T) {

	var expected byte

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.ClearBit(0x00)

	out := buf.ReadBit(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferClearBitNext(t *testing.T) {

	var expected byte

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.ClearBitNext()

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

	_ = buf.ReadBit(0x20)

}

func TestBufferReadbitPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadBit(-0x01)

}

func TestBufferSetbitPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.SetBit(0x20)

}

func TestBufferSetbitPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.SetBit(-0x01)

}

func TestBufferClearbitPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.ClearBit(0x20)

}

func TestBufferClearbitPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.ClearBit(-0x01)

}

func TestBufferFlipbitPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.FlipBit(0x20)

}

func TestBufferFlipbitPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.FlipBit(-0x01)

}

func TestBufferWriteBytesPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteBytes(0x04, []byte{0x01, 0x01})

}

func TestBufferWriteBytesPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteBytes(-0x01, []byte{0x01, 0x01})

}

func TestBufferReadBytesPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadBytes(0x20, 1)

}

func TestBufferReadBytesPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadBytes(-0x01, 1)

}

func TestBufferWriteU16LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteU16LE(0x00, []uint16{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteU16LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteU16LE(-0x01, []uint16{0x01, 0x02})

}

func TestBufferWriteU16BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteU16BE(0x00, []uint16{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteU16BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteU16BE(-0x01, []uint16{0x01, 0x02})

}

func TestBufferWriteU32LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteU32LE(0x00, []uint32{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteU32LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteU32LE(-0x01, []uint32{0x01})

}

func TestBufferWriteU32BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteU32BE(0x00, []uint32{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteU32BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteU32BE(-0x01, []uint32{0x01})

}

func TestBufferWriteU64LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteU64LE(0x00, []uint64{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteU64LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	buf.WriteU64LE(-0x01, []uint64{0x01})

}

func TestBufferWriteU64BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteU64BE(0x00, []uint64{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteU64BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	buf.WriteU64BE(-0x01, []uint64{0x01})

}

func TestBufferWriteI16LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteI16LE(0x00, []int16{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteI16LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteI16LE(-0x01, []int16{0x01, 0x02})

}

func TestBufferWriteI16BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteI16BE(0x00, []int16{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteI16BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteI16BE(-0x01, []int16{0x01, 0x02})

}

func TestBufferWriteI32LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteI32LE(0x00, []int32{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteI32LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteI32LE(-0x01, []int32{0x01})

}

func TestBufferWriteI32BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteI32BE(0x00, []int32{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteI32BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteI32BE(-0x01, []int32{0x01})

}

func TestBufferWriteI64LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteI64LE(0x00, []int64{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteI64LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	buf.WriteI64LE(-0x01, []int64{0x01})

}

func TestBufferWriteI64BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteI64BE(0x00, []int64{0x01, 0x02, 0x03, 0x04})

}

func TestBufferWriteI64BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderwriteError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	buf.WriteI64BE(-0x01, []int64{0x01})

}

func TestBufferReadU16LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadU16LE(0x04, 1)

}

func TestBufferReadU16LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadU16LE(-0x01, 1)

}

func TestBufferReadU16BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadU16BE(0x04, 1)

}

func TestBufferReadU16BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadU16BE(-0x01, 1)

}

func TestBufferReadU32LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadU32LE(0x04, 1)

}

func TestBufferReadU32LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadU32LE(-0x01, 1)

}

func TestBufferReadU32BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadU32BE(0x04, 1)

}

func TestBufferReadU32BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadU32BE(-0x01, 1)

}

func TestBufferReadU64LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadU64LE(0x04, 1)

}

func TestBufferReadU64LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadU64LE(-0x01, 1)

}

func TestBufferReadU64BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadU64BE(0x04, 1)

}

func TestBufferReadU64BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadU64BE(-0x01, 1)

}

func TestBufferReadI16LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadI16LE(0x04, 1)

}

func TestBufferReadI16LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadI16LE(-0x01, 1)

}

func TestBufferReadI16BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadI16BE(0x04, 1)

}

func TestBufferReadI16BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadI16BE(-0x01, 1)

}

func TestBufferReadI32LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadI32LE(0x04, 1)

}

func TestBufferReadI32LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadI32LE(-0x01, 1)

}

func TestBufferReadI32BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadI32BE(0x04, 1)

}

func TestBufferReadI32BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadI32BE(-0x01, 1)

}

func TestBufferReadI64LEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadI64LE(0x04, 1)

}

func TestBufferReadI64LEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadI64LE(-0x01, 1)

}

func TestBufferReadI64BEPanic1(t *testing.T) {

	defer panicChecker(t, BufferOverreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadI64BE(0x04, 1)

}

func TestBufferReadI64BEPanic2(t *testing.T) {

	defer panicChecker(t, BufferUnderreadError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	_ = buf.ReadI64BE(-0x01, 1)

}

func TestBufferGrowPanic(t *testing.T) {

	defer panicChecker(t, BufferInvalidByteCountError)

	buf := NewBuffer([]byte{0x00, 0x00})

	buf.Grow(-1)

}

func TestBufferTruncateLeftPanic(t *testing.T) {

	defer panicChecker(t, BufferInvalidByteCountError)

	buf := NewBuffer([]byte{0x00, 0x00})

	buf.TruncateLeft(-1)

}

func TestBufferTruncateRightPanic(t *testing.T) {

	defer panicChecker(t, BufferInvalidByteCountError)

	buf := NewBuffer([]byte{0x00, 0x00})

	buf.TruncateRight(-1)

}

/*

benchmarks

*/

func BenchmarkBufferWriteBytes(b *testing.B) {

	b.ReportAllocs()

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {

		buf.WriteBytes(0x00, []byte{0x01, 0x02})

	}

}

func BenchmarkBufferReadBytes(b *testing.B) {

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

func BenchmarkBufferReadBit(b *testing.B) {

	b.ReportAllocs()

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {

		_ = buf.ReadBit(0x00)

	}

}

func BenchmarkBufferReadBits(b *testing.B) {

	b.ReportAllocs()

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {

		_ = buf.ReadBits(0x00, 2)

	}

}

func BenchmarkBufferSetBit(b *testing.B) {

	b.ReportAllocs()

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {

		buf.SetBit(0x00)

	}

}

func BenchmarkBufferClearBit(b *testing.B) {

	b.ReportAllocs()

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {

		buf.ClearBit(0x00)

	}

}

func BenchmarkBufferGrow(b *testing.B) {

	b.ReportAllocs()

	buf := NewBuffer()

	for n := 0; n < b.N; n++ {

		buf.Grow(1)
		buf.Reset()

	}

}
