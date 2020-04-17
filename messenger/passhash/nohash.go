package passhash

import "bytes"

var _ PasswordHasher = NoHash{}

// NoHash is a non-hashing pass-through mock
type NoHash struct{}

// Hash salts and hashes the given password returning the resulting hash
func (h NoHash) Hash(password []byte) ([]byte, error) {
	cp := make([]byte, len(password))
	copy(cp, password)
	return cp, nil
}

// Compare returns true if the given password corresponds to the given hash,
// otherwise returns false
func (h NoHash) Compare(password, hash []byte) (bool, error) {
	return bytes.Equal(password, hash), nil
}
