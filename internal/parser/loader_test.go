package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	tmp := t.TempDir()
	path := filepath.Join(tmp, ".env")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestLoadFile_Success(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	ef, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(ef.Entries))
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := LoadFile("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFiles_MultipleFiles(t *testing.T) {
	p1 := writeTempEnv(t, "A=1\n")
	p2 := writeTempEnv(t, "B=2\n")
	result, err := LoadFiles([]string{p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 files, got %d", len(result))
	}
}

func TestEnvFile_Get(t *testing.T) {
	path := writeTempEnv(t, "MY_KEY=hello\n")
	ef, _ := LoadFile(path)

	val, ok := ef.Get("MY_KEY")
	if !ok || val != "hello" {
		t.Errorf("expected MY_KEY=hello, got %q ok=%v", val, ok)
	}

	_, ok = ef.Get("MISSING")
	if ok {
		t.Error("expected MISSING key to not be found")
	}
}

func TestEnvFile_Keys(t *testing.T) {
	path := writeTempEnv(t, "X=1\nY=2\nZ=3\n")
	ef, _ := LoadFile(path)
	keys := ef.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "X" || keys[1] != "Y" || keys[2] != "Z" {
		t.Errorf("unexpected key order: %v", keys)
	}
}
