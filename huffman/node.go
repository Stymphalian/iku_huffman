package huffman

type Node struct {
	symbol byte
	freq   float64
	parent *Node
	left   *Node
	right  *Node
}

type NodePQ []*Node

func (this Node) IsLeaf() bool {
	return this.left == nil && this.right == nil
}

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
