package uid_test

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode"

	"github.com/hoodie-ninja/uid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBytes(t *testing.T) {
	check := func(in, expectedS string, expectedB []byte) {
		actual, ok := uid.Parse(in)
		assert.True(t, ok)
		assert.Exactly(t, expectedS, actual.String())
		assert.Exactly(t, expectedB, actual.Bytes())
	}
	check(string(uid.Nil().Bytes()), uid.NilCanonical, uid.Nil().Bytes())
	check(string(uid.Max().Bytes()), uid.MaxCanonical, uid.Max().Bytes())
	check(string(ref4Bytes), ref4, ref4Bytes)
	check(string(ref7Bytes), ref7, ref7Bytes)
}

func TestParseBytesBad(t *testing.T) {
	checkFail := func(ref []byte) {
		// corrupt in version
		bad := make([]byte, len(ref))
		copy(bad, ref)
		bad[6] = 0x01 // uuid v1 not supported
		id, ok := uid.Parse(string(bad))
		assert.Exactly(t, uid.Nil(), id)
		assert.False(t, ok)
		// bad variant
		bad = make([]byte, len(ref)) // reset
		copy(bad, ref)
		bad[8] = 0x11 // ms variant not supported
		id, ok = uid.Parse(string(bad))
		assert.Exactly(t, uid.Nil(), id)
		assert.False(t, ok)
	}
	checkFail(ref4Bytes)
	checkFail(ref7Bytes)
	checkFail(uid.Nil().Bytes())
	checkFail(uid.Max().Bytes())
}

func checkString2Bytes(t *testing.T, s string, b []byte) {
	t.Helper()
	actual, ok := uid.Parse(s)
	assert.True(t, ok)
	assert.Exactly(t, b, actual.Bytes())
}

func TestParseCanonical(t *testing.T) {
	checkString2Bytes(t, uid.NilCanonical, uid.Nil().Bytes())
	checkString2Bytes(t, uid.MaxCanonical, uid.Max().Bytes())
	checkString2Bytes(t, strings.ToUpper(uid.MaxCanonical), uid.Max().Bytes())
	checkString2Bytes(t, ref4, ref4Bytes)
	checkString2Bytes(t, strings.ToUpper(ref4), ref4Bytes)
	checkString2Bytes(t, ref7, ref7Bytes)
	checkString2Bytes(t, strings.ToUpper(ref7), ref7Bytes)
	rfj, err := json.Marshal(ref4)
	require.NoError(t, err)
	checkString2Bytes(t, string(rfj), ref4Bytes)
}

func assertBadTxt(t *testing.T, bad []rune) {
	t.Helper()
	id, ok := uid.Parse(string(bad))
	assert.Exactly(t, uid.Nil(), id)
	assert.False(t, ok)
}

func TestParseCanonicalBad(t *testing.T) {
	checkFail := func(ref string) {
		// bad sep
		bad := []rune(ref) // set
		bad[8] = '_'       // invalid separator
		assertBadTxt(t, bad)
		// bad version
		bad = []rune(ref) // reset
		bad[14] = '1'     // uuid v1 not supported
		assertBadTxt(t, bad)
		// bad variant
		bad = []rune(ref) // reset
		bad[19] = 'C'     // ms variant not supported
		assertBadTxt(t, bad)
		// bad hex
		bad = []rune(ref) // reset
		bad[0] = 'g'      // invalid hex rune
		assertBadTxt(t, bad)
	}
	checkFail(ref4)
	checkFail(ref7)
	checkFail(uid.NilCanonical)
	checkFail(uid.MaxCanonical)
}

func TestParseCompact32(t *testing.T) {
	checkString2Bytes(t, uid.NilCompact32, uid.Nil().Bytes())
	checkString2Bytes(t, strings.ToUpper(uid.NilCompact32), uid.Nil().Bytes())
	checkString2Bytes(t, uid.MaxCompact32, uid.Max().Bytes())
	checkString2Bytes(t, strings.ToUpper(uid.MaxCompact32), uid.Max().Bytes())
	checkString2Bytes(t, ref4b32, ref4Bytes)
	checkString2Bytes(t, strings.ToUpper(ref4b32), ref4Bytes)
	checkString2Bytes(t, ref7b32, ref7Bytes)
	checkString2Bytes(t, strings.ToUpper(ref7b32), ref7Bytes)
	rfj32, err := json.Marshal(ref4b32)
	require.NoError(t, err)
	checkString2Bytes(t, string(rfj32), ref4Bytes)
}

