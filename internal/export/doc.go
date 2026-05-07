// Package export serializes environment variable maps from envchain chains
// into shell-consumable formats.
//
// Supported formats:
//
//	FormatShell  — POSIX shell "export KEY='VALUE'" lines, suitable for eval
//	FormatDotenv — Plain "KEY=VALUE" lines, compatible with dotenv tooling
//	FormatJSON   — A JSON object mapping keys to string values
//
// Example usage:
//
//	ex := export.NewExporter(os.Stdout, export.FormatShell)
//	if err := ex.Write(vars); err != nil {
//	    log.Fatal(err)
//	}
package export
