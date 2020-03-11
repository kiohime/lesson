package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	subDirToSkip := "lesson2"

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == subDirToSkip {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}
		fmt.Printf("visited file or dir: %q\n", path)
		return nil
	}
	//
	fmt.Println("On Unix:")
	err := filepath.Walk(".", walkFunc)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", "test", err)
		return
	}
}
