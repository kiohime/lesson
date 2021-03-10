package main

import (
	"errors"
	"fmt"
	"os"

	walk "github.com/kiohime/lesson/filefind"
)

// ////////////////////////////////////////////////

// Получает переменную пути, проверет - файл или каталог. Если файл, то выдает ошибку
func rootCheck(r string) error {
	fmt.Println("### rootcheck")
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

// // Сканирование данных в пути и добавление их в кэш отрисовки
// func startWalk() error {
// 	fmt.Println("### startWalk")

// 	var walkError error

// 	// настройка сканирования данных
// 	walkFunc := func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			fmt.Printf("walkfunc : %v\n", walkError)
// 			walkError = err
// 			return walkError
// 		}

// 		// тэги данных "каталог" и "файл"
// 		isDir := info.Mode().IsDir()
// 		isFile := info.Mode().IsRegular()

// 		// в зависимости от режима отрисовки, составлять кэш отрисовки
// 		switch scanMode {

// 		// только каталоги
// 		case 1:
// 			// fmt.Println("scanDir is", scanDir)
// 			// fmt.Println("isDir is", isDir)

// 			if scanDir && isDir {
// 				// fmt.Printf("visited : %q\n", path)
// 				// fmt.Println("11111111111111", argCache)
// 				argCache = append(argCache, path)
// 			}

// 		// только файлы
// 		case 2:
// 			if scanFile && isFile {
// 				// fmt.Printf("visited : %q\n", path)
// 				argCache = append(argCache, path)
// 			}
// 		// по дефолту добавляет всё
// 		default:
// 			// fmt.Printf("visited : %q\n", path)
// 			argCache = append(argCache, path)
// 		}
// 		return walkError
// 	}
// 	// сканирование данных в переменной пути
// 	fmt.Printf("start walking\n")
// 	walkError = filepath.Walk(rootDir, walkFunc)
// 	fmt.Println(walkError)
// 	return nil
// }

////////////////////////////////////////////////

//устанавливает режим отрисовки скинированного
func writeBaser() error {
	fmt.Println("### writebaser")

	// проверка переменной пути на то, является ли та настоящим путем, если нет - остановить программу
	err := rootCheck(rootDir)
	// err = errors.New("TEST_ERROR")
	if err != nil {
		fmt.Println("root checking error")
		return err
	}
	// установка очереди отработки флагов и режимов отрисовки по флагам. по умолчанию считывает и каталоги, и файлы
	switch {
	case scanDir && !scanFile:
		// -d
		scanMode = 1
		fmt.Println("scanDir set")
	case !scanDir && scanFile:
		// -f
		scanMode = 2
		fmt.Println("scanFile set")
	default:
		// no args
		fmt.Println("default operators set")
		scanMode = 1
	}

	exportFileName := ""
	switch scanMode {
	case 1:
		exportFileName = baseNameDirs
	case 2:
		exportFileName = baseNameFiles
	default:
		exportFileName = baseNameDefault
	}

	exportFullPath := workDir + exportFileName

	// парсит
	res, errs := walk.StartWalk(
		[]string{rootDir},
		walk.Options{
			Mode:     scanMode,
			SkipDir:  !scanDir,
			SkipFile: !scanFile,
		},
	)

	argCache = res
	// ничего не делает
	// err = printer(dataForPrinter...)
	// if err != nil {
	// 	fmt.Printf("error in caching : %q\n", err)
	// 	return err
	// }
	// создает файл
	err = bahniFile(exportFullPath, &argCache)
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
	return nil
}

// ///////////////////////////////////////////////
