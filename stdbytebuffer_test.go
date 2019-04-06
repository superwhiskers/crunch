package crunch

import (
	"bytes"
	"testing"
)

func BenchmarkStdByteBufferWrite(b *testing.B) {

	b.ReportAllocs()

	buf := bytes.NewBuffer([]byte{})

	var err error
	for n := 0; n < b.N; n++ {

		_, err = buf.Write([]byte{0x01, 0x02})
		if err != nil {

			b.Fatal(err)

		}
		buf.Reset() // needs to be done because we'll overwrite otherwise

	}

}

func BenchmarkStdByteBufferRead(b *testing.B) {

	b.ReportAllocs()

	buf := bytes.NewBuffer([]byte{0x00, 0x00})

	var err error
	for n := 0; n < b.N; n++ {

		_, err = buf.ReadByte()
		if err != nil {

			b.Fatal(err)

		}
		err = buf.UnreadByte()
		if err != nil {

			b.Fatal(err)

		}

	}

}
