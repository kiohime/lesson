package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	inputMode := os.Args[1]
	rootDir := "c:\\__downloads\\Godot 3 Complete Developer Course - 2D and 3D"
	subDirToSkip := ".lool"

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == subDirToSkip {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}

		tolkoFayly := info.Mode().IsRegular()
		tolkoPuti := info.Mode().IsDir()
		switch inputMode {
		case "-tf":
			if tolkoFayly {
				fmt.Printf("visited file: %q\n", path)
			}
		case "-tp":
			if tolkoPuti {
				fmt.Printf("visited dir: %q\n", path)
			}
			// default:
			// 	fmt.Printf("visited file or path: %q\n", path)
		}
		return nil
	}
	//

	fmt.Println("###########SEARCHING###########")
	err := filepath.Walk(rootDir, walkFunc)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", "test", err)
		return
	}
}
