package uid

import (
	"encoding/binary"
	"math"
	"time"
)

// NewV7 constructs a new v7 UUID. Enforces method 3 of monotonicity.
func NewV7() UUID { return make7() }

const m = 1_000_000 // ns per ms

// Time returns the embedded timestamp of UUID. For non-V7, or timestamps beyond the unix nanosecond range
// (year 2262+), zero(time.Time) is returned. If you don't pre-check version use `.IsZero()` to ensure time is "real".
//
//nolint:mnd // locality of behavior
func (u UUID) Time() time.Time {
	if u.Version() != Version7 {
		return time.Time{}
	}
	// rebuild unix_ts_ms
	ms := int64(u.b[0])<<40 | int64(u.b[1])<<32 | int64(u.b[2])<<24 | int64(u.b[3])<<16 | int64(u.b[4])<<8 | int64(u.b[5])
	if ms > math.MaxInt64/m-1 { // ms*m+unslot(ra) would overflow int64
		return time.Time{}
	}
	ra := uint16(u.b[6]&0x0f)<<8 | // top 4 of rand_a
		uint16(u.b[7]) // bottom 8 of rand_a
	return time.Unix(0, ms*m+unslot(ra))
}

//nolint:mnd,gosec // locality of behavior, false positive index out of range
func make7() UUID {
	var b [16]byte
	ns, ra := tick()
	if ns < 0 {
		panic("v7 UUID does not support time before epoch")
	}
	// set unix_ts_ms
	ms := ns / m
	b[0], b[1], b[2], b[3], b[4], b[5] = byte(ms>>40), byte(ms>>32), byte(ms>>24), byte(ms>>16), byte(ms>>8), byte(ms)
	// set rand_a
	b[6] = byte(ra>>8) & 0x0f // top 4 bits of rand_a
	b[7] = byte(ra)
	// fill rand_b
	binary.LittleEndian.PutUint64(b[8:16], rand64())
	// version, variant
	b[6], b[8] = (b[6]&0x0f)|0x70, (b[8]&0x3f)|0x80
	return UUID{b}
}

func tick() (int64, uint16) {
	ns := now().UnixNano()
	return ns, slot(ns)
}

// unslot returns the first nanosecond of slot randA. Exact inverse of slot: slot(unslot(k))==k for all k.
//
//nolint:mnd // lob
func unslot(randA uint16) int64 { return (int64(randA)*15625 + 63) / 64 }

// slot returns rand_a for a given unixnano: 4096 slots per ms, each 15625/64 (244.140625) ns wide.
//
//nolint:mnd,gosec // lob; result < 4096
func slot(ns int64) uint16 { return uint16(ns % m * 64 / 15625) }
