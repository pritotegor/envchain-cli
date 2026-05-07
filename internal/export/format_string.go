package export

import "fmt"

// ParseFormat converts a string to a Format constant.
// Returns an error if the string does not match a known format.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatShell, FormatDotenv, FormatJSON:
		return Format(s), nil
	default:
		return "", fmt.Errorf("export: unknown format %q; valid values are shell, dotenv, json", s)
	}
}

// String returns the string representation of the Format.
func (f Format) String() string {
	return string(f)
}

// KnownFormats returns all supported format names as a slice.
func KnownFormats() []string {
	return []string{
		string(FormatShell),
		string(FormatDotenv),
		string(FormatJSON),
	}
}
