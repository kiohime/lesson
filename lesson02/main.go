package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("c:\\_working\\filesearch_files.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "menu") {
			fmt.Println(scanner.Text())
			// os.Exit(1)
		}

	}

	if err := scanner.Err(); err != nil {
		// Handle the error
	}

}
