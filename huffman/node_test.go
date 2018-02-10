package huffman

import (
	"bytes"
	"container/heap"
	"testing"
)

func TestHeap(t *testing.T) {
	n := make(NodePQ, 0)
	heap.Init(&n)
	a := &Node{0x00, 0.1, nil, nil, nil}
	b := &Node{0x01, 0.2, nil, nil, nil}
	c := &Node{0x02, 0.3, nil, nil, nil}
	d := &Node{0x03, 0.4, nil, nil, nil}
	e := &Node{0x04, 0.5, nil, nil, nil}

	heap.Push(&n, d)
	heap.Push(&n, c)
	heap.Push(&n, a)
	heap.Push(&n, e)
	heap.Push(&n, b)

	want := []byte{0x00, 0x01, 0x02, 0x03, 0x04}
	got := make([]byte, 0)

	for n.Len() > 0 {
		got = append(got, heap.Pop(&n).(*Node).symbol)
	}
	if bytes.Compare(want, got) != 0 {
		t.Errorf("Priority queue failed")
	}
}
