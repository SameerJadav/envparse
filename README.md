# GoEnvParse

[![Go Reference](https://pkg.go.dev/badge/github.com/SameerJadav/go-envparse.svg)](https://pkg.go.dev/github.com/SameerJadav/go-envparse) [![CI](https://github.com/SameerJadav/go-envparse/actions/workflows/ci.yml/badge.svg)](https://github.com/SameerJadav/go-envparse/actions/workflows/ci.yml)

GoEnvParse is a Go package for parsing environment variables from `.env` files. It provides a simple and efficient way to load environment variables from `.env` file into your Go applications.

## Features

- Parse environment files from any `io.Reader` source
- Handling of quoted values (double quotes, single quotes, and backticks)
- Variable expansion in non-quoted and double-quoted values
- Error reporting with line numbers for invalid syntax

## Parsing Details

- Double-quoted values are unescaped, including unicode characters
- Single-quoted and backtick-quoted values are treated as literal strings
- Variable expansion is performed in non-quoted and double-quoted values
- Does not support multiline values

## Installation

```shell
go get github.com/SameerJadav/go-envparse
```

## Usage

```go
package main

import (
	"log"
	"os"

	"github.com/SameerJadav/go-envparse"
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

## Contributing

Contributions are welcome. Please open an issue or submit a pull request.

## License

GoEnvParse is open-source and available under the [MIT License](./LICENSE).
