package huffman

import (
	"log"
	"testing"
)

func TestHufflepuffEncodeDecode(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    []byte
		Expected []byte
	}{
		{
			Name:     "empty string",
			Input:    []byte(""),
			Expected: []byte(""),
		},
		{
			Name:     "1 char",
			Input:    []byte("1"),
			Expected: []byte("1"),
		},
		{
			Name:     "2 chars same",
			Input:    []byte("11"),
			Expected: []byte("11"),
		},
		{
			Name:     "2 chars different",
			Input:    []byte("12"),
			Expected: []byte("12"),
		},
		{
			Name:     "gopher koan",
			Input:    []byte("Don't communicate by sharing memory; share memory by communicating."),
			Expected: []byte("Don't communicate by sharing memory; share memory by communicating."),
		},
		{
			Name:     "gopher koan utf8",
			Input:    []byte("通过通信共享内存"),
			Expected: []byte("通过通信共享内存"),
		},
		{
			Name:     "emoji",
			Input:    []byte("😺️😸️😹️😻️😼️😽️🙀️😿️😾️🐈️"),
			Expected: []byte("😺️😸️😹️😻️😼️😽️🙀️😿️😾️🐈️"),
		},
		{
			Name:     "multiline \n",
			Input:    []byte("😺️😸️😹️😻️😼️😽️🙀️😿️😾️🐈️\n😺️😸️😹️😻️😼️😽️🙀️😿️😾️🐈️\n"),
			Expected: []byte("😺️😸️😹️😻️😼️😽️🙀️😿️😾️🐈️\n😺️😸️😹️😻️😼️😽️🙀️😿️😾️🐈️\n"),
		},
		{
			Name:     "multiline \r\n",
			Input:    []byte("😺️😸️😹️😻️😼️😽️🙀️😿️😾️🐈️\r\n😺️😸️😹️😻️😼️😽️🙀️😿️😾️🐈️\r\n"),
			Expected: []byte("😺️😸️😹️😻️😼️😽️🙀️😿️😾️🐈️\r\n😺️😸️😹️😻️😼️😽️🙀️😿️😾️🐈️\r\n"),
		},
		{
			Name:     "hex",
			Input:    []byte("\x01\x02\x03\x04\x05\x06\x07\x08\x09\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x20"),
			Expected: []byte("\x01\x02\x03\x04\x05\x06\x07\x08\x09\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x20"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			H, err := NewHufflepuffInitBytes(tc.Input)
			if err != nil {
				t.Errorf("got Init err %s want nil", err)
			}
			enc, err := H.Encode()
			if err != nil {
				t.Errorf("got Encode err %s want nil", err)
			}
			dec, err := H.DecodeBytes(enc)
			if err != nil {
				t.Errorf("got Decode err %s want nil", err)
			}
			if string(dec) != string(tc.Expected) {
				t.Errorf("mismatch wanted \n%s got\n%s", string(tc.Expected), string(dec))
			}
			log.Println(string(dec))
		})
	}
}
