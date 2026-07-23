package app

import (
	"os"
)

const defaultTailLines = 200
const defaultTailBytes = 256 * 1024

type LogReader struct {
	MaxBytes int64
	MaxLines int
}

func DefaultLogReader() *LogReader {
	return &LogReader{MaxBytes: defaultTailBytes, MaxLines: defaultTailLines}
}

func (r *LogReader) ReadTail(path string) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.Size() <= r.MaxBytes {
		return r.tailLines(path, r.MaxLines)
	}

	return r.tailBytes(path, r.MaxBytes)
}

func (r *LogReader) tailLines(path string, maxLines int) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lineCount := 0
	for _, b := range data {
		if b == '\n' {
			lineCount++
		}
	}

	if lineCount <= maxLines {
		return data, nil
	}

	skipLines := lineCount - maxLines
	pos := 0
	for skipLines > 0 && pos < len(data) {
		if data[pos] == '\n' {
			skipLines--
		}
		pos++
	}
	return data[pos:], nil
}

func (r *LogReader) tailBytes(path string, maxBytes int64) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := info.Size()
	seekPos := fileSize - maxBytes
	if seekPos < 0 {
		seekPos = 0
	}

	buf := make([]byte, fileSize-seekPos)
	_, err = f.ReadAt(buf, seekPos)
	if err != nil {
		return nil, err
	}

	// Skip to first complete line
	if seekPos > 0 {
		for i := 0; i < len(buf); i++ {
			if buf[i] == '\n' {
				buf = buf[i+1:]
				break
			}
		}
	}

	return buf, nil
}
