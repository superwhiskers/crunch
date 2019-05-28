package crunch

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

func TestMiniBufferReadUNEN(t *testing.T) {

	var (
		expected1 = []uint16{0x01}
		expected2 = []uint16{0x100}
		expected3 = []uint32{0x01}
		expected4 = []uint32{0x1000000}
		expected5 = []uint64{0x01}
		expected6 = []uint64{0x100000000000000}
	)

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	out1 := []uint16{0x00}

	buf.ReadU16LE(&out1, 0x00, 1)
	if !cmp.Equal(out1, expected1) {

		t.Fatalf("expected uint16 array does not match the one gotten (got %#v, expected %#v)", out1, expected1)

	}

	buf.ReadU16BE(&out1, 0x00, 1)
	if !cmp.Equal(out1, expected2) {

		t.Fatalf("expected uint16 array does not match the one gotten (got %#v, expected %#v)", out1, expected2)

	}

	out2 := []uint32{0x00}

	buf.ReadU32LE(&out2, 0x00, 1)
	if !cmp.Equal(out2, expected3) {

		t.Fatalf("expected uint32 array does not match the one gotten (got %#v, expected %#v)", out2, expected3)

	}

	buf.ReadU32BE(&out2, 0x00, 1)
	if !cmp.Equal(out2, expected4) {

		t.Fatalf("expected uint32 array does not match the one gotten (got %#v, expected %#v)", out2, expected4)

	}

	out3 := []uint64{0x00}

	buf.ReadU64LE(&out3, 0x00, 1)
	if !cmp.Equal(out3, expected5) {

		t.Fatalf("expected uint64 array does not match the one gotten (got %#v, expected %#v)", out3, expected5)

	}

	buf.ReadU64BE(&out3, 0x00, 1)
	if !cmp.Equal(out3, expected6) {

		t.Fatalf("expected uint64 array does not match the one gotten (got %#v, expected %#v)", out3, expected6)

	}

}

func TestMiniBufferReadUNENNext(t *testing.T) {

	var (
		expected1 = []uint16{0x01}
		expected2 = []uint16{0x100}
		expected3 = []uint32{0x01}
		expected4 = []uint32{0x1000000}
		expected5 = []uint64{0x01}
		expected6 = []uint64{0x100000000000000}
	)

	var off int64

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	out1 := []uint16{0x00}

	buf.ReadU16LENext(&out1, 1)
	if !cmp.Equal(out1, expected1) {

		t.Fatalf("expected uint16 array does not match the one gotten (got %#v, expected %#v)", out1, expected1)

	}

	buf.ByteOffset(&off)
	if off != 2 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	buf.ReadU16BENext(&out1, 1)
	if !cmp.Equal(out1, expected2) {

		t.Fatalf("expected uint16 array does not match the one gotten (got %#v, expected %#v)", out1, expected2)

	}

	buf.ByteOffset(&off)
	if off != 2 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	out2 := []uint32{0x00}

	buf.ReadU32LENext(&out2, 1)
	if !cmp.Equal(out2, expected3) {

		t.Fatalf("expected uint32 array does not match the one gotten (got %#v, expected %#v)", out2, expected3)

	}

	buf.ByteOffset(&off)
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	buf.ReadU32BENext(&out2, 1)
	if !cmp.Equal(out2, expected4) {

		t.Fatalf("expected uint32 array does not match the one gotten (got %#v, expected %#v)", out2, expected4)

	}

	buf.ByteOffset(&off)
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	out3 := []uint64{0x00}

	buf.ReadU64LENext(&out3, 1)
	if !cmp.Equal(out3, expected5) {

		t.Fatalf("expected uint64 array does not match the one gotten (got %#v, expected %#v)", out3, expected5)

	}

	buf.ByteOffset(&off)
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	buf.ReadU64BENext(&out3, 1)
	if !cmp.Equal(out3, expected6) {

		t.Fatalf("expected uint64 array does not match the one gotten (got %#v, expected %#v)", out3, expected6)

	}

	buf.ByteOffset(&off)
	if off != 8 {

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

func TestMiniBufferWriteUNEN(t *testing.T) {

	var (
		expected byte = 0x01

		out []byte
	)

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	buf.WriteU16LE(0x00, []uint16{0x01, 0x01})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ReadBytes(&out, 0x02, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.WriteU16BE(0x00, []uint16{0x100, 0x100})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ReadBytes(&out, 0x02, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.WriteU32LE(0x00, []uint32{0x01})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.WriteU32BE(0x00, []uint32{0x1000000})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.WriteU64LE(0x00, []uint64{0x01})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.WriteU64BE(0x00, []uint64{0x100000000000000})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

}

func TestMiniBufferWriteUNENNext(t *testing.T) {

	var (
		expected byte = 0x01

		out []byte
		off int64
	)

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	buf.WriteU16LENext([]uint16{0x01})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 2 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	buf.WriteU16BENext([]uint16{0x100})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 2 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	buf.WriteU32LENext([]uint32{0x01})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	buf.WriteU32BENext([]uint32{0x1000000})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 4 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	buf.WriteU64LENext([]uint64{0x01})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 8 {

		t.Fatalf("incorrect offset: %d", off)

	}
	buf.SeekByte(0x00, false)

	buf.WriteU64BENext([]uint64{0x100000000000000})

	buf.ReadBytes(&out, 0x00, 1)
	if expected != out[0] {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out[0], expected)

	}

	buf.ByteOffset(&off)
	if off != 8 {

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

	buf.SetBit(0x00, 1)

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

	buf.SetBitNext(1)

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

func TestMiniBufferSetbitPanic(t *testing.T) {

	defer panicChecker(t, BufferInvalidBitError)

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.SetBit(0x00, 2)

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

	out := []uint32{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	for n := 0; n < b.N; n++ {

		buf.ReadU32LE(&out, 0x00, 2)

	}

	_ = out

}
