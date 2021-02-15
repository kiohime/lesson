package main
import (
    "fmt"
    "os"
    "strings"
)

func main() {
	a := os.Getenv("Path")
    // fmt.Println(a)
    
    if !strings.Contains(a, ";c:\\totalcmd\\tools\\bin") {
        a += ";c:\\totalcmd\\tools\\bin"
    }
    a = fmt.Sprintf("setx Path \"%v\"\npause", a)


    f, err := os.Create("install.bat")
        if err != nil {
            panic(err)
        }
    defer f.Close()

    _, err = f.Write([]byte(a))
        if err != nil {
            panic(err)
        }
}