package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
//
//gocyclo:ignore
func GenerateComplex(oldFiles map[string][]byte) (files map[string][]byte, e error) {
	magicCommentRegex := regexp.MustCompile("(?m)^\\/\\/generator:complex ([A-z]{1,}) ([A-z]{1,}) ([A-z]{1,}) ([0-9]{1,}) ([A-z]{1,})$")

	files = oldFiles

	for name := range files {
		fmt.Println("* scanning", name)

		files[name] = magicCommentRegex.ReplaceAllFunc(files[name], func(comment []byte) []byte {
			arguments := magicCommentRegex.FindAllStringSubmatch(string(comment), -1)[0][1:]
			fmt.Printf("* invocation. args: %v\n", arguments)

			/* argument verification */

			if arguments[0] != "Buffer" && arguments[0] != "MiniBuffer" {
				fmt.Println("! invalid argument for position 0:", arguments[0])
				return []byte(fmt.Sprint("// invalid argument provided in position zero:", arguments[0]))
			}

			if arguments[1] != "Read" && arguments[1] != "Write" {
				fmt.Println("! invalid argument for position 1:", arguments[1])
				return []byte(fmt.Sprint("// invalid argument provided in position one:", arguments[1]))
			}

			if arguments[2] != "I" && arguments[2] != "U" && arguments[2] != "F" {
				fmt.Println("! invalid argument for position 2:", arguments[2])
				return []byte(fmt.Sprint("// invalid argument provided in position two:", arguments[2]))
			}

			if arguments[3] != "16" && arguments[3] != "32" && arguments[3] != "64" {
				fmt.Println("! invalid argument for position 3:", arguments[3])
				return []byte(fmt.Sprint("// invalid argument provided in position three:", arguments[3]))
			}

			if arguments[4] != "BE" && arguments[4] != "LE" {
				fmt.Println("! invalid argument for position 4:", arguments[4])
				return []byte(fmt.Sprint("// invalid argument provided in position four:", arguments[4]))
			}

			/* convenience definitions */

			intType := strings.Join(
				[]string{
					map[string]string{
						"I": "int",
						"U": "uint",
						"F": "float",
					}[arguments[2]],
					arguments[3],
				},
				"",
			)

			intBits, err := strconv.Atoi(arguments[3])
			if err != nil {
				fmt.Println("! unable to convert string to integer:", err)
				return []byte(fmt.Sprint("// conversion error:", err))
			}
			intBytes := intBits / 8
			functionName := strings.Join([]string{arguments[1], arguments[2], arguments[3], arguments[4]}, "")
			functionNameNext := strings.Join([]string{functionName, "Next"}, "")

			/* code generation */

			builder := &jen.Group{}
			builder.Comment(strings.Join([]string{
				"// ",
				functionName,
				" ",
				map[string]string{
					"Read":  "reads",
					"Write": "writes",
				}[arguments[1]],
				" a slice of ",
				intType,
				"s ",
				map[string]string{
					"Read":  "from",
					"Write": "to",
				}[arguments[1]],
				" the buffer at the\n",
			}, ""))
			builder.Comment(strings.Join([]string{
				"// specified offset in ",
				map[string]string{
					"BE": "big-endian",
					"LE": "little-endian",
				}[arguments[4]],
				" without modifying the internal\n",
			}, ""))
			builder.Comment("// offset value\n")

			function := builder.Func().Params(jen.Id("b").Op("*").Id(arguments[0])).Id(functionName)
			if arguments[1] == "Write" {
				function.Params(
					jen.Id("off").Id("int64"),
					jen.Id("data").Index().Id(intType))
			} else if arguments[1] == "Read" {
				if arguments[0] == "Buffer" {
					function.Params(
						jen.Id("off"),
						jen.Id("n").Id("int64")).Params(jen.Id("out").Index().Id(intType))
				} else {
					function.Params(
						jen.Id("out").Op("*").Index().Id(intType),
						jen.Id("off"),
						jen.Id("n").Id("int64"))
				}
			}

			function.BlockFunc(func(body *jen.Group) {
				if arguments[0] == "Buffer" {
					if arguments[1] == "Read" {
						body.If(jen.Parens(jen.Id("off").Op("+").Id("n").Op("*").Lit(intBytes)).Op(">").Id("b").Dot("cap")).
							Block(jen.Panic(jen.Id("BufferOverreadError")))
						body.If(jen.Id("off").Op("<").Lit(0x00)).
							Block(jen.Panic(jen.Id("BufferUnderreadError")))
					} else {
						body.If(jen.Parens(jen.Id("off").Op("+").Id("int64").
							Call(jen.Len(jen.Id("data"))).Op("*").Lit(intBytes)).Op(">").Id("b").Dot("cap")).
							Block(jen.Panic(jen.Id("BufferOverwriteError")))
						body.If(jen.Id("off").Op("<").Lit(0x00)).
							Block(jen.Panic(jen.Id("BufferUnderwriteError")))
					}
				}

				if arguments[1] == "Read" {
					// read generation code
					if arguments[0] == "Buffer" {
						body.Id("out").Op("=").Id("make").
							Call(
								jen.Index().Id(intType),
								jen.Id("n"))
					}
					body.Id("i").Op(":=").Id("int64").
						Call(jen.Lit(0))

					if arguments[2] == "F" {
						body.Var().Id("u").Id(strings.Join([]string{"uint", arguments[3]}, ""))
					}

					body.BlockFunc(func(loop *jen.Group) {
						loop.Id("read_loop:")

						if arguments[2] == "F" {
							// this is necessary for precedence
							v := loop.Id("u").Op("=")
							if arguments[4] == "BE" {
								v.CallFunc(func(g *jen.Group) {
									s := g.Null()
									for i := intBytes - 1; i > 0; i-- {
										s = s.Id(strings.Join([]string{"uint", arguments[3]}, "")).
											Call(jen.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Parens(jen.Lit(i).Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))))

										if i < intBytes-1 {
											// subsequent iterations
											s = s.Op("<<").Lit((intBytes - i - 1) * 8)
										}
										s = s.Op("|")
									}
									s = s.Id(strings.Join([]string{"uint", arguments[3]}, "")).
										Call(jen.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))).Op("<<").Lit((intBytes - 1) * 8)
								})
							} else if arguments[4] == "LE" {
								v.CallFunc(func(g *jen.Group) {
									s := g.Id(strings.Join([]string{"uint", arguments[3]}, "")).
										Call(jen.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))).Op("|")
									for i := 1; i < intBytes; i++ {
										s = s.Id(strings.Join([]string{"uint", arguments[3]}, "")).
											Call(jen.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Parens(jen.Lit(i).Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes))))).Op("<<").Lit(i * 8)

										if i < intBytes-1 {
											// all operations except last
											s = s.Op("|")
										}
									}
								})
							}
						}

						var orChain *jen.Statement
						if arguments[0] == "Buffer" {
							orChain = loop.Id("out").Index(jen.Id("i"))
						} else {
							orChain = loop.Parens(jen.Op("*").Id("out")).Index(jen.Id("i"))
						}
						orChain = orChain.Op("=")

						if arguments[4] == "BE" {
							// double check that we don't need to special-case for the float* types
							if arguments[2] == "U" || arguments[2] == "I" {
								for i := intBytes - 1; i > 0; i-- {
									orChain = orChain.Id(intType).
										Call(jen.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Parens(jen.Lit(i).Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))))

									if i < intBytes-1 {
										// subsequent iterations
										orChain = orChain.Op("<<").Lit((intBytes - i - 1) * 8)
									}
									orChain = orChain.Op("|")
								}
								orChain = orChain.Id(intType).
									Call(jen.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))).Op("<<").Lit((intBytes - 1) * 8)
							} else if arguments[2] == "F" {
								// this doesn't really represent the or chain, but this is necessary so i'm not having more special cases
								orChain = orChain.Op("*").Parens(jen.Op("*").Id(intType)).Parens(jen.Id("unsafe").Dot("Pointer").
									Call(jen.Op("&").Id("u")))
							}
						} else {
							// check if we're reading a float* type
							if arguments[2] == "U" || arguments[2] == "I" {
								orChain = orChain.Id(intType).
									Call(jen.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))).Op("|")
								for i := 1; i < intBytes; i++ {
									orChain = orChain.Id(intType).
										Call(jen.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Parens(jen.Lit(i).Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes))))).Op("<<").Lit(i * 8)
									if i < intBytes-1 {
										// all operations except last
										orChain = orChain.Op("|")
									}
								}
							} else if arguments[2] == "F" {
								// again, this doesn't really represent the or chain
								orChain = orChain.Op("*").Parens(jen.Op("*").Id(intType)).Parens(jen.Id("unsafe").Dot("Pointer").
									Call(jen.Op("&").Id("u")))
							}
						}
						loop.Id("i").Op("++")
						loop.If(jen.Id("i").Op("<").Id("n")).
							Block(jen.Goto().Id("read_loop"))
					})
				} else {
					// write generation code
					body.Id("i").Op(":=").Lit(0)
					body.Id("n").Op(":=").Len(jen.Id("data"))
					body.BlockFunc(func(loop *jen.Group) {
						loop.Id("write_loop:")
						if arguments[4] == "BE" {
							if arguments[2] == "U" || arguments[2] == "I" {
								loop.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Id("int64").
									Call(jen.Id("i").Op("*").Lit(intBytes))).Op("=").Id("byte").
									Call(jen.Id("data").Index(jen.Id("i")).Op(">>").Lit((intBytes - 1) * 8))
								for i := intBytes - 1; i > 1; i-- {
									loop.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Id("int64").
										Call(jen.Lit(intBytes - i).Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))).Op("=").Id("byte").
										Call(jen.Id("data").Index(jen.Id("i")).Op(">>").Lit((i - 1) * 8))
								}
								loop.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Id("int64").
									Call(jen.Lit(intBytes - 1).Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))).Op("=").Id("byte").
									Call(jen.Id("data").Index(jen.Id("i")))
							} else if arguments[2] == "F" {
								loop.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Id("int64").
									Call(jen.Id("i").Op("*").Lit(intBytes))).Op("=").Id("byte").
									Call(jen.Op("*").Parens(jen.Op("*").Id(strings.Join([]string{"uint", arguments[3]}, ""))).Parens(jen.Id("unsafe").Dot("Pointer").
										Call(jen.Op("&").Id("data").Index(jen.Id("i")))).Op(">>").Lit((intBytes - 1) * 8))
								for i := intBytes - 1; i > 1; i-- {
									loop.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Id("int64").
										Call(jen.Lit(intBytes - i).Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))).Op("=").Id("byte").
										Call(jen.Op("*").Parens(jen.Op("*").Id(strings.Join([]string{"uint", arguments[3]}, ""))).Parens(jen.Id("unsafe").Dot("Pointer").
											Call(jen.Op("&").Id("data").Index(jen.Id("i")))).Op(">>").Lit((i - 1) * 8))
								}
								loop.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Id("int64").
									Call(jen.Lit(intBytes - 1).Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))).Op("=").Id("byte").
									Call(jen.Op("*").Parens(jen.Op("*").Id(strings.Join([]string{"uint", arguments[3]}, ""))).Parens(jen.Id("unsafe").Dot("Pointer").
										Call(jen.Op("&").Id("data").Index(jen.Id("i")))))
								/*
									looking into reworking it to be like this:

									*(*uintBITS)(unsafe.Pointer(b.obuf + (i * 8))) = (or chain)
								*/
							}
						} else {
							if arguments[2] == "U" || arguments[2] == "I" {
								loop.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Id("int64").
									Call(jen.Id("i").Op("*").Lit(intBytes))).Op("=").Id("byte").
									Call(jen.Id("data").Index(jen.Id("i")))
								for i := 1; i < intBytes; i++ {
									loop.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Id("int64").
										Call(jen.Lit(i).Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))).Op("=").Id("byte").
										Call(jen.Id("data").Index(jen.Id("i")).Op(">>").Lit(i * 8))
								}
							} else if arguments[2] == "F" {
								loop.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Id("int64").
									Call(jen.Id("i").Op("*").Lit(intBytes))).Op("=").Id("byte").
									Call(jen.Op("*").Parens(jen.Op("*").Id(strings.Join([]string{"uint", arguments[3]}, ""))).Parens(jen.Id("unsafe").Dot("Pointer").
										Call(jen.Op("&").Id("data").Index(jen.Id("i")))))
								for i := 1; i < intBytes; i++ {
									loop.Id("b").Dot("buf").Index(jen.Id("off").Op("+").Id("int64").
										Call(jen.Lit(i).Op("+").Parens(jen.Id("i").Op("*").Lit(intBytes)))).Op("=").Id("byte").
										Call(jen.Op("*").Parens(jen.Op("*").Id(strings.Join([]string{"uint", arguments[3]}, ""))).Parens(jen.Id("unsafe").Dot("Pointer").
											Call(jen.Op("&").Id("data").Index(jen.Id("i")))).Op(">>").Lit(i * 8))
								}
							}
						}
						loop.Id("i").Op("++")
						loop.If(jen.Id("i").Op("<").Id("n")).
							Block(jen.Goto().Id("write_loop"))
					})
				}

				if arguments[0] == "Buffer" && arguments[1] == "Read" {
					body.Return()
				}
			})

			outputBuffer := bytes.NewBuffer([]byte{})
			err = builder.Render(outputBuffer)
			if err != nil {
				fmt.Println("! unable to render code:", err)
				return []byte("// render failure")
			}

			_, _ = outputBuffer.Write([]byte("\n\n"))

			builder = &jen.Group{}
			builder.Comment(strings.Join([]string{
				"// ",
				functionNameNext,
				" ",
				map[string]string{
					"Read":  "reads",
					"Write": "writes",
				}[arguments[1]],
				" a slice of ",
				intType,
				"s ",
				map[string]string{
					"Read":  "from",
					"Write": "to",
				}[arguments[1]],
				" the buffer at the\n",
			}, ""))
			builder.Comment(strings.Join([]string{
				"// current offset in ",
				map[string]string{
					"BE": "big-endian",
					"LE": "little-endian",
				}[arguments[4]],
				" and moves the offset forward the\n",
			}, ""))
			builder.Comment("// amount of bytes written\n")
			function = builder.Func().Params(jen.Id("b").Op("*").Id(arguments[0])).Id(functionNameNext)
			if arguments[1] == "Write" {
				function.Params(jen.Id("data").Index().Id(intType))
			} else if arguments[1] == "Read" {
				if arguments[0] == "Buffer" {
					function.Params(jen.Id("n").Id("int64")).Params(jen.Id("out").Index().Id(intType))
				} else {
					function.Params(jen.Id("out").Op("*").Index().Id(intType), jen.Id("n").Id("int64"))
				}
			}

			function.BlockFunc(func(body *jen.Group) {
				if arguments[1] == "Write" {
					body.Id("b").Dot(functionName).
						Call(jen.Id("b").Dot("off"), jen.Id("data"))
					body.Id("b").Dot("SeekByte").
						Call(jen.Id("int64").
							Call(jen.Len(jen.Id("data"))).Op("*").Lit(intBytes), jen.Lit(true))
				} else if arguments[1] == "Read" {
					if arguments[0] == "Buffer" {
						body.Id("out").Op("=").Id("b").Dot(functionName).
							Call(jen.Id("b").Dot("off"), jen.Id("n"))
					} else {
						body.Id("b").Dot(functionName).
							Call(jen.Id("out"), jen.Id("b").Dot("off"), jen.Id("n"))
					}

					body.Id("b").Dot("SeekByte").
						Call(jen.Id("n").Op("*").Lit(intBytes), jen.Lit(true))

					if arguments[0] == "Buffer" {
						body.Return()
					}
				}
			})

			err = builder.Render(outputBuffer)
			if err != nil {
				fmt.Println("! unable to render code:", err)
				return []byte("// render failure")
			}

			return outputBuffer.Bytes()
		})
	}
	return
}
