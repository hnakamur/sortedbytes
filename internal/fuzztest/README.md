fuzztest
========

## Usage

Install `go-fuzz` and `go-fuzz-build`

```
go get -u github.com/dvyukov/go-fuzz/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz-build
```

Build the fuzz tests in this `fuzztest` directory.

```
go-buzz-build
```

Run fuzz tests.

```
go-fuzz -func FuzzTakeString -workdir work/TakeString
```

```
go-fuzz -func FuzzTakeNullString -workdir work/TakeNullString
```

```
go-fuzz -func FuzzTakeInt32 -workdir work/TakeInt32
```

```
go-fuzz -func FuzzTakeNullInt32 -workdir work/TakeNullInt32
```

```
go-fuzz -func FuzzTakeInt64 -workdir work/TakeInt64
```

```
go-fuzz -func FuzzTakeNullInt64 -workdir work/TakeNullInt64
```

```
go-fuzz -func FuzzTakeFloat64 -workdir work/TakeFloat64
```

```
go-fuzz -func FuzzTakeNullFloat64 -workdir work/TakeNullFloat64
```

```
go-fuzz -func FuzzTakeBool -workdir work/TakeBool
```

```
go-fuzz -func FuzzTakeNullBool -workdir work/TakeNullBool
```
