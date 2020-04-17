package sessid_test

import (
	"fmt"
	"testing"

	"github.com/romshark/messenger-sim/messenger/sessid"

	"github.com/stretchr/testify/require"
)

func TestNewGeneratorErr(t *testing.T) {
	// Too short
	l := makeRange(0, 15)

	// Non-divisible by 4 without remainder
	l = append(l, makeRange(17, 19)...)
	l = append(l, makeRange(21, 23)...)

	for _, i := range l {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			g, err := sessid.NewGenerator(i)
			require.Error(t, err)
			require.Nil(t, g)
		})
	}
}

func TestGenerator(t *testing.T) {
	for _, i := range []uint{
		16, 20, 24, 28, 32,
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			g, err := sessid.NewGenerator(i)
			require.NoError(t, err)

			sid, err := g.New()
			require.NoError(t, err)
			require.Len(t, sid, int(i))
		})
	}
}

// makeRange returns a slice of unsigned integers
// beginning with first and ending with the last values
func makeRange(first, last uint) []uint {
	if first > last {
		panic(fmt.Errorf("invalid range (%d..%d)", first, last))
	}
	s := make([]uint, 0, last-first+1)
	for ; first <= last; first++ {
		s = append(s, first)
	}
	return s
}
