package huffman

import (
	"container/heap"
	"errors"
	"fmt"
	"io"
	"log"
)

type byteSeqWriter struct {
	out    io.Writer
	carry  byte
	carryN uint
}

// read 'bits' off the 'seq' putting them in the lower bit places
// 'n' denotes how many bits were actually read off
// 'b' is where the data is stored
// mutates the ByteSeq.
func readBits(seq uint64, len uint, bits uint) (n uint, b byte) {
	if bits > 8 {
		log.Fatal("Cannot read off more than 8 bits at a time")
	}
	if bits > len {
		b = byte(seq)
		n = len
	} else {
		b = byte(seq >> (len - bits))
		n = bits
	}
	return
}

// Writes the entire byte sequence to the output writer
// Return the number of BITS successfully taken off the byte sequence.
// we get a copy of ByteSeq so it should be safe to mutate seq
func (this *byteSeqWriter) Write(seq ByteSeq) (n int, err error) {
	// First handle any carry bits from the last Write
	if this.carryN > 0 {
		room := 8 - this.carryN
		nBits, b := readBits(seq.Pattern, seq.Len, room)
		n += int(nBits)

		// clear the left most 'nBits' from the pattern
		left := 64 - seq.Len - nBits
		seq.Pattern = ((seq.Pattern << left) >> left)
		seq.Len = seq.Len - nBits

		var tmp byte = this.carry | (b << (room - nBits))
		if this.carryN+nBits == 8 {
			_, err := this.out.Write([]byte{tmp})
			if err != nil {
				return n, err
			}
			this.carry = 0x00
			this.carryN = 0
		} else {
			this.carry = tmp
			this.carryN += nBits
		}
	}

	for seq.Len > 0 {
		nBits, bits := readBits(seq.Pattern, seq.Len, 8)
		n += int(nBits)

		// clear the left most 'nBits' from the pattern
		left := 64 - (seq.Len - nBits)
		seq.Pattern = (seq.Pattern << left) >> left
		seq.Len -= nBits

		if nBits < 8 {
			// at the end of th sequence but didn't have a full byte of data
			this.carry = bits << (8 - nBits)
			this.carryN = nBits
		} else {
			_, err := this.out.Write([]byte{bits})
			if err != nil {
				return n, err
			}
		}
	}

	return n, nil
}

// Flush any last remaining bits as a single byte.
// Returns n - the number of 'real' bits that were in the byte. Those bits
// are placed in the most-significant digits of the byte and the
// least-signiicant digits are set to zero.
// YOU MUST CALL THIS FUNCTION.
func (this *byteSeqWriter) Flush() (n int, err error) {
	n = int(this.carryN)
	if this.carryN > 0 {
		_, err := this.out.Write([]byte{this.carry})
		if err != nil {
			return 0, err
		}
		this.carry = 0x00
		this.carryN = 0
	}
	return n, nil
}

type Writer struct {
	out  *byteSeqWriter
	dict map[byte]ByteSeq
}

func NewWriter(w io.Writer, dict map[byte]float64) (io.Writer, error) {
	this := new(Writer)
	this.out = &byteSeqWriter{w, 0, 0}
	this.Reset(dict)
	return this, nil
}

// Compress the provided 'p' bytes into the provided buffer
// return 'n' number of bytes from p which were compresseed.
func (this *Writer) Write(p []byte) (n int, err error) {
	// traverse through the src and output the bits
	n = 0

	for i := 0; i < len(p); i++ {
		seq, ok := this.dict[p[i]]
		if !ok {
			return n, errors.New(fmt.Sprintf("Unrecognized symbol %v", p[i]))
		}

		nBits, err := this.out.Write(seq)
		if err != nil {
			return n, err
		}
		if nBits != int(seq.Len) {
			return n, errors.New(fmt.Sprintf("Failed to write bits from sequence %v", seq))
		}

		n += 1
	}

	return n, nil
}

func (this *Writer) Reset(dict map[byte]float64) error {
	leafNodes, err := createTreeReturnLeafs(dict)
	if err != nil {
		return err
	}

	this.dict, err = buildMap(leafNodes)
	if err != nil {
		return err
	}

	return nil
}

func BuildFrequencyDict(src []byte) map[byte]float64 {
	dict := make(map[byte]float64)
	for _, b := range src {
		dict[b] += 1
	}

	total := float64(len(src))
	for k, _ := range dict {
		dict[k] /= total
	}
	return dict
}

// -----------------------------------------------------------------------------
// Internals
// -----------------------------------------------------------------------------

// Return tree + all the base nodes given the mapping of prioriy of symbols
func createTreeReturnLeafs(dict map[byte]float64) ([]*Node, error) {
	pq := make(NodePQ, 0)
	leafNodes := make([]*Node, 0)
	heap.Init(&pq)

	for k, v := range dict {
		n := &Node{k, v, nil, nil, nil}
		heap.Push(&pq, n)
		leafNodes = append(leafNodes, n)
	}

	for pq.Len() > 1 {
		a := heap.Pop(&pq).(*Node)
		b := heap.Pop(&pq).(*Node)
		n := &Node{0x00, a.freq + b.freq, nil, a, b}
		a.parent = n
		b.parent = n
		heap.Push(&pq, n)
	}

	if pq.Len() != 1 {
		return nil, errors.New("Failed to make tree")
	}

	return leafNodes, nil
}

// Construct a map from 'symbol' -> byte sequence given the set of
// leaf nodes. This should be called from the result of 'createTreeReturnLeafs'
func buildMap(leafNodes []*Node) (map[byte]ByteSeq, error) {
	dict := make(map[byte]ByteSeq)
	for i := 0; i < len(leafNodes); i++ {

		// traverse to root
		var byteSeq ByteSeq
		n := leafNodes[i]
		for n != nil {
			p := n.parent
			if p == nil {
				// we are at the root, stop processing
				break
			}

			// add a bit depending on the if we are the left or right child
			if p.left == n {
				byteSeq.AddBit(0)
			} else if p.right == n {
				byteSeq.AddBit(1)
			} else {
				return nil, errors.New("current node is not a child of its parent")
			}

			// keep going up
			n = p
		}

		dict[leafNodes[i].symbol] = byteSeq
	}

	return dict, nil
}
