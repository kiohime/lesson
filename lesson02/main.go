package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	f, err := os.Open("e:\\filesearch_files.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		lineFile := filepath.Base(line)
		if strings.Contains(lineFile, "mama") {
			fmt.Println(line)
			// os.Exit(1)
		}
	}

	if err := scanner.Err(); err != nil {
		// Handle the error
	}

}
