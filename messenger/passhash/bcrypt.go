package passhash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var _ PasswordHasher = Bcrypt{}

// Bcrypt implements the PasswordHasher interface using bcrypt
type Bcrypt struct {
	Cost int
}

// NewBcrypt creates a new bcrypt-based password hasher
func NewBcrypt() Bcrypt {
	return Bcrypt{
		Cost: bcrypt.DefaultCost,
	}
}

// Hash salts and hashes the given password returning the resulting hash
func (h Bcrypt) Hash(password []byte) ([]byte, error) {
	if len(password) < 1 {
		return nil, errors.New("invalid password input (empty)")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.Cost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// Compare returns true if the given password corresponds to the given hash,
// otherwise returns false
func (h Bcrypt) Compare(password, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, password)
	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return false, nil
	case err != nil:
		return false, err
	}
	return true, nil
}
