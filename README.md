# Yet another UUID library!?

This library exists to free UUIDs from "Too Much Crypto" https://eprint.iacr.org/2019/1492.pdf and enhance usage across
the many places where UUIDs are commonly found (JSON, SQL, etc).

This library agrees with stdlib's opionion about what UUIDs are worthwhile (v4 and v7), and that creating UUIDs should
never fail (no constructor errors). However, even the standard library returns useless parse errors as sentinels where
simple `false` would force developers to actuall handle instead of relay the sentinel.

## But the crypto!

This library follows Go's `math/rand/v2` and Linux's `/dev/random` changes to use ChaCha20-based cryptographic
pseudorandom number generators to ensure error-free generation and speed. Randomness is drawn from `math/rand/v2`'s
lock-free, per-CPU ChaCha8 generators, runtime-seeded from OS entropy. UUIDs are not cryptographic keys or secrets.

## But the errors!

Errors returned from unmarshalling functions are the constant sentinel `ErrInvalid`. With no dynamic content to
translate or sanitize it is functionally boolean: `nil` or not.

Boolean success and sentinel error returns free (require) you to handle parsing/unmarshalling failures your way.

```go
id, ok := uid.Parse(r.PathValue("id"))
if !ok {
    // observe it your way
    slog.Warn("bad ID", "id", sanitizeForLog(input))
    badIDCounter.Inc()
    // translate responses your way
    http.Error(w, messagePrinter.Sprint("invalid ID"), http.StatusBadRequest)
    return
}
```

# How To

New Random UUID (v4)...
```go
id := uid.NewV4()
```

New Sortable UUID (v7 with "method 3", extended precision monotonicity)
```go
id := uid.NewV7()
```

Note: `uid.NewV7Strict` is now deprecated. If you need single-node monotonically sortable ids, use
`uid.FromStdlib(uuid.NewV7())`

## Databases

`UUID` implements `database/sql/driver.Valuer` and `database/sql.Scanner`, it drops straight into
`database/sql` and `pgx` as `uuid` typed columns.

```go
_, err := db.Exec(`INSERT INTO widgets (id) VALUES ($1)`, uid.NewV7())

var id uid.UUID
err = db.QueryRow(`SELECT id FROM widgets WHERE ...`).Scan(&id)
```

### sqlc

To use this type with `sqlc`, map `uuid` columns to `uid.UUID` with an override:

```yaml
overrides:
  - db_type: "uuid"
    go_type: "github.com/hoodie-ninja/uid.UUID"
```

## Short Serializations

The "hex-and-dash" encoding of a canonical UUID is already URL-safe and contains no ambiguous characters. Omitting the
dashes (which are positional anyway) gives you a short (32-runes), case-insensitive, URL-safe identifier string.

Sometimes an even shorter (but still non-binary) string is helpful. `uid` supports Compact UUIDs and ShortUUIDs.

### Compact UUIDs for Constrained Grammars (NCName)

`Parse` supports automatic detection and decoding of `UUID-NCName-32` and `UUID-NCName-64` compact encodings for
constrained grammars.

`UUID.Compact64()` and `UUID.Compact32()` return the Base64 and Base32 NCName encoded values, respectively.

Per the draft, `Compact32()` emits lower-case, and `Parse` accepts any Base32. Base64 is case-sensitive with upper-case
bookends.

These formats achieve or preserve the goals of compaction, URL-safety, and CSS/DOM identifier safety.

More info: https://datatracker.ietf.org/doc/draft-taylor-uuid-ncname/

### ShortUUID support

Python ShortUUID is problematic in multiple ways.

1. The common implementation accepts ANY alphabet (and padding) endangering transferability.
2. The encoding algorithm does not encode standard alphabets using standard mappings. If you "ShortUUID" encode using
the Base64 alphabet, you cannot Base64 decode the result back into the original bytes.
3. Base57 (default alphabet) has no other usage.
4. Optimizing for "manual human entry" is problematic in itself but Base57 still includes the `o` rune. The more
commonly used Base56 omits `o`.
5. Base57 alphabet ShortUUIDs may contain leading digits (often due to left-padding with `2`) making them unsuitable for
DOM and CSS identifiers without escaping.

Despite all this, it's a popular library and you may be interacting with a system that already uses them so the
following helpers

`FromPythonShort` enables decoding of Python ShortUUID encoded UUIDs using the default alphabet (Base57) and padding
(22).

`ToPythonShort` encodes a given `UUID` into a Python ShortUUID using the default alphabet (Base57) and padding (22).

# What about `uuid`? #

This library was originally written to remove the penalties of crypto/rand costs to non-cryptographic ID generation and
the ergonomics of fallible constructors. Now that uuid has been added to Go (1.27) the constructor issue is solved, but
the implementors still leaned on the "too much crypto" solution. Moreover, the purely sentinel error that standard lib
Parse returns is unexported so callers cannot errors.Is against it even if they did want it.

This library will continue to provide an opinionated parser and a richer data-type tool. `ToStdlib` and `FromStdlib`
move values between the two by copy. Use whatever generator you want, import uid to add codec/sql support.

```go
id := uid.FromStdlib(uuid.NewV7()) // stdlib generation, uid features
su := uid.ToStdlib(uid.NewV4())    // uid generation, stdlib API
```

Note: Stdlib v7 monotonic still uses "too much crypto" AND has a slower implementation than I'd like. If there's
interest I can bring back `NewV7Strict()` as `NewV7Monotonic()` with a `crypto/rand`-free implementation that's about
3x faster than stdlib. But you're already wrong if you think you need more than `NewV7()`.
