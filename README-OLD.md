## Errors with context

[![Go Report Card](https://goreportcard.com/badge/github.com/rzajac/zrr)](https://goreportcard.com/report/github.com/rzajac/zrr)
[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rzajac/zrr)

The package `zrr` provides a way to add and inspect, in type safe manner, 
context for errors.  

Most importantly `errors.Is`, `errors.As` and `errors.Unwrap` work with 
`zrr.Error` as expected.

## Installation

```
go get github.com/rzajac/zrr
```

## When is it useful?

Imagine somewhere deep in your call tree one of the methods returns an error 
which bubbles up to the top where it can be logged using for example `zerolog`. 
Wouldn't it be great if you could log not only the error message, but its 
context? Especially that some context information might not be available at the 
logging level.

This is where `zrr` might be useful.  

## Package level errors.

By definition package level errors must be immutable. With `zrr` you can 
create package level errors with `Imm` constructor function.

```
var ErrPackageLevel = zrr.Imm("package level error", "ECode")
```

Errors created this way will be immutable, but you still can add 
context to them. When context adding methods are called on immutable error 
it's first cloned, and the keys are added on the cloned instance. The cool
thing is that `errors.Is` still works as expected on the cloned instance:   

```
// Create immutable error.
var ErrPackageLevel = zrr.Imm("package level error", "ECode")

// Somewhere in the code use ErrPackageLevel and add context to it.
err := ErrPackageLevel.Str("path", "/path/to/file").Str("code", "ENewCode")

// Notice the error code has not been changed.
fmt.Println(ErrPackageLevel) 
fmt.Println(err)
fmt.Println(errors.Is(err, ErrPackageLevel))
fmt.Println(zrr.GetCode(err))
fmt.Println(zrr.GetStr(err, "path"))

// Output:
// package level error"
// package level error"
// true
// ENewCode true
// /path/to/file true
```

## Wrapping other errors

You can easily decorate other error instances with context fields and keep 
ability to unwrap them with `errors.As(e1, err)` and test 
them with `errors.Is(e1, err)`. 

```
err := errors.New("some error")

e1 := zrr.Wrap(err).Str("key", "value")

fmt.Println(e1)
fmt.Println(errors.Is(e1, err))
fmt.Println(zrr.GetStr(e1, "key"))

// Output:
// some error
// true
// value true
```

## Wrapping `zrr` errors.

```
err := zrr.New("my error")

e1 := fmt.Errorf("zrr wrapped: %w", err)

fmt.Println(e1)
fmt.Println(errors.Is(e1, err))

// Output:
// zrr wrapped: my error
// true
```

## Key value example

```
// Create an error and add bunch of context fields to it.
err := zrr.Wrap(errors.New("std error")).
    Code("ECode").
    Str("str", "string").
    Int("int", 5).
    Float64("float64", 1.23).
    Time("time", time.Date(2020, time.October, 7, 23, 47, 0, 0, time.UTC)).
    Bool("bool", true)

// Somewhere else (maybe during logging extract the context fields.
iter := err.Fields()
for iter.Next() {
    key, val := iter.Get()
    fmt.Printf("%s = %v\n", key, val)
}

// Output:
// bool = true
// code = ECode
// float64 = 1.23
// int = 5
// str = string
// time = 2020-10-07 23:47:00 +0000 UTC
```

For more examples visit [pkg.go.dev](https://pkg.go.dev/mod/github.com/rzajac/zrr).

## Inspecting error metadata 

```
var err error
err = zrr.New("message", "ECode").Int("retry", 5)

fmt.Println(zrr.GetInt(err, "retry"))   // 5 true
fmt.Println(zrr.GetInt(err, "not_set")) // 0 false
fmt.Println(zrr.HasKey(err, "retry"))   // true
fmt.Println(zrr.HasKey(err, "not_set")) // false

// Output:
// 5 true
// 0 false
// true
// false
```

# License

BSD-2-Clause