package uid_test

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/byron-janrain/uid"
	"github.com/stretchr/testify/assert"
)

func TestV7(t *testing.T) {
	freezeNow := time.Now()
	defer uid.ReseedPRNG()()
	defer uid.SetNowFunc(func() time.Time { return freezeNow })()
	id := uid.NewV7()
	// identity
	assert.Exactly(t, uid.Version7, id.Version())
	assert.False(t, id.IsMax())
	assert.False(t, id.IsNil())
	assert.False(t, id.Time().IsZero())
	// randomness preserved
	assert.True(t, strings.HasSuffix(id.String(), "-9987-7ece6d368aac"))
	// time
	assert.Exactly(t, freezeNow.UnixMilli(), id.Time().UnixMilli())
	// ensure microsecond is within truncation tolerance
	assert.InDelta(t, freezeNow.UnixMicro(), id.Time().UnixMicro(), 1)
	// ensure microsecond is within rounding tolerance 1/4096th of a ms (245ns) + call stack time
	assert.InDelta(t, freezeNow.UnixNano(), id.Time().UnixNano(), 300)
	id2, ok := uid.Parse(id.String())
	assert.True(t, ok)
	assert.Exactly(t, id, id2)
}

func TestPreEpochPanic(t *testing.T) {
	// check that negative times panic (impossible outside of this unit test)
	defer uid.SetNowFunc(func() time.Time { return time.Unix(-10, 0) })()
	assert.Panics(t, func() {
		_ = uid.NewV7()
	})
}

func TestEpochEdgeCase(t *testing.T) {
	defer uid.SetNowFunc(func() time.Time { return time.Unix(0, 0) })()
	assert.Exactly(t, int64(0), uid.NewV7().Time().UnixNano())
}

func TestSanity(t *testing.T) {
	// create matching lists as fast as we can across 2ms to ensure capture 1 full ms
	ts1, ts2 := []uid.UUID{}, []uid.UUID{}
	for start := time.Now(); time.Since(start) < time.Millisecond; {
		id := uid.NewV7()
		ts1, ts2 = append(ts1, id), append(ts2, id) // fill both arrays instead of cloning later
	}
	mss := map[int64]bool{}
	for _, i := range ts1 {
		mss[i.Time().UnixMilli()] = true
	}
	// verify setup
	assert.NotEmpty(t, ts1)
	assert.NotEmpty(t, ts2)
	assert.Exactly(t, ts1, ts2)
	assert.Len(t, mss, 2) // breaking across ms should only have 2 different ms values
	// test that times were generated in order
	assert.True(t, slices.IsSortedFunc(ts1, uid.Compare))
	// test uuids are unique (includes randomness)
	assert.Len(t, ts1, len(slices.Compact(ts1)))
}

func TestV7StrictIsV7(t *testing.T) {
	freezeNow := time.Now()
	defer uid.ReseedPRNG()()
	defer uid.SetNowFunc(func() time.Time { return freezeNow })()
	id := uid.NewV7Strict()
	// identity
	assert.Exactly(t, uid.Version7, id.Version())
	assert.False(t, id.IsMax())
	assert.False(t, id.IsNil())
	assert.False(t, id.Time().IsZero())
	// randomness preserved
	assert.True(t, strings.HasSuffix(id.String(), "-9987-7ece6d368aac"))
	// time
	assert.Exactly(t, freezeNow.UnixMilli(), id.Time().UnixMilli())
	id2, ok := uid.Parse(id.String())
	assert.True(t, ok)
	assert.Exactly(t, id, id2)
}

func TestSanityBatching(t *testing.T) {
	// create matching lists as fast as we can across 2ms to ensure capture 1 full ms
	ts1, ts2 := []uid.UUID{}, []uid.UUID{}
	for start := time.Now(); time.Since(start) < time.Millisecond; {
		id := uid.NewV7Strict()
		ts1, ts2 = append(ts1, id), append(ts2, id) // fill both arrays instead of cloning later
	}
	mss := map[int64]bool{}
	for _, i := range ts1 {
		mss[i.Time().UnixMilli()] = true
	}
	ts := map[time.Time]bool{}
	for _, i := range ts1 {
		ts[i.Time()] = true
	}
	// verify setup
	assert.NotEmpty(t, ts1)
	assert.NotEmpty(t, ts2)
	assert.Exactly(t, ts1, ts2)
	assert.Len(t, mss, 2) // breaking across ms should only have 2 different ms values
	// test that times were generated in order
	assert.True(t, slices.IsSortedFunc(ts1, uid.Compare))
	// test uuids are unique (includes randomness)
	assert.Len(t, ts1, len(slices.Compact(ts1)))
	// assert times are strictly monotonic
	assert.Len(t, ts1, len(ts))
}
