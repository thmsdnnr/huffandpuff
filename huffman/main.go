package huffman

import (
	"bufio"
	"bytes"
	"container/heap"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	bitstream "github.com/dgryski/go-bitstream"
)

const runeEOF = '见'

// Coder represents a Huffman encoding object.
type Coder interface {
	Init() error
	Encode() ([]byte, error)
	DecodeBytes([]byte) ([]byte, error)
}

// keyCount represents the number of times key is seen.
type keyCount struct {
	key   rune
	count int
}

// TreeNode represents a binary tree node having a value and left and right children.
type TreeNode struct {
	value *keyCount
	left  *TreeNode
	right *TreeNode
}

// NewTreeNode creates a new TreeNode having value k.
func NewTreeNode(k keyCount) *TreeNode {
	return &TreeNode{
		value: &k,
		left:  nil,
		right: nil,
	}
}

// max returns the greater of a and b
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Height retursn the height of tree rooted at node n.
func Height(n *TreeNode) int {
	if n == nil {
		return 0
	}
	leftHeight := Height(n.left)
	rightHeight := Height(n.right)
	return max(leftHeight, rightHeight) + 1
}

// pathSum stores the sum of all paths in pList using DFS. pList is a map of
// character in path to the directions taken (0 = left, 1 = right) to reach
// that character in the tree.
func pathSum(n *TreeNode, p string, pList *[]map[rune]string) {
	if n.left == nil && n.right == nil {
		mapPath := map[rune]string{n.value.key: p}
		*pList = append(*pList, mapPath)
		return
	}
	pathSum(n.left, p+"0", pList)
	pathSum(n.right, p+"1", pList)
	return
}

// keyCountList is a list of keyCount structs and implements the .Sort interface
type keyCountList []keyCount

func (k keyCountList) Len() int           { return len(k) }
func (k keyCountList) Less(i, j int) bool { return k[i].count < k[j].count }
func (k keyCountList) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }

// buildCodeWordDict
func (h *Huffandpuff) buildCodeWordDict() {
	pq := make(PriorityQueue, len(h.freqDict))
	i := 0
	for _, v := range h.keyCountList {
		pq[i] = &Item{
			value:    NewTreeNode(v),
			priority: v.count,
			index:    i,
		}
		i++
	}
	for pq.Len() > 1 {
		item1 := heap.Pop(&pq).(*Item)
		item2 := heap.Pop(&pq).(*Item)
		ntn := NewTreeNode(keyCount{
			key:   item1.value.value.key + item2.value.value.key,
			count: item1.priority + item2.priority,
		})
		ntn.left = item1.value
		ntn.right = item2.value
		ni := &Item{
			value:    ntn,
			priority: item1.priority + item2.priority,
		}
		pq.Push(ni)
		pq.update(ni, ni.value, ni.priority)
	}
	if pq.Len() == 0 {
		return
	}
	rootNode := heap.Pop(&pq).(*Item)
	pathList := []map[rune]string{}
	pathSum(rootNode.value, "", &pathList)
	for _, s := range pathList {
		for k, v := range s {
			h.decodingDict[v] = k
			h.encodingDict[k] = v
		}
	}
}

// buildFrequencyDict populates h.freqDict with the count of occurences of each rune.
func (h *Huffandpuff) buildFrequencyDict() error {
	if h.freqDict == nil {
		h.freqDict = map[rune]int{}
	}
	kvs := map[rune]keyCount{}
	r, err := h.getReader()
	if err != nil {
		return err
	}
	h.freqDict[runeEOF] = 1
	kvs[runeEOF] = keyCount{key: runeEOF, count: 1}
	c, _, err := r.ReadRune()
	for err != io.EOF {
		if _, ok := h.freqDict[c]; !ok {
			h.freqDict[c] = 1
			kvs[c] = keyCount{key: c, count: 1}
		} else {
			h.freqDict[c]++
			kv := kvs[c]
			kv.count++
			kvs[c] = kv
		}
		c, _, err = r.ReadRune()
	}
	return nil
}

var magicBytesHeader = []byte("=^･ｪ･^=")
var byteDelimiter = []byte("_(ツ)_")

type frequencyDict map[rune]int

type compressedHeader struct {
	HuffmanDictionary map[string]rune `json:"hd"`
}

func (h *Huffandpuff) getHeader() ([]byte, error) {
	writeBuf := bytes.NewBuffer(magicBytesHeader)
	header, err := h.headerToBytes()
	if err != nil {
		return nil, err
	}
	writeBuf.Write(header)
	writeBuf.Write(byteDelimiter)
	return writeBuf.Bytes(), nil
}

// ToFile saves the huffman encoded data to file at filepath.
func (h *Huffandpuff) ToFile(filepath string) error {
	header, err := h.getHeader()
	if err != nil {
		return err
	}
	fInfo, _ := os.Stat(filepath)
	if fInfo != nil {
		if fInfo.IsDir() {
			return fmt.Errorf("cannot save file as a directory")
		}
	}
	fPtr, err := os.Create(filepath)
	if err != nil {
		return err
	}

	writeCt, err := fPtr.Write(header)
	if err != nil {
		return err
	}
	if writeCt != len(header) {
		return fmt.Errorf("header truncated, wanted %d only wrote %d", len(header), writeCt)
	}

	// write the body
	return h.EncodeToFile(fPtr)
}

