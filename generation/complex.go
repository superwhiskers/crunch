package main

import (
	"fmt"
	"regexp"
	"strings"
	"bytes"
	"strconv"

	"github.com/dave/jennifer/jen"
)

// GenerateComplex takes an array of file names and the prefix for all of them
// it runs over all of the provided files and searches for "magic comments"
// that look like this:
//
// 	//generator:complex <receiver: Buffer | MiniBuffer> <rw: [read | write]> <sn: [signed | unsigned]> <is: [intsize (restricted to *8 atm)]> <en: [big | little]>
//
// if it finds one, it generates two functions in this pattern:
//
// 	// <naming> <reads | writes> a slice of <integer type>s to the buffer at the
// 	// specified offset in <big-endian | little-endian> without modifying the internal
// 	// offset value
// 	func (b *<receiver>) <naming>(off int64, data []<integer type>) {
//
// 		/* standard buffer method prelude omitted for brevity */
//
// 		var (
// 			i = 0
// 			n = len(data)
// 		)
// 		{
// 		<read | write>_loop:
// 			/** generate lines that look like the following for the number of bits to convert **/
// 			b.buf[off+int64(i*<number of bits / 8>)] = byte(data[i] <add >> <(line of this section - 1) * 8> for each line after line 0>)
//
// 			/** the above lines are reversed and some bits are changed around for big-endian **/
//
// 			i++
// 			if i < n {
//
// 				goto <read | write>_loop
//
// 			}
// 		}
// 	}
//
// 	// <naming>Next <reads | writes> a slice of <integer type>s to the buffer at the
// 	// current offset in <big-endian | little-endian> and moves the offset forward the
// 	// amount of bytes written
// 	func (b *<receiver>) <naming>Next(data []<integer type>) {
//
// 		b.<naming>(b.off, data)
// 		b.SeekByte(int64(len(data))*<number of bits / 8>, true)
//
// 	}
//
// that is (mostly) it. after source tweaking, it outputs each modifed file into a new file with a name
// like this:
//
// 	<filename w/o leading underscore>.generated.go
//
// after that, it goes to the next file and does the same.
func GenerateComplex(oldFiles map[string][]byte) (files map[string][]byte, e error) {
	r := regexp.MustCompile("(?m)^\\/\\/generator:complex ([A-z]{1,}) ([A-z]{1,}) ([A-z]{1,}) ([0-9]{1,}) ([A-z]{1,})$")

	files = oldFiles

	for name, _ := range files {
		fmt.Println("* scanning", name)

		files[name] = r.ReplaceAllFunc(files[name], func(m []byte) []byte {
			a := r.FindAllStringSubmatch(string(m), -1)[0][1:]
			fmt.Printf("* found match. args: %v\n", a)

			// verify the arguments
			if a[0] != "Buffer" && a[0] != "MiniBuffer" {
				fmt.Println("! invalid argument for position 0:", a[0])
				return []byte{}
			}

			if a[1] != "Read" && a[1] != "Write" {
				fmt.Println("! invalid argument for position 1:", a[1])
				return []byte{}
			}

			if a[2] != "I" && a[2] != "U" && a[2] != "F" {
				fmt.Println("! invalid argument for position 2:", a[2])
				return []byte{}
			}

			if a[3] != "16" && a[3] != "32" && a[3] != "64" {
				fmt.Println("! invalid argument for position 3:", a[3])
				return []byte{}
			}

			if a[4] != "BE" && a[4] != "LE" {
				fmt.Println("! invalid argument for position 4:", a[4])
				return []byte{}
			}

			intType := strings.Join([]string{map[string]string{"I": "int", "U": "uint", "F": "float"}[a[2]], a[3]}, "")

			intBits, e := strconv.Atoi(a[3])
			if e != nil {
				fmt.Println("! unable to convert string to integer:", e)
				return []byte{}
			}
			intBytes := intBits/8

			g := &jen.Group{}
			g.Comment(strings.Join([]string{"// ", a[1], a[2], a[3], a[4], " ", map[string]string{"Read": "reads", "Write": "writes"}[a[1]], " a slice of ", intType, "s ", map[string]string{"Read": "from", "Write": "to"}[a[1]], " the buffer at the\n"}, ""))
			g.Comment(strings.Join([]string{"// specified offset in ", map[string]string{"BE": "big-endian", "LE": "little-endian"}[a[4]], " without modifying the internal\n"}, ""))
			g.Comment("// offset value\n")
			f := g.Func().
				Params(jen.Id("b").Op("*").Id(a[0])).
				Id(strings.Join([]string{a[1], a[2], a[3], a[4]}, ""))
			if a[1] == "Write" {
				f.Params(jen.Id("off").Id("int64"), jen.Id("data").Index().Id(intType))
			} else if a[1] == "Read" {
				if a[0] == "Buffer" {
					f.Params(jen.Id("off"), jen.Id("n").Id("int64")).
						Params(jen.Id("out").Index().Id(intType))
				} else {
					f.Params(jen.Id("out").Op("*").Index().Id(intType), jen.Id("off"), jen.Id("n").Id("int64"))
				}
			}

			f.BlockFunc(func(g *jen.Group) {
				if a[0] == "Buffer" {
					if a[1] == "Read" {
						g.If(jen.Parens(jen.Id("off").Op("+").Id("n").Op("*").Lit(intBytes)).Op(">").Id("b").Dot("cap")).Block(jen.Panic(jen.Id("BufferOverreadError")))
						g.If(jen.Id("off").Op("<").Lit(0x00)).Block(jen.Panic(jen.Id("BufferUnderreadError")))
					} else {
						g.If(jen.Parens(jen.Id("off").Op("+").Id("int64").Call(jen.Len(jen.Id("data"))).Op("*").Lit(intBytes)).Op(">").Id("b").Dot("cap")).Block(jen.Panic(jen.Id("BufferOverwriteError")))
						g.If(jen.Id("off").Op("<").Lit(0x00)).Block(jen.Panic(jen.Id("BufferUnderwriteError")))
					}
				}

				if a[1] == "Read" {
					// read generation code
					if a[0] == "Buffer" {
						g.Id("out").Op("=").Id("make").Call(jen.Index().Id(intType), jen.Id("n"))
					}
					g.Id("i").Op(":=").Id("int64").Call(jen.Lit(0))
					g.BlockFunc(func(g *jen.Group) {
						g.Id("read_loop:")
						var w *jen.Statement
						if a[0] == "Buffer" {
							w = g.Id("out").Index(jen.Id("i"))
						} else {
							w = g.Parens(jen.Op("*").Id("out")).Index(jen.Id("i"))
						}
						w = w.Op("=")

						if a[4] == "BE" {
							for i := intBytes - 1; i > 0; i-- {
								w = w.Id(intType).Call(jen.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Parens(jen.Lit(i).Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))))

								if i < intBytes - 1 {
									// subsequent iterations
									w = w.Op("<<").Lit((intBytes - i - 1)*8)
								}
								w = w.Op("|")
							}
							w = w.Id(intType).Call(jen.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))).Op("<<").Lit((intBytes - 1)*8)
						} else {
							w = w.Id(intType).Call(jen.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))).Op("|")
							for i := 1; i < intBytes; i++ {
								w = w.Id(intType).Call(jen.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Parens(jen.Lit(i).Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes))))).Op("<<").Lit(i*8)

								if i < intBytes - 1 {
									// all operations except last
									w = w.Op("|")
								}
							}
						}
						g.Id("i").Op("++")
						g.If(jen.Id("i").Op("<").Id("n")).Block(jen.Goto().Id("read_loop"))
					})
				} else {
					// write generation code
					g.Id("i").Op(":=").Lit(0)
					g.Id("n").Op(":=").Len(jen.Id("data"))
					g.BlockFunc(func(g *jen.Group) {
						g.Id("write_loop:")
						if a[4] == "BE" {
							g.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Id("int64").Call(jen.Id("i").Op("*").Lit(intBytes))).Op("=").Id("byte").Call(jen.Id("data").Index(jen.Id("i")).Op(">>").Lit((intBytes - 1)*8))
							for i := intBytes - 1; i > 1; i-- {
								g.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Id("int64").Call(jen.Lit(intBytes-i).Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))).Op("=").Id("byte").Call(jen.Id("data").Index(jen.Id("i")).Op(">>").Lit((i-1)*8))
							}
							g.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Id("int64").Call(jen.Lit(intBytes-1).Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))).Op("=").Id("byte").Call(jen.Id("data").Index(jen.Id("i")))
						} else {
							g.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Id("int64").Call(jen.Id("i").Op("*").Lit(intBytes))).Op("=").Id("byte").Call(jen.Id("data").Index(jen.Id("i")))
							for i := 1; i < intBytes; i++ {
								g.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Id("int64").Call(jen.Lit(i).Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))).Op("=").Id("byte").Call(jen.Id("data").Index(jen.Id("i")).Op(">>").Lit(i*8))
							}
						}
						g.Id("i").Op("++")
						g.If(jen.Id("i").Op("<").Id("n")).Block(jen.Goto().Id("write_loop"))
					})
				}

				if a[0] == "Buffer" && a[1] == "Read" {
					g.Return()
				}
			})

			w := bytes.NewBuffer([]byte{})
			if g.Render(w) != nil {
				fmt.Println("! unable to render code:", e)
				return []byte("// render failiure")
			}

			_, _ = w.Write([]byte("\n"))

			g = &jen.Group{}
			g.Comment(strings.Join([]string{"// ", a[1], a[2], a[3], a[4], "Next ", map[string]string{"Read": "reads", "Write": "writes"}[a[1]], " a slice of ", intType, "s ", map[string]string{"Read": "from", "Write": "to"}[a[1]], " the buffer at the\n"}, ""))
			g.Comment(strings.Join([]string{"// current offset in ", map[string]string{"BE": "big-endian", "LE": "little-endian"}[a[4]], " and moves the offset forward the\n"}, ""))
			g.Comment("// amount of bytes written\n")
			f = g.Func().
				Params(jen.Id("b").Op("*").Id(a[0])).
				Id(strings.Join([]string{a[1], a[2], a[3], a[4], "Next"}, ""))
			if a[1] == "Write" {
				f.Params(jen.Id("data").Index().Id(intType))
			} else if a[1] == "Read" {
				if a[0] == "Buffer" {
					f.Params(jen.Id("n").Id("int64")).
						Params(jen.Id("out").Index().Id(intType))
				} else {
					f.Params(jen.Id("out").Op("*").Index().Id(intType), jen.Id("n").Id("int64"))
				}
			}

			f.BlockFunc(func(g *jen.Group) {
				if a[1] == "Write" {
					g.Id("b").Dot(strings.Join([]string{a[1], a[2], a[3], a[4]}, "")).Call(jen.Id("b").Dot("off"), jen.Id("data"))
					g.Id("b").Dot("SeekByte").Call(jen.Id("int64").Call(jen.Len(jen.Id("data"))).Op("*").Lit(intBytes), jen.Lit(true))
				} else if a[1] == "Read" {
					if a[0] == "Buffer" {
						g.Id("out").Op("=").Id("b").Dot(strings.Join([]string{a[1], a[2], a[3], a[4]}, "")).Call(jen.Id("b").Dot("off"), jen.Id("n"))
					} else {
						g.Id("b").Dot(strings.Join([]string{a[1], a[2], a[3], a[4]}, "")).Call(jen.Id("out"), jen.Id("b").Dot("off"), jen.Id("n"))
					}

					g.Id("b").Dot("SeekByte").Call(jen.Id("n").Op("*").Lit(intBytes), jen.Lit(true))

					if a[0] == "Buffer" {
						g.Return()
					}
				}
			})

			if g.Render(w) != nil {
				fmt.Println("! unable to render code:", e)
				return []byte("// render failiure")
			}

			return w.Bytes()
		})
	}
	return
}
