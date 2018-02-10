package huffman

import (
	"testing"
)

func TestWriteBit(t *testing.T) {
	var s ByteSeq
	s.AddBit(1)
	s.AddBit(0)
	s.AddBit(1)
	s.AddBit(1)
	s.AddBit(1)
	s.AddBit(1)
	s.AddBit(0)
	s.AddBit(1)
	s.AddBit(1)

	if s.Len != 9 {
		t.Error("Len should be 9")
	}
	if s.Pattern != 0x01bd {
		t.Errorf("Pattern should be 0x01bd but got %#x", s.Pattern)
	}
}

func TestByteSeqString(t *testing.T) {
	var s ByteSeq
	s.AddBit(1)
	s.AddBit(0)
	s.AddBit(1)
	s.AddBit(1)
	s.AddBit(1)
	s.AddBit(1)
	s.AddBit(0)
	s.AddBit(1)
	s.AddBit(1)

	if s.String() != "110111101" {
		t.Errorf("ByteSeq as string does not match")
	}
}
