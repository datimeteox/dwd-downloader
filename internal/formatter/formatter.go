package formatter

import (
	"fmt"
	"io"
	"strings"
)

// ExtendedFormatter implements fmt.Formatter with extended conversion symbols 'U' (uppercase) and 'L' (lowercase)
type ExtendedFormatter struct{}

func (f *ExtendedFormatter) Format(w io.Writer, format string, args ...interface{}) {
	// Simplified formatter that handles the custom verbs !U and !L
	// This is a simplified version that doesn't fully implement fmt.Formatter interface
	// but handles the specific use case in the patterns

	if len(args) == 0 {
		io.WriteString(w, format)
		return
	}

	// For the specific patterns used, we know the format string structure
	// We'll handle the custom conversions inline
	result := format
	for i, arg := range args {
		var replacement string
		if i < len(args) {
			switch v := arg.(type) {
			case string:
				replacement = v
			default:
				replacement = fmt.Sprint(v)
			}
		}
		// Replace {param!U} or {param!L} patterns
		result = strings.ReplaceAll(result, fmt.Sprintf("{param!%c}", 'U'), strings.ToUpper(replacement))
		result = strings.ReplaceAll(result, fmt.Sprintf("{param!%c}", 'L'), strings.ToLower(replacement))
	}

	io.WriteString(w, result)
}

// FormatString applies the formatter to a string with the given format
func (f *ExtendedFormatter) FormatString(format string, args ...interface{}) string {
	result := format
	for _, arg := range args {
		s := fmt.Sprint(arg)
		result = strings.ReplaceAll(result, "{param!U}", strings.ToUpper(s))
		result = strings.ReplaceAll(result, "{param!L}", strings.ToLower(s))
	}
	return result
}
