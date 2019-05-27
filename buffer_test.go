package crunch

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

/*
func TestBufferReadComplex(t *testing.T) {

	var (
		expected1 = []byte{0x01}
		expected2 = []uint16{0x01}
		expected3 = []uint16{0x100}
		expected4 = []uint32{0x01}
		expected5 = []uint64{0x01}
	)

	buf := NewBuffer([]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	out1 := buf.ReadComplex(0x00, 1, Unsigned8, LittleEndian).([]byte)
	if !cmp.Equal(out1, expected1) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out1, expected1)

	}

	out2 := buf.ReadComplex(0x00, 1, Unsigned16, LittleEndian).([]uint16)
	if !cmp.Equal(out2, expected2) {

		t.Fatalf("expected uint16 array does not match the one gotten (got %#v, expected %#v)", out2, expected2)

	}

	out2 = buf.ReadComplex(0x00, 1, Unsigned16, BigEndian).([]uint16)
	if !cmp.Equal(out2, expected3) {

		t.Fatalf("expected uint16 array does not match the one gotten (got %#v, expected %#v)", out2, expected3)

	}

	out3 := buf.ReadComplex(0x00, 1, Unsigned32, LittleEndian).([]uint32)
	if !cmp.Equal(out3, expected4) {

		t.Fatalf("expected uint32 array does not match the one gotten (got %#v, expected %#v)", out3, expected4)

	}

	out4 := buf.ReadComplex(0x00, 1, Unsigned64, LittleEndian).([]uint64)
	if !cmp.Equal(out4, expected5) {

		t.Fatalf("expected uint64 array does not match the one gotten (got %#v, expected %#v)", out4, expected5)

	}

}

func TestBufferReadComplexNext(t *testing.T) {

	var expected = []byte{0x01, 0x02}

	buf := NewBuffer([]byte{0x01, 0x02, 0x03, 0x04})

	out := buf.ReadComplexNext(2, Unsigned8, LittleEndian).([]byte)

	if !cmp.Equal(expected, out) {

		t.Fatalf("expected byte array does not match the one gotten (got %#v, expected %#v)", out, expected)

	}

}*/

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

/*
func TestBufferWriteComplex(t *testing.T) {

	var expected byte = 0x01

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	buf.WriteComplex(0x00, []byte{0x01}, Unsigned8, LittleEndian)

	out := buf.ReadByte(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

	buf.WriteComplex(0x00, []uint16{0x01}, Unsigned16, LittleEndian)

	out = buf.ReadByte(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

	buf.WriteComplex(0x00, []uint16{0x100}, Unsigned16, BigEndian)

	out = buf.ReadByte(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

	buf.WriteComplex(0x00, []uint32{0x01}, Unsigned32, LittleEndian)

	out = buf.ReadByte(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

	buf.WriteComplex(0x00, []uint64{0x01}, Unsigned64, LittleEndian)

	out = buf.ReadByte(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}

func TestBufferWriteComplexNext(t *testing.T) {

	var expected byte = 0x01

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	buf.WriteComplexNext([]byte{0x01, 0x02}, Unsigned8, LittleEndian)

	out := buf.ReadByte(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

	buf.Seek(0x00, false)
	buf.WriteComplexNext([]uint16{0x01}, Unsigned16, LittleEndian)

	out = buf.ReadByte(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

	buf.Seek(0x00, false)
	buf.WriteComplexNext([]uint32{0x01}, Unsigned32, LittleEndian)

	out = buf.ReadByte(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

	buf.Seek(0x00, false)
	buf.WriteComplexNext([]uint64{0x01}, Unsigned64, LittleEndian)

	out = buf.ReadByte(0x00)
	if expected != out {

		t.Fatalf("expected byte does not match the one gotten (got %d, expected %d)", out, expected)

	}

}*/

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

/*func TestBufferReadComplexPanic(t *testing.T) {

	defer panicChecker(t, BufferInvalidIntegerSizeError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	_ = buf.readComplex(0x00, 1, -1, LittleEndian).([]byte)

}

func TestBufferWriteComplexPanic(t *testing.T) {

	defer panicChecker(t, BufferInvalidIntegerSizeError)

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	buf.writeComplex(0x00, 1, -1, LittleEndian)

}*/

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

/*
func BenchmarkBufferWriteComplex(b *testing.B) {

	b.ReportAllocs()

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {

		buf.WriteComplex(0x00, []uint32{0x01, 0x02}, Unsigned32, LittleEndian)

	}

}

func BenchmarkBufferReadComplex(b *testing.B) {

	b.ReportAllocs()

	buf := NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	var out interface{}
	for n := 0; n < b.N; n++ {

		out = buf.ReadComplex(0x00, 2, Unsigned32, LittleEndian)

	}

	_ = out

}*/
