package main

import (
	"io"
	"log"
	"os"
)

// Mock implementation of gh. It just prints whatever is in the file pointed to
// by the GH_RESPONSE environment variable.
func main() {
	filename := os.Getenv("GH_RESPONSE")
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("opening file %q: %v", filename, err)
	}
	defer file.Close()

	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		log.Fatalf("writing response: %v", err)
	}
}
