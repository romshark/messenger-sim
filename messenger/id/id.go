package id

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

// ID represents a universally unique
// lexicographically sortable identifier
type ID struct{ ulid.ULID }

// New creates a new unique ID
func New() (ID, error) {
	now := time.Now().UTC()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(now.UnixNano())), 0)
	id, err := ulid.New(ulid.Timestamp(now), entropy)
	if err != nil {
		return ID{}, err
	}
	return ID{ULID: id}, nil
}

// FromString tries to parse the ID from a string
func FromString(s string) (ID, error) {
	id, err := ulid.Parse(s)
	return ID{id}, err
}
