package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"fyne.io/fyne/v2/widget"
)

var (
	argCache       []string
	dataForPrinter []string
)
var (
	progressBar   *widget.ProgressBarInfinite
	entry_widget  *modifiedEntry
	screen_widget *widget.Label
)

// ///////////////СТРУКТУРА НАСТРОЙКИ ПРИЛОЖЕНИЯ/////////////////////////////

// AppSet - глобальная переменная с настройками приложения
var AppSet *AppSettings

// AppSettings - глобальная структура с настройками приложения
type AppSettings struct {
	AppMode        int
	ScanMode       int
	RootDir        string
	WorkDir        string
	BaseNameFiles  string
	BaseNameDirs   string
	ResultFileName string
}

// блок инициализации: установка рабочего пути для файлов базы и поиска
func initialize(aset *AppSettings) error {
	fmt.Println("### initialize")

	aset.BaseNameFiles = "filesearch_files.txt"
	aset.BaseNameDirs = "filesearch_directory.txt"
	aset.ResultFileName = "result.txt"

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("UserHomeDir error: %v", err)
	}
	filePath := filepath.Join(homeDir, ".config", "kiohime", "file.txt")
	filePathDir := filepath.Dir(filePath)
	aset.WorkDir = filePathDir + "\\"
	err = os.MkdirAll(filePathDir, 0777)
	// 0666 for files
	if err != nil {
		return fmt.Errorf("Database error: %v", err)
	}
	fmt.Println("### END initialize")
	return nil
}

//bahniFile - создает файл с именем inputName из массива данных inputData, данные добавляются через ньюлайн
func bahniFile(inputName string, inputData *[]string) error {
	fmt.Println("#### bahnifile")
	fmt.Println(inputName)
	// fmt.Println(inputData)

	// err := os.Remove(inputName)
	// создание файла по полному пути, вставка значений из кэша отрисовки с обрезкой лишняка
	file, err := os.OpenFile(inputName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println("#### END bahnifile")
		return fmt.Errorf("bahni file 1 : %v", err)
	}
	s := strings.Join(*inputData, "\n")
	_, e := file.WriteString(s)

	if e != nil {
		err = e
		fmt.Println("#### END bahnifile")
		return fmt.Errorf("bahni file 2 : %v", err)
	}

	// закрытие файла
	err = file.Close()
	if err != nil {
		fmt.Println("#### END bahnifile")
		return fmt.Errorf("bahni file 3 : %v", err)
	}
	return nil
}

// decider - executes searching mode, depending on switch
func decider(input *modifiedEntry, scr *widget.Label, aset *AppSettings) {
	assert(input, scr)

	fmt.Printf("on decider AppMode is %v\n", aset.AppMode)

	progressBar.Start()
	progressBar.Show()
	argCache = nil
	switch aset.AppMode {
	case 0:
		argCache = append(argCache, input.Text)
	case 1:
		aset.RootDir = input.Text
	}
	result := executer(aset)
	// dataBox.Add(widget.NewLabel(result))
	scr.Text = result
	input.Text = ""
	input.Refresh()
	progressBar.Hide()
	progressBar.Stop()
}

//executer - запускает программу в устновленном режиме
func executer(aset *AppSettings) string {
	result := ""
	fmt.Printf("on executer AppMode is %v\n", aset.AppMode)
	switch aset.AppMode {
	case 0:
		base, e := readBaser(aset)
		if e != nil {
			a := fmt.Errorf("error in reading base : %v", e)
			fmt.Println(a)
			// os.Exit(1)
		}
		result = strings.Join(base, "\n")
	case 1:
		e := writeBaser(aset)
		if e != nil {
			a := fmt.Errorf("error in writing base : %v", e)
			fmt.Println(a)
		}
		// result = "000000000000000000000000000000"
		fmt.Println("THE END")
	}
	return result
}

// assert - проверяет на наличие критических ошибок
func assert(args ...interface{}) {
	okGlobal := true
	msg := "assertion failed:\n"
	for i, arg := range args {
		// fmt.Printf("____%v: %v\n", i, arg)
		ok := true
		switch t := arg.(type) {
		case bool:
			ok = t
		case error:
			ok = t == nil
		case nil:
			ok = false
			fmt.Printf("NILNIL\n")
		}
		if reflect.ValueOf(arg).IsNil() {
			ok = false
		}

		if !ok {
			_, file, line, _ := runtime.Caller(1)
			msg += fmt.Sprintf("\t(%v, %T) %v: %v\n", i, arg, file, line)
		}
		okGlobal = okGlobal && ok
	}
	if !okGlobal {
		panic(msg)
	}
}

func main() {
	AppSet = &AppSettings{}
	err := initialize(AppSet)
	if err != nil {
		fmt.Printf("Error in initialisation : %q\n", err)
		keyWait()
		os.Exit(1)
	}
	gui(AppSet)
}
