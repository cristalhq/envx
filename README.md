# envx

[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]
[![coverage-img]][coverage-url]
[![version-img]][version-url]

Go environment utils

## Features

* Simple API.
* Dependency-free.
* Clean and tested code.
* Fully compatible with `env` package.

See [GUIDE.md](https://github.com/cristalhq/envx/blob/main/GUIDE.md) for more details.

## Install

Go version 1.17+

```
go get github.com/cristalhq/envx
```

## Example

```go
envs := []string{"ENVX_TIMEOUT=20s"} // or os.Environ()

var d time.Duration
eset := envx.NewEnvSet("ENVX")
eset.Duration(&d, "TIMEOUT", 10*time.Second, "just a timeout")

err := eset.Parse(envs)
if err != nil {
	panic(err)
}

fmt.Println(d)

// Output: 20s
```

Also see examples: [examples_test.go](https://github.com/cristalhq/envx/blob/main/example_test.go).

## Documentation

See [these docs][pkg-url].

## License

[MIT License](LICENSE).

[build-img]: https://github.com/cristalhq/envx/workflows/build/badge.svg
[build-url]: https://github.com/cristalhq/envx/actions
[pkg-img]: https://pkg.go.dev/badge/cristalhq/envx
[pkg-url]: https://pkg.go.dev/github.com/cristalhq/envx
[reportcard-img]: https://goreportcard.com/badge/cristalhq/envx
[reportcard-url]: https://goreportcard.com/report/cristalhq/envx
[coverage-img]: https://codecov.io/gh/cristalhq/envx/branch/main/graph/badge.svg
[coverage-url]: https://codecov.io/gh/cristalhq/envx
[version-img]: https://img.shields.io/github/v/release/cristalhq/envx
[version-url]: https://github.com/cristalhq/envx/releases
