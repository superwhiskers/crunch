package crunch

import "testing"

func BenchmarkUnsafeBufferWrite(b *testing.B) {

	b.ReportAllocs()

	buf := &UnsafeBuffer{}
	NewUnsafeBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {
		
		buf.WriteBytes(0x00, []byte{0x01, 0x02})

	}

}

func BenchmarkUnsafeBufferRead(b *testing.B) {

	b.ReportAllocs()
	
	buf := &UnsafeBuffer{}
	NewUnsafeBuffer(&buf, []byte{0x00, 0x00, 0x00, 0x00})

	out := []byte{0x00, 0x00}
	
	for n := 0; n < b.N; n++ {

		buf.ReadBytes(&out, 0x00, 2)

	}

}
