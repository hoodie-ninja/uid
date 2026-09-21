package uid

import "uuid"

// ToStdlib returns u as a standard library uuid.UUID.
func ToStdlib(u UUID) uuid.UUID { return u.b }

// FromStdlib returns su as a UUID. The uuid constructors enforce version and variant, so su is taken as valid.
func FromStdlib(su uuid.UUID) UUID { return UUID{su} }

// NewV7Strict returns a v7 UUID from stdlib.
//
// Deprecated: Just use stdlib `uuid.NewV7()` if you want strictness/monotonicity.
func NewV7Strict() UUID { return FromStdlib(uuid.NewV7()) }
