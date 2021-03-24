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
func readBaser(aset *Settings, adata *Data) ([]string, error) {
	result := []string{}
	fmt.Println("### readbaser")

	// проверка на кол-во аргументов на 24 марта 2021 аргумент всегда один,
	argLen := len(adata.Cache)
	// fmt.Println("11111  ", argLen)
	fmt.Println()
	if argLen == 0 {
		err := errors.New("no search arguments was inputed")
		return nil, err
	}
	// так что эта часть кода вообщето не нужна?

	base := aset.WorkDir + aset.TargetFileName
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

		for _, a := range adata.Cache {
			// fmt.Println("argument is ", a)
			if strings.Contains(lineFile, a) {
				result = append(result, line)
			}
		}
	}

	// read := aset.WorkDir + aset.ResultFileName

	// err = bahniFile(read, &adata.PrintData)

	// if err != nil {
	// 	fmt.Printf("error in making file : %q\n", err)
	// 	return nil, err
	// }

	return result, err
}
