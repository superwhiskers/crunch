package crunch

import "testing"

/*

tests

*/

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
