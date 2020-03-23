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
	rootDir    string
	searchDir  bool
	searchFile bool
	searchMode int
	helpMode   bool
	pathCache  []string
)

////////////////////////////////////////////////
func mainFunc() error {
	// fmt.Println("=== mainFunc called")
	return nil
}

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

// Проверка флагов и устанавливает режим отрисовки скинированного
func initiate() error {
	// составная переменная для корневого каталога
	// rootDir1 := ""
	// rootDir2 := " : totalcommander mode"
	// rootDir = rootDir1 + rootDir2
	// установка алиасов и значений флагов
	flagSet := cli.New("!PROG! сканирует пути и найденное кладет в файл", mainFunc)
	flagSet.Elements(
		cli.Flag("-od -dir : показывает только пути", &searchDir),
		cli.Flag("-of -file : показывает только файлы", &searchFile),
		cli.Flag("-h -help -? /? : справка", &helpMode),
		cli.Flag(": file pathes", &rootDir),
	)
	// парсинг введеных аргументов на предмет флагов
	args := os.Args
	err := flagSet.Parse(args)

	// проверка переменной пути на то, является ли та настоящим путем, если нет - остановить программу
	errRoot := rootCheck(rootDir)

	// установка очереди отработки флагов и режимов отрисовки по флагам. по умолчанию считывает и каталоги, и файлы
	switch {
	case helpMode:
		err = flagSet.PrintHelp()
		keyWait()
		defer os.Exit(9)
	case errRoot != nil:
		return errRoot
	case searchDir && !searchFile:
		// -od
		searchMode = 1
		fmt.Println("searchDir set")
	case !searchDir && searchFile:
		// -of
		searchMode = 2
		fmt.Println("searchFile set")
	default:
		// no args
		fmt.Println("default operators set")
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
		//
		// тэги данных "каталог" и "файл"
		isDir := info.Mode().IsDir()
		isFile := info.Mode().IsRegular()

		// в зависимости от режима отрисовки, составлять кэш отрисовки
		switch searchMode {

		// только каталоги
		case 1:
			if searchDir && isDir {
				// fmt.Printf("visited : %q\n", path)
				pathCache = append(pathCache, path)
			}

		// только файлы
		case 2:
			if searchFile && isFile {
				// fmt.Printf("visited : %q\n", path)
				pathCache = append(pathCache, path)
			}
		// по дефолту добавляет всё
		default:
			// fmt.Printf("visited : %q\n", path)
			pathCache = append(pathCache, path)
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
	switch searchMode {
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
	for i, data := range pathCache {
		// fmt.Printf("%q\n", data)
		if i < len(pathCache)-1 {
			_, _ = file.WriteString(data + "\n")
		} else {
			_, _ = file.WriteString(data)
		}
	}
	// закрытие файла
	file.Close()

	return nil
}

////////////////////////////////////////////////
func printer() error {
	// fmt.Println(pathCache)
	return nil
}

// пауза
func keyWait() {
	fmt.Printf("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

////////////////////////////////////////////////
func main() {
	// устанавливает режимы
	err := initiate()
	if err != nil {
		fmt.Printf("error in initiation : %q\n", err)
		keyWait()
		os.Exit(1)
	}
	// парсит
	err = startWalk()
	if err != nil {
		fmt.Printf("error in walking : %q\n", err)
		keyWait()
		os.Exit(1)
	}
	// ничего не делает
	err = printer()
	if err != nil {
		fmt.Printf("error in caching : %q\n", err)
		keyWait()
		os.Exit(1)
	}
	// создает файл
	err = makeFile()
	if err != nil {
		fmt.Printf("error in making file : %q\n", err)
		keyWait()
		os.Exit(1)
	}

}