// headerToBytes writes Huffandpuff header to bytes.
func (h *Huffandpuff) headerToBytes() ([]byte, error) {
	ch := &compressedHeader{
		HuffmanDictionary: h.decodingDict,
	}
	htb, err := json.Marshal(ch)
	return htb, err
}

func (h *Huffandpuff) getDictionaryJSON() ([]byte, error) {
	return json.Marshal(h.decodingDict)
}

func (h *Huffandpuff) writeCodeword(w *bitstream.BitWriter, cw string) error {
	for _, c := range cw {
		if c == '1' {
			if err := w.WriteBit(true); err != nil {
				return err
			}
		} else {
			if err := w.WriteBit(false); err != nil {
				return err
			}
		}
	}
	return nil
}

// EncodeToFile encodes to a file directly.
func (h *Huffandpuff) EncodeToFile(f *os.File) error {
	defer f.Close()
	if !h.hasInit {
		return fmt.Errorf("Encode called without initialization")
	}
	r, err := h.getReader()
	if err != nil {
		return err
	}
	bnw := bitstream.NewWriter(f)
	ch, _, err := r.ReadRune()
	for err != io.EOF {
		enc, ok := h.encodingDict[ch]
		if !ok {
			log.Fatalf("no codeword for %x", string(ch))
		}
		h.writeCodeword(bnw, enc)
		ch, _, err = r.ReadRune()
		if err == io.EOF {
			h.writeCodeword(bnw, h.encodingDict[runeEOF])
			bnw.Flush(true)
			return nil
		}
	}
	return nil
}

func (h *Huffandpuff) getFilePtr() (*os.File, error) {
	if h.filename == "" {
		return nil, fmt.Errorf("filename not set")
	}
	fInfo, _ := os.Stat(h.filename)
	if fInfo != nil {
		if fInfo.IsDir() {
			return nil, fmt.Errorf("filename is a directory")
		}
	}
	return os.Open(h.filename)
}

// FromFile encodes to a file directly.
func (h *Huffandpuff) FromFile(filename string) error {
	h.filename = filename
	fPtr, err := h.getFilePtr()
	if err != nil {
		return err
	}
	defer fPtr.Close()

	h.file = fPtr
	r, err := h.getReader()
	if err != nil {
		return err
	}
	firstLine, _, err := r.ReadLine()
	if err != nil {
		return err
	}
	firstLineStr := string(firstLine)
	if !strings.HasPrefix(firstLineStr, string(magicBytesHeader)) {
		return fmt.Errorf("missing magic bytes header")
	}
	afterBytes := strings.SplitAfter(firstLineStr, string(magicBytesHeader))
	dict := strings.Split(afterBytes[1], string(byteDelimiter))[0]
	var ch compressedHeader
	offset := len(magicBytesHeader) + len(dict) + len(byteDelimiter)
	err = json.Unmarshal([]byte(dict), &ch)
	if err != nil {
		return err
	}
	h.decodingDict = ch.HuffmanDictionary
	h.encodedOffset = offset
	return nil
}

// GetDict returns the decoding dictionary
func (h *Huffandpuff) GetDict() map[string]rune {
	return h.decodingDict
}

// init initializes the Huffandpuff.
func (h *Huffandpuff) init() error {
	h.codeDict = map[string]int64{}
	h.decodingDict = map[string]rune{}
	h.encodingDict = map[rune]string{}
	if err := h.buildFrequencyDict(); err != nil {
		return err
	}
	var kpl keyCountList
	for k, v := range h.freqDict {
		kpl = append(kpl, keyCount{key: k, count: v})
	}
	sort.Sort(kpl)
	h.keyCountList = kpl
	h.buildCodeWordDict()
	h.hasInit = true
	return nil
}

// InitBytes initializes the encoder/decoder from bytes b.
func (h *Huffandpuff) InitBytes(b []byte) error {
	h.encStr = b
	return h.init()
}

// InitFile initializes the encoder/decoder from file f.
func (h *Huffandpuff) InitFile(f *os.File) error {
	h.file = f
	return h.init()
}

func (h *Huffandpuff) getReader() (*bufio.Reader, error) {
	if h.encStr != nil {
		return bufio.NewReader(bytes.NewBuffer(h.encStr)), nil
	}
	if h.file != nil {
		if _, err := h.file.Seek(0, 0); err != nil {
			return nil, err
		}
		return bufio.NewReader(h.file), nil
	}
	return nil, fmt.Errorf("no reader source")
}

