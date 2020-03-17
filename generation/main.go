package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	directoryName  *string
	generationKind *string
)

func init() {
	directoryName = flag.String("directory", ".", "the directory to run code generation on the files in")
	generationKind = flag.String("generation", "", "what kind of generation to run on the files")
}

func main() {
	flag.Parse()

	var err error

	*directoryName, err = filepath.Abs(*directoryName)
	if err != nil {
		fmt.Printf("! unable to get absolute path of provided directory. (%v)\n", err)
		return
	}

	directory, err := os.Open(*directoryName)
	if err != nil {
		fmt.Printf("! unable to open source directory. (%v)\n", err)
		return
	}
	defer directory.Close()

	directoryInfo, err := directory.Stat()
	if err != nil {
		fmt.Printf("! unable to get source directory information. (%v)\n", err)
		return
	}

	if !directoryInfo.IsDir() {
		fmt.Println("! provided source directory is not a directory")
		return
	}

	// we can be sure that this will not error
	// we know that it is a directory
	_ = directory.Chdir()

	directoryFiles, err := directory.Readdir(0)
	if err != nil {
		fmt.Printf("! unable to read directory information. (%v)\n", err)
	}

	files := map[string][]byte{}
	for _, fileInfo := range directoryFiles {
		if filepath.Ext(fileInfo.Name()) == ".go" && !fileInfo.IsDir() && strings.HasPrefix(fileInfo.Name(), "_") {
			file, err := os.Open(fileInfo.Name())
			if err != nil {
				fmt.Printf("! unable to open source file. (%v)\n", err)
				continue
			}

			files[fileInfo.Name()], err = ioutil.ReadAll(file)
			if err != nil {
				fmt.Printf("! unable to read file contents. (%v)\n", err)
				continue
			}

			_ = file.Close()
		}
	}

	var generated map[string][]byte

	switch *generationKind {
	case "complex":
		generated, err = GenerateComplex(files)
		if err != nil {
			fmt.Printf("! @ GenerateComplex() %v\n", err)
			return
		}
	default:
		fmt.Println("! no valid generation mode specified")
		return
	}

	for name, contents := range generated {
		file, err := os.OpenFile(strings.Join([]string{strings.TrimSuffix(strings.TrimPrefix(name, "_"), ".go"), ".generated.go"}, ""), os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("! unable to open source file. (%v)\n", err)
			continue
		}

		_, err = file.Write(contents)
		if err != nil {
			fmt.Printf("! unable to update contents of generated file. (%v)\n", err)
			continue
		}

		_ = file.Close()
	}
}
