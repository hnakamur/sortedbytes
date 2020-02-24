keybytes
========

Package keybytes provides feature for encoding and decoding key bytes.
Two encoded key bytes of two different keys keeps the order
so that you can use encoded key bytes for keys in a key-value-store
which is capable to do range scans.

Supported types are
bool, int32, int64, float64, string,
sql.NulBool, sql.NullInt32, sql.NullInt64, sql.NullFloat64, and sql.NullString.

Note time.Time and sql.NullTime are not supported.
You can use int64 or sql.NullInt64 for timestamps with time.Time.Unix() or
time.Time.UnixNano().

The encoding in this package is a subset of the FDB Tuple layer typecodes encoding.
https://github.com/apple/foundationdb/blob/92b41e3562e639e16dbe0142cc479a3304e9c08a/design/tuple.md
https://activesphere.com/blog/2018/08/17/order-preserving-serialization
