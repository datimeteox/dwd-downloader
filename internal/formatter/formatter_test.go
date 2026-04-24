package formatter

import (
	"bytes"
	"testing"
)

func TestFormatString(t *testing.T) {
	f := &ExtendedFormatter{}

	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "no placeholders",
			format:   "hello world",
			args:     []interface{}{},
			expected: "hello world",
		},
		{
			name:     "single uppercase placeholder",
			format:   "{param!U}",
			args:     []interface{}{"hello"},
			expected: "HELLO",
		},
		{
			name:     "single lowercase placeholder",
			format:   "{param!L}",
			args:     []interface{}{"HELLO"},
			expected: "hello",
		},
		{
			name:     "mixed case input uppercase",
			format:   "{param!U}",
			args:     []interface{}{"HeLLo"},
			expected: "HELLO",
		},
		{
			name:     "mixed case input lowercase",
			format:   "{param!L}",
			args:     []interface{}{"HeLLo"},
			expected: "hello",
		},
		{
			name:     "multiple placeholders same arg",
			format:   "{param!U}_{param!L}",
			args:     []interface{}{"Test"},
			expected: "TEST_test",
		},
		{
			name:     "numeric value uppercase",
			format:   "{param!U}",
			args:     []interface{}{123},
			expected: "123",
		},
		{
			name:     "numeric value lowercase",
			format:   "{param!L}",
			args:     []interface{}{123},
			expected: "123",
		},
		{
			name:     "empty string",
			format:   "{param!U}",
			args:     []interface{}{""},
			expected: "",
		},
		{
			name:     "prefix and suffix",
			format:   "prefix-{param!U}-suffix",
			args:     []interface{}{"middle"},
			expected: "prefix-MIDDLE-suffix",
		},
		{
			name:     "multiple args both placeholders",
			format:   "{param!U}-{param!L}",
			args:     []interface{}{"First", "Second"},
			expected: "FIRST-first",
		},
		{
			name:     "empty format",
			format:   "",
			args:     []interface{}{"hello"},
			expected: "",
		},
		{
			name:     "only uppercase in format",
			format:   "value={param!U}",
			args:     []interface{}{"test"},
			expected: "value=TEST",
		},
		{
			name:     "only lowercase in format",
			format:   "value={param!L}",
			args:     []interface{}{"TEST"},
			expected: "value=test",
		},
		{
			name:     "both placeholders with one arg",
			format:   "{param!U}-{param!L}",
			args:     []interface{}{"Mixed"},
			expected: "MIXED-mixed",
		},
		{
			name:     "special characters",
			format:   "{param!U}",
			args:     []interface{}{"test@#$%"},
			expected: "TEST@#$%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := f.FormatString(tt.format, tt.args...)
			if result != tt.expected {
				t.Errorf("FormatString(%q, %v) = %q, want %q", tt.format, tt.args, result, tt.expected)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	f := &ExtendedFormatter{}

	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "empty args",
			format:   "hello",
			args:     []interface{}{},
			expected: "hello",
		},
		{
			name:     "with args",
			format:   "{param!U}",
			args:     []interface{}{"test"},
			expected: "TEST",
		},
		{
			name:     "complex",
			format:   "prefix_{param!U}_suffix",
			args:     []interface{}{"middle"},
			expected: "prefix_MIDDLE_suffix",
		},
		{
			name:     "no args with format",
			format:   "https://example.com/path",
			args:     []interface{}{},
			expected: "https://example.com/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			f.Format(&buf, tt.format, tt.args...)
			result := buf.String()
			if result != tt.expected {
				t.Errorf("Format(%q, %v) = %q, want %q", tt.format, tt.args, result, tt.expected)
			}
		})
	}
}

func TestFormatStringEdgeCases(t *testing.T) {
	f := &ExtendedFormatter{}

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "nil pointer receiver",
			testFunc: func(t *testing.T) {
				// Should not panic with nil receiver
				var f *ExtendedFormatter
				result := f.FormatString("test")
				if result != "test" {
					t.Errorf("expected 'test', got %q", result)
				}
			},
		},
		{
			name: "unicode uppercase",
			testFunc: func(t *testing.T) {
				result := f.FormatString("{param!U}", "hello 世界")
				// Note: Go's strings.ToUpper may not handle all unicode perfectly
				if result == "" {
					t.Error("should not produce empty result for unicode")
				}
			},
		},
		{
			name: "unicode lowercase",
			testFunc: func(t *testing.T) {
				result := f.FormatString("{param!L}", "HELLO 世界")
				if result == "" {
					t.Error("should not produce empty result for unicode")
				}
			},
		},
		{
			name: "multiple args replaces all placeholders",
			testFunc: func(t *testing.T) {
				// Each arg will replace both {param!U} and {param!L}
				result := f.FormatString("{param!U}-{param!L}", "First", "Second")
				// First arg: First -> {param!U} = FIRST, {param!L} = first
				// After first arg: "FIRST-first"
				// Second arg: Second -> tries to replace {param!U} and {param!L} but they're already replaced
				// So final: FIRST-first (placeholders consumed by first arg)
				expected := "FIRST-first"
				if result != expected {
					t.Errorf("expected %q, got %q", expected, result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}
