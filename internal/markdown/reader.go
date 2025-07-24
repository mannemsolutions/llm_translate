package markdown

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Reader is an object that can read (e.a. stdin) and return one part at a time
type Reader struct {
	reader *bufio.Reader
}

// NewReader returns an initialized Reader
func NewReader(stream *bufio.Reader) Reader {
	return Reader{
		reader: stream,
	}
}

// NewFromStdin returns an initialized Reader, which runs from stdin
func NewFromStdin() Reader {
	return NewReader(bufio.NewReader(os.Stdin))
}

// Read reads line by line until it found an empty line, and returns a Part
// from these lines
func (r Reader) Read() (*Part, error) {
	var (
		line    string
		lines   []string
		readErr error
	)
	if r.reader == nil {
		return nil, io.EOF
	}
	for {
		line, readErr = r.reader.ReadString('\n')
		if readErr != nil {
			// last part
			if readErr == io.EOF {
				r.reader = nil
				if len(line) > 0 {
					lines = append(lines, line)
				}
				break
			}
			return nil, readErr
		} else if line == "" {
			part := Part(strings.Join(lines, "\n"))
			return &part, readErr
		}
		lines = append(lines, line)
		break
	}
	part := Part(strings.Join(lines, "\n"))
	return &part, readErr
}
