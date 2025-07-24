package markdown

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Reader is an object that can read (e.a. stdin) and return one part at a time
type Reader struct {
	scanner *bufio.Scanner
}

// NewReader returns an initialized Reader
func NewReader(stream *bufio.Scanner) Reader {
	return Reader{
		scanner: stream,
	}
}

// NewFromStdin returns an initialized Reader, which runs from stdin
func NewFromStdin() Reader {
	return NewReader(bufio.NewScanner(os.Stdin))
}

// Read reads line by line until it found an empty line, and returns a Part
// from these lines
func (r Reader) Read() (p Part, err error) {
	var (
		lines   []string
		readErr error
	)
	if r.scanner == nil {
		return p, io.EOF
	}
	for r.scanner.Scan() {
		line := r.scanner.Text()
		if line == "" {
			part := Part(strings.Join(lines, "\n"))
			return part, readErr
		}
		lines = append(lines, line)
	}
	lines = append(lines, r.scanner.Text())
	part := Part(strings.Join(lines, "\n"))
	return part, readErr
}
