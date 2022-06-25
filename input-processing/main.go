package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	filterTerm = "error"
	blockSize  = uint64(4096)
)

func main() {
	flag.Parse()

	fmt.Println("SP// Backend Developer Test - Input Processing")
	fmt.Println()

	// Handle the special case for no filter simply
	if len(filterTerm) == 0 {
		if _, err := io.Copy(os.Stdout, os.Stdin); err != nil {
			log.Fatalln("Error while transferring Stdin to Stdout")
		}
		return
	}

	buffer := make([]byte, blockSize)
	possibleMatchPosition := 0
	matchedLine := false
	var inProgressLine bytes.Buffer

	for {
		sizeRead, err := readBuffer(os.Stdin, buffer)
		for _, b := range buffer[:sizeRead] {
			if err := inProgressLine.WriteByte(b); err != nil {
				log.Fatalln("Error storing in progress line.", err)
			}

			if b == '\n' {
				if matchedLine {
					if _, err := inProgressLine.WriteTo(os.Stdout); err != nil {
						log.Fatalln("Error writing matching line to standard out.", err.Error())
					}
				} else {
					inProgressLine.Reset()
					possibleMatchPosition = 0
					matchedLine = false
				}
			}

			if !matchedLine {
				// Check if the byte matches the next possible character in the filter term if we haven't fully matched yet.
				if possibleMatchPosition < len(filterTerm) && b == filterTerm[possibleMatchPosition] {
					// We are matching the filter term still, increment the counter
					possibleMatchPosition++
					if possibleMatchPosition == len(filterTerm) {
						matchedLine = true
					}
				} else if possibleMatchPosition > 0 {
					// Reset the possible match counter, byte did not match
					possibleMatchPosition = 0
				}
			}
		}

		if err != nil {
			break
		}
	}
}

func readBuffer(reader io.Reader, buffer []byte) (int, error) {
	sizeRead := 0
	for {
		n, err := reader.Read(buffer[sizeRead:])
		sizeRead += n
		if err != nil {
			return sizeRead, err
		}
		if sizeRead == len(buffer) {
			break
		}
	}
	return sizeRead, nil
}

func init() {
	flag.StringVar(&filterTerm, "filter", filterTerm, "Sets the string to filter standard input lines to standard output. If empty everything is copied over")
	flag.Uint64Var(&blockSize, "bs", blockSize, "Sets the input buffer block size")
}
