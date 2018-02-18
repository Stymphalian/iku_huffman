package huffman

import (
	"sort"
	"testing"
)

const modelTestText = "A_DEAD_DAD_CEDED_A_BAD_BABE_A_BEADED_ABACA_BED"

func TestModel_CreateModelFromText(t *testing.T) {
	src := []byte(modelTestText)
	_, err := CreateModelFromText(src)
	if err != nil {
		t.Errorf("Failed to create a basic model from the text")
	}
}

func TestModel_GetPattern(t *testing.T) {
	src := []byte(modelTestText)
	m, err := CreateModelFromText(src)
	if err != nil {
		t.Errorf("Failed to create a basic model from the text")
	}

	got, err := m.GetPattern(byte('A'))
	want := &ByteSeq{0x02, 2}
	if err != nil {
		t.Error(err)
	}
	if got.Len != 2 {
		t.Errorf("Expected %d len got %d len for symbol %v",
			got.Len, want.Len, 'A')
	}
	if got.Pattern != 0x00 {
		t.Errorf("Expected %#x pattern got %#x pattern for symbol %v",
			got.Pattern, want.Pattern, 'A')
	}

	_, err = m.GetPattern(byte('P'))
	if err == nil {
		t.Errorf("Expecting to fail to get pattern buit instead got a pattern.")
	}
}

func TestModel_TestSort(t *testing.T) {
	src := []byte(modelTestText)
	m, err := CreateModelFromText(src)
	if err != nil {
		t.Errorf("Failed to create a basic model from the text")
	}

	ps := make(symbolByteSeqPairNameLenSort, 0)
	for k, v := range m.patternDict {
		ps = append(ps, symbolByteSeqPair{k, v})
	}

	sort.Stable(ps)
	sort.Stable(symbolByteSeqPairLenNameSort(ps))
}

func TestModel_MarhsalBinary(t *testing.T) {
	src := []byte(modelTestText)
	m, err := CreateModelFromText(src)
	if err != nil {
		t.Errorf("Failed to create a basic model from the text")
	}

	b, err := m.MarshalBinary()
	if err != nil {
		t.Errorf("Failed to marshal model to binary format %v", err)
	}
	if len(b) != 6 {
		t.Errorf("Not enough numbers %v", len(b))
	}
}

func TestModel_UnmarhsalBinary(t *testing.T) {
	alphabet := []byte("abcdefghijklmnopqrstuvwxyz")
	binary := []byte{
		5, 5, 5, 5, 5,
		4, 5, 5, 5, 5,
		5, 4, 5, 5, 5,
		5, 4, 4, 4, 5,
		4, 5, 5, 5, 5,
		5,
	}
	m := &Model{}
	err := m.UnmarshalBinary(alphabet, binary)
	if err != nil {
		t.Errorf("Failed to unmarshal model to binary format %v", err)
	}
	if len(m.patternDict) != 26 {
		t.Errorf("pattern dict should have 26 elements but only have %v",
			len(m.patternDict))
	}
	for i := 0; i < len(alphabet); i++ {
		if m.patternDict[alphabet[i]].Len != uint(binary[i]) {
			t.Errorf("Symbol %c is pattern dict does not have the correct length. got %d, want %d",
				alphabet[i], m.patternDict[alphabet[i]].Len, uint(binary[i]))
		}
	}

	symbols := make([]byte, 0)
	m.tree.InOrderTraversal(func(n *Node) {
		symbols = append(symbols, n.symbol)
	})
	for i := 0; i < len(alphabet); i++ {
		found := false
		for j := 0; j < len(symbols); j++ {
			if alphabet[i] == symbols[j] {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Missing symbol %c in huffman tree", alphabet[i])
		}
	}
}
