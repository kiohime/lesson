package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/macroblock/imed/pkg/cli"
)

//flags
var (
	searchDir  bool
	searchFile bool
	searchMode int
	helpMode   bool
	pathCache  []string
)

func mainFunc() error {
	// fmt.Println("=== mainFunc called")
	return nil
}

func initiate() error {
	flagSet := cli.New("!PROG! сканирует пути и найденное кладет в файл", mainFunc)
	flagSet.Elements(
		cli.Flag("-od -tp : показывает только пути", &searchDir),
		cli.Flag("-of -tf : показывает только файлы", &searchFile),
		cli.Flag("-h -help -? /? : справка", &helpMode),
	)

	// err := flagSet.PrintHelp()
	// if err != nil {
	// 	fmt.Println("print help error: ", err)
	// }
	args := os.Args
	err := flagSet.Parse(args)
	// if err != nil {
	// 	// fmt.Println("parse error: ", err)
	// 	// fmt.Println("hint: ", flagSet.GetHint())
	// }

	// fmt.Println(searchDir, searchFile)
	// SETTING MODES
	switch {
	case helpMode:
		err = flagSet.PrintHelp()
		defer os.Exit(1)
	case searchDir && !searchFile:
		// -od
		searchMode = 1
		fmt.Println("searchDir set")
	case !searchDir && searchFile:
		// -of
		searchMode = 2
		fmt.Println("searchFile set")

	default:
		// no args or -od+of
		fmt.Println("default operators set")
	}
	return err
}

func printer() error {
	// fmt.Println(pathCache)
	return nil
}

func makeFile() error {
	file, err := os.OpenFile("filesearch.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	// dataWriter := bufio.NewWriter(file)
	for _, data := range pathCache {
		// fmt.Println(data)
		_, _ = file.WriteString(data + "\n")
	}
	// dataWriter.Flush()
	file.Close()
	return err
}

func startWalk() error {
	var walkError error
	rootDir := "c:\\_working\\_godot_projects\\asd\\"

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			walkError = err
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
	if walkError != nil {
		fmt.Printf("error walking the path : %q\n", walkError)
	}
	return walkError
}

func main() {
	// setting filemodes
	err := initiate()
	if err != nil {
		fmt.Printf("error in initiation : %q\n", err)
		os.Exit(1)
	}
	// start parcing
	err = startWalk()
	if err != nil {
		fmt.Printf("error in walking : %q\n", err)
		os.Exit(2)
	}

	err = printer()
	if err != nil {
		fmt.Printf("error in caching : %q\n", err)
		os.Exit(3)
	}

	err = makeFile()
	if err != nil {
		fmt.Printf("error in making file : %q\n", err)
		os.Exit(4)
	}

}
