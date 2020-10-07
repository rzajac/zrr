# zrr

Package `zrr` gives ability to add key value pair metadata to errors. 

# Installing

```
go get github.com/rzajac/zrr
```

# Adding metadata

```
err := zrr.Wrap(errors.New("std error")).
    Code("ECode").
    Str("str", "here").
    Int("int", 5).
    Float64("float64", 1.23).
    Time("time", time.Date(2020, time.October, 7, 23, 47, 0, 0, time.UTC)).
    Bool("bool", true)

fmt.Println(err.Error()) // std error :: bool=true code="ECode" float64=1.23 int=5 str="here" time=2020-10-07T23:47:00Z
```

# Inspecting metadata 

```
err := zrr.New("message", "ECode").Int("retry", 5)

fmt.Println(zrr.GetInt(err, "retry"))   // 5 true
fmt.Println(zrr.GetInt(err, "not_set")) // 0 false
fmt.Println(zrr.HasKey(err, "retry"))   // true
fmt.Println(zrr.HasKey(err, "not_set")) // false
```

# License

BSD-2-Clause