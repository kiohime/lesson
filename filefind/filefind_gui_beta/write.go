package main

import (
	"errors"
	"fmt"
	"os"

	walk "github.com/kiohime/lesson/filefind"
)

// ////////////////////////////////////////////////

//rootCheck - Получает переменную пути, проверет - файл или каталог. Если файл, то выдает ошибку
func rootCheck(r string) error {
	fmt.Println("## rootcheck")
	file, fileOpenErr := os.Open(r)
	if fileOpenErr != nil {
		return fileOpenErr
	}
	defer file.Close()
	fi, err := file.Stat()
	switch {
	case !fi.IsDir():
		err = errors.New("root var is file")
	default:
		fmt.Println("root var is directory")
	}
	return err
}

////////////////////////////////////////////////

//writeBaser - пишет базу
func writeBaser(aset *Settings, adata *Data) error {
	fmt.Println("# writebaser")

	// проверка переменной пути на то, является ли та настоящим путем, если нет - остановить программу
	err := rootCheck(aset.RootDir)
	// err = errors.New("TEST_ERROR")
	if err != nil {
		fmt.Println("root checking error")
		return err
	}
	fmt.Printf("ScanMode is %v", aset.ScanMode)
	exportFullPath := aset.WorkDir + aset.TargetFileName

	// парсит
	res, errs := walk.StartWalk(
		[]string{aset.RootDir},
		walk.Options{
			WalkScanMode: aset.ScanMode,
		},
	)
	// adata.Cache = res

	// создает файл
	err = bahniFile(exportFullPath, &res)
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		var err string
		// fmt.Printf("error in walking : %q\n", err)
		for i, e := range errs {
			err += fmt.Sprintf("	%v : %v\n", i, e)
		}
		return fmt.Errorf("%v", err)
	}
	fmt.Println("# END writebaser")
	return nil
}

// ///////////////////////////////////////////////
