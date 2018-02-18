package huffman

import (
	"bytes"
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

func TestByteSeqWriter_WriteSimple(t *testing.T) {
	dest := bytes.NewBuffer([]byte{})
	w := NewByteSeqWriter(dest)

	w.Write(ByteSeq{0x02, 2})
	w.Write(ByteSeq{0x02, 2})
	w.Write(ByteSeq{0x03, 2})
	w.Write(ByteSeq{0x03, 2})
	w.Write(ByteSeq{0x00, 1})
	w.Write(ByteSeq{0x00, 1})
	w.Write(ByteSeq{0x02, 2})
	w.Write(ByteSeq{0x02, 2})
	w.Write(ByteSeq{0x00, 1})
	w.Write(ByteSeq{0x00, 1})
	w.Flush()
	if bytes.Compare(dest.Bytes(), []byte{0xaf, 0x28}) != 0 {
		t.Errorf("Did not get all bytes %#v", dest.Bytes())
	}
}

func TestByteSeqWriter_WriteFlush(t *testing.T) {
	dest := bytes.NewBuffer([]byte{})
	w := NewByteSeqWriter(dest)

	// Write out the sequence, but because it is only 13 bits we will have
	// 5 bits to carry over
	// 1111 0000 1011 0
	seq := ByteSeq{0x1e16, 13}
	n, err := w.Write(seq)
	if err != nil {
		t.Error(err)
	}
	if n != 8 {
		t.Errorf("Was expecting to write 8 bits out of 13 but wrote %d", n)
	}
	if bytes.Compare(dest.Bytes(), []byte{0xf0}) != 0 {
		t.Errorf("Did not get all bytes %#v", dest.Bytes())
	}
	if seq.Pattern != 0x1e16 || seq.Len != 13 {
		t.Error("Calling Write should not mutate the seq")
	}

	// Write a second sequence, but it doesn't have enough bits so nothing get
	// output to the stream. we still take the bits of the sequence though
	// 11
	seq2 := ByteSeq{0x03, 2}
	n, err = w.Write(seq2)
	if err != nil {
		t.Error(err)
	}
	if n != 0 {
		t.Errorf("2: Was not expecting any bits to be written but wrote %d", n)
	}
	if bytes.Compare(dest.Bytes(), []byte{0xf0}) != 0 {
		t.Errorf("2: Did not get all bytes %#v", dest.Bytes())
	}

	// Write a 3rd sequence, take one bit to finish the carry then write out the
	// remainder of the sequence.
	// 10011 1111 1
	seq3 := ByteSeq{0x27f, 10}
	n, err = w.Write(seq3)
	if err != nil {
		t.Error(err)
	}
	if n != 16 {
		t.Errorf("3: expecting to write 16 bits of data but wrote %d", n)
	}
	if bytes.Compare(dest.Bytes(), []byte{0xf0, 0xb7, 0x3f}) != 0 {
		t.Errorf("3: Did not get all bytes %#v", dest.Bytes())
	}

	// Flush the stream to get the last byte
	n, err = w.Flush()
	if err != nil {
		t.Error(err)
	}
	if n != 1 {
		t.Error("4. Flush did not report the correct number of bits left")
	}
	// 1111 0000 1011 0111 0011 1111 1
	if bytes.Compare(dest.Bytes(), []byte{0xf0, 0xb7, 0x3f, 0x80}) != 0 {
		t.Errorf("3: Did not get all bytes %#v", dest.Bytes())
	}
}

// func TestReadBits(t *testing.T) {
// 	var pattern uint64 = 0x51de
// 	var len uint = 15

// 	nBits, bits := readBits(pattern, len, 3)
// 	if nBits != 3 {
// 		t.Fail()
// 	}
// 	if bits != 0x05 {
// 		t.Fail()
// 	}
// }

// func TestReadBitsNotEnough(t *testing.T) {
// 	var pattern uint64 = 0x05
// 	var len uint = 3

// 	nBits, bits := readBits(pattern, len, 8)
// 	if nBits != 3 {
// 		t.Fail()
// 	}
// 	if bits != 0x05 {
// 		t.Fail()
// 	}
// }
