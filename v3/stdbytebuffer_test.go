/*

crunch - utilities for taking bytes out of things
copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>

this source code form is subject to the terms of the mozilla public
license, v. 2.0. if a copy of the mpl was not distributed with this
file, you can obtain one at http://mozilla.org/MPL/2.0/.

*/

package v3

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

func BenchmarkStdByteBufferGrow(b *testing.B) {

	b.ReportAllocs()

	buf := &bytes.Buffer{}

	for n := 0; n < b.N; n++ {

		buf.Grow(1)
		buf.Reset()

	}

}
