package huffman

import (
	"bytes"
	"errors"
	"io"
	"log"

	"github.com/Stymphalian/iku_bits/bitwriter"
)

type ByteSeq struct {
	Pattern uint64
	Len     uint
}

func (this *ByteSeq) AddBit(one int) {
	if this.Len >= 64 {
		panic(errors.New("Can't add more than 64 bits to ByteSeq"))
	}
	if one == 1 {
		this.Pattern |= (0x01 << this.Len)
	}
	this.Len += 1
}

func (this ByteSeq) String() string {
	buf := bytes.NewBufferString("")
	for i := int(this.Len) - 1; i >= 0; i-- {
		if this.Pattern&(0x01<<uint(i)) > 0 {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
	}
	return buf.String()
}

type ByteSeqWriter struct {
	// The writer interface in which to write the Bytes
	w bitwriter.Interface
}

func NewByteSeqWriter(w io.Writer) *ByteSeqWriter {
	b, err := bitwriter.NewBitWriter(w)
	if err != nil {
		log.Fatal("Failed to create a bit writer", err)
	}
	return &ByteSeqWriter{b}
}

// Writes the byte sequence into the writer stream
// Return[int] the number of bits written to the stream
// Return[error] nil if okay, otherwise error object
func (this *ByteSeqWriter) Write(seq ByteSeq) (n int, err error) {
	n = this.w.Remain()
	for i := int(seq.Len - 1); i >= 0; i-- {
		n += 1
		var err error
		if seq.Pattern&(1<<uint(i)) > 0 {
			err = this.w.WriteBit(1)
		} else {
			err = this.w.WriteBit(0)
		}
		if err != nil {
			return (n / 8) * 8, err
		}
	}
	return (n / 8) * 8, nil
}

// Closes the byte seq writer flushing the remaining bits as a single byte into
// the stream.
// Returns int - The number of bits written
// Returns error - nil if okay, error otherwise
func (this *ByteSeqWriter) Flush() (n int, err error) {
	return this.w.Flush()
}
