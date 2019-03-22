package crunch

import "testing"

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
