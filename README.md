# Errors with metadata

[![Go Report Card](https://goreportcard.com/badge/github.com/rzajac/zrr)](https://goreportcard.com/report/github.com/rzajac/zrr)
[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rzajac/zrr)

The package `zrr` provides a way to add and inspect type safe error context.  

```
go get github.com/rzajac/zrr
```

# When is it useful?

The error context might be useful for example when logging errors which were 
created in some deeper parts of your code.   
 
# Examples

```
if err := somepackage.DoStuff(); err != nil {
    err = zrr.Wrap(err).
        Code("ECode").
        Str("str", "here").
        Int("int", 5).
        Float64("float64", 1.23).
        Time("time", time.Date(2020, time.October, 7, 23, 47, 0, 0, time.UTC)).
        Bool("bool", true)
    return err
}

fmt.Println(err.Error()) // std error :: bool=true code="ECode" float64=1.23 int=5 str="here" time=2020-10-07T23:47:00Z
```

For more examples visit [pkg.go.dev](https://pkg.go.dev/mod/github.com/rzajac/zrr).

# Inspecting error metadata 

```
err := zrr.New("message", "ECode").Int("retry", 5)

fmt.Println(zrr.GetInt(err, "retry"))   // 5 true
fmt.Println(zrr.GetInt(err, "not_set")) // 0 false
fmt.Println(zrr.HasKey(err, "retry"))   // true
fmt.Println(zrr.HasKey(err, "not_set")) // false
```

# License

BSD-2-Clause