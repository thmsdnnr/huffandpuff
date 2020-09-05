package huffman

import "io"

type BitTwiddler interface {
	WriteBit(bool) error
	ReadBit(io.Reader) (byte, error)
}

type bitTwid struct {
	buffer    byte
	bufferIdx int
	bitBuffer byte
	writer    io.Writer
	reader    io.Reader
	inBuffer  []byte
}

func (b *bitTwid) Init(w io.Writer, r io.Reader) error {
	b.writer = w
	b.reader = r
	return nil
}

// WriteBit writes a 0 or 1 based on the value of b.
func (b *bitTwid) WriteBit(bit bool) error {
	if b.bufferIdx == 8 {
		// write byte.
		b.writer.Write([]byte{b.buffer})
		b.bufferIdx = 0
		b.buffer = 0
	} else {
		b.buffer <<= 1
		if bit {
			b.buffer |= 1
		}
	}
	return nil
}

// func (b *bitTwid) ReadBit(r io.Reader) (byte, error) {
// 	if len(b.inBuffer) == 0 {
// 		b.buffer, err = r.ReadLine()
// 		if err != io.EOF {
// 			return 0, err
// 		}
// 	}
// 	if len(b.inBuffer) == 0 {
// 		return 0, nil
// 	}
// 	b.buffer = b.inBuffer[0]
// }
