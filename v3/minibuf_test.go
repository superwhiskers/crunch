/*

crunch - utilities for taking bytes out of things
Copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.

*/

package v3

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

/*

utilities

*/

var MiniBufferComparer = cmp.Comparer(func(x, y *MiniBuffer) bool {

	return cmp.Equal(x.buf, y.buf) &&
		x.off == y.off &&
		x.cap == y.cap &&
		x.boff == y.boff &&
		x.bcap == y.bcap

})

/*

tests

*/

func TestNewMiniBuffer(t *testing.T) {

	var (
		expected1 = &MiniBuffer{
			buf:  []byte{0x00, 0x00, 0x00, 0x00},
			off:  0x00,
			cap:  4,
			boff: 0x00,
			bcap: 32,
		}
		expected2 = &MiniBuffer{
			buf:  []byte{},
			off:  0x00,
			cap:  0,
			boff: 0x00,
			bcap: 0,
		}
		expected3 = &MiniBuffer{
			buf:  []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			off:  0x00,
			cap:  8,
			boff: 0x00,
			bcap: 64,
		}
	)

	out := &MiniBuffer{}

	NewMiniBuffer(&out, []byte{0x00, 0x00, 0x00, 0x00})
	if !cmp.Equal(expected1, out, MiniBufferComparer) {

		t.Fatalf("expected minibuffer does not match the one gotten (got %#v, expected %#v)", out, expected1)

	}

	NewMiniBuffer(&out)
	if !cmp.Equal(expected2, out, MiniBufferComparer) {

		t.Fatalf("expected minibuffer does not match the one gotten (got %#v, expected %#v)", out, expected2)

	}

	NewMiniBuffer(&out, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00, 0x00, 0x00})
	if !cmp.Equal(expected3, out, MiniBufferComparer) {

		t.Fatalf("expected minibuffer does not match the one gotten (got %#v, expected %#v)", out, expected3)

	}

}

func TestMiniBufferBytes(t *testing.T) {

	var expected = []byte{0x01, 0x02, 0x03, 0x04}

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x01, 0x02, 0x03, 0x04})

	out := []byte{}
	buf.Bytes(&out)
	if !cmp.Equal(expected, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

}

func TestMiniBufferByteCapacity(t *testing.T) {

	var expected int64 = 4

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	var out int64
	buf.ByteCapacity(&out)
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferBitCapacity(t *testing.T) {

	var expected int64 = 32

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	var out int64
	buf.BitCapacity(&out)
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferByteOffset(t *testing.T) {

	var expected int64 = 0x02

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.SeekByte(0x02, false)

	var out int64
	buf.ByteOffset(&out)
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferBitOffset(t *testing.T) {

	var expected int64 = 0x0f

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.SeekBit(0x0f, false)

	var out int64
	buf.BitOffset(&out)
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferRefresh(t *testing.T) {

	var (
		expected1 int64 = 4
		expected2 int64 = 32
	)

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.Refresh()

	if expected1 != buf.cap || expected2 != buf.bcap {

		t.Fatalf("expected int64(s) do not match the ones gotten (got %d and %d, expected %d and %d)", buf.cap, buf.bcap, expected1, expected2)

	}

}

func TestMiniBufferReset(t *testing.T) {

	var expected = []byte{}

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.Reset()

	if !cmp.Equal(expected, buf.buf) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", buf.buf, expected)

	}

}

func TestMiniBufferTruncateLeft(t *testing.T) {

	var (
		expected1 byte  = 0x02
		expected2 int64 = 1
	)

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x01, 0x02})

	buf.TruncateLeft(1)

	var out1 []byte
	buf.ReadBytes(&out1, 0x00, 1)
	if expected1 != out1[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out1, expected1)

	}

	var out2 int64
	buf.ByteCapacity(&out2)
	if expected2 != out2 {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out2, expected2)

	}

}

func TestMiniBufferTruncateRight(t *testing.T) {

	var (
		expected1 byte  = 0x01
		expected2 int64 = 1
	)

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x01, 0x02})

	buf.TruncateRight(1)

	var out1 []byte
	buf.ReadBytes(&out1, 0x00, 1)
	if expected1 != out1[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out1, expected1)

	}

	var out2 int64
	buf.ByteCapacity(&out2)
	if expected2 != out2 {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out2, expected2)

	}

}

