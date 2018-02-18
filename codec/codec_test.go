package codec

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestCodec_Simple(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "test-")
	if err != nil {
		log.Fatal(err)
	}
	// defer os.Remove(tmpfile.Name())
	fmt.Println(tmpfile.Name())

	src := []byte("hello world")
	encoder, err := NewEncoder(tmpfile)
	if err != nil {
		t.Error(err)
	}
	n, err := encoder.Write(src, HAS_MODEL)
	if err != nil {
		t.Error(err)
	}
	if n != len(src) {
		t.Errorf("Failed to write payload, only wrote %d bytes out of %d",
			n, len(src))
	}

	// tmpfile.Seek(0, 0)
	// b, _ := ioutil.ReadAll(tmpfile)
	// for i := 0; i < len(b); i++ {
	// 	fmt.Printf("%#x ", b[i])
	// }
	// fmt.Println()

	tmpfile.Seek(0, 0)
	decoder, err := NewDecoder(tmpfile)
	if err != nil {
		t.Error(err)
	}
	got, err := decoder.Read()
	if err != nil {
		t.Error(err)
	}
	if got == nil {
		t.Errorf("Expected to have the input but got nil")
	}

	if bytes.Compare(src, got) != 0 {
		t.Errorf("payload was not retrieved. got = %v, want = %v\n", got, src)
	}
}
