package eventdb

import (
	"testing"
)

func TestDialAndClose(t *testing.T) {
	db, err := Dial("tcp", ":6379")
	if err != nil {
		t.Error(err)
	}

	defer db.Close()
}
