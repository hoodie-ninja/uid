//nolint:gochecknoglobals,gochecknoinits // test data
package uid_test

import (
	"embed"
	"encoding/csv"
	"io"
	"testing"

	"github.com/hoodie-ninja/uid"
	"github.com/stretchr/testify/assert"
)

//go:embed samples.csv
var files embed.FS

var sampleData []sampleid

type sampleid struct{ Canonical, B32, B64 string }

func init() {
	csvFile, err := files.Open("samples.csv")
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(csvFile)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		entry := sampleid{row[1], row[2], row[4]}
		sampleData = append(sampleData, entry)
	}
}

const (
	ref4    = "01867b2c-a0dd-459c-98d7-89e545538d6c"
	ref4b32 = "eagdhwlfa3vm4rv4j4vcvhdlmj" // draft-taylor-uuid-ncname-01 appendix A
	ref4b64 = "EAYZ7LKDdWcjXieVFU41sJ"
	ref7    = "0191e843-b452-7ac4-b853-8ee3953a28af"
	ref7b32 = "hagi6qq5ukkwequ4o4oktukfpl"
	ref7b64 = "HAZHoQ7RSrEhTjuOVOiivL"
)

var (
	ref4Bytes = []byte{0x01, 0x86, 0x7b, 0x2c, 0xa0, 0xdd, 0x45, 0x9c, 0x98, 0xd7, 0x89, 0xe5, 0x45, 0x53, 0x8d, 0x6c}
	ref7Bytes = []byte{0x01, 0x91, 0xE8, 0x43, 0xB4, 0x52, 0x7A, 0xC4, 0xB8, 0x53, 0x8E, 0xE3, 0x95, 0x3A, 0x28, 0xAF}
)

func TestParseCompactSamples(t *testing.T) {
	tested := 0
	for _, sample := range sampleData {
		c, ok := uid.Parse(sample.Canonical)
		assert.True(t, ok)
		b32, ok := uid.Parse(sample.B32)
		assert.True(t, ok)
		b64, ok := uid.Parse(sample.B64)
		assert.True(t, ok)
		ids := [3]uid.UUID{c, b32, b64}
		// ensure equivalent parsing
		assert.Exactly(t, ids[0], ids[1])
		assert.Exactly(t, ids[1], ids[2])
		for _, id := range ids {
			assert.Exactly(t, uid.Version4, id.Version())
			assert.Exactly(t, uid.Variant9562, id.Variant())
			assert.Exactly(t, sample.Canonical, id.String())
			assert.Exactly(t, sample.B32, id.Compact32())
			assert.Exactly(t, sample.B64, id.Compact64())
		}
		tested++
	}
	assert.Exactly(t, 1000, tested)
}
