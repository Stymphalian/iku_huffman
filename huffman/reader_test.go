package huffman

import (
	"bytes"
	"testing"
)

func TestReader_SimpleRead(t *testing.T) {
	// The encoded source test is
	// abcabcabc
	// encoded in bytes as
	// 0101 1010 1101 0110 <- last bit here is empty.
	src := bytes.NewBuffer([]byte{0x5a, 0xc6})
	m, err := CreateModelFromText([]byte("aaaaaaaaaabbbbbccccc"))
	if err != nil {
		t.Fatal("Failed to create model")
	}
	r, err := NewReader(src, m)
	if err != nil {
		t.Fatal("Failed to create reader")
	}

	dest := make([]byte, 9)
	n, err := r.Read(dest)
	if err != nil {
		t.Errorf("Failed to decode symbols: %v", err)
	}
	if n != 9 {
		t.Errorf("Failed to read in expected 9 bytes, got %d", n)
	}
}
