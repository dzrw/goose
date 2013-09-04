package eventdb

import (
	"log"
	"testing"
)

func TestRandomKey(t *testing.T) {
	key, err := RandomKey("goose:", 16)
	if err != nil {
		t.Error(err)
		return
	}

	log.Println(key)
}
