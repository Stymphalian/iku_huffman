package huffman

import (
	"bytes"
	"testing"
)

func TestWriter_WriteSimple(t *testing.T) {
	dest := bytes.NewBuffer([]byte{})
	m, err := CreateModelFromText([]byte("aaaaaaaaaabbbbbccccc"))
	if err != nil {
		t.Fatal("Failed to create model")
	}
	w, err := NewWriter(dest, m)
	if err != nil {
		t.Fatal("Failed to create writer")
	}

	src := []byte("abc")
	n, err := w.Write(src)
	if err != nil || n != len(src) {
		t.Errorf("Failed to write bytes: %v", err)
	}
	err = w.Close()
	if err != nil {
		t.Errorf("Failed to close the stream: %v", err)
	}
	if bytes.Compare(dest.Bytes(), []byte{0x58}) != 0 {
		t.Errorf("Failed to write bit sequence %#v, got = %#v",
			[]byte{0x58}, dest.Bytes())
	}
	if w.BitsWritten() != 5 {
		t.Errorf("Expected 5 bits to be written but got %d", w.BitsWritten())
	}
}

func TestWriter_Write(t *testing.T) {
	src := []byte(loremText)
	dest := bytes.NewBuffer([]byte{})
	m, err := CreateModelFromText(src)
	if err != nil {
		t.Fatal("Failed to create model")
	}
	w, err := NewWriter(dest, m)
	if err != nil {
		t.Fatal("Failed to create writer")
	}

	n, err := w.Write(src)
	if err != nil || n != len(src) {
		t.Errorf("Failed to write bytes: %v", err)
	}
	err = w.Close()
	if err != nil {
		t.Errorf("Failed to close the stream: %v", err)
	}

	if len(src) < len(dest.Bytes()) {
		t.Errorf("Huffman did not actually compress the data %d, %d", len(src), len(dest.Bytes()))
	}
}

const loremText = typicalDefaultText