// Encode encodes the string b using huffman coding.
func (h *Huffandpuff) Encode() ([]byte, error) {
	if !h.hasInit {
		return nil, fmt.Errorf("Encode called without initialization")
	}
	r, err := h.getReader()
	if err != nil {
		return nil, err
	}
	var writeBuf bytes.Buffer
	bnw := bitstream.NewWriter(&writeBuf)
	ch, _, err := r.ReadRune()
	for err != io.EOF {
		enc, ok := h.encodingDict[ch]
		if !ok {
			log.Fatalf("no codeword for %s", string(ch))
		}
		h.writeCodeword(bnw, enc)
		ch, _, err = r.ReadRune()
		if err == io.EOF {
			h.writeCodeword(bnw, h.encodingDict[runeEOF])
			bnw.Flush(true)
			break
		}
	}
	return writeBuf.Bytes(), nil
}

// DecodeFromFile decodes the huffman coded string d given h.decodingDict.
func (h *Huffandpuff) DecodeFromFile() ([]byte, error) {
	filePtr, err := h.getFilePtr()
	if err != nil {
		return nil, err
	}
	_, err = filePtr.Seek(int64(h.encodedOffset), 0)
	if err != nil {
		return nil, err
	}
	h.file = filePtr
	bnw := bitstream.NewReader(h.file)
	var buf bytes.Buffer
	match := ""
	b, err := bnw.ReadBit()
	if err != nil && err != io.EOF {
		log.Fatalf("ReadBit: %s", err)
	}
	for err != io.EOF {
		if b == false {
			match += "0"
		} else {
			match += "1"
		}
		out, ok := h.decodingDict[match]
		if !ok {
			b, err = bnw.ReadBit()
			if err != nil && err != io.EOF {
				log.Fatalf("ReadBit: %s", err)
			}
			continue
		}
		if out == runeEOF {
			break
		}

		_, err := buf.WriteRune(out)
		if err != nil {
			log.Fatalf("write err: %s", err)
		}

		match = ""
		b, err = bnw.ReadBit()
		if err != nil && err != io.EOF {
			log.Fatalf("ReadBit: %s", err)
		}
	}
	return buf.Bytes(), nil
}

// DecodeBytes decodes the huffman coded string d given h.decodingDict.
func (h *Huffandpuff) DecodeBytes(d []byte) ([]byte, error) {
	var buf bytes.Buffer
	match := ""
	bnw := bitstream.NewReader(bytes.NewBuffer(d))
	b, err := bnw.ReadBit()
	if err != nil && err != io.EOF {
		log.Fatalf("ReadBit: %s", err)
	}
	for err != io.EOF {
		if b == false {
			match += "0"
		} else {
			match += "1"
		}
		out, ok := h.decodingDict[match]
		if !ok {
			b, err = bnw.ReadBit()
			if err != nil && err != io.EOF {
				log.Fatalf("ReadBit: %s", err)
			}
			continue
		}
		if out == runeEOF {
			break
		}

		_, err := buf.WriteRune(out)
		if err != nil {
			log.Fatalf("write err: %s", err)
		}

		// got match
		match = ""

		b, err = bnw.ReadBit()
		if err != nil && err != io.EOF {
			log.Fatalf("ReadBit: %s", err)
		}
	}
	return buf.Bytes(), nil
}

// Huffandpuff represents an instance of a Huffman Decoder.
type Huffandpuff struct {
	encodedOffset int
	filename      string
	hasInit       bool
	file          *os.File
	reader        io.Reader
	keyCountList  []keyCount
	freqDict      map[rune]int
	codeDict      map[string]int64
	encodingDict  map[rune]string
	decodingDict  map[string]rune
	encStr        []byte
	bytesWritten  int64
}

// NewHuffandpuffInitBytes returns an initialized Huffandpuff with bytes s.
func NewHuffandpuffInitBytes(s []byte) (*Huffandpuff, error) {
	var th Huffandpuff
	if err := th.InitBytes(s); err != nil {
		return nil, err
	}
	return &th, nil
}

// NewHuffandpuffInitFile initializes Huffandpuff with a file handle.
func NewHuffandpuffInitFile(f *os.File) (*Huffandpuff, error) {
	if f == nil {
		return nil, fmt.Errorf("NewHuffandpuffInitReader called with nil file")
	}
	var h Huffandpuff
	if err := h.InitFile(f); err != nil {
		return nil, err
	}
	return &h, nil
}

// NewHuffandpuffFromFile loads a Huffandpuff encoded file from file filepath
func NewHuffandpuffFromFile(filepath string) (*Huffandpuff, error) {
	if filepath == "" {
		return nil, fmt.Errorf("NewHuffandpuffFromFile called with empty filepath")
	}
	var h Huffandpuff
	if err := h.FromFile(filepath); err != nil {
		return nil, err
	}
	return &h, nil
}

func main() {
	f, err := os.Open("./test.txt")
	if err != nil {
		log.Fatalf("file could not be opened: %s", err)
	}

	H, err := NewHuffandpuffInitFile(f)
	if err != nil {
		log.Fatalf("could not init Huffandpuff: %s", err)
	}

	enc, err := H.Encode()
	if err != nil {
		log.Fatalf("encoding err: %s", err)
	}

	dec, err := H.DecodeBytes(enc)
	if err != nil {
		log.Fatalf("encoding err: %s", err)
	}

	fmt.Print(string(dec))

	if err := f.Close(); err != nil {
		log.Fatalf("could not close file %s", err)
	}
}
