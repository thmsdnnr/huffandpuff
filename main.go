// This package implements a Huffman Coding CLI.
// https://en.wikipedia.org/wiki/Huffman_coding

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/thmsdnnr/hufflepuff/v2/huffman"
)

// ReadFromEncodeTo reads from file encoding to to file.
func ReadFromEncodeTo(from string, to string) error {
	f, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("file could not be opened: %s", err)
	}
	H, err := huffman.NewHufflepuffInitFile(f)
	if err != nil {
		return fmt.Errorf("could not init hufflepuff: %s", err)
	}
	if err := H.ToFile(to); err != nil {
		return err
	}
	return nil
}

// ReadFromDecodeTo reads encoded from from file, decoding to to file.
func ReadFromDecodeTo(from string, to string) error {
	h, err := huffman.NewHufflepuffFromFile(from)
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
	// inputFile := "./100-0.txt"
	outputFile := "./testingfoo.txt"
	stdOut := "/dev/stdout"
	// fooOut := "./foo.txt"

	// err := ReadFromEncodeTo(inputFile, outputFile)
	// if err != nil {
	// 	log.Fatalf("rfwt err: %s", err)
	// }

	err := ReadFromDecodeTo(outputFile, stdOut)
	if err != nil {
		log.Fatalf("rfwt err: %s", err)
	}

	// err = ReadFromDecodeTo(outputFile, fooOut)
	// if err != nil {
	// 	log.Fatalf("rfwt err: %s", err)
	// }

	// , err := os.Open("./test.txt")
	// if err != nil {
	// 	log.Fatalf("file could not be opened: %s", err)
	// }

	// // H, err := huffman.NewHufflepuffInitFile(f)
	// // if err != nil {
	// // 	log.Fatalf("could not init hufflepuff: %s", err)
	// // }

	// h, err := huffman.NewHufflepuffFromFile("./testingfoo.txt")
	// if err != nil {
	// 	log.Fatalf("from file err: %s", err)
	// }
	// dec, err := h.DecodeFromFile()
	// if err != nil {
	// 	log.Fatalf("decode from file err: %s", err)
	// }
	// fmt.Print(string(dec))

	// enc, err := H.Encode()
	// if err != nil {
	// 	log.Fatalf("encoding err: %s", err)
	// }

	// fmt.Printf("%s", string(enc))

	// for _, c := range enc {
	// 	fmt.Printf("%b", c)
	// }
	// fmt.Printf("%x ", string(enc))
	// v := fmt.Sprintf("%03b", enc)
	// log.Println(v)
	// if err := H.ToFile("./testingfoo.txt"); err != nil {
	// 	log.Fatalf("tf err: %s", err)
	// }

	// dec, err := H.DecodeBytes(enc)
	// if err != nil {
	// 	log.Fatalf("encoding err: %s", err)
	// }

	// fmt.Print(string(dec))

	// log.Printf("%+v", H.GetDict())
}
