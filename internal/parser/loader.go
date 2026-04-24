package parser

import (
	"fmt"
	"os"
	"path/filepath"
)

// LoadFile opens a .env file at the given path and parses it.
func LoadFile(path string) (*EnvFile, error) {
	clean := filepath.Clean(path)
	f, err := os.Open(clean)
	if err != nil {
		return nil, fmt.Errorf("could not open file %q: %w", clean, err)
	}
	defer f.Close()

	ef, err := Parse(f)
	if err != nil {
		return nil, fmt.Errorf("could not parse file %q: %w", clean, err)
	}

	return ef, nil
}

// LoadFiles loads multiple .env files and returns them keyed by path.
func LoadFiles(paths []string) (map[string]*EnvFile, error) {
	result := make(map[string]*EnvFile, len(paths))
	for _, p := range paths {
		ef, err := LoadFile(p)
		if err != nil {
			return nil, err
		}
		result[p] = ef
	}
	return result, nil
}

// Keys returns all keys present in the EnvFile in order.
func (ef *EnvFile) Keys() []string {
	keys := make([]string, 0, len(ef.Entries))
	for _, e := range ef.Entries {
		keys = append(keys, e.Key)
	}
	return keys
}

// Get returns the value for a key, or empty string if not found.
func (ef *EnvFile) Get(key string) (string, bool) {
	e, ok := ef.Index[key]
	if !ok {
		return "", false
	}
	return e.Value, true
}
