package store

import "testing"

func redisTestStore(t *testing.T) *redisStore {
	rs, err := NewRedisStore()

	if err != nil {
		t.Fatal(err)
	}

	return rs
}

func TestRedis(t *testing.T) {
	performStoreTest(t, redisTestStore(t))
}
