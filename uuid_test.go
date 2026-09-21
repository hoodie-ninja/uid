package uid_test

import (
	"encoding"
	"encoding/json"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/hoodie-ninja/uid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsts(t *testing.T) {
	assert.Exactly(t, uid.MaxCanonical, strings.ToLower(uid.MaxCanonical))
	assert.Exactly(t, uid.MaxCompact32, strings.ToLower(uid.MaxCompact32))
	assert.Exactly(t, uid.MaxCompact64, strings.ToUpper(uid.MaxCompact64))
	assert.Exactly(t, uid.NilCanonical, strings.ToLower(uid.NilCanonical))
	assert.Exactly(t, uid.NilCompact32, strings.ToLower(uid.NilCompact32))
	assert.Exactly(t, uid.NilCompact64, strings.ToUpper(uid.NilCompact64))
}

func TestCommonAccessors(t *testing.T) {
	// v4 max is a good surrogate since none of it's values are zero AND unshifting has an edge case
	id := uid.Max()
	// check methods
	assert.Exactly(t, uid.VersionMax, id.Version())
	assert.Exactly(t, uid.VariantMax, id.Variant())
	// check marshal-unmarshal equivalences
	var id2 uid.UUID
	// text
	txt, err := id.MarshalText()
	require.NoError(t, err)
	assert.Exactly(t, id.String(), string(txt))
	require.NoError(t, id2.UnmarshalText(txt))
	assert.Exactly(t, id, id2)
	// json
	data, err := id.MarshalJSON()
	require.NoError(t, err)
	require.NoError(t, id2.UnmarshalJSON(data))
	assert.Exactly(t, id, id2)
	// binary
	data, err = id.MarshalBinary()
	require.NoError(t, err)
	assert.Exactly(t, id.Bytes(), data)
	require.NoError(t, id2.UnmarshalBinary(data))
	assert.Exactly(t, id, id2)
	// check compact forms in text unmarshaling
	for _, txt := range []string{id.Compact32(), id.Compact64()} {
		var id2 uid.UUID
		// txt
		err := id2.UnmarshalText([]byte(txt))
		require.NoError(t, err)
		assert.Exactly(t, id, id2)
		// json
		data, err := json.Marshal(txt)
		require.NoError(t, err)
		err = json.Unmarshal(data, &id2)
		require.NoError(t, err)
		assert.Exactly(t, id, id2)
	}
}

func TestBytesImmutable(t *testing.T) {
	id := uid.Max()   // use max for non-zero values
	id.Bytes()[0] = 0 // should be copy
	assert.Exactly(t, id.Bytes(), uid.Max().Bytes())
}

// AppendText exists to satisfy encoding.TextAppender; nothing else pins its signature.
var _ encoding.TextAppender = uid.UUID{}

func TestAppendText(t *testing.T) {
	id := uid.Max() // non-zero in every position

	appended, err := id.AppendText(nil)
	require.NoError(t, err)
	assert.Exactly(t, uid.MaxCanonical, string(appended))

	// a full slice forces append to move the backing array; the hex destinations must follow it
	full := []byte("abc")
	appended, err = id.AppendText(slices.Clip(full))
	require.NoError(t, err)
	assert.Exactly(t, "abc"+uid.MaxCanonical, string(appended))
	assert.Exactly(t, "abc", string(full), "input must be untouched")

	// the reuse idiom: one buffer, many UUIDs, no growth after the first
	buf := make([]byte, 0, 64)
	for _, id := range []uid.UUID{uid.Nil(), uid.Max(), uid.NewV4(), uid.NewV7()} {
		buf, err = id.AppendText(buf[:0])
		require.NoError(t, err)
		assert.Exactly(t, id.String(), string(buf))
	}
}

// textSink escapes anything stored in it, so AllocsPerRun below measures the returned buffer.
//
//nolint:gochecknoglobals // must outlive the closures to defeat escape analysis
var textSink []byte

func TestAppendTextAllocations(t *testing.T) {
	id := uid.NewV4()

	buf := make([]byte, 0, len(uid.NilCanonical))
	assert.Zero(t, testing.AllocsPerRun(100, func() {
		buf, _ = id.AppendText(buf[:0]) //nolint:errcheck // never errors
	}), "AppendText into a buffer with room must not allocate")

	// one buffer, retained: a String round-trip inside MarshalText would double this
	allocs := testing.AllocsPerRun(100, func() {
		textSink, _ = id.MarshalText() //nolint:errcheck // never errors
	})
	assert.InDelta(t, 1, allocs, 0)
}

func TestUnmarshalBinaryFail(t *testing.T) {
	var id uid.UUID
	require.ErrorIs(t, id.UnmarshalBinary([]byte{}), uid.ErrInvalid)
}

func TestUnmarshalTextFail(t *testing.T) {
	var id uid.UUID
	require.ErrorIs(t, id.UnmarshalText([]byte{}), uid.ErrInvalid)
}

func TestUnmarshalJSONFail(t *testing.T) {
	var id uid.UUID
	require.ErrorIs(t, id.UnmarshalJSON([]byte{}), uid.ErrInvalid)
}

func TestNil(t *testing.T) {
	assert.Exactly(t, uid.VersionNil, uid.Nil().Version()) // version
	assert.True(t, uid.Nil().IsNil())                      // identity
	// binary equivalence
	assert.Exactly(t, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, uid.Nil().Bytes())
	// negatives
	assert.False(t, uid.Nil().IsMax())
	assert.True(t, uid.Nil().Time().IsZero())
	// stringifications
	assert.Exactly(t, uid.NilCanonical, uid.Nil().String())
	assert.Exactly(t, uid.NilCompact32, uid.Nil().Compact32())
	assert.Exactly(t, uid.NilCompact64, uid.Nil().Compact64())
	actual, err := uid.Nil().MarshalJSON()
	require.NoError(t, err)
	assert.Exactly(t, `"`+uid.NilCanonical+`"`, string(actual))
}

func TestMax(t *testing.T) {
	assert.Exactly(t, uid.VersionMax, uid.Max().Version()) // version
	// identity
	assert.True(t, uid.Max().IsMax())
	assert.False(t, uid.Max().IsNil())
	assert.True(t, uid.Max().Time().IsZero())
	// binary equivalence
	expected := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	assert.Exactly(t, expected, uid.Max().Bytes())
	// stringifications
	assert.Exactly(t, uid.MaxCanonical, uid.Max().String())
	assert.Exactly(t, uid.MaxCompact32, uid.Max().Compact32())
	assert.Exactly(t, uid.MaxCompact64, uid.Max().Compact64())
	actual, err := uid.Max().MarshalJSON()
	require.NoError(t, err)
	assert.Exactly(t, `"`+uid.MaxCanonical+`"`, string(actual))
}

func TestCompare(t *testing.T) {
	assert.Exactly(t, 0, uid.Compare(uid.Nil(), uid.Nil()))
	assert.Exactly(t, 0, uid.Compare(uid.Max(), uid.Max()))
	id1 := uid.NewV7()
	time.Sleep(time.Microsecond)
	id2 := uid.NewV7()
	assert.Exactly(t, -1, uid.Compare(id1, id2))
}

func TestSanityCollisions(t *testing.T) {
	const count = 2 * 1000 * 1000 // fill 2,000,000 ids to ensure we don't have a glaring issue
	ids := map[uid.UUID]bool{}    // doubles as check that uid.UUID can be used as map key
	// v4
	for range count {
		id := uid.NewV4()
		_, ok := ids[id]
		assert.False(t, ok)
		ids[id] = true
	}
	// v7 in the same set
	for range count {
		id := uid.NewV7()
		_, ok := ids[id]
		assert.False(t, ok)
		ids[id] = true
	}
}
