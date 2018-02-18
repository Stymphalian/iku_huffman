package codec

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/Stymphalian/iku_huffman/huffman"
)

const (
	VERSION   = uint16(0x53)
	HAS_MODEL = 0x0001
)

type Encoder struct {
	w  io.Writer
	m  *huffman.Model
	hw *huffman.Writer
}

func NewEncoder(w io.Writer) (*Encoder, error) {
	m := huffman.DefaultModel()
	hw, err := huffman.NewWriter(w, m)
	if err != nil {
		return nil, err
	}
	return &Encoder{w, m, hw}, nil
}

func (this *Encoder) Write(p []byte, flags uint16) (int, error) {
	err := binary.Write(this.w, binary.LittleEndian, VERSION)
	if err != nil {
		return 0, err
	}
	err = binary.Write(this.w, binary.LittleEndian, flags)
	if err != nil {
		return 0, err
	}

	err = binary.Write(this.w, binary.LittleEndian, uint64(len(p)))
	if err != nil {
		return 0, err
	}

	// Optionally write the huffman tree model
	// not needed assuming that the Decoder know what model to use.
	if flags&HAS_MODEL > 0 {
		bs, err := this.m.MarshalBinary()
		if err != nil {
			return 0, nil
		}
		this.w.Write(bs)
	}

	// Write the paylaod
	n, err := this.hw.Write(p)
	if err != nil {
		return n, err
	}
	err = this.hw.Close()
	return n, err
}

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) (*Decoder, error) {
	return &Decoder{r}, nil
}

func (this *Decoder) Read() ([]byte, error) {
	var version uint16
	var flags uint16
	var payloadLen uint64

	err := binary.Read(this.r, binary.LittleEndian, &version)
	if err != nil {
		return nil, err
	}
	err = binary.Read(this.r, binary.LittleEndian, &flags)
	if err != nil {
		return nil, err
	}
	err = binary.Read(this.r, binary.LittleEndian, &payloadLen)
	if err != nil {
		return nil, err
	}

	var m *huffman.Model
	// Optionally write the huffman tree
	if flags&HAS_MODEL > 0 {
		m = &huffman.Model{}
		alphabet := huffman.DefaultAlphabet()
		tree := make([]byte, 0)
		for i := 0; i < len(alphabet); i++ {
			var b byte
			err := binary.Read(this.r, binary.LittleEndian, &b)
			if err != nil {
				return nil, err
			}
			tree = append(tree, b)
		}
		err := m.UnmarshalBinary(alphabet, tree)
		if err != nil {
			return nil, err
		}
	} else {
		m = huffman.DefaultModel()
	}

	hr, err := huffman.NewReader(this.r, m)
	if err != nil {
		return nil, err
	}

	p := make([]byte, payloadLen)
	n, err := hr.Read(p)
	if err != nil {
		return nil, err
	}
	if uint64(n) != payloadLen {
		return nil, fmt.Errorf("Failed to read %d symbols from the stream.", n)
	}
	return p, nil
}
