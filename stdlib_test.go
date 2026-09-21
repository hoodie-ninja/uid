package uid_test

import (
	"testing"
	"time"
	"uuid"

	"github.com/hoodie-ninja/uid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestV7StrictIsV7(t *testing.T) {
	id := uid.NewV7Strict()
	assert.Exactly(t, uid.Version7, id.Version())
	assert.Exactly(t, uid.Variant9562, id.Variant())
	assert.False(t, id.IsMax())
	assert.False(t, id.IsNil())
	// uuid's counter carries the timestamp past the clock under load, so only a loose bound holds
	assert.WithinDuration(t, time.Now(), id.Time(), time.Minute)
	id2, ok := uid.Parse(id.String())
	assert.True(t, ok)
	assert.Exactly(t, id, id2)
}

func TestV7StrictIsOrdered(t *testing.T) {
	ids := make([]uid.UUID, 4000)
	for i := range ids {
		ids[i] = uid.NewV7Strict()
	}
	for i := 1; i < len(ids); i++ {
		if uid.Compare(ids[i-1], ids[i]) >= 0 {
			t.Fatalf("not strictly increasing at %d: %s >= %s", i, ids[i-1], ids[i])
		}
	}
}

func TestStdlibRoundTrip(t *testing.T) {
	check := func(id uid.UUID) {
		t.Helper()
		assert.Exactly(t, id, uid.FromStdlib(uid.ToStdlib(id)))
		// two independent implementations must agree on the canonical form
		assert.Exactly(t, id.String(), uid.ToStdlib(id).String())
		su, err := uuid.Parse(id.String())
		require.NoError(t, err)
		assert.Exactly(t, id, uid.FromStdlib(su))
	}
	check(uid.Nil())
	check(uid.Max())
	check(mustParse(t, ref4))
	check(mustParse(t, ref7))
	check(uid.NewV4())
	check(uid.NewV7())
}

func TestStdlibParsesUID(t *testing.T) {
	// uid accepts encodings uuid does not; the canonical form is the shared contract
	for _, id := range []uid.UUID{uid.Nil(), uid.Max(), uid.NewV4(), uid.NewV7()} {
		su, err := uuid.Parse(id.String())
		require.NoError(t, err)
		assert.Exactly(t, id.Bytes(), su[:])
	}
}

func mustParse(t *testing.T, s string) uid.UUID {
	t.Helper()
	id, ok := uid.Parse(s)
	require.True(t, ok)
	return id
}
