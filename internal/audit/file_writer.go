package audit

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileWriter is an io.WriteCloser that appends audit events to a file on disk.
type FileWriter struct {
	f *os.File
}

// OpenFile opens (or creates) the audit log file at path for appending.
// Intermediate directories are created automatically.
func OpenFile(path string) (*FileWriter, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("audit: create log dir: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &FileWriter{f: f}, nil
}

// Write implements io.Writer.
func (fw *FileWriter) Write(p []byte) (int, error) {
	return fw.f.Write(p)
}

// Close closes the underlying file.
func (fw *FileWriter) Close() error {
	return fw.f.Close()
}

// NewFileLogger is a convenience constructor that opens path and returns a
// ready-to-use Logger backed by that file. The caller must Close the returned
// FileWriter when finished.
func NewFileLogger(path string) (*Logger, *FileWriter, error) {
	fw, err := OpenFile(path)
	if err != nil {
		return nil, nil, err
	}
	return New(fw), fw, nil
}
