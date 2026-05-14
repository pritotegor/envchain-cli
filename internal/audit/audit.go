// Package audit provides a simple append-only audit log that records
// mutations made to environment chains (add, delete, rename, copy, merge).
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Chain     string    `json:"chain"`
	Key       string    `json:"key,omitempty"`
	Extra     string    `json:"extra,omitempty"`
}

// Logger writes audit entries to a newline-delimited JSON file.
type Logger struct {
	path string
}

// NewLogger creates a Logger that appends to the file at path.
// The parent directory is created if it does not exist.
func NewLogger(path string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("audit: create log dir: %w", err)
	}
	return &Logger{path: path}, nil
}

// Record appends a single entry to the log file.
func (l *Logger) Record(action, chain, key, extra string) error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	entry := Entry{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Chain:     chain,
		Key:       key,
		Extra:     extra,
	}
	if err := json.NewEncoder(f).Encode(entry); err != nil {
		return fmt.Errorf("audit: write entry: %w", err)
	}
	return nil
}

// ReadAll reads and returns all entries from the log file.
// Returns an empty slice if the file does not exist.
func (l *Logger) ReadAll() ([]Entry, error) {
	f, err := os.Open(l.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
