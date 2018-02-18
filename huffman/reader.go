package huffman

import (
	"fmt"
	"io"
)

type Reader struct {
	r *BitReader
	m *Model
}

func NewReader(r io.Reader, m *Model) (*Reader, error) {
	return &Reader{NewBitReader(r), m}, nil
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
