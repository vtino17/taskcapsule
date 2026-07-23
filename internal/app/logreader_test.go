package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestLog(t *testing.T, dir, name string, lines int) string {
	t.Helper()
	path := filepath.Join(dir, name+".log")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for i := 0; i < lines; i++ {
		_, err := f.WriteString("line " + strings.Repeat("x", 64) + "\n")
		if err != nil {
			t.Fatal(err)
		}
	}
	return path
}

func TestReadTailSmallFile(t *testing.T) {
	dir := t.TempDir()
	writeTestLog(t, dir, "test", 10)

	reader := &LogReader{MaxBytes: 256 * 1024, MaxLines: 200}
	data, err := reader.ReadTail(filepath.Join(dir, "test.log"))
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Count(string(data), "\n")
	if lines != 10 {
		t.Errorf("expected 10 lines, got %d", lines)
	}
}

func TestReadTailTruncateLines(t *testing.T) {
	dir := t.TempDir()
	writeTestLog(t, dir, "test", 500)

	reader := &LogReader{MaxBytes: 256 * 1024, MaxLines: 10}
	data, err := reader.ReadTail(filepath.Join(dir, "test.log"))
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Count(string(data), "\n")
	if lines > 11 {
		t.Errorf("expected at most 11 lines, got %d", lines)
	}
}

func TestReadTailLargeByteLimit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "big.log")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	// Write 1MB of data
	for i := 0; i < 10240; i++ {
		f.WriteString(strings.Repeat("x", 100) + "\n")
	}
	f.Close()

	reader := &LogReader{MaxBytes: 1024, MaxLines: 200}
	data, err := reader.ReadTail(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) > 2048 {
		t.Errorf("expected data under 2048 bytes, got %d", len(data))
	}
}

func TestReadTailEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.log")
	os.WriteFile(path, []byte{}, 0644)

	reader := DefaultLogReader()
	data, err := reader.ReadTail(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 0 {
		t.Errorf("expected empty data, got %d bytes", len(data))
	}
}

func TestReadTailNoTrailingNewline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nolf.log")
	err := os.WriteFile(path, []byte("line without newline"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	reader := DefaultLogReader()
	data, err := reader.ReadTail(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("expected data")
	}
}

func TestReadTailMissingFile(t *testing.T) {
	reader := DefaultLogReader()
	_, err := reader.ReadTail("/nonexistent/path.log")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
