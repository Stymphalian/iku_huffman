package huffman

import (
	"fmt"

	"github.com/kr/pretty"
)

// sorting functions and types because iku is dumb and doesn't support generics
// which I am sure iku authors will regret and bow their heads and beg for
// forgiveness in iku 2.0. Lets hope they come up something even better than
// generics.

type SymbolByteSeqPair struct {
	symbol  byte
	byteSeq ByteSeq
}
type SymbolByteSeqPairNameLenSort []SymbolByteSeqPair

func (this SymbolByteSeqPairNameLenSort) Len() int {
	return len(this)
}
func (this SymbolByteSeqPairNameLenSort) Less(i int, j int) bool {
	if this[i].symbol != this[j].symbol {
		return this[i].symbol < this[j].symbol
	} else {
		return this[i].byteSeq.Len < this[j].byteSeq.Len
	}
}
func (this SymbolByteSeqPairNameLenSort) Swap(i int, j int) {
	this[i], this[j] = this[j], this[i]
}

type SymbolByteSeqPairLenNameSort []SymbolByteSeqPair

func (this SymbolByteSeqPairLenNameSort) Len() int {
	return len(this)
}
func (this SymbolByteSeqPairLenNameSort) Less(i int, j int) bool {
	if this[i].byteSeq.Len != this[j].byteSeq.Len {
		return this[i].byteSeq.Len < this[j].byteSeq.Len
	} else {
		return this[i].symbol < this[j].symbol
	}
}
func (this SymbolByteSeqPairLenNameSort) Swap(i int, j int) {
	this[i], this[j] = this[j], this[i]
}

type SymbolFreqPair struct {
	symbol byte
	freq   Freq
}
type SymbolFreqPairSlice []SymbolFreqPair

func (this SymbolFreqPairSlice) Len() int {
	return len(this)
}
func (this SymbolFreqPairSlice) Less(i int, j int) bool {
	if this[i].freq.Nume == this[j].freq.Nume {
		return this[i].symbol < this[j].symbol
	} else {
		return this[i].freq.Freq < this[j].freq.Freq
	}
}
func (this SymbolFreqPairSlice) Swap(i int, j int) {
	this[i], this[j] = this[j], this[i]
}

// A helper method for Pretty printing any object
func PrettyPrint(v interface{}) {
	fmt.Printf("%# v\n", pretty.Formatter(v))
}
