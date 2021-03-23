package walk

import (
	"fmt"
	"os"
	"path/filepath"
)

type Options struct {
	Mode     int
	SkipDir  bool
	SkipFile bool
}

// StartWalk - Сканирование данных в пути и добавление их в кэш отрисовки
func StartWalk(paths []string, options Options) ([]string, []error) {
	fmt.Println("## StartWalk")
	var ret []string
	var errors []error
	for _, p := range paths {
		result, errs := startWalk(p, options)
		ret = append(ret, result...)
		errors = append(errors, errs...)
	}
	fmt.Println("## END StartWalk")
	return ret, errors
}

func startWalk(rootDir string, options Options) ([]string, []error) {
	fmt.Printf("### start walking\n")
	var ret []string
	var errors []error

	// настройка сканирования данных
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("walkfunc : %v\n", err)
			errors = append(errors, err)
			if info == nil {
				return nil
			}
		}

		// тэги данных "каталог" и "файл"
		isDir := info.Mode().IsDir()
		isFile := info.Mode().IsRegular()

		// в зависимости от режима отрисовки, составлять кэш отрисовки
		switch options.Mode {

		// только каталоги
		case 0:
			// fmt.Println("scanDir is", scanDir)
			// fmt.Println("isDir is", isDir)

			if !options.SkipDir && isDir {
				// fmt.Printf("visited : %q\n", path)
				// fmt.Println("11111111111111", argCache)
				ret = append(ret, path)
			}

		// только файлы
		case 1:
			if !options.SkipFile && isFile {
				// fmt.Printf("visited : %q\n", path)
				ret = append(ret, path)
			}
		// по дефолту добавляет всё
		default:
			// fmt.Printf("visited : %q\n", path)
			ret = append(ret, path)
		}
		return nil
	}
	// сканирование данных в переменной пути

	err := filepath.Walk(rootDir, walkFunc)
	if err != nil {
		errors = append(errors, err)
	}
	fmt.Printf("### END walking\n")
	// fmt.Println(walkError)
	return ret, errors
}
