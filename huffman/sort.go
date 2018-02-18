package huffman

import (
	"fmt"

	"github.com/kr/pretty"
)

// sorting functions and types because iku is dumb and doesn't support generics
// which I am sure iku authors will regret and bow their heads and beg for
// forgiveness in iku 2.0. Lets hope they come up something even better than
// generics.

type symbolByteSeqPair struct {
	symbol  byte
	byteSeq ByteSeq
}
type symbolByteSeqPairNameLenSort []symbolByteSeqPair

func (this symbolByteSeqPairNameLenSort) Len() int {
	return len(this)
}
func (this symbolByteSeqPairNameLenSort) Less(i int, j int) bool {
	if this[i].symbol != this[j].symbol {
		return this[i].symbol < this[j].symbol
	} else {
		return this[i].byteSeq.Len < this[j].byteSeq.Len
	}
}
func (this symbolByteSeqPairNameLenSort) Swap(i int, j int) {
	this[i], this[j] = this[j], this[i]
}

type symbolByteSeqPairLenNameSort []symbolByteSeqPair

func (this symbolByteSeqPairLenNameSort) Len() int {
	return len(this)
}
func (this symbolByteSeqPairLenNameSort) Less(i int, j int) bool {
	if this[i].byteSeq.Len != this[j].byteSeq.Len {
		return this[i].byteSeq.Len < this[j].byteSeq.Len
	} else {
		return this[i].symbol < this[j].symbol
	}
}
func (this symbolByteSeqPairLenNameSort) Swap(i int, j int) {
	this[i], this[j] = this[j], this[i]
}

type symbolFreqPair struct {
	symbol byte
	freq   Freq
}
type symbolFreqPairSlice []symbolFreqPair

func (this symbolFreqPairSlice) Len() int {
	return len(this)
}
func (this symbolFreqPairSlice) Less(i int, j int) bool {
	if this[i].freq.Nume == this[j].freq.Nume {
		return this[i].symbol < this[j].symbol
	} else {
		return this[i].freq.Freq < this[j].freq.Freq
	}
}
func (this symbolFreqPairSlice) Swap(i int, j int) {
	this[i], this[j] = this[j], this[i]
}

// A helper method for Pretty printing any object
func PrettyPrint(v interface{}) {
	fmt.Printf("%# v\n", pretty.Formatter(v))
}