func TestParseCompact32Bad(t *testing.T) {
	checkFail := func(ref string) {
		// bad version
		bad := []rune(ref) // reset
		bad[0] = 'b'       // uuid v1 not supported
		assertBadTxt(t, bad)
		// bad variant
		bad = []rune(ref) // reset
		bad[25] = 'b'     // ms variant not supported
		assertBadTxt(t, bad)
		// bad b32
		bad = []rune(ref) // reset
		bad[1] = '_'      // invalid b32 rune
		assertBadTxt(t, bad)
	}
	checkFail(ref4b32)
	checkFail(ref7b32)
	checkFail(uid.NilCompact32)
	checkFail(uid.MaxCompact32)
}

func TestParseCompact64(t *testing.T) {
	checkString2Bytes(t, uid.NilCompact64, uid.Nil().Bytes())
	checkString2Bytes(t, uid.MaxCompact64, uid.Max().Bytes())
	checkString2Bytes(t, ref4b64, ref4Bytes)
	checkString2Bytes(t, ref7b64, ref7Bytes)
	rfj64, err := json.Marshal(ref4b64)
	require.NoError(t, err)
	checkString2Bytes(t, string(rfj64), ref4Bytes)
}

func TestParseCompact64Bad(t *testing.T) {
	checkFail := func(ref string) {
		// bad version
		bad := []rune(ref)
		bad[0] = 'B' // uuid v1 not supported
		assertBadTxt(t, bad)
		// bad variant
		bad = []rune(ref)
		bad[21] = 'B' // ms variant not supported
		assertBadTxt(t, bad)
		// bad b64
		bad = []rune(ref) // reset
		bad[1] = '$'      // invalid b64 rune
		assertBadTxt(t, bad)
	}
	checkFail(ref4b64)
	checkFail(ref7b64)
	checkFail(uid.NilCompact64)
	checkFail(uid.MaxCompact64)
}

func TestParseBoundaries(t *testing.T) {
	checkString2Bytes(t, "{"+ref4+"}", ref4Bytes) // ms-style
	assertBadTxt(t, []rune("x"+ref4+"y"))
	assertBadTxt(t, []rune("{"+ref4+`"`)) // mismatched
	assertBadTxt(t, []rune("x"+ref4b32+"y"))
	assertBadTxt(t, []rune("{"+ref4b32+"}")) // braces are canonical-only
	assertBadTxt(t, []rune("x"+ref4b64+"y"))
	assertBadTxt(t, []rune("{"+ref4b64+"}"))
}

func TestParseCompactCRLF(t *testing.T) {
	// base32/base64 decoders silently ignore \r and \n; swallowed chars must not part-decode
	assertBadTxt(t, []rune("E000000000000\r00\r0000J")) // found by FuzzParse
	b32 := []rune(ref4b32)
	b32[5] = '\n'
	assertBadTxt(t, b32)
}

func TestParseCompact32Multibyte(t *testing.T) {
	// 26 bytes but 25 runes ('ſ' uppercases to ASCII 'S'): must fail, not part-decode into an invalid UUID
	assertBadTxt(t, []rune("eſ"+ref4b32[3:]))
}

func TestParseCompact64Multibyte(t *testing.T) {
	// 22 bytes but 21 runes: must fail cleanly, not panic
	assertBadTxt(t, []rune("Eé"+strings.Repeat("A", 18)+"I"))
}

func TestParseCompact64CaseSensitive(t *testing.T) {
	checkFail := func(ref string, b []byte) {
		r := []rune(ref)
		if unicode.IsLower(r[2]) {
			r[2] = unicode.ToUpper(r[2])
		} else {
			r[2] = unicode.ToLower(r[2])
		}
		id, ok := uid.Parse(string(r))
		assert.True(t, ok)
		assert.NotEqual(t, b, id.Bytes())
	}
	checkFail(ref4b64, ref4Bytes)
	checkFail(ref7b64, ref7Bytes)
}

func TestParseBadLen(t *testing.T) { assertBadTxt(t, []rune{}) }
