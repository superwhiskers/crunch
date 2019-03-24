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

	expected := &MiniBuffer{
		buf:  []byte{0x00, 0x00, 0x00, 0x00},
		off:  0x00,
		cap:  4,
		boff: 0x00,
		bcap: 32,
	}

	out := &MiniBuffer{}
	NewMiniBuffer(&out, []byte{0x00, 0x00, 0x00, 0x00})

	if !cmp.Equal(expected, out, MiniBufferComparer) {

		t.Fatalf("expected minibuffer does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

}

func TestMiniBufferBytes(t *testing.T) {

	expected := []byte{0x01, 0x02, 0x03, 0x04}

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x01, 0x02, 0x03, 0x04})

	out := []byte{}
	buf.Bytes(&out)
	if !cmp.Equal(expected, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

}

func TestMiniBufferCapacity(t *testing.T) {

	var expected int64 = 4

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	var out int64
	buf.Capacity(&out)
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

func TestMiniBufferOffset(t *testing.T) {

	var expected int64 = 0x02

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.Seek(0x02, false)

	var out int64
	buf.Offset(&out)
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

	// there is no actual way to test this because
	// the statistics are updated every time Grow is called

}

func TestMiniBufferGrow(t *testing.T) {

	var expected int64 = 4

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00})

	buf.Grow(2)

	var out int64
	buf.Capacity(&out)
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferSeek(t *testing.T) {

	var expected int64 = 0x04

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.Seek(2, true)
	buf.Seek(2, true)

	var out int64
	buf.Offset(&out)
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

	buf.Seek(0x04, false)
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
	buf.Offset(&out)
	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestMiniBufferAfter(t *testing.T) {

	var expected int64 = 2

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	var out int64

	buf.Seek(0x01, false)
	buf.After(&out)

	if expected != out {

		t.Fatalf("expected int64 does not match the one gotten (got %d, expected %d)", out, expected)

	}

	buf.After(&out, 0x01)

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

	var expected []byte = []byte{0x03, 0x04}

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x01, 0x02, 0x03, 0x04})

	var out []byte
	buf.ReadBytes(&out, 0x02, 2)
	if !cmp.Equal(expected, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

}

func TestMiniBufferReadBytesNext(t *testing.T) {

	var expected []byte = []byte{0x01, 0x02}

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x01, 0x02, 0x03, 0x04})

	var out []byte
	buf.ReadBytesNext(&out, 2)
	if !cmp.Equal(expected, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

}

func TestMiniBufferWriteBytes(t *testing.T) {

	var expected []byte = []byte{0x01, 0x02}

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

	var expected []byte = []byte{0x01, 0x02}

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	buf.WriteBytesNext([]byte{0x01, 0x02})

	var out []byte
	buf.ReadBytes(&out, 0x00, 2)
	if !cmp.Equal(expected, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

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

	var expected byte = 0

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

func BenchmarkMiniBufferWrite(b *testing.B) {

	b.ReportAllocs()

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {

		buf.WriteBytes(0x00, []byte{0x01, 0x02})

	}

}

func BenchmarkMiniBufferRead(b *testing.B) {

	b.ReportAllocs()

	buf := &MiniBuffer{}
	NewMiniBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	out := []byte{}

	for n := 0; n < b.N; n++ {

		buf.ReadBytes(&out, 0x00, 2)

	}

}
