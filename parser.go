package envparse

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Parse reads an env file from io.Reader, returning a map of keys-value pairs and an error.
//
// Only double-quoted values are escaped. Single-quoted and backquoted values
// are treated as literal strings. Variable expansion (${...} and $...)
// is performed in non-quoted and double-quoted values.
//
// Note: This function does not support multiline values.
// Each key-value pair must be on a single line.
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
