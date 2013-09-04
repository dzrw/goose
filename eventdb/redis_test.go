package eventdb

import (
	"testing"
)

func TestRedisEventProvider(t *testing.T) {
	testData := []int{0, 1, 15}

	for _, db := range testData {
		ep := RedisEventProvider("tcp", ":6379", db)
		err := ep.Dial()

		switch {
		case err == ErrCouldNotSelectDb:
			ep.Close()
			t.Error(err)
			return
		case err != nil:
			t.Error(err)
			return
		default:
			ep.Close()
		}
	}
}

func TestRedisWithBadDb(t *testing.T) {
	testData := []int{-1, 16}

	for _, db := range testData {
		ep := RedisEventProvider("tcp", ":6379", db)
		err := ep.Dial()
		if err == nil {
			ep.Close()
			t.Errorf("should not have been able to select db %d", db)
			return
		}
	}
}
