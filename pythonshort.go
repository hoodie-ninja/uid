package uid

import (
	"math/big"
	"strings"
)

const (
	fiftySeven     = 57
	pythonShortLen = 22
	b57decRef      = "23456789" + "ABCDEFGH" + "JKLMN" + "PQRSTUVWXYZ" + "abcdefghijk" + "mnopqrstuvwxyz"
)

//nolint:gochecknoglobals // wtb const arrays and (u)int128
var (
	b57encRef = [fiftySeven]rune{
		'2', '3', '4', '5', '6', '7', '8', '9', // 8/57
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', // 16/57
		'J', 'K', 'L', 'M', 'N', // 21/57
		'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', // 32/57
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', // 43/57
		'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', // 57/57
	}
	big57 = big.NewInt(fiftySeven)
)

// ToPythonShort returns the Python ShortUUID encoding of u. See https://pypi.org/project/shortuuid.
func ToPythonShort(u UUID) string {
	out, q, r := pythonShortBase(), new(big.Int).SetBytes(u.b[:]), new(big.Int)
	for i := pythonShortLen - 1; i > -1; i-- {
		q.QuoRem(q, big57, r)
		out[i] = b57encRef[r.Int64()]
		if q.Int64() == 0 {
			break
		}
	}
	return string(out[:])
}

// FromPythonShort parses a UUID from Python ShortUUID encoded ps.
func FromPythonShort(ps string) (UUID, bool) {
	ps = strings.TrimSpace(ps)
	if len(ps) != pythonShortLen {
		return UUID{}, false
	}
	if ps == MaxPythonShort {
		return Max(), true
	}
	if ps == NilPythonShort {
		return Nil(), true
	}
	n := new(big.Int)
	for _, r := range ps {
		i := int64(strings.IndexRune(b57decRef, r))
		if i == -1 {
			return UUID{}, false
		}
		n.Mul(n, big57).Add(n, big.NewInt(i))
	}
	out := UUID{}
	n.FillBytes(out.b[:])
	return out, true
}

func pythonShortBase() [22]rune {
	return [22]rune{
		'2', '2', '2', '2', '2', // 5/22
		'2', '2', '2', '2', '2', // 10/22
		'2', '2', '2', '2', '2', // 15/22
		'2', '2', '2', '2', '2', // 20/22
		'2', '2', // 22/22
	}
}
