package store

import (
	"errors"
	"time"
)

var (
	NeverExpire     = time.Time{}
	ErrItemNotFound = errors.New("not found")
)

type Namespace interface {
	// JSON encoded structs
	FindByID(id string, out interface{}) error
	Save(id string, item interface{}, expiration time.Time) error
	Delete(id string) error
}

type Store interface {
	Namespace(name string) Namespace
}

func ExpiresIn(d time.Duration) time.Time {
	return time.Now().Add(d)
}
