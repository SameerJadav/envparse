package envparse

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Parse reads an environment variables file from the provided io.Reader and
// returns a map of key-value pairs. The function also returns an error if
// any issues are encountered during parsing.
func Parse(r io.Reader) (map[string]string, error) {
	result := make(map[string]string)
	scanner := bufio.NewScanner(r)
	var lineNumber int
	var err error

	for scanner.Scan() {
		lineNumber++

		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' || !strings.Contains(line, "=") {
			continue
		}

		line = strings.TrimPrefix(line, "export ")

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		value = strings.TrimSpace(value)

		if quote, start, end, ok := isQuoted(value); ok {
			switch quote {
			case '"':
				value, err = strconv.Unquote(value[start : end+1])
				if err != nil {
					return nil, fmt.Errorf("failed to unquote value at line %d: %w", lineNumber, err)
				}
				value = expandVariables(value, result)
			case '`':
				value, err = strconv.Unquote(value[start : end+1])
				if err != nil {
					return nil, fmt.Errorf("failed to unquote value at line %d: %w", lineNumber, err)
				}
			case '\'':
				value = value[start+1 : end]
			}
		} else {
			if i := strings.IndexByte(value, '#'); i >= 0 {
				value = strings.TrimSpace(value[:i])
			}
			value = expandVariables(value, result)
		}

		result[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	return result, nil
}

// ParseFile reads an environment variables file from the specified filename
// and returns a map of key-value pairs. The function also returns an error
// if any issues are encountered during file reading or parsing.
//
// This function uses the [Parse] function internally to process the file contents.
func ParseFile(filename string) (map[string]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read env file %s: %w", filename, err)
	}
	return Parse(bytes.NewReader(content))
}

func isQuoted(value string) (byte, int, int, bool) {
	if len(value) < 2 {
		return 0, -1, -1, false
	}

	if quote := value[0]; quote == '"' || quote == '\'' || quote == '`' {
		for i := len(value) - 1; i > 0; i-- {
			if value[i] == quote && value[i-1] != '\\' {
				return quote, 0, i, true
			}
		}
		return quote, 0, -1, false
	}

	return 0, -1, -1, false
}

func expandVariables(value string, result map[string]string) string {
	return os.Expand(value, func(key string) string {
		if val, ok := result[key]; ok {
			return val
		}
		return os.Getenv(key)
	})
}
