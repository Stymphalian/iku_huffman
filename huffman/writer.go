package huffman

import (
	"io"
)

type Writer struct {
	w           *ByteSeqWriter
	m           *Model
	bitsWritten uint64
}

func NewWriter(w io.Writer, m *Model) (*Writer, error) {
	return &Writer{NewByteSeqWriter(w), m, 0}, nil
}

func (this *Writer) Write(p []byte) (int, error) {
	numBytesWritten := 0
	for i := 0; i < len(p); i++ {
		symbol, err := this.m.GetPattern(p[i])
		if err != nil {
			return numBytesWritten, err
		}

		bitsWritten, err := this.w.Write(symbol)
		if err != nil {
			return numBytesWritten, err
		}
		this.bitsWritten += uint64(bitsWritten)
		numBytesWritten += 1
	}
	return numBytesWritten, nil
}

func (this *Writer) Close() error {
	bitsWritten, err := this.w.Flush()
	if err != nil {
		return err
	}

	this.bitsWritten += uint64(bitsWritten)
	return nil
}

func (this *Writer) BitsWritten() uint64 {
	return this.bitsWritten
}
