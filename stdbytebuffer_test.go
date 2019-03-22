package crunch

import (
	"testing"
	"bytes"
)

func BenchmarkStdByteBufferWrite(b *testing.B) {

	b.ReportAllocs()

	buf := bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	for n := 0; n < b.N; n++ {
		
		buf.Write([]byte{0x01, 0x02})
		buf.Reset() // needs to be done because we'll overwrite otherwise

	}

}

func BenchmarkStdByteBufferRead(b *testing.B) {

	b.ReportAllocs()
	
	buf := bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	out := []byte{0x00, 0x00}
	
	for n := 0; n < b.N; n++ {

		buf.Read(out)
		buf.Reset()

	}

}
