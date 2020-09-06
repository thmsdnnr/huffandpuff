// This package implements a Huffman Coding CLI.
// https://en.wikipedia.org/wiki/Huffman_coding

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/thmsdnnr/huffandpuff/v2/huffman"
)

// ReadFromEncodeTo reads from file encoding to to file.
func ReadFromEncodeTo(from string, to string) error {
	f, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("file could not be opened: %s", err)
	}
	H, err := huffman.NewHuffandpuffInitFile(f)
	if err != nil {
		return fmt.Errorf("could not init Huffandpuff: %s", err)
	}
	if err := H.ToFile(to); err != nil {
		return err
	}
	return nil
}

// ReadFromDecodeTo reads encoded from from file, decoding to to file.
func ReadFromDecodeTo(from string, to string) error {
	h, err := huffman.NewHuffandpuffFromFile(from)
	if err != nil {
		return err
	}
	dec, err := h.DecodeFromFile()
	if err != nil {
		log.Fatalf("decode from file err: %s", err)
	}
	fPtr, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fPtr.Close()

	_, err = fPtr.Write(dec)
	return err
}

func main() {
	var inputFile, outputFile string
	var compress, decompress bool
	flag.StringVar(&inputFile, "in", "", "(required) name of input file")
	flag.StringVar(&outputFile, "out", "/dev/stdout", "(optional) name of output file, default stdout")
	flag.BoolVar(&compress, "c", false, "compress infile")
	flag.BoolVar(&decompress, "d", false, "decompress infile")
	flag.Parse()
	// inputFile := "./100-0.txt"
	// outputFile := "./testingfoo.txt"
	// stdOut := "/dev/stdout"
	// fooOut := "./foo.txt"
	if inputFile == "" {
		flag.Usage()
		os.Exit(1)
	}
	if outputFile == "" {
		log.Printf("no output file specified, using stdout")
	}

	if compress && decompress || !(compress || decompress) {
		flag.Usage()
		os.Exit(1)
	}

	if compress {
		err := ReadFromEncodeTo(inputFile, outputFile)
		if err != nil {
			log.Fatalf("ReadFromEncodeTo err: %s", err)
		}
	} else if decompress {
		err := ReadFromDecodeTo(inputFile, outputFile)
		if err != nil {
			log.Fatalf("ReadFromDecodeTo err: %s", err)
		}
	}
}
