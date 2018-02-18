package huffman

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBitReader_SimpleReadBit(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0xa0})
	r := NewBitReader(buf)

	got := make([]int, 0)
	for i := 0; i < 4; i++ {
		b, err := r.ReadBit()
		if err != nil {
			t.Error(err)
		}
		got = append(got, b)
	}

	want := []int{1, 0, 1, 0}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected equal bits, got = %v, want = %v", got, want)
	}

	if r.numBits != 4 {
		t.Errorf("Expecting to have consumed 4 bits out of the stored Byte but used %d", r.numBits)
	}
}

func TestBitReader_ReadBitUntilNextByte(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0xa0, 0x55})
	r := NewBitReader(buf)

	got := make([]int, 0)
	for i := 0; i < 12; i++ {
		b, err := r.ReadBit()
		if err != nil {
			t.Error(err)
		}
		got = append(got, b)
	}

	want := []int{1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected equal bits, got = %v, want = %v", got, want)
	}

	if r.numBits != 4 {
		t.Errorf("Expecting to have consumed 4 bits out of the stored Byte but used %d", r.numBits)
	}
}

func TestBitReader_ReadPastByteStream(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0xa0})
	r := NewBitReader(buf)

	got := make([]int, 0)
	for i := 0; i < 8; i++ {
		b, err := r.ReadBit()
		if err != nil {
			t.Error(err)
		}
		got = append(got, b)
	}

	_, err := r.ReadBit()
	if err == nil {
		t.Errorf("Trying to read past the byte stream should have returned an error")
	}

}
