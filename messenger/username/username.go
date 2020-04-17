package username

import (
	"errors"
)

// Username represents a unique lower-case username
// which may only consist of 'a'-'z', underscores and '0'-'9'
type Username string

const (
	// MinLength defines the minimum username length
	MinLength = 3

	// MaxLength defines the maximum username length
	MaxLength = 30
)

// Errors
var (
	ErrTooShort             = errors.New("username too short")
	ErrTooLong              = errors.New("username too long")
	ErrContainsIllegalChars = errors.New(
		"username contains illegal characters",
	)
	ErrContainsIllegalSequences = errors.New(
		"username contains illegal sequences",
	)
)

// Validate returns an error if the username is invalid
func (u Username) Validate() error {
	switch {
	case u == "":
		return ErrTooShort
	case len(u) < MinLength:
		return ErrTooShort
	case len(u) > MaxLength:
		return ErrTooLong
	}
	for i, c := range u {
		if !isLowerAlpha(c) && !isDigit(c) && c != '_' {
			return ErrContainsIllegalChars
		}
		if i > 0 && c == '_' && u[i-1] == '_' {
			// Sequence of multiple underscores
			return ErrContainsIllegalSequences
		}
	}
	return nil
}

func isLowerAlpha(r rune) bool { return r >= 'a' && r <= 'z' }
func isDigit(r rune) bool      { return r >= '0' && r <= '9' }
