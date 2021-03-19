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

	appMode        int
	workDir        = ""
	baseNameFiles  = "filesearch_files.txt"
	baseNameDirs   = "filesearch_directory.txt"
	resultFileName = "result.txt"
)
var (
	input_widget  *modifiedEntry
	screen_widget *widget.Label
)

//bahniFile - создает файл с именем inputName из массива данных inputData, данные добавляются через ньюлайн
func bahniFile(inputName string, inputData *[]string) error {
	fmt.Println("#### bahnifile")
	// fmt.Println(inputName)
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

//executer - запускает программу в устновленном режиме
func executer() string {
	result := ""
	switch appMode {
	case 0:
		base, e := readBaser()
		if e != nil {
			a := fmt.Errorf("error in reading base : %v", e)
			fmt.Println(a)
			// os.Exit(1)
		}
		result = strings.Join(base, "\n")
	case 1:

		// Временно! пока галочек нет!!!!
		scanDir = true

		e := writeBaser()
		if e != nil {
			a := fmt.Errorf("error in writing base : %v", e)
			fmt.Println(a)
			// os.Exit(1)
		}
		// result = "000000000000000000000000000000"
		fmt.Println("THE END")
	}
	return result

}

// findBtn - executes searching mode, depending on switch
func findBtn(input *modifiedEntry, scr *widget.Label) func() {

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

type modifiedEntry struct {
	widget.Entry
}

// onEsc - clears entry
func (e *modifiedEntry) onEsc() {
	fmt.Println(e.Entry.Text)
	e.Entry.SetText("")
}

func (e *modifiedEntry) onEnter() {
	fn := findBtn(input_widget, screen_widget)
	fn()
}

// newEscapeEntry - rewriting basic entry widget
func newModifiedEntry() *modifiedEntry {
	entry := &modifiedEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

// TypedKey - overriding default TypedKey method in fyne.Focusable, adding switch
// ТУТ ПРЕСЕТЫ ДЛЯ КЛАВИАТУРЫ
func (e *modifiedEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyEscape:
		e.onEsc()
	case fyne.KeyEnter, fyne.KeyReturn:
		e.onEnter()
	default:
		e.Entry.TypedKey(key)
	}
}

// makeModeSelection_widget
func makeModeSelection_widget() *widget.RadioGroup {
	s := widget.NewRadioGroup([]string{"Поиск", "Сканирование"}, func(s string) {
		switch s {
		case "Поиск":
			appMode = 0
		case "Сканирование":
			appMode = 1
		}
	})
	s.SetSelected("Поиск")
	s.Refresh()
	return s
}

// func makeContainerTree(mode *widget.RadioGroup, ) *fyne.Container {

// 	return
// }

type Input_widget struct {
	Mode_widget  *widget.RadioGroup
	Input_widget *modifiedEntry
}

func makeForm(i *modifiedEntry, s *widget.Label) *widget.Form {
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Data", Widget: i, HintText: "input data"},
		},
		OnCancel: nil,
		OnSubmit: findBtn(i, s),
	}
	i.Validator = validation.NewRegexp(`.+`, "input smthing")
	return form
}

func assert(label string, args ...bool) {
	for _, arg := range args {
		if !arg {
			panic(label)
		}
	}
}

func makeContainerTree(modeWidget *widget.RadioGroup, entry *modifiedEntry, form *widget.Form) *fyne.Container {
	assert("makeContainerTree", modeWidget != nil, entry != nil, form != nil)
	dataBox := container.NewWithoutLayout(screen_widget)
	dataBoxScroll := container.NewScroll(dataBox)
	allContainers := container.NewBorder(
		// top
		container.NewVBox(
			container.NewBorder(
				nil, nil, modeWidget, nil, container.NewBorder(
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
	return allContainers
	// return nil
}

func gui() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Notepad")
	myWindow.Resize(fyne.NewSize(1000, 500))

	mode := makeModeSelection_widget()
	entry := newModifiedEntry()

	i := Input_widget{
		Mode_widget:  mode,
		Input_widget: entry,
	}
	// fmt.Printf("debug %v\n", i)
	screen_widget = widget.NewLabel("")

	form := makeForm(i.Input_widget, screen_widget)
	fmt.Printf("debug form %v\n", form)

	_ = form
	tree := makeContainerTree(i.Mode_widget, i.Input_widget, nil)
	fmt.Printf("debug tree %v\n", tree)

	myWindow.SetContent(tree)
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