func TestMiniBufferGrow(t *testing.T) {

	var expected int64 = 4

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00})

	buf.Grow(2)

	var out int64
	buf.ByteCapacity(&out)
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferGrow2(t *testing.T) {

	var expected int64 = 4

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, make([]byte, 2, 4))

	buf.Grow(2)

	var out int64
	buf.ByteCapacity(&out)
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferSeekByte(t *testing.T) {

	var expected int64 = 0x04

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.SeekByte(2, true)
	buf.SeekByte(2, true)

	var out int64
	buf.ByteOffset(&out)
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferSeekBit(t *testing.T) {

	var expected int64 = 0x20

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.SeekBit(16, true)
	buf.SeekBit(16, true)

	var out int64
	buf.BitOffset(&out)
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferAlignBit(t *testing.T) {

	var expected int64 = 0x20

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.SeekByte(0x04, false)
	buf.AlignBit()

	var out int64
	buf.BitOffset(&out)
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferAlignByte(t *testing.T) {

	var expected int64 = 0x04

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.SeekBit(0x20, false)
	buf.AlignByte()

	var out int64
	buf.ByteOffset(&out)
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferAfterByte(t *testing.T) {

	var expected int64 = 2

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	var out int64

	buf.SeekByte(0x01, false)
	buf.AfterByte(&out)

	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

	buf.AfterByte(&out, 0x01)

	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferAfterBit(t *testing.T) {

	var expected int64 = 16

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	var out int64

	buf.SeekBit(0x0f, false)
	buf.AfterBit(&out)

	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

	buf.AfterBit(&out, 0x0f)

	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferReadBytes(t *testing.T) {

	var expected = []byte{0x03, 0x04}

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x01, 0x02, 0x03, 0x04})

	var out []byte
	buf.ReadBytes(&out, 0x02, 2)
	if !cmp.Equal(expected, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

}

func TestMiniBufferReadBytesNext(t *testing.T) {

	var expected = []byte{0x01, 0x02}

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x01, 0x02, 0x03, 0x04})

	var out []byte
	buf.ReadBytesNext(&out, 2)
	if !cmp.Equal(expected, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

}

//gocyclo:ignore
func TestMiniBufferReadSNEN(t *testing.T) {

	var (
		expectedReadu16le = []uint16{0x01, 0x00}
		expectedReadu16be = []uint16{0x100, 0x00}
		expectedReadu32le = []uint32{0x01, 0x00}
		expectedReadu32be = []uint32{0x1000000, 0x00}
		expectedReadu64le = []uint64{0x01, 0x01}
		expectedReadu64be = []uint64{0x100000000000000, 0x100000000000000}

		expectedReadi16le = []int16{0x01, 0x00}
		expectedReadi16be = []int16{0x100, 0x00}
		expectedReadi32le = []int32{0x01, 0x00}
		expectedReadi32be = []int32{0x1000000, 0x00}
		expectedReadi64le = []int64{0x01, 0x01}
		expectedReadi64be = []int64{0x100000000000000, 0x100000000000000}

		off = int64(0)

		out1 = []uint16{0x00, 0x00}
		out2 = []uint32{0x00, 0x00}
		out3 = []uint64{0x00, 0x00}

		out4 = []int16{0x00, 0x00}
		out5 = []int32{0x00, 0x00}
		out6 = []int64{0x00, 0x00}

		out7 = []float32{0x00, 0x00}
		out8 = []float64{0x00, 0x00}
	)

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	//
	// u16le
	//

	buf.ReadU16LENext(&out1, 2)
	if !cmp.Equal(out1, expectedReadu16le) {

		t.Fatalf("expected uint16 array does not match the one gotten (got %#v, expected %#v)", out1, expectedReadu16le)

	}

	buf.ByteOffset(&off)
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u16be
	//

	buf.ReadU16BENext(&out1, 2)
	if !cmp.Equal(out1, expectedReadu16be) {

		t.Fatalf("expected uint16 array does not match the one gotten (got %#v, expected %#v)", out1, expectedReadu16be)

	}

	buf.ByteOffset(&off)
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u32le
	//

	buf.ReadU32LENext(&out2, 2)
	if !cmp.Equal(out2, expectedReadu32le) {

		t.Fatalf("expected uint32 array does not match the one gotten (got %#v, expected %#v)", out2, expectedReadu32le)

	}

	buf.ByteOffset(&off)
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u32be
	//

	buf.ReadU32BENext(&out2, 2)
	if !cmp.Equal(out2, expectedReadu32be) {

		t.Fatalf("expected uint32 array does not match the one gotten (got %#v, expected %#v)", out2, expectedReadu32be)

	}

	buf.ByteOffset(&off)
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u64le
	//

	buf.ReadU64LENext(&out3, 2)
	if !cmp.Equal(out3, expectedReadu64le) {

		t.Fatalf("expected uint64 array does not match the one gotten (got %#v, expected %#v)", out3, expectedReadu64le)

	}

	buf.ByteOffset(&off)
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u64be
	//

	buf.ReadU64BENext(&out3, 2)
	if !cmp.Equal(out3, expectedReadu64be) {

		t.Fatalf("expected uint64 array does not match the one gotten (got %#v, expected %#v)", out3, expectedReadu64be)

	}

	buf.ByteOffset(&off)
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i16le
	//

	buf.ReadI16LENext(&out4, 2)
	if !cmp.Equal(out4, expectedReadi16le) {

		t.Fatalf("expected int16 array does not match the one gotten (got %#v, expected %#v)", out4, expectedReadi16le)

	}

	buf.ByteOffset(&off)
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i16be
	//

	buf.ReadI16BENext(&out4, 2)
	if !cmp.Equal(out4, expectedReadi16be) {

		t.Fatalf("expected int16 array does not match the one gotten (got %#v, expected %#v)", out4, expectedReadi16be)

	}

	buf.ByteOffset(&off)
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i32le
	//

	buf.ReadI32LENext(&out5, 2)
	if !cmp.Equal(out5, expectedReadi32le) {

		t.Fatalf("expected int32 array does not match the one gotten (got %#v, expected %#v)", out5, expectedReadi32le)

	}

	buf.ByteOffset(&off)
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i32be
	//

	buf.ReadI32BENext(&out5, 2)
	if !cmp.Equal(out5, expectedReadi32be) {

		t.Fatalf("expected int32 array does not match the one gotten (got %#v, expected %#v)", out5, expectedReadi32be)

	}

	buf.ByteOffset(&off)
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i64le
	//

	buf.ReadI64LENext(&out6, 2)
	if !cmp.Equal(out6, expectedReadi64le) {

		t.Fatalf("expected int64 array does not match the one gotten (got %#v, expected %#v)", out6, expectedReadi64le)

	}

	buf.ByteOffset(&off)
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i64be
	//

	buf.ReadI64BENext(&out6, 2)
	if !cmp.Equal(out6, expectedReadi64be) {

		t.Fatalf("expected int64 array does not match the one gotten (got %#v, expected %#v)", out6, expectedReadi64be)

	}

	buf.ByteOffset(&off)
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// floats (we use a slightly different buffer to have higher confidence here)
	//

	var (
		expectedReadf32le = []float32{-45.318115, -8.536945e+32}
		expectedReadf32be = []float32{-3.081406, -1.0854726e-29}
		expectedReadf64le = []float64{-1.498274907009045e+261, 5.278146837874356e-99}
		expectedReadf64be = []float64{-42.42, 3.621}
	)

	NewMiniBuffer(&buf, []byte{0xc0, 0x45, 0x35, 0xc2, 0x8f, 0x5c, 0x28, 0xf6, 0x40, 0xc, 0xf7, 0xce, 0xd9, 0x16, 0x87, 0x2b})

	//
	// f32le
	//

	buf.ReadF32LENext(&out7, 2)
	if !cmp.Equal(out7, expectedReadf32le) {

		t.Fatalf("expected float32 array does not match the one gotten (got %#v, expected %#v)", out7, expectedReadf32le)

	}
	buf.ByteOffset(&off)
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// f32be
	//

	buf.ReadF32BENext(&out7, 2)
	if !cmp.Equal(out7, expectedReadf32be) {

		t.Fatalf("expected float32 array does not match the one gotten (got %#v, expected %#v)", out7, expectedReadf32be)

	}
	buf.ByteOffset(&off)
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// f64le
	//

	buf.ReadF64LENext(&out8, 2)
	if !cmp.Equal(out8, expectedReadf64le) {

		t.Fatalf("expected float64 array does not match the one gotten (got %#v, expected %#v)", out8, expectedReadf64le)

	}
	buf.ByteOffset(&off)
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// f64be
	//

	buf.ReadF64BENext(&out8, 2)
	if !cmp.Equal(out8, expectedReadf64be) {

		t.Fatalf("expected float64 array does not match the one gotten (got %#v, expected %#v)", out8, expectedReadf64be)

	}
	buf.ByteOffset(&off)
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

}

func TestMiniBufferWriteBytes(t *testing.T) {

	var expected = []byte{0x01, 0x02}

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteBytes(0x02, []byte{0x01, 0x02})

	var out []byte
	buf.ReadBytes(&out, 0x02, 2)
	if !cmp.Equal(expected, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

}

func TestMiniBufferWriteBytesNext(t *testing.T) {

	var expected = []byte{0x01, 0x02}

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteBytesNext([]byte{0x01, 0x02})

	var out []byte
	buf.ReadBytes(&out, 0x00, 2)
	if !cmp.Equal(expected, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

}

//gocyclo:ignore
func TestMiniBufferWriteSNEN(t *testing.T) {

	var (
		out           = []byte{0x00}
		off           = int64(0)
		expected byte = 0x01
	)

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	//
	// u16le
	//

	buf.WriteU16LENext([]uint16{0x01, 0x01})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ReadBytes(&out, 0x02, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u16be
	//

	buf.WriteU16BENext([]uint16{0x100, 0x100})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ReadBytes(&out, 0x02, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u32le
	//

	buf.WriteU32LENext([]uint32{0x01, 0x01})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ReadBytes(&out, 0x04, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u32be
	//

	buf.WriteU32BENext([]uint32{0x1000000, 0x1000000})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ReadBytes(&out, 0x04, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u64le
	//

	buf.WriteU64LENext([]uint64{0x01, 0x01})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ReadBytes(&out, 0x08, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// u64be
	//

	buf.WriteU64BENext([]uint64{0x100000000000000, 0x100000000000000})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ReadBytes(&out, 0x08, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i16le
	//

	buf.WriteI16LENext([]int16{0x01, 0x01})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ReadBytes(&out, 0x02, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i16be
	//

	buf.WriteI16BENext([]int16{0x100, 0x100})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ReadBytes(&out, 0x02, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i32le
	//

	buf.WriteI32LENext([]int32{0x01, 0x01})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ReadBytes(&out, 0x04, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i32be
	//

	buf.WriteI32BENext([]int32{0x1000000, 0x1000000})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ReadBytes(&out, 0x04, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i64le
	//

	buf.WriteI64LENext([]int64{0x01, 0x01})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ReadBytes(&out, 0x08, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// i64be
	//

	buf.WriteI64BENext([]int64{0x100000000000000, 0x100000000000000})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ReadBytes(&out, 0x08, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// floats (we use a slightly different expected integer array to have higher confidence here)
	//

	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	expected2 := []byte{0xc0, 0x45, 0x35, 0xc2, 0x8f, 0x5c, 0x28, 0xf6, 0x40, 0xc, 0xf7, 0xce, 0xd9, 0x16, 0x87, 0x2b}

	//
	// f32le
	//

	buf.WriteF32LENext([]float32{-45.318115, -8.536945e+32})

	buf.ReadBytes(&out, 0x00, 8)
	if !cmp.Equal(expected2[0:8], out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected2)

	}

	buf.ByteOffset(&off)
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// f32be
	//

	buf.WriteF32BENext([]float32{-3.081406, -1.0854726e-29})

	buf.ReadBytes(&out, 0x00, 8)
	if !cmp.Equal(expected2[0:8], out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected2)

	}

	buf.ByteOffset(&off)
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// f64le
	//

	buf.WriteF64LENext([]float64{-1.498274907009045e+261, 5.278146837874356e-99})

	buf.ReadBytes(&out, 0x00, 16)
	if !cmp.Equal(expected2, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected2)

	}

	buf.ByteOffset(&off)
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	//
	// f64be
	//

	buf.WriteF64BENext([]float64{-42.42, 3.621})

	buf.ReadBytes(&out, 0x00, 16)
	if !cmp.Equal(expected2, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected2)

	}

	buf.ByteOffset(&off)
	if off != 16 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)
}

func TestMiniBufferReadBit(t *testing.T) {

	var expected byte = 1

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x01, 0x00, 0x00, 0x00})

	var out byte
	buf.ReadBit(&out, 0x07)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferReadBitNext(t *testing.T) {

	var expected byte = 1

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0xff, 0x00, 0x00, 0x00})

	var out byte
	buf.ReadBitNext(&out)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferReadBits(t *testing.T) {

	var expected uint64 = 5

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x0d, 0x00, 0x00, 0x00})

	var out uint64
	buf.ReadBits(&out, 0x05, 3)
	if expected != out {

		t.Fatalf("expected uint64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferReadBitsNext(t *testing.T) {

	var expected uint64 = 13

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x0d, 0x00, 0x00, 0x00})

	var out uint64
	buf.ReadBitsNext(&out, 8)
	if expected != out {

		t.Fatalf("expected uint64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferSetBit(t *testing.T) {

	var expected byte = 1

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.SetBit(0x00)

	var out byte
	buf.ReadBit(&out, 0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferSetBitNext(t *testing.T) {

	var expected byte = 1

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.SetBitNext()

	var out byte
	buf.ReadBit(&out, 0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferClearBit(t *testing.T) {

	var expected byte

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.ClearBit(0x00)

	var out byte
	buf.ReadBit(&out, 0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferClearBitNext(t *testing.T) {

	var expected byte

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.ClearBitNext()

	var out byte
	buf.ReadBit(&out, 0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferSetBits(t *testing.T) {

	var expected uint64 = 5

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.SetBits(0x06, 5, 3)

	var out uint64
	buf.ReadBits(&out, 0x06, 3)
	if expected != out {

		t.Fatalf("expected uint64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferSetBitsNext(t *testing.T) {

	var expected uint64 = 13

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.SetBitsNext(13, 8)

	var out uint64
	buf.ReadBits(&out, 0x00, 8)
	if expected != out {

		t.Fatalf("expected uint64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferFlipBit(t *testing.T) {

	var expected byte

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0xff, 0xff, 0xff, 0xff})

	buf.FlipBit(0x01)

	var out byte
	buf.ReadBit(&out, 0x01)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferFlipBitNext(t *testing.T) {

	var expected byte = 1

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.FlipBitNext()

	var out byte
	buf.ReadBit(&out, 0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferClearAllBits(t *testing.T) {

	var expected = []byte{0x00, 0x00, 0x00, 0x00}

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0xff, 0x00, 0xff, 0x00})

	buf.ClearAllBits()

	if !cmp.Equal(expected, buf.buf) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", buf.buf, expected)

	}

}

func TestMiniBufferSetAllBits(t *testing.T) {

	var expected = []byte{0xff, 0xff, 0xff, 0xff}

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0xff, 0x00, 0xff})

	buf.SetAllBits()

	if !cmp.Equal(expected, buf.buf) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", buf.buf, expected)

	}

}

func TestMiniBufferFlipAllBits(t *testing.T) {

	var expected = []byte{0xff, 0x00, 0xff, 0x00}

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0xff, 0x00, 0xff})

	buf.FlipAllBits()

	if !cmp.Equal(expected, buf.buf) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", buf.buf, expected)

	}

}

/*

benchmarks

*/

func BenchmarkMiniBufferWriteBytes(b *testing.B) {

	b.ReportAllocs()

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {

		buf.WriteBytes(0x00, []byte{0x01, 0x02})

	}

}

func BenchmarkMiniBufferReadBytes(b *testing.B) {

	b.ReportAllocs()

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	out := []byte{}

	for n := 0; n < b.N; n++ {

		buf.ReadBytes(&out, 0x00, 2)

	}

}

func BenchmarkMiniBufferWriteU32LE(b *testing.B) {

	b.ReportAllocs()

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {

		buf.WriteU32LE(0x00, []uint32{0x01, 0x02})

	}

}

func BenchmarkMiniBufferReadU32LE(b *testing.B) {

	b.ReportAllocs()

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	out := []uint32{0x00, 0x00}
	for n := 0; n < b.N; n++ {

		buf.ReadU32LE(&out, 0x00, 2)

	}

	_ = out

}

func BenchmarkMiniBufferWriteF32LE(b *testing.B) {

	b.ReportAllocs()

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {

		buf.WriteF32LE(0x00, []float32{0.01, 0.02})

	}

}

func BenchmarkMiniBufferReadF32LE(b *testing.B) {

	b.ReportAllocs()

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	out := []float32{0x00, 0x00}
	for n := 0; n < b.N; n++ {

		buf.ReadF32LE(&out, 0x00, 2)

	}

	_ = out

}

func BenchmarkMiniBufferReadBit(b *testing.B) {

	b.ReportAllocs()

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	var out byte
	for n := 0; n < b.N; n++ {

		buf.ReadBit(&out, 0x00)

	}

	_ = out

}

func BenchmarkMiniBufferReadBits(b *testing.B) {

	b.ReportAllocs()

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	var out uint64
	for n := 0; n < b.N; n++ {

		buf.ReadBits(&out, 0x00, 2)

	}

	_ = out

}

func BenchmarkMiniBufferSetBit(b *testing.B) {

	b.ReportAllocs()

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {

		buf.SetBit(0x00)

	}

}

func BenchmarkMiniBufferClearBit(b *testing.B) {

	b.ReportAllocs()

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {

		buf.ClearBit(0x00)

	}

}

func BenchmarkMiniBufferGrow(b *testing.B) {

	b.ReportAllocs()

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf)

	for n := 0; n < b.N; n++ {

		buf.Grow(1)
		buf.Reset()

	}

}
