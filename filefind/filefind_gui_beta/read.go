package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//readBaser - reading mode : сканирует существующую базу
func readBaser() ([]string, error) {
	dataForPrinter = nil
	fmt.Println("### readbaser")

	argLen := len(argCache)
	if argLen == 0 {
		err := errors.New("no search arguments was inputed")
		return nil, err
	}
	// fmt.Println(argCache)
	// readBaseName := basePath +

	// exportFileName := ""
	// switch scanMode {
	// case 1:
	// 	exportFileName = baseNameDirs
	// case 2:
	// 	exportFileName = baseNameFiles
	// default:
	// 	exportFileName = baseNameDefault
	// }

	// exportFullPath := workDir + exportFileName

	base := workDir + baseNameDirs
	fmt.Println(base)
	f, err := os.Open(base)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)
	// countScan := 0
	for scanner.Scan() {
		// fmt.Println("######################################")
		// fmt.Println("counter is ", countScan)
		line := scanner.Text()
		line = strings.TrimSpace(line)
		lineFile := filepath.Base(line)

		for _, a := range argCache {
			// fmt.Println("argument is ", a)
			if strings.Contains(lineFile, a) {
				dataForPrinter = append(dataForPrinter, line)
				// fmt.Println(line)
				// os.Exit(1)
			}
		}
		// countScan++
	}

	read := workDir + resultFileName

	err = bahniFile(read, &dataForPrinter)

	if err != nil {
		fmt.Printf("error in making file : %q\n", err)
		return nil, err
	}

	// err = printer(dataForPrinter...)
	// if err != nil {
	// 	fmt.Printf("error in printing : %q\n", err)
	// 	return err
	// }

	return dataForPrinter, err
}
