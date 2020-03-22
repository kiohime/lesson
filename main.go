package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/macroblock/imed/pkg/cli"
)

var (
	rootDir    string
	searchDir  bool
	searchFile bool
	// tcMode     bool
	searchMode int
	helpMode   bool
	pathCache  []string
)

func mainFunc() error {
	// fmt.Println("=== mainFunc called")
	return nil
}

func initiate() error {

	rootDir1 := ""
	rootDir2 := " : totalcommander mode"
	rootDir = rootDir1 + rootDir2

	flagSet := cli.New("!PROG! сканирует пути и найденное кладет в файл", mainFunc)
	flagSet.Elements(
		cli.Flag("-od -dir : показывает только пути", &searchDir),
		cli.Flag("-of -file : показывает только файлы", &searchFile),
		cli.Flag("-h -help -? /? : справка", &helpMode),
		cli.Flag(rootDir1, &rootDir),
	)

	args := os.Args
	err := flagSet.Parse(args)

	// SETTING MODES
	switch {
	case helpMode:
		err = flagSet.PrintHelp()
		defer os.Exit(1)
	// case tcMode:
	// 	fmt.Printf("tcMode set, %v\n", tcMode)
	// 	defer os.Exit(1)
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

func printer() error {
	// fmt.Println(pathCache)
	return nil
}

func makeFile() error {
	exportFileName := ""
	switch searchMode {
	case 1:
		exportFileName = "filesearch_directory.txt"
	case 2:
		exportFileName = "filesearch_files.txt"
	default:
		exportFileName = "filesearch_default.txt"
	}

	exportFilePath := ""
	exportFullPath := exportFilePath + exportFileName
	file, err := os.OpenFile(exportFullPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	for i, data := range pathCache {
		// fmt.Println(data)
		if i < len(pathCache)-1 {
			_, _ = file.WriteString(data + "\n")
		} else {
			_, _ = file.WriteString(data)
		}

	}
	file.Close()
	return err
}

func startWalk() error {
	var walkError error
	// rootDir = "c:\\__downloads\\"

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			walkError = err
			return walkError
		}

		isDir := info.Mode().IsDir()
		isFile := info.Mode().IsRegular()

		switch searchMode {
		case 1:
			if searchDir && isDir {
				// fmt.Printf("visited : %q\n", path)
				pathCache = append(pathCache, path)
			}
		case 2:
			if searchFile && isFile {
				// fmt.Printf("visited : %q\n", path)
				pathCache = append(pathCache, path)
			}
		default:
			// fmt.Printf("visited : %q\n", path)
			pathCache = append(pathCache, path)
		}
		return walkError
	}

	walkError = filepath.Walk(rootDir, walkFunc)
	return walkError
}

func keyWait() {
	fmt.Printf("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func main() {
	// setting filemodes
	err := initiate()
	if err != nil {
		fmt.Printf("error in initiation : %q\n", err)
		keyWait()
	}
	// start parcing
	err = startWalk()
	if err != nil {
		fmt.Printf("error in walking : %q\n", err)
		keyWait()
	}

	err = printer()
	if err != nil {
		fmt.Printf("error in caching : %q\n", err)
		keyWait()
	}

	err = makeFile()
	if err != nil {
		fmt.Printf("error in making file : %q\n", err)
		keyWait()
	}

}
