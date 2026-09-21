/*
Package uid generates and parses RFC 9562 v4 and v7 UUIDs without returning errors.

Constructors cannot fail: randomness comes from a ChaCha8 CSPRNG seeded once from crypto/rand. Parse detects and
decodes canonical (plain, json quoted, or ms-style braced), raw byte, and NCName compact (Base32/Base64) encodings.
Python ShortUUID conversion is provided for interoperability.

UUID V7 uses Method 3 (Replace Leftmost Random Bits with Increased Clock Precision) to implement single-node
monotonicity.
*/
package uid

const (
	// MaxCanonical is the canonical RFC9562 "Max" UUID.
	MaxCanonical = "ffffffff-ffff-ffff-ffff-ffffffffffff"

	// MaxCompact32 is the canonical NCName Compact Base32 "Max" UUID.
	MaxCompact32 = "p777777777777777777777777p"

	// MaxCompact64 is the canonical NCName Compact Base64 "Max" UUID.
	MaxCompact64 = "P____________________P"

	// MaxPythonShort is the canonical "Max" Python ShortUUID.
	MaxPythonShort = "oZEq7ovRbLq6UnGMPwc8B5"

	// NilCanonical is the canonical RFC9562 "Nil" UUID.
	NilCanonical = "00000000-0000-0000-0000-000000000000"

	// NilCompact32 is the canonical NCName Compact Base32 "Nil" UUID.
	NilCompact32 = "aaaaaaaaaaaaaaaaaaaaaaaaaa"

	// NilCompact64 is the canonical NCName Compact Base64 "Nil" UUID.
	NilCompact64 = "AAAAAAAAAAAAAAAAAAAAAA"

	// NilPythonShort is the canonical "Nil" Python ShortUUID.
	NilPythonShort = "2222222222222222222222"
)

// Version is the RFC9562 UUID Version.
type Version byte

const (
	// VersionNil is the Nil UUID version.
	VersionNil = Version(0b0000_0000)

	// Version4 is the version of random UUIDs.
	Version4 = Version(0b0000_0100)

	// Version7 is the version of time-sortable UUIDs.
	Version7 = Version(0b0000_0111)

	// VersionMax is the Max UUID version.
	VersionMax = Version(0b0000_1111)

	versionBad = Version(0b1111_1111) // sentinel version to identify malformed UUIDs.
)

// Variant is the RFC9562 variant.
type Variant byte

const (
	// Variant9562 is the value of the variant bits of v4 or v7.
	Variant9562 = Variant(0b0000_0010)

	// VariantNil is the value of the variant bits of a Nil UUID.
	VariantNil = Variant(0b0000_0000)

	// VariantMax is the value of the variant bits of a Max UUID.
	VariantMax = Variant(0b0000_0111)
)

//nolint:gochecknoglobals // wtb const arrays
var (
	bytesNil = [16]byte{}
	bytesMax = [16]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)
