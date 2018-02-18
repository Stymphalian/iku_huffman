package huffman

import (
	"bytes"
	"container/heap"
	"testing"
)

const codebookTestText = "A_DEAD_DAD_CEDED_A_BAD_BABE_A_BEADED_ABACA_BED"

func TestCodebook_BuildFrequencyDict(t *testing.T) {
	d := BuildFrequencyDict([]byte("ABACAAACBC"))
	testcases := []struct {
		s byte
		n uint64
		d uint64
		f float64
	}{
		{byte('A'), 5, 10, 0.5},
		{byte('B'), 2, 10, 0.2},
		{byte('C'), 3, 10, 0.3},
	}

	for _, tc := range testcases {
		got := d[tc.s]
		if got.Nume != tc.n {
			t.Errorf("Numerator mismtach: got %v, want %v for symbol %v",
				got.Nume, tc.n, tc.s)
		}
		if got.Deno != tc.d {
			t.Errorf("Denominator mismtach: got %v, want %v for symbol %v",
				got.Deno, tc.d, tc.s)
		}
		if got.Freq != tc.f {
			t.Errorf("Frequency mismtach: got %v, want %v for symbol %v",
				got.Freq, tc.f, tc.s)
		}
	}
}

func TestCodebook_TestFreqAddsTo1(t *testing.T) {
	d := BuildFrequencyDict([]byte(codebookTestText))
	sum := float64(0.0)
	for _, v := range d {
		sum += v.Freq
	}
	if sum-1.0 > 0.0000001 {
		t.Errorf("Frequency of characters do not add to 1")
	}
}

func TestCodebook_TestNodeIsLeaf(t *testing.T) {
	nilNode := &Node{0, 0, nil, nil, nil}
	for _, tc := range []struct {
		n    *Node
		want bool
	}{
		{nilNode, true},
		{&Node{0, 0, nil, nilNode, nil}, false},
		{&Node{0, 0, nil, nil, nilNode}, false},
		{&Node{0, 0, nil, nilNode, nilNode}, false},
	} {
		if tc.n.IsLeaf() != tc.want {
			t.Errorf("Expected node %v.IsLeaf() == %v, but got %v",
				tc.n, tc.want, tc.n.IsLeaf())
		}
	}
}

func TestCodebook_TestHeap(t *testing.T) {
	n := make(NodePQ, 0)
	heap.Init(&n)
	a := &Node{byte('a'), 2, nil, nil, nil}
	b := &Node{byte('b'), 6, nil, nil, nil}
	c := &Node{byte('c'), 7, nil, nil, nil}
	d := &Node{byte('d'), 10, nil, nil, nil}
	e := &Node{byte('e'), 10, nil, nil, nil}
	f := &Node{byte('f'), 11, nil, nil, nil}

	heap.Push(&n, d)
	heap.Push(&n, c)
	heap.Push(&n, a)
	heap.Push(&n, f)
	heap.Push(&n, e)
	heap.Push(&n, b)

	want := []byte("abcdef")
	got := make([]byte, 0)

	for n.Len() > 0 {
		got = append(got, heap.Pop(&n).(*Node).symbol)
	}
	if bytes.Compare(want, got) != 0 {
		t.Errorf("Priority queue failed")
	}
}

func TestCodebook_BuildHuffmanTree(t *testing.T) {
	d := map[byte]*Freq{
		0x00: &Freq{0, 0, 0.4},
		0x01: &Freq{0, 0, 0.35},
		0x02: &Freq{0, 0, 0.2},
		0x03: &Freq{0, 0, 0.05},
	}
	root, err := buildHuffmanTree(d)
	if err != nil {
		t.Error(err)
	}

	dict, err := buildPatternDict(root)
	if err != nil {
		t.Error(err)
	}
	dict, err = canonicalCodebook(dict)
	if err != nil {
		t.Error(err)
	}

	wantDict := map[byte]ByteSeq{
		0x00: ByteSeq{0x00, 1},
		0x01: ByteSeq{0x02, 2},
		0x02: ByteSeq{0x06, 3},
		0x03: ByteSeq{0x07, 3},
	}
	for k, v := range wantDict {
		if dict[k] != v {
			t.Errorf("byteSeq don't match %v, got = %v, want = %v", k, dict[k], v)
		}
	}
}

