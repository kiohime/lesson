package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/macroblock/imed/pkg/cli"
)

var (
	rootDir        string
	scanDir        bool
	scanFile       bool
	scanMode       int
	argCache       []string
	dataForPrinter []string

	workDir         = ""
	baseNameDefault = "filesearch_default.txt"
	baseNameFiles   = "filesearch_files.txt"
	baseNameDirs    = "filesearch_directory.txt"
	resultFileName  = "result.txt"
)

////////////////////////////////////////////////

// Получает переменную пути, проверет - файл или каталог. Если файл, то выдает ошибку
func rootCheck(r string) error {
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

// Сканирование данных в пути и добавление их в кэш отрисовки
func startWalk() error {
	var walkError error

	// настройка сканирования данных
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("walkfunc : %v\n", walkError)
			walkError = err
			return walkError
		}

		// тэги данных "каталог" и "файл"
		isDir := info.Mode().IsDir()
		isFile := info.Mode().IsRegular()

		// в зависимости от режима отрисовки, составлять кэш отрисовки
		switch scanMode {

		// только каталоги
		case 1:
			if scanDir && isDir {
				// fmt.Printf("visited : %q\n", path)
				argCache = append(argCache, path)
			}

		// только файлы
		case 2:
			if scanFile && isFile {
				// fmt.Printf("visited : %q\n", path)
				argCache = append(argCache, path)
			}
		// по дефолту добавляет всё
		default:
			// fmt.Printf("visited : %q\n", path)
			argCache = append(argCache, path)
		}
		return walkError
	}
	// сканирование данных в переменной пути
	fmt.Printf("start walking\n")
	walkError = filepath.Walk(rootDir, walkFunc)
	fmt.Println(walkError)
	return nil
}

func bahniFile(inputName string, inputData *[]string) error {
	// создание файла по полному пути, вставка значений из кэша отрисовки с обрезкой лишняка

	file, err := os.OpenFile(inputName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	s := strings.Join(*inputData, "\n")
	_, e := file.WriteString(s)
	if e != nil {
		err = e
		return err
	}

	// закрытие файла
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

////////////////////////////////////////////////
func printer(result ...string) error {

	if len(result) == 0 {
		fmt.Println("no results were found")
		return nil
	}

	for _, r := range result {
		i, e := fmt.Println(r)
		if e != nil {
			fmt.Println(i)
			return e
		}
	}
	// fmt.Println(argCache)
	return nil
}

// пауза
func keyWait() {
	fmt.Printf("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

////////////////////////////////////////////////
func readBaser() error {
	argLen := len(argCache)
	if argLen == 0 {
		err := errors.New("no search arguments was inputed")
		return err
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

	base := workDir + "filelist_folders_upload.txt"
	fmt.Println(base)
	f, err := os.Open(base)
	if err != nil {
		fmt.Println(err)
		keyWait()
		os.Exit(1)
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
		return err
	}

	err = printer(dataForPrinter...)
	if err != nil {
		fmt.Printf("error in printing : %q\n", err)
		return err
	}

	return nil
}

////////////////////////////////////////////////

//устанавливает режим отрисовки скинированного
func writeBaser() error {
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
	err = startWalk()
	if err != nil {
		fmt.Printf("error in walking : %q\n", err)
		return err
	}
	// ничего не делает
	// err = printer(dataForPrinter...)
	// if err != nil {
	// 	fmt.Printf("error in caching : %q\n", err)
	// 	return err
	// }
	// создает файл
	err = bahniFile(exportFullPath, &argCache)
	if err != nil {
		fmt.Printf("error in making file : %q\n", err)
		return err
	}
	return nil
}

///////////////////////////////////////////////

func initialize() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("UserHomeDir error: %v", err)
	}
	filePath := filepath.Join(homeDir, ".config", "kiohime")
	workDir = filePath + "\\"
	err = os.MkdirAll(filepath.Dir(filePath), 0777)
	// 0666 for files
	if err != nil {
		return fmt.Errorf("Database error: %v", err)
	}
	return nil
}

////////////////////////////////////////////////
func main() {
	// checking database folders
	err := initialize()
	if err != nil {
		fmt.Printf("Error in initialisation : %q\n", err)
		keyWait()
		os.Exit(1)
	}

	// установка алиасов и значений флагов
	flagSet := cli.New("!PROG! сканирует пути и найденное кладет в файл", writeBaser)
	flagSet.Elements(
		// cli.Flag("--path : показывает только пути", &path),
		cli.Flag("-d -dir : сканирует только пути", &scanDir),
		cli.Flag("-f -file : сканирует только файлы", &scanFile),
		cli.Flag(": file pathes", &rootDir),
		cli.Command("sd search : режим поиска по имеющейся базе", readBaser,
			cli.Flag(": search arguments", &argCache),
			cli.Flag("-d -dir : ищет в базе только пути", &scanDir),
			cli.Flag("-f -file : ищет в базе только файлы", &scanFile)),
		cli.Flag("-h -help -? /? : справка", flagSet.PrintHelp).Terminator(),
	)

	// парсинг введеных аргументов на предмет флагов
	args := os.Args
	err = flagSet.Parse(args)
	// err = errors.New("TEST_ERROR")

	// fmt.Println("#########" + workDir)

	if err != nil {
		fmt.Printf("error in parsing arguments : %q\n", err)
		keyWait()
		os.Exit(1)
	}
	keyWait()
}
