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
func readBaser(aset *AppSettings) ([]string, error) {
	dataForPrinter = nil
	fmt.Println("### readbaser")

	// проверка на кол-во аргументов на 24 марта 2021 аргумент всегда один,
	// так что эта часть кода вообщето не нужна?
	argLen := len(argCache)
	fmt.Println("11111  ", argLen)
	fmt.Println()
	if argLen == 0 {
		err := errors.New("no search arguments was inputed")
		return nil, err
	}

	base := aset.WorkDir + aset.BaseNameDirs
	// fmt.Println(base)
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

	read := aset.WorkDir + aset.ResultFileName

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
