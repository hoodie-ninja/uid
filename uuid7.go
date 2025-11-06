package uid

import (
	"sync"
	"time"
)

// NewV7 constructs a new v7 UUID. Enforces method 3 of monotonicity.
func NewV7() UUID { return make7(tick) }

const scale, m, mf64, slot2ns, ns2slot = 4096, 1_000_000, float64(m), mf64 / float64(scale), float64(scale) / mf64

// Time returns the embedded timestamp of UUID. For non-V7 zero(time.Time) is returned. If you don't pre-check version
// use `.IsZero()` to ensure time is "real".
//
//nolint:mnd // locality of behavior
func (u UUID) Time() time.Time {
	if u.Version() != Version7 {
		return time.Time{}
	}
	// rebuild unix_ts_ms
	ms := int64(u.b[0])<<40 | int64(u.b[1])<<32 | int64(u.b[2])<<24 | int64(u.b[3])<<16 | int64(u.b[4])<<8 | int64(u.b[5])
	ra := uint16(u.b[6]&0x0f)<<8 | // top 4 of rand_a
		uint16(u.b[7]) // bottom 8 of rand_a
	return time.Unix(0, ms*m+unslot(ra))
}

//nolint:mnd,gosec // locality of behavior, falst positive index out of range
func make7(tickFn func() int64) UUID {
	var b [16]byte
	ns := tickFn()
	if ns < 0 {
		panic("v7 UUID does not support time before epoch")
	}
	// set unix_ts_ms
	ms := ns / m
	b[0], b[1], b[2], b[3], b[4], b[5] = byte(ms>>40), byte(ms>>32), byte(ms>>24), byte(ms>>16), byte(ms>>8), byte(ms)
	// set rand_a
	ra := slot(ns)
	b[6] = byte((ra >> 8)) & 0x0f // set top 4 bytes of rand_a
	b[7] = byte(ra)
	// fill rand_b
	_, _ = rng.Read(b[8:]) //nolint:errcheck // never returns errors
	// version, variant
	b[6], b[8] = (b[6]&0x0f)|0x70, (b[8]&0x3f)|0x80
	return UUID{b}
}

func tick() int64 { return now().UnixNano() }

/*
NewV7Strict returns a v7 UUID with guaranteed (beyond RFC method 3) local monotonicity.
You don't need this, if you think you need finer than sub-millisecond precision in IDs, what you really need is a
sequence generator and not more accurate timekeeping.
*/
func NewV7Strict() UUID { return make7(tickBatch) }

//nolint:gochecknoglobals // unexported
var (
	mux      = new(sync.Mutex)
	lastTime = slottedNow()
)

func tickBatch() int64 {
	defer mux.Unlock()
	mux.Lock()
	n := slottedNow()
	for !n.After(lastTime) {
		n = slottedNow()
	}
	lastTime = n
	return lastTime.UnixNano()
}

// returns the Time of t's slot (1/4096 of ms).
func slottedNow() time.Time {
	n := now()
	return time.Unix(0, n.UnixMilli()*m+unslot(slot(n.UnixNano())))
}

// returns ns from a given slot (rand_a).
func unslot(randA uint16) int64 {
	ns := float64(randA) * slot2ns
	return int64(ns)
}

// return slot (rand_a) for a given unixnano (ns).
func slot(ns int64) uint16 {
	ms := ns / m
	nsr := ns - ms*m
	s := float64(nsr) * ns2slot
	return uint16(s)
}
