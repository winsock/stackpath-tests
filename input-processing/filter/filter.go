package filter

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// Filter A helper that allows filtering text lines from a source to a destination if it contains the filter term
// This class is not thread safe, only one thread may access filter at a time.
type Filter struct {
	filterTerm            string
	possibleMatchPosition int
	matchedLine           bool
	buffer                []byte
	inProgressLine        bytes.Buffer
}

func New(filterTerm string, blockSize uint64) *Filter {
	return &Filter{
		filterTerm:            filterTerm,
		possibleMatchPosition: 0,
		matchedLine:           false,
		buffer:                make([]byte, blockSize),
	}
}

// Process Reads from source and if a line matches the filterTerm writes to dest. Will block until there is an error or EOF is hit
func (f *Filter) Process(dest io.Writer, source io.Reader) error {
	// Make sure we reset for possible future runs
	defer f.Reset()

	for {
		sizeRead, err := f.readBuffer(source)
		for _, b := range f.buffer[:sizeRead] {
			if err := f.inProgressLine.WriteByte(b); err != nil {
				return fmt.Errorf("error writting to in progress line buffer, %w", err)
			}

			if b == '\n' {
				if f.matchedLine {
					if _, err := f.inProgressLine.WriteTo(dest); err != nil {
						return fmt.Errorf("error writting to destination, %w", err)
					}
				}
				f.Reset()
				continue
			}

			if !f.matchedLine {
				// Check if the byte matches the next possible character in the filter term if we haven't fully matched yet.
				if f.possibleMatchPosition < len(f.filterTerm) && b == f.filterTerm[f.possibleMatchPosition] {
					// We are matching the filter term still, increment the counter
					f.possibleMatchPosition++
					if f.possibleMatchPosition == len(f.filterTerm) {
						f.matchedLine = true
					}
				} else if f.possibleMatchPosition > 0 {
					// Reset the possible match counter, byte did not match
					f.possibleMatchPosition = 0
				}
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("error while reading byte from source, %w", err)
		}
	}
}

func (f *Filter) Reset() {
	f.inProgressLine.Reset()
	f.possibleMatchPosition = 0
	f.matchedLine = false
}

func (f *Filter) readBuffer(reader io.Reader) (int, error) {
	sizeRead := 0
	for {
		n, err := reader.Read(f.buffer[sizeRead:])
		sizeRead += n
		if err != nil {
			return sizeRead, err
		}
		if sizeRead == len(f.buffer) {
			break
		}
	}
	return sizeRead, nil
}
