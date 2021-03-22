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
	progressBar   *widget.ProgressBarInfinite
	inputTab      Input_widget
	input_widget  *modifiedEntry
	screen_widget *widget.Label
)

type modifiedEntry struct {
	widget.Entry
}

type Input_widget struct {
	Mode  *widget.RadioGroup
	Entry *modifiedEntry
}

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

// findBtn - executes searching mode, depending on switch
func findBtn(input *modifiedEntry, scr *widget.Label) func() {

	assert("findBtn", input != nil, scr != nil)
	return func() {
		progressBar.Start()
		progressBar.Show()
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
		progressBar.Hide()
		progressBar.Stop()
	}
}

// onEsc - clears entry
func (e *modifiedEntry) onEsc() {
	fmt.Println(e.Entry.Text)
	e.Entry.SetText("")
}

func (e *modifiedEntry) onEnter() {

	fn := findBtn(inputTab.Entry, screen_widget)
	fn()

}

func makeForm(i *modifiedEntry, s *widget.Label) *widget.Form {
	assert("makeForm", i != nil, s != nil)
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Data", Widget: i, HintText: "input data"},
		},
		OnCancel: nil,
		// OnCancel: doCancel(),
		OnSubmit: findBtn(i, s),
	}
	i.Validator = validation.NewRegexp(`.+`, "input smthing")
	return form
}

// func doCancel() func() {
// 	return func() {}
// }

// assert - выдает сообщения об ошибке если есть nil
func assert(label string, args ...bool) {
	for i, arg := range args {
		if !arg {
			fmt.Printf("arg %v\narg %v\n", i, arg)
			panic(label)
		}
	}
}

func makeContainerTree(modeWidget *widget.RadioGroup, entry *modifiedEntry, form *widget.Form) *fyne.Container {
	assert("makeContainerTree", modeWidget != nil, entry != nil, form != nil)
	screen_container := container.NewWithoutLayout(screen_widget)
	screen_scrool_container := container.NewScroll(screen_container)
	all_container := container.NewBorder(
		// top
		container.NewVBox(
			container.NewBorder(
				nil, nil, modeWidget, nil, container.NewBorder(
					nil, nil, widget.NewSeparator(), nil, form,
				),
			),
			widget.NewSeparator(),
			progressBar,
		),
		nil,
		nil,
		nil,
		// other
		screen_scrool_container,
	)
	return all_container
}

func gui() {
	the_app := app.New()
	app_window := the_app.NewWindow("Notepad")
	app_window.Resize(fyne.NewSize(1000, 500))

	mode_widget := makeModeSelection_widget()
	entry_widget := newModifiedEntry()

	inputTab = Input_widget{
		Mode:  mode_widget,
		Entry: entry_widget,
	}
	screen_widget = widget.NewLabel("")

	form_widget := makeForm(inputTab.Entry, screen_widget)
	progressBar = widget.NewProgressBarInfinite()
	progressBar.Hide()

	tree_container := makeContainerTree(inputTab.Mode, inputTab.Entry, form_widget)

	app_window.SetContent(tree_container)
	app_window.ShowAndRun()
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
