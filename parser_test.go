package envparse

import (
	"bytes"
	"os"
	"testing"
)

var expected = map[string]string{
	"SIMPLE_VAR":                          "value",
	"EMPTY":                               "",
	"EMPTY_SINGLE_QUOTES":                 "",
	"EMPTY_DOUBLE_QUOTES":                 "",
	"EMPTY_BACKTICKS":                     "",
	"QUOTED_VAR":                          "value",
	"SINGLE_QUOTED_VAR":                   "value",
	"BACKQUOTED_VAR":                      "value",
	"DOUBLE_QUOTES_SPACED":                "    double quotes    ",
	"SINGLE_QUOTES_SPACED":                "    single quotes    ",
	"BACKQUOTE_SPACED":                    "    back quotes    ",
	"UNQUOTED_WITH_SPACES":                "this has spaces",
	"NESTED_IN_DOUBLE":                    "This is a 'single quote' and a `backtick` inside double quotes",
	"NESTED_IN_SINGLE":                    "This is a \"double quote\" and a `backtick` inside single quotes",
	"NESTED_IN_BACKTICK":                  "This is a \"double quote\" and a 'single quote' inside backticks",
	"ESCAPED_IN_DOUBLE":                   "This has \"escaped double quotes\"",
	"ESCAPED_IN_SINGLE":                   "This has \\'escaped single quotes\\'",
	"MIXED_QUOTES":                        "Double with 'single' and \"escaped double\" and `backtick`",
	"UNQUOTED_NEWLINE":                    "This is a line\\nAnd this is another line",
	"DOUBLE_QUOTED_NEWLINE":               "This is in double quotes\nThis is a new line in double quotes",
	"SINGLE_QUOTED_NEWLINE":               "This is in single quotes\\nThis is a new line in single quotes",
	"BACKTICK_QUOTED_NEWLINE":             "This is in backticks\\nThis is a new line in backticks",
	"VARIABLE_EXPANSION":                  "value/expansion",
	"NESTED_EXPANSION":                    "value/expansion/nested",
	"VARIABLE_EXPANSION_ALT":              "value/expansion",
	"NESTED_EXPANSION_ALT":                "value/expansion/nested",
	"DOUBLEQUOTED_VARIABLE_EXPANSION":     "value/expansion",
	"DOUBLEQUOTED_NESTED_EXPANSION":       "value/expansion/nested",
	"DOUBLEQUOTED_VARIABLE_EXPANSION_ALT": "value/expansion",
	"DOUBLEQUOTED_NESTED_EXPANSION_ALT":   "value/expansion/nested",
	"SINGLEQUOTED_VARIABLE_EXPANSION":     "${SIMPLE_VAR}/expansion",
	"SINGLEQUOTED_NESTED_EXPANSION":       "${VARIABLE_EXPANSION}/nested",
	"SINGLEQUOTED_VARIABLE_EXPANSION_ALT": "$SIMPLE_VAR/expansion",
	"SINGLEQUOTED_NESTED_EXPANSION_ALT":   "$VARIABLE_EXPANSION/nested",
	"BACKQUOTED_VARIABLE_EXPANSION":       "${SIMPLE_VAR}/expansion",
	"BACKQUOTED_NESTED_EXPANSION":         "${VARIABLE_EXPANSION}/nested",
	"BACKQUOTED_VARIABLE_EXPANSION_ALT":   "$SIMPLE_VAR/expansion",
	"BACKQUOTED_NESTED_EXPANSION_ALT":     "$VARIABLE_EXPANSION/nested",
	"UNMATCHED_DOUBLEQUOTE":               "\"value",
	"UNMATCHED_SINGLEQUOTE":               "'value",
	"UNMATCHED_BACKQUOTE":                 "`value",
	"INLINE_COMMENTS":                     "value",
	"INLINE_COMMENTS_DOUBLE_QUOTES":       "inline comments outside of #doublequotes",
	"INLINE_COMMENTS_SINGLE_QUOTES":       "inline comments outside of #singlequotes",
	"INLINE_COMMENTS_BACKQUOTES":          "inline comments outside of #backticks",
	"EXPORTED_VAR":                        "value",
	"EQUAL_SIGNS":                         "equals==",
	"SPACED_KEY":                          "value",
}

func TestParse(t *testing.T) {
	content, err := os.ReadFile("test.env")
	if err != nil {
		t.Fatal(err)
	}

	result, err := Parse(bytes.NewReader(content))
	if err != nil {
		t.Fatalf("Parse returned unexpected error: %v", err)
	}

	compareMaps(t, result, expected)
}

func TestParseFile(t *testing.T) {
	result, err := ParseFile("test.env")
	if err != nil {
		t.Fatalf("Parse returned unexpected error: %v", err)
	}

	compareMaps(t, result, expected)
}

func BenchmarkParse(b *testing.B) {
	content, err := os.ReadFile("test.env")
	if err != nil {
		b.Fatal(err)
	}
	reader := bytes.NewReader(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = Parse(reader)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseFile(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseFile("test.env")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func compareMaps(t *testing.T, result, expected map[string]string) {
	t.Helper()

	if len(result) != len(expected) {
		t.Errorf("Map size mismatch: got %d entries, want %d entries", len(result), len(expected))
	}

	for expectedKey, expectedValue := range expected {
		actualValue, ok := result[expectedKey]
		if !ok {
			t.Errorf("Missing key: %s is not present in the result map, but it was expected", expectedKey)
		}
		if actualValue != expectedValue {
			t.Errorf("Value mismatch for key %s: got %s, want %s", expectedKey, actualValue, expectedValue)
		}
	}

	for key, value := range result {
		if _, ok := expected[key]; !ok {
			t.Errorf("Unexpected key-value pair: %s=%s is present in the result map, but it was not expected.", key, value)
		}
	}
}
