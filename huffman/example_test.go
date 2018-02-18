package huffman

import (
	"bytes"
	"testing"
)

func TestExample_WriteAndRead(t *testing.T) {
	src := []byte(loremText)
	m, _ := CreateModelFromText(src)

	encoded := bytes.NewBuffer([]byte{})
	w, _ := NewWriter(encoded, m)
	w.Write(src)
	_ = w.Close()
	r, _ := NewReader(encoded, m)

	decoded := make([]byte, len(src))
	r.Read(decoded)

	if bytes.Compare(decoded, src) != 0 {
		t.Errorf("Failed to encode/decode the text. \ngot  = %s\nwant = %s", decoded, src)
	}
}
