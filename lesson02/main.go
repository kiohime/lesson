package main

import (
	"fmt"
	"os"
)

func main() {

	// dir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(dir)

	// err = os.Chdir("c:\\Users\\park-\\go\\src\\github.com\\kiohime\\lesson\\filesearch\\")
	// if err != nil {
	// 	panic(err)
	// }

	// dir, err = os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(dir)

	fmt.Println(os.Getenv("FILE_SEARCH_PATH"))

}
