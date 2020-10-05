# zrr

Package `zrr` provides errors with metadata as key value pairs.

The `errors.Wrap` function returns a new error to which we can add metadata.
For example:

```
_, err := ioutil.ReadAll(r)
if err != nil {
        return zrr.Wrap(err).Str("origin", "67b75223ad8c2183")
}

// ...

fmt.Println(zrr.Wrap(err).GetStr("origin))
// Output: 67b75223ad8c2183

```

# License

BSD-2-Clause