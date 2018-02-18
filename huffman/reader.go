package huffman

import (
	"fmt"
	"io"

	"github.com/Stymphalian/iku_bits/bitreader"
)

type Reader struct {
	r bitreader.Interface
	m *Model
}

func NewReader(r io.Reader, m *Model) (*Reader, error) {
	b, err := bitreader.NewBitReader(r)
	if err != nil {
		return nil, err
	}
	return &Reader{b, m}, nil
}

func (this *Reader) Read(p []byte) (int, error) {
	numBytes := 0
	node := this.m.tree
	for numBytes < len(p) {
		if node.IsLeaf() {
			p[numBytes] = node.symbol
			node = this.m.tree
			numBytes += 1
		}

		b, err := this.r.ReadBit()
		if err != nil {
			return numBytes, err
		}

		if b == 1 {
			if node.right != nil {
				node = node.right
			} else {
				return numBytes, fmt.Errorf(
					"Invalid huffman tree, expecting a right child but found nil")
			}
		} else {
			if node.left != nil {
				node = node.left
			} else {
				return numBytes, fmt.Errorf(
					"Invalid huffman tree, expecting a left child but found nil")
			}
		}
	}
	return numBytes, nil
}
