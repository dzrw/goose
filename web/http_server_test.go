package web

import (
	"testing"
)

const ListenAddr = ":9001"

func TestDialNonResponsiveServer(t *testing.T) {
	err := CheckConnection(ListenAddr)
	if err == nil {
		t.Error("should not be able to connect")
	}
}

func TestStartHttpServerAndStop(t *testing.T) {
	s, err := StartHttpServer(ListenAddr, &fakeHandler{})
	if err != nil {
		t.Error(err)
	}

	err = CheckConnection(ListenAddr)
	if err != nil {
		t.Error("should be able to connect: ", err)
	}

	s.Stop()

	err = CheckConnection(ListenAddr)
	if err == nil {
		t.Error("should not be able to connect")
	}
}
