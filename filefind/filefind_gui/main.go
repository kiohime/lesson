package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
)

var (
	rootDir        string
	scanDir        bool
	scanFile       bool
	scanMode       int
	argCache       []string
	dataForPrinter []string

	appMode         int
	workDir         = ""
	baseNameDefault = "filesearch_default.txt"
	baseNameFiles   = "filesearch_files.txt"
	baseNameDirs    = "filesearch_directory.txt"
	resultFileName  = "result.txt"
)

// создает файл с именем inputName из массива данных inputData, данные добавляются через ньюлайн
func bahniFile(inputName string, inputData *[]string) error {
	fmt.Println("### bahnifile")
	// fmt.Println(inputName)
	// fmt.Println(inputData)

	err := os.Remove(inputName)
	// создание файла по полному пути, вставка значений из кэша отрисовки с обрезкой лишняка
	file, err := os.OpenFile(inputName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("bahni file 1 : %v", err)
	}
	s := strings.Join(*inputData, "\n")
	_, e := file.WriteString(s)
	if e != nil {
		err = e
		return fmt.Errorf("bahni file 2 : %v", err)
	}

	// закрытие файла
	err = file.Close()
	if err != nil {
		return fmt.Errorf("bahni file 3 : %v", err)
	}
	return nil
}

// ////////////////////////////////////////////////
// func printer(result ...string) error {
// 	fmt.Println("### printer")

// 	if len(result) == 0 {
// 		fmt.Println("no results were found")
// 		return nil
// 	}

// 	for _, r := range result {
// 		dataBox.Add(widget.NewLabel(r))
// 		// i, e := fmt.Println(r)
// 		// if e != nil {
// 		// 	fmt.Println(i)
// 		// 	return e
// 		// }
// 	}
// 	// fmt.Println(argCache)
// 	return nil
// }

// пауза
func keyWait() {
	fmt.Printf("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// блок инициализации: установка рабочего пути для файлов базы и поиска
func initialize() error {
	fmt.Println("### initialize")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("UserHomeDir error: %v", err)
	}
	filePath := filepath.Join(homeDir, ".config", "kiohime", "file.txt")
	filePathDir := filepath.Dir(filePath)
	workDir = filePathDir + "\\"
	err = os.MkdirAll(filePathDir, 0777)
	// 0666 for files
	if err != nil {
		return fmt.Errorf("Database error: %v", err)
	}
	return nil
}

// ////////////////////////////////////////////////
// func main() {
// 	// checking database folders
// 	err := initialize()
// 	if err != nil {
// 		fmt.Printf("Error in initialisation : %q\n", err)
// 		keyWait()
// 		os.Exit(1)
// 	}

// 	// установка алиасов и значений флагов
// 	flagSet := cli.New("!PROG! сканирует пути и найденное кладет в файл", writeBaser)
// 	flagSet.Elements(
// 		// cli.Flag("--path : показывает только пути", &path),
// 		cli.Flag("-d -dir : сканирует только пути", &scanDir),
// 		cli.Flag("-f -file : сканирует только файлы", &scanFile),
// 		cli.Flag(": file pathes", &rootDir),
// 		cli.Command("sd search : режим поиска по имеющейся базе", readBaser,
// 			cli.Flag(": search arguments", &argCache),
// 			cli.Flag("-d -dir : ищет в базе только пути", &scanDir),
// 			cli.Flag("-f -file : ищет в базе только файлы", &scanFile)),
// 		cli.Flag("-h -help -? /? : справка", flagSet.PrintHelp).Terminator(),
// 	)

// 	// парсинг введеных аргументов на предмет флагов
// 	args := os.Args
// 	err = flagSet.Parse(args)
// 	// err = errors.New("TEST_ERROR")

// 	// fmt.Println("#########" + workDir)

// 	if err != nil {
// 		fmt.Printf("error in parsing arguments : %q\n", err)
// 		keyWait()
// 		os.Exit(1)
// 	}
// 	keyWait()
// }

// запускает программу в устновленном режиме
func executer() string {
	result := ""
	switch appMode {
	case 0:
		base, e := readBaser()
		if e != nil {
			a := fmt.Errorf("error in reading base : %v", e)
			fmt.Println(a)
			os.Exit(1)
		}
		result = strings.Join(base, "\n")
	case 1:

		// Временно! пока галочек нет!!!!
		scanDir = true

		e := writeBaser()
		if e != nil {
			a := fmt.Errorf("error in writing base : %v", e)
			fmt.Println(a)
			os.Exit(1)
		}
		// result = "000000000000000000000000000000"
	}
	return result

}

func findBtn(input *widget.Entry, scr *widget.Label) func() {
	return func() {
		argCache = nil
		rootDir = ""
		switch appMode {
		case 0:
			argCache = append(argCache, input.Text)
		case 1:
			rootDir = input.Text
		}
		result := executer()
		// dataBox.Add(widget.NewLabel(result))
		scr.Text = result
		input.Text = ""
		input.Refresh()
	}
}

func gui() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Notepad")
	myWindow.Resize(fyne.NewSize(1000, 500))

	setter := widget.NewRadioGroup([]string{"Поиск", "Сканирование"}, func(s string) {
		switch s {
		case "Поиск":
			appMode = 0
		case "Сканирование":
			appMode = 1
		}
	})
	setter.SetSelected("Поиск")
	setter.Refresh()

	curScreen := widget.NewLabel("")
	// screen.Text = "mama"
	dataBox := container.NewWithoutLayout(curScreen)
	dataBoxScroll := container.NewScroll(dataBox)
	// dataBoxScroll.Move()
	// dataBoxScroll.Resize(fyne.NewSize(400, 400))

	input := widget.NewEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Data", Widget: input, HintText: "input data"},
		},
		OnCancel: nil,
		OnSubmit: findBtn(input, curScreen),
	}
	input.Validator = validation.NewRegexp(`.+`, "input smthing")

	allContainers := container.NewBorder(
		// top
		container.NewVBox(
			container.NewBorder(
				nil, nil, setter, nil, container.NewBorder(
					nil, nil, widget.NewSeparator(), nil, form,
				),
			),
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		// other
		dataBoxScroll,
	)

	myWindow.SetContent(allContainers)
	myWindow.ShowAndRun()
}

func main() {

	err := initialize()
	if err != nil {
		fmt.Printf("Error in initialisation : %q\n", err)
		keyWait()
		os.Exit(1)
	}
	gui()

}
