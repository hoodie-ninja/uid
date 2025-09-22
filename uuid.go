package uid

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
)

// UUID is a UUID as defined by RFC...
// Underlying array is unexported for immutability. UUID is comparable using `==`.
// The zero value is Nil UUID.
//
//nolint:recvcheck // only unserializers should have (temporarily) have pointers.
type UUID struct{ b [16]byte }

// Version returns u's version.
func (u UUID) Version() Version { return Version(u.b[6] >> 4) } //nolint:mnd // lob

// Variant is u's variant.
func (u UUID) Variant() Variant {
	if u.Version() == VersionMax {
		return VariantMax
	}
	return Variant(u.b[8] >> 6) //nolint:mnd // lob
}

// Bytes returns a copy of u's raw bytes.
func (u UUID) Bytes() []byte { return u.b[:] } // copy

// MarshalBinary implements encoding.BinaryMarshaler. Never returns errors.
func (u UUID) MarshalBinary() ([]byte, error) { return u.b[:], nil }

// UnmarshalBinary implement encoding.BinaryUnmarshaler.
func (u *UUID) UnmarshalBinary(b []byte) error {
	if id, ok := Parse(string(b)); ok {
		*u = id
		return nil
	}
	return errors.New("") //nolint:err113 // non-nil sentinel
}

// String implements fmt.Stringer. Returns canonical RFC-4122 representation.
func (u UUID) String() string {
	buf := make([]byte, 36) //nolint:mnd // lob
	buf[8], buf[13], buf[18], buf[23] = '-', '-', '-', '-'
	hex.Encode(buf[0:8], u.b[0:4])
	hex.Encode(buf[9:13], u.b[4:6])
	hex.Encode(buf[14:18], u.b[6:8])
	hex.Encode(buf[19:23], u.b[8:10])
	hex.Encode(buf[24:], u.b[10:])
	return string(buf)
}

// MarshalText implements encoding.TextMarshaler. Never returns errors.
func (u UUID) MarshalText() ([]byte, error) { return []byte(u.String()), nil }

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *UUID) UnmarshalText(b []byte) error {
	if id, ok := Parse(string(b)); ok {
		*u = id
		return nil
	}
	return errors.New("") //nolint:err113 // non-nil sentinel
}

// MarshalJSON implements encoding/json.Marshaler. Never returns errors.
func (u UUID) MarshalJSON() ([]byte, error) { return []byte(`"` + u.String() + `"`), nil }

// UnmarshalJSON implements encoding/json.Unmarshaler.
func (u *UUID) UnmarshalJSON(b []byte) error {
	if id, ok := Parse(string(b)); ok {
		*u = id
		return nil
	}
	return errors.New("") //nolint:err113 // non-nil sentinel
}

// Compact32 returns NCName Base32 representation.
func (u UUID) Compact32() string {
	b := u.shifted()
	b[15] >>= 1
	//nolint:mnd // lob
	return string(u.Version()+65) + base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b[:])[0:25]
}

// Compact64 returns NCName Base64 representation.
func (u UUID) Compact64() string {
	b := u.shifted()
	b[15] >>= 2
	//nolint:mnd // lob
	return string(u.Version()+65) + base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b[:])[0:21]
}

// Nil constructs a Nil UUID (all 0).
func Nil() UUID { return UUID{bytesNil /*copy*/} }

// IsNil returns true when u is the Nil UUID.
func (u UUID) IsNil() bool { return u.b == bytesNil } // compare to zero array is highly optimized

// Max constructs a Max UUID (all F).
func Max() UUID { return UUID{bytesMax /*copy*/} }

// IsMax returns true when u is the Max UUID.
func (u UUID) IsMax() bool { return u.b == bytesMax }

//nolint:nonamedreturns,mnd // golf, locality of behavior
func (u UUID) shifted() (out [16]byte) {
	ints := [4]uint32{
		binary.BigEndian.Uint32(u.b[0:4]),
		binary.BigEndian.Uint32(u.b[4:8]),
		binary.BigEndian.Uint32(u.b[8:12]),
		binary.BigEndian.Uint32(u.b[12:16]),
	}
	variant := (ints[2] & 0xf0000000) >> 24
	ints[1] = (ints[1] & 0xffff0000) | ((ints[1] & 0x00000fff) << 4) | (ints[2] & 0x0fffffff >> 24)
	ints[2] = (ints[2]&0x00ffffff)<<8 | ints[3]>>24
	ints[3] = (ints[3] << 8) | variant
	binary.BigEndian.PutUint32(out[0:4], ints[0])
	binary.BigEndian.PutUint32(out[4:8], ints[1])
	binary.BigEndian.PutUint32(out[8:12], ints[2])
	binary.BigEndian.PutUint32(out[12:16], ints[3])
	return //nolint:gofumpt // covered by nonamedreturns
}

// Compare is a helper for sorting/deduping by monotonic time. Note: Sorting non-v7 IDs is a design flaw.
func Compare(a, b UUID) int { return bytes.Compare(a.b[:8], b.b[:8]) } // unix_ms_ts and rand_a (monotonic times)
