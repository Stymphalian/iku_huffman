package huffman

import (
	"bytes"
	"errors"
	"io"
	"log"
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
	w io.Writer
	// The remaining bits to be written
	carry byte
	// How many bits have yet to be flushed to the stream
	carryN uint
}

func NewByteSeqWriter(w io.Writer) *ByteSeqWriter {
	return &ByteSeqWriter{w, 0, 0}
}

// Writes the byte sequence into the writer stream
// Return[int] the number of bits written to the stream
// Return[error] nil if okay, otherwise error object
func (this *ByteSeqWriter) Write(seq ByteSeq) (n int, err error) {
	// First handle any carry bits from the last Write
	if this.carryN > 0 {
		room := 8 - this.carryN
		nBits, b := readBits(seq.Pattern, seq.Len, room)

		// clear the left most 'nBits' from the pattern
		left := 64 - seq.Len - nBits
		seq.Pattern = ((seq.Pattern << left) >> left)
		seq.Len = seq.Len - nBits

		var tmp byte = this.carry | (b << (room - nBits))
		if this.carryN+nBits == 8 {
			_, err := this.w.Write([]byte{tmp})
			if err != nil {
				return n, err
			}
			n += 8
			this.carry = 0x00
			this.carryN = 0
		} else {
			this.carry = tmp
			this.carryN += nBits
		}
	}

	for seq.Len > 0 {
		nBits, bits := readBits(seq.Pattern, seq.Len, 8)

		// clear the left most 'nBits' from the pattern
		left := 64 - (seq.Len - nBits)
		seq.Pattern = (seq.Pattern << left) >> left
		seq.Len -= nBits

		if nBits < 8 {
			// at the end of th sequence but didn't have a full byte of data
			this.carry = bits << (8 - nBits)
			this.carryN = nBits
		} else {
			_, err := this.w.Write([]byte{bits})
			if err != nil {
				return n, err
			}
			n += int(nBits)
		}
	}

	return n, nil
}

// Closes the byte seq writer flushing the remaining bits as a single byte into
// the stream.
// Returns int - The number of bits written
// Returns error - nil if okay, error otherwise
func (this *ByteSeqWriter) Flush() (n int, err error) {
	n = int(this.carryN)
	if this.carryN > 0 {
		_, err := this.w.Write([]byte{this.carry})
		if err != nil {
			return 0, err
		}
		this.carry = 0x00
		this.carryN = 0
	}
	return n, nil
}

// read 'bits' off the 'seq' putting them in the lower bit places
// 'n' denotes how many bits were actually read off
// 'b' is where the data is stored
// mutates the ByteSeq.
func readBits(seq uint64, len uint, bits uint) (n uint, b byte) {
	if bits > 8 {
		log.Fatal("Cannot read off more than 8 bits at a time")
	}
	if bits > len {
		b = byte(seq)
		n = len
	} else {
		b = byte(seq >> (len - bits))
		n = bits
	}
	return
}
