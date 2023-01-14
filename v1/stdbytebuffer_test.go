/*

crunch - utilities for taking bytes out of things
Copyright (c) 2019-2020 superwhiskers <whiskerdev@protonmail.com>

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.

*/

package v1

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
