package main

import (
	"flag"
	"fmt"
	"github.com/stackpath/backend-developer-tests/input-processing/filter"
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

	textFilter := filter.New(filterTerm, blockSize)
	if err := textFilter.Process(os.Stdout, os.Stdin); err != nil {
		log.Fatalln("Error encountered while filtering", err)
	}
}

func init() {
	flag.StringVar(&filterTerm, "filter", filterTerm, "Sets the string to filter standard input lines to standard output. If empty everything is copied over")
	flag.Uint64Var(&blockSize, "bs", blockSize, "Sets the input buffer block size")
}
