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

	var out *Buffer

	out = NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})
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
