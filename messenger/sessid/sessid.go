// Package sessid provides a random cryptographically secure
// session identifier generator
package sessid

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// SessionID represents a cryptographically secure session identifier
type SessionID string

// Generator represents a non-thread-safe
// cryptographically secure session id generator
type Generator struct {
	buf []byte
}

// NewGenerator creates a new session id generator instance
func NewGenerator(idLength uint) (*Generator, error) {
	if idLength < MinIDLength {
		return nil, fmt.Errorf(
			"id length too short (%d/%d)",
			idLength, MinIDLength,
		)
	}

	if !isDivisibleBy4WithoutRemainder(idLength) {
		return nil, fmt.Errorf("invalid identifier length (%d) "+
			", length must be divisible by 4 without remainder", idLength)
	}
	g := &Generator{
		buf: make([]byte, base64SrcBufLen(idLength)),
	}

	// Ensure that a cryptographically secure PRNG is available
	if _, err := io.ReadFull(rand.Reader, g.buf); err != nil {
		return nil, fmt.Errorf(
			"crypto/rand is unavailable: Read() failed with %w", err,
		)
	}
	return g, nil
}

// New generates a new session identifier
//
// WARNING: this method isn't thread safe!
func (g *Generator) New() (SessionID, error) {
	if _, err := rand.Read(g.buf); err != nil {
		return "", fmt.Errorf("reading rand: %w", err)
	}
	return SessionID(base64.StdEncoding.EncodeToString(g.buf)), nil
}

// MinIDLength defines the minimum identifier length
const MinIDLength = 16

func isDivisibleBy4WithoutRemainder(i uint) bool { return !(i%4 > 0) }

func base64SrcBufLen(base64IDLength uint) uint { return base64IDLength * 6 / 8 }
