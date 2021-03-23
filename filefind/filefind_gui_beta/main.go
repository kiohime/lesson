package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
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
	quitProcess    chan bool

	appMode        int
	workDir        = ""
	baseNameFiles  = "filesearch_files.txt"
	baseNameDirs   = "filesearch_directory.txt"
	resultFileName = "result.txt"
)
var (
	progressBar   *widget.ProgressBarInfinite
	entry_widget  *modifiedEntry
	screen_widget *widget.Label
)

// /////////////////////////////////////////////////////////////////////////
type modifiedEntry struct {
	widget.Entry
	input *Input_widget
}

// newEscapeEntry - rewriting basic entry widget
func newModifiedEntry() *modifiedEntry {
	entry := &modifiedEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func newEntry() *modifiedEntry {
	return newModifiedEntry()
}

// onEsc - clears entry
func (e *modifiedEntry) onEsc() {
	// fmt.Println(e.Entry.Text)
	e.Entry.SetText("")
}

func (e *modifiedEntry) onEnter() {
	assert(e)
	assert(e.input)
	// findBtn(input_widget.Entry, screen_widget)
	fmt.Printf("on onEnter AppSet is %v\n", AppSet.AppMode)
	decider(e.input.Entry, screen_widget, AppSet)

}

func (e *modifiedEntry) setInput(input *Input_widget) {
	e.input = input
}

// /////////////////////////////////////////////////////////////
type modifiedSelect struct {
	widget.Select
	settings *Input_widget
}

func newModifiedSelect() *modifiedSelect {
	selEntry := &modifiedSelect{}
	selEntry.ExtendBaseWidget(selEntry)
	return selEntry
}

func (s *modifiedSelect) setSettings(settings *Input_widget) {
	assert(s)
	s.settings = settings
}

func newSelect() *modifiedSelect {
	s := newModifiedSelect()

	return s
}

func newSettings(m *widget.RadioGroup) *modifiedSelect {
	check := newSelect()
	t := []string{"каталоги", "файлы", "И ТО, И ДРУГОЕ"}
	check.Options = t
	check.OnChanged = settingsChanged
	// check.Disable()
	return check
}

func settingsChanged(c string) {
	fmt.Println(c)
	switch c {
	case "каталоги":
		AppSet.ScanMode = 0
	case "файлы":
		AppSet.ScanMode = 1
		// case "И ТО, И ДРУГОЕ":
	}
}

// //////////////////////////////////////////////////////////
type AppSettings struct {
	AppMode  int
	ScanMode int
}

var AppSet *AppSettings

// ////////////////////////////////////////////////////////////////////////
type Input_widget struct {
	Mode   *widget.RadioGroup
	Entry  *modifiedEntry
	Form   *widget.Form
	Option *modifiedSelect
}

func newInputWidget(m *widget.RadioGroup, e *modifiedEntry, f *widget.Form, s *modifiedSelect) *Input_widget {
	assert(m, e, f, s)
	return &Input_widget{Mode: m, Entry: e, Form: f, Option: s}
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

// пауза
func keyWait() {
	fmt.Printf("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// блок инициализации: установка рабочего пути для файлов базы и поиска
func initialize() error {
	fmt.Println("### initialize")
	AppSet = &AppSettings{
		AppMode:  appMode,
		ScanMode: scanMode,
	}
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
func executer(aset *AppSettings) string {
	result := ""
	fmt.Printf("on executer AppMode is %v\n", aset.AppMode)
	switch aset.AppMode {
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
		// scanDir = true

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

// newModeWidget
func newModeWidget(aset *AppSettings) *widget.RadioGroup {
	s := widget.NewRadioGroup([]string{"Поиск", "Сканирование"}, func(s string) {
		switch s {
		case "Поиск":
			aset.AppMode = 0
		case "Сканирование":
			aset.AppMode = 1
		}
		fmt.Printf("set AppMode %v\n", aset.AppMode)
	})
	s.SetSelected("Поиск")
	s.Refresh()
	return s
}

// decider - executes searching mode, depending on switch
func decider(input *modifiedEntry, scr *widget.Label, aset *AppSettings) {
	assert(input, scr)

	fmt.Printf("on decider AppMode is %v\n", aset.AppMode)

	progressBar.Start()
	progressBar.Show()
	argCache = nil
	rootDir = ""
	switch aset.AppMode {
	case 0:
		argCache = append(argCache, input.Text)
	case 1:
		rootDir = input.Text
	}
	result := executer(aset)
	// dataBox.Add(widget.NewLabel(result))
	scr.Text = result
	input.Text = ""
	input.Refresh()
	progressBar.Hide()
	progressBar.Stop()

}

func newForm(i *modifiedEntry, s *widget.Label) *widget.Form {
	assert(i, s)
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Data", Widget: i, HintText: "input data"},
		},
		OnCancel: nil,
		OnSubmit: i.onEnter,

		// OnSubmit: findBtn(i, s),
	}
	i.Validator = validation.NewRegexp(`.+`, "input smthing")
	return form
}

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

func makeContainerTree(i *Input_widget) *fyne.Container {
	screen_container := container.NewWithoutLayout(screen_widget)
	screen_scroll_container := container.NewScroll(screen_container)
	mode_settings_container := container.NewVBox(
		i.Mode,
		widget.NewSeparator(),
		i.Option,
		// widget.NewSeparator(),
	)
	form_container := container.NewBorder(
		nil, nil, nil, nil, i.Form,
	)
	all_container := container.NewBorder(
		// top
		container.NewVBox(
			container.NewBorder(
				nil, nil, mode_settings_container, nil, form_container,
			),
			widget.NewSeparator(),
			progressBar,
		),
		nil,
		nil,
		nil,
		// other
		screen_scroll_container,
	)
	return all_container
}

func gui() {
	the_app := app.New()
	app_window := the_app.NewWindow("Notepad")
	app_window.Resize(fyne.NewSize(1000, 500))

	progressBar = widget.NewProgressBarInfinite()
	progressBar.Hide()

	screen_widget = widget.NewLabel("")
	mode_widget := newModeWidget(AppSet)
	entry_widget := newEntry()
	form_widget := newForm(entry_widget, screen_widget)
	settings_widget := newSettings(mode_widget)
	settings_widget.Selected = "каталоги"
	settings_widget.OnChanged("каталоги")

	input_widget := newInputWidget(mode_widget, entry_widget, form_widget, settings_widget)
	entry_widget.setInput(input_widget)
	settings_widget.setSettings(input_widget)

	tree_container := makeContainerTree(input_widget)

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
