//nolint:errcheck // benchmarks
package uid_test

import (
	"testing"

	"github.com/byron-janrain/uid"
	gofrsuuid "github.com/gofrs/uuid"
	googleuuid "github.com/google/uuid"
)

func BenchmarkV4(b *testing.B) {
	for b.Loop() {
		_ = uid.NewV4()
	}
}

func BenchmarkGoogleV4(b *testing.B) {
	for b.Loop() {
		_, _ = googleuuid.NewRandom() // ignoring error is best-case for performance comparison but don't
	}
}

func BenchmarkGofrsV4(b *testing.B) {
	for b.Loop() {
		_, _ = gofrsuuid.NewV4() // ignoring error is best-case for performance comparison but don't
	}
}

func BenchmarkV7(b *testing.B) {
	for b.Loop() {
		_ = uid.NewV7()
	}
}

func BenchmarkV7Strict(b *testing.B) {
	for b.Loop() {
		_ = uid.NewV7Strict()
	}
}

func BenchmarkGoogleV7(b *testing.B) {
	for b.Loop() {
		_, _ = googleuuid.NewV7() // ignoring error is unrealistic but best-case for performance comparison
	}
}

func BenchmarkGofrsV7(b *testing.B) {
	for b.Loop() {
		_, _ = gofrsuuid.NewV7() // ignoring error is unrealistic but best-case for performance comparison
	}
}

func BenchmarkParseNil(b *testing.B) {
	for b.Loop() {
		_, _ = uid.Parse(uid.NilCanonical)
	}
}

func BenchmarkParseMax(b *testing.B) {
	for b.Loop() {
		_, _ = uid.Parse(uid.MaxCanonical)
	}
}

func BenchmarkParse4(b *testing.B) {
	for b.Loop() {
		_, _ = uid.Parse(ref4)
	}
}

func BenchmarkParseGoogle4(b *testing.B) {
	for b.Loop() {
		_, _ = googleuuid.Parse(ref4)
	}
}

func BenchmarkParseGofrs4(b *testing.B) {
	for b.Loop() {
		_, _ = gofrsuuid.FromString(ref4)
	}
}

func BenchmarkParse7(b *testing.B) {
	for b.Loop() {
		_, _ = uid.Parse(ref7)
	}
}

func BenchmarkParseGoogle7(b *testing.B) {
	for b.Loop() {
		_, _ = googleuuid.Parse(ref7)
	}
}

func BenchmarkParseGofrs7(b *testing.B) {
	for b.Loop() {
		_, _ = gofrsuuid.FromString(ref7)
	}
}
