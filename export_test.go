package uid

import (
	"math/rand/v2"
	"time"
)

// ReseedPRNG swaps the PRNG for a zero-seeded ChaCha8 to make output testably predictable.
func ReseedPRNG() func() {
	old := rand64
	rand64 = rand.NewChaCha8([32]byte{}).Uint64
	return func() { rand64 = old }
}

// SetNowFunc replaces the internal time.Now for unit testing returns a deferrable that undoes this change.
func SetNowFunc(f func() time.Time) func() {
	now = f
	return func() {
		now = time.Now
	}
}
