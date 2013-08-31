package eventdb

import (
	"testing"
)

func TestNopEventProvider(t *testing.T) {
	ep := NopEventProvider()
	err := ep.Dial()
	if err != nil {
		t.Error(err)
		return
	}

	ep.Close()
}

func TestRedisEventProvider(t *testing.T) {
	ep := RedisEventProvider("tcp", ":6379", 0)
	err := ep.Dial()
	if err != nil {
		t.Error(err)
		return
	}

	ep.Close()
}
