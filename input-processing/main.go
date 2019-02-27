package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Println("SP// Backend Developer Test - Input Processing")
	fmt.Println()

	// Make sure input is being piped to this program. STDIN should be a
	// named pipe.
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeNamedPipe != 0 {
		fmt.Println("Please pipe input to this program.")
		return
	}

	// Read STDIN into a new buffered reader
	reader := bufio.NewReader(os.Stdin)

	// TODO: Look for lines in the STDIN reader that contain "error" and output them.
}
