package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/macroblock/imed/pkg/cli"
)

var (
	rootDir  string
	scanDir  bool
	scanFile bool
	scanMode int
	argCache []string
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

// Создание файлов на экспорт
func makeFile() error {
	// установка имени файла зависит от режима отрисовки
	exportFileName := ""
	switch scanMode {
	case 1:
		exportFileName = "filesearch_directory.txt"
	case 2:
		exportFileName = "filesearch_files.txt"
	default:
		exportFileName = "filesearch_default.txt"
	}

	// установка пути файла, можно забить хардкодом для конкретного места
	exportFilePath := ""
	exportFullPath := exportFilePath + exportFileName

	// создание файла по полному пути, вставка значений из кэша отрисовки с обрезкой лишняка
	file, err := os.OpenFile(exportFullPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	for i, data := range argCache {
		// fmt.Printf("%q\n", data)
		if i < len(argCache)-1 {
			_, _ = file.WriteString(data + "\n")
		} else {
			_, _ = file.WriteString(data)
		}
	}
	// закрытие файла
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

////////////////////////////////////////////////
func printer() error {
	// fmt.Println(argCache)
	return nil
}

// пауза
func keyWait() {
	fmt.Printf("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

////////////////////////////////////////////////
func reader() error {
	if len(argCache) == 0 {
		err := errors.New("no search arguments was inputed")
		return err
	}
	fmt.Println(argCache)
	return nil
}

////////////////////////////////////////////////

// Проверка флагов и устанавливает режим отрисовки скинированного
func mainFunc() error {
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

	// парсит
	err = startWalk()
	if err != nil {
		fmt.Printf("error in walking : %q\n", err)
		return err
	}
	// ничего не делает
	err = printer()
	if err != nil {
		fmt.Printf("error in caching : %q\n", err)
		return err
	}
	// создает файл
	err = makeFile()
	if err != nil {
		fmt.Printf("error in making file : %q\n", err)
		return err
	}
	return nil
}

////////////////////////////////////////////////
func main() {
	// установка алиасов и значений флагов
	flagSet := cli.New("!PROG! сканирует пути и найденное кладет в файл", mainFunc)
	flagSet.Elements(
		cli.Flag("-d -dir : показывает только пути", &scanDir),
		cli.Flag("-f -file : показывает только файлы", &scanFile),
		cli.Flag("-h -help -? /? : справка", flagSet.PrintHelp).Terminator(),
		cli.Flag(": file pathes", &rootDir),
		cli.Command("sd search : режим поиска по имеющейся базе", reader,
			cli.Flag(": search arguments", &argCache)),
	)

	// парсинг введеных аргументов на предмет флагов
	args := os.Args
	err := flagSet.Parse(args)
	// err = errors.New("TEST_ERROR")
	if err != nil {
		fmt.Printf("error in parsing arguments : %q\n", err)
		keyWait()
		os.Exit(1)
	}
	keyWait()
}
