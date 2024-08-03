# EnvParse

[![Go Reference](https://pkg.go.dev/badge/github.com/SameerJadav/envparse.svg)](https://pkg.go.dev/github.com/SameerJadav/envparse) [![CI](https://github.com/SameerJadav/envparse/actions/workflows/ci.yml/badge.svg)](https://github.com/SameerJadav/envparse/actions/workflows/ci.yml)

EnvParse is a Go package designed for efficiently parsing environment variables from `.env` files. It provides a straightforward and performant way to load environment variables into your Go applications.

## Features

- Parse environment files from any `io.Reader` source
- Parse environment files directly from a file
- Handling of quoted values (double quotes, single quotes, and backticks)
- Variable expansion in non-quoted and double-quoted values
- Error reporting with line numbers for invalid syntax

## Parsing Details

- Double-quoted values are unescaped, including unicode characters
- Single-quoted and backtick-quoted values are treated as literal strings
- Variable expansion is performed in non-quoted and double-quoted values
- Any empty keys, commented lines (lines prefixed with `#`), and invalid lines (lines that are not comments and do not have an `=` sign) will not be parsed
- Inline comments (e.g., `KEY=VALUE#inline comment`) are removed from the value; if a value must contain a `#`, the value must be quoted (e.g., `KEY="VALUE#with hash"`)
- Does not support multiline values

## Installation

```shell
go get github.com/SameerJadav/envparse
```

## Usage

Parse from `io.Reader`

```go
package main

import (
	"log"
	"os"

	"github.com/SameerJadav/envparse"
)

func main() {
	file, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	env, err := envparse.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range env {
		os.Setenv(key, value)
	}
}
```

Parse from File

```go
package main

import (
	"log"
	"os"

	"github.com/SameerJadav/envparse"
)

func main() {
	env, err := envparse.ParseFile(".env")
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range env {
		os.Setenv(key, value)
	}
}
```

## Contributing

Contributions are welcome. Please open an issue or submit a pull request.

## License

GoEnvParse is open-source and available under the [MIT License](./LICENSE).