func TestCodebook_BuildTreeAndFreqDict(t *testing.T) {
	src := []byte(codebookTestText)
	weights := BuildFrequencyDict(src)
	root, err := buildHuffmanTree(weights)
	if err != nil {
		t.Error(err)
	}

	dict, err := buildPatternDict(root)
	if err != nil {
		t.Error(err)
	}
	dict, err = canonicalCodebook(dict)
	if err != nil {
		t.Error(err)
	}

	wantDict := map[byte]ByteSeq{
		'_': ByteSeq{0x02, 2},
		'D': ByteSeq{0x01, 2},
		'A': ByteSeq{0x00, 2},
		'E': ByteSeq{0x06, 3},
		'C': ByteSeq{0x0f, 4},
		'B': ByteSeq{0x0e, 4},
	}
	for k, v := range wantDict {
		if dict[k] != v {
			t.Errorf("byteSeq don't match %c, got = %v, want = %v", k, dict[k], v)
		}
	}
}

func TestCodebook_canonicalCodebook(t *testing.T) {
	codebook := map[byte]ByteSeq{
		byte('a'): ByteSeq{0x03, 2},
		byte('b'): ByteSeq{0x00, 1},
		byte('c'): ByteSeq{0x05, 3},
		byte('d'): ByteSeq{0x04, 3},
	}
	newbook, err := canonicalCodebook(codebook)
	if err != nil {
		t.Errorf("Failed to create canonical codebook: %v", err)
	}
	wantbook := map[byte]ByteSeq{
		byte('b'): ByteSeq{0x00, 1},
		byte('a'): ByteSeq{0x02, 2},
		byte('c'): ByteSeq{0x06, 3},
		byte('d'): ByteSeq{0x07, 3},
	}

	for k, v := range wantbook {
		if newbook[k] != v {
			t.Errorf("Did not get canonical pattern for %c, got=%v, want=%v", k, newbook[k], v)
		}
	}
}

func TestCodebook_SimpleCanonicalHuffmanTreeSimple(t *testing.T) {
	codebook := map[byte]ByteSeq{
		byte('a'): ByteSeq{0x00, 1},
		byte('b'): ByteSeq{0x02, 2},
		byte('c'): ByteSeq{0x03, 2},
	}
	newbook, err := canonicalCodebook(codebook)
	if err != nil {
		t.Errorf("Failed to create canonical codebook: %v", err)
	}

	tree, err := canonicalHuffmanTree(newbook)
	if err != nil {
		t.Error(err)
	}

	got, err := buildPatternDict(tree)
	if err != nil {
		t.Error(err)
	}

	for k, v := range got {
		if got[k] != newbook[k] {
			t.Errorf(`New huffman tree does not produce the correct bit pattern
for symbol %c. got %v , want %v`, k, v, newbook[k])
		}
	}
}

func TestCodebook_SimpleCanonicalHuffmanTree(t *testing.T) {
	codebook := map[byte]ByteSeq{
		byte('a'): ByteSeq{0x00, 1},
		byte('b'): ByteSeq{0x02, 2},
		byte('c'): ByteSeq{0x03, 2},
	}
	newbook, err := canonicalCodebook(codebook)
	if err != nil {
		t.Errorf("Failed to create canonical codebook: %v", err)
	}

	tree, err := canonicalHuffmanTree(newbook)
	if err != nil {
		t.Error(err)
	}

	got, err := buildPatternDict(tree)
	if err != nil {
		t.Error(err)
	}

	for k, v := range got {
		if got[k] != newbook[k] {
			t.Errorf(`New huffman tree does not produce the correct bit pattern
for symbol %c. got %v , want %v`, k, v, newbook[k])
		}
	}
}
