package username_test

import (
	"simulator/messenger/username"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	maxLen := makeString('a', username.MaxLength)
	minLen := makeString('a', username.MinLength)

	for _, t1 := range []string{
		minLen,
		maxLen,
		"username",
		"user_name",
	} {
		t.Run(t1, func(t *testing.T) {
			u := username.Username(t1)
			require.NoError(t, u.Validate())
		})
	}
}

func TestNewErr(t *testing.T) {
	for _, t1 := range []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			"empty",
			"",
			username.ErrTooShort,
		},
		{
			"too short",
			makeString('a', username.MinLength-1),
			username.ErrTooShort,
		},
		{
			"too long",
			makeString('a', username.MaxLength+1),
			username.ErrTooLong,
		},
		{
			"contains upper-case characters",
			"UserName",
			username.ErrContainsIllegalChars,
		},
		{
			"contains illegal characters",
			"username$",
			username.ErrContainsIllegalChars,
		},
		{
			"contains illegal characters",
			"user.name",
			username.ErrContainsIllegalChars,
		},
		{
			"contains spaces",
			"user name",
			username.ErrContainsIllegalChars,
		},
		{
			"contains sequences of multiple underscores",
			"user__name",
			username.ErrContainsIllegalSequences,
		},
	} {
		t.Run(t1.name, func(t *testing.T) {
			u := username.Username(t1.input)
			err := u.Validate()
			require.Error(t, err)
			require.Equal(t, t1.expectedError, err)
		})
	}
}

func makeString(r rune, l int) string {
	b := strings.Builder{}
	b.Grow(l)
	for i := 0; i < l; i++ {
		b.WriteRune(r)
	}
	return b.String()
}
