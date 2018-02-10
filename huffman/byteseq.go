package huffman

import (
	"bytes"
	"errors"
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
