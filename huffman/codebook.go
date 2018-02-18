package huffman

import (
	"bytes"
	"container/heap"
	"errors"
	"fmt"
	"sort"
)

type Freq struct {
	Nume uint64
	Deno uint64
	Freq float64
}

func BuildFrequencyDict(src []byte) map[byte]*Freq {
	dict := make(map[byte]*Freq)

	for _, b := range src {
		_, ok := dict[b]
		if !ok {
			dict[b] = &Freq{0, 0, 0}
		}
		dict[b].Nume += 1
	}

	total := uint64(len(src))
	for k, _ := range dict {
		dict[k].Deno = total
		dict[k].Freq = float64(dict[k].Nume) / float64(dict[k].Deno)
	}
	return dict
}

type Node struct {
	symbol byte
	freq   float64
	parent *Node
	left   *Node
	right  *Node
}

func (this *Node) IsLeaf() bool {
	return this.left == nil && this.right == nil
}

func (this *Node) Size() int {
	num := 0
	if this.left != nil {
		num += this.left.Size()
	}
	if this.right != nil {
		num += this.right.Size()
	}
	return num + 1
}

func (this *Node) InOrderTraversal(fn func(*Node)) {
	// In order traversal
	if this.left != nil {
		this.left.InOrderTraversal(fn)
	}
	fn(this)
	if this.right != nil {
		this.right.InOrderTraversal(fn)
	}
}

func (this *Node) InOrderTraversalDepth(fn func(*Node, int), depth int) {
	// In order traversal
	if this.left != nil {
		this.left.InOrderTraversalDepth(fn, depth+1)
	}
	fn(this, depth)
	if this.right != nil {
		this.right.InOrderTraversalDepth(fn, depth+1)
	}
}

func (this *Node) String() string {
	s := bytes.NewBuffer([]byte{})
	this.InOrderTraversalDepth(func(n *Node, depth int) {
		for i := 0; i < depth; i++ {
			s.WriteString(" ")
		}
		s.WriteString(fmt.Sprintf("%c:%.2f\n", n.symbol, n.freq))
	}, 0)
	return s.String()
}

// Type used for priority queue
type NodePQ []*Node

// sort.Interface methods
func (this NodePQ) Len() int {
	return len(this)
}
func (this NodePQ) Less(i int, j int) bool {
	return this[i].freq < this[j].freq
}
func (this NodePQ) Swap(i int, j int) {
	this[i], this[j] = this[j], this[i]

}

// heap.Interface methods
func (this *NodePQ) Push(x interface{}) {
	*this = append(*this, x.(*Node))
}
func (this *NodePQ) Pop() interface{} {
	old := *this
	n := len(old)
	item := old[n-1]
	*this = old[0 : n-1]
	return item
}

func buildHuffmanTree(dict map[byte]*Freq) (*Node, error) {

	// Return tree + all the base nodes given the mapping of prioriy of symbols
	pq := make(NodePQ, 0)
	heap.Init(&pq)

	sortedKeys := make(SymbolFreqPairSlice, 0)
	for k, v := range dict {
		sortedKeys = append(sortedKeys, SymbolFreqPair{k, *v})
	}
	sort.Stable(sortedKeys)

	for i := 0; i < len(sortedKeys); i++ {
		k := byte(sortedKeys[i].symbol)
		n := &Node{k, dict[k].Freq, nil, nil, nil}
		heap.Push(&pq, n)
	}

	for pq.Len() > 1 {
		a := heap.Pop(&pq).(*Node)
		b := heap.Pop(&pq).(*Node)
		n := &Node{0x00, a.freq + b.freq, nil, a, b}

		if a.Size() > b.Size() {
			// We always want the larger sub-tree on the right side
			n.left, n.right = n.right, n.left
		}
		a.parent = n
		b.parent = n
		heap.Push(&pq, n)
	}

	if pq.Len() != 1 {
		return nil, errors.New("Failed to make tree")
	}
	root := pq.Pop().(*Node)
	return root, nil
}

// Construct a map from 'symbol' -> byte sequence given the set of
// leaf nodes. This should be called from the result of 'createTreeReturnLeafs'
func buildPatternDict(root *Node) (map[byte]ByteSeq, error) {
	leafNodes := make([]*Node, 0)
	root.InOrderTraversal(func(n *Node) {
		if n.IsLeaf() {
			leafNodes = append(leafNodes, n)
		}
	})

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

// Given a codebook mapping symbol -> length of byte pattern
// Create a canonical encoding of the codebook
func canonicalCodebook(codebook map[byte]ByteSeq) (map[byte]ByteSeq, error) {
	ps := make(SymbolByteSeqPairLenNameSort, 0)
	for k, v := range codebook {
		ps = append(ps, SymbolByteSeqPair{k, v})
	}
	sort.Stable(ps)

	newCodeBook := make(map[byte]ByteSeq)
	patLen := uint(0)
	pat := uint64(0)
	for i := 0; i < len(ps); i++ {
		patLen = ps[i].byteSeq.Len
		newCodeBook[ps[i].symbol] = ByteSeq{pat, patLen}

		pat += 1
		if i < len(ps)-1 {
			if ps[i+1].byteSeq.Len > ps[i].byteSeq.Len {
				diff := ps[i+1].byteSeq.Len - ps[i].byteSeq.Len
				pat <<= diff
			}
		}
	}

	return newCodeBook, nil
}

func canonicalHuffmanTree(codebook map[byte]ByteSeq) (*Node, error) {
	ps := make(SymbolByteSeqPairLenNameSort, 0)
	for k, v := range codebook {
		ps = append(ps, SymbolByteSeqPair{k, v})
	}
	sort.Stable(ps)

	root := &Node{}
	for i := 0; i < len(ps); i++ {

		n := root
		for j := int(ps[i].byteSeq.Len - 1); j >= 0; j-- {
			if ps[i].byteSeq.Pattern&(1<<uint(j)) > 0 {
				if n.right == nil {
					n.right = &Node{0, 0, n, nil, nil}
				}
				n = n.right
			} else {
				if n.left == nil {
					n.left = &Node{0, 0, n, nil, nil}
				}
				n = n.left
			}

			// is at root node
			if j == 0 {
				n.symbol = ps[i].symbol
				n.freq = 0.0
			}
		}
	}
	return root, nil
}
