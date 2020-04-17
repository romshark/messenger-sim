package passhash_test

import (
	"simulator/messenger/passhash"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBcrypt(t *testing.T) {
	var h PasswordCompareHasher = passhash.NewBcrypt()
	input := []byte("some password")
	hash, err := h.Hash(input)
	require.NoError(t, err)
	require.NotNil(t, hash)

	ok, err := h.Compare(input, hash)
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = h.Compare([]byte("wrong password"), hash)
	require.NoError(t, err)
	require.False(t, ok)
}

func TestBcryptEmpty(t *testing.T) {
	var h PasswordCompareHasher = passhash.NewBcrypt()
	input := []byte{}
	hash, err := h.Hash(input)
	require.Error(t, err)
	require.Nil(t, hash)
}

type PasswordCompareHasher interface {
	passhash.PasswordHasher
	passhash.PasswordComparer
}
