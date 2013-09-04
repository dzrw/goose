package eventdb

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var ErrUnderflow = errors.New("not enough bits")

func RandomKey(prefix string, c int) (str string, err error) {
	if c < 16 {
		c = 16
	}

	buf := make([]byte, c)
	n, err := io.ReadFull(rand.Reader, buf)

	switch {
	case n != len(buf):
		err = ErrUnderflow
		fallthrough
	case err != nil:
		return
	}

	str, err = EncodeBuffer(buf)
	if err != nil {
		return
	}

	str = fmt.Sprintf("%s%s", prefix, str)
	return
}

func EncodeBuffer(input []byte) (str string, err error) {
	var buf bytes.Buffer

	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write(input)
	encoder.Close()

	str = string(buf.Bytes())
	return
}
