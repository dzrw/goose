package eventdb

import (
	"fmt"
	"net/http"
	"time"
)

type Message struct {
	Tag       string
	Request   *http.Request
	MoreInfo  map[string]string
	Timestamp int32
}

func NewMessage(tag string, req *http.Request, moreInfo map[string]string) *Message {
	if moreInfo == nil {
		moreInfo = make(map[string]string)
	}

	return &Message{tag, req, moreInfo, int32(time.Now().Unix())}
}

func (self *Message) Tagged() bool {
	return self.Tag != ""
}

func (self *Message) String() string {
	code := "UFO"
	if self.Tagged() {
		code = "HIT"
	}

	req := self.Request
	str := fmt.Sprintf("%s: %s %s from %s",
		code, req.Method, req.URL.Path, req.RemoteAddr)

	return str
}
