package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiohime/findfiles/pkg/fzscreen"
)

type fzScreenWidget struct {
	Entry  *widget.Entry
	Screen *fyne.Container
}

func (c *fzScreenWidget) update(in []string) {
	if c.Entry.Text != "" {
		go func() {
			result.Objects = []fyne.CanvasObject{c.Entry}
			fzres := fzscreen.Render(in, c.Entry.Text)
			for _, obj := range fzres.Objects {
				// fmt.Println("%%%", obj)
				result.Objects = append(result.Objects, obj)
			}
		}()
	}
}

var result *fyne.Container
var prevText string

func inputSlice(path string) []string {
	// a := []string{"fyoobar111", "fayzaza2r", "fofbyar"}
	fmt.Println("inputslice")

	file, err := os.Open(path)
	if err != nil {
		fmt.Errorf("error : %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var line string
	result := []string{}
	for {
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		result = append(result, line)
	}
	return result

}

func main() {

	theApp := app.New()
	app_window := theApp.NewWindow("")
	app_window.Resize(fyne.NewSize(1000, 500))
	entry_widget := widget.NewEntry()

	dummy_screen := fyne.NewContainer()
	// entryText := entry_widget.Text

	texter := &fzScreenWidget{
		Entry:  entry_widget,
		Screen: dummy_screen,
	}
	tester := inputSlice("c:\\Users\\apak\\.config\\kiohime\\root.txt")
	// fmt.Println(inputSlice("c:\\Users\\apak\\.config\\kiohime\\sysvol.txt"))
	entry_widget.OnCursorChanged = func() {
		text := entry_widget.Text
		fmt.Printf("%v/%v\n", text, prevText)
		if text != prevText && len(text) >= 3 {
			texter.update(tester)
		}
		prevText = text
	}

	result = container.NewBorder(
		texter.Entry,
		nil,
		nil,
		nil,
		fyne.NewContainer(texter.Screen),
	)

	app_window.SetContent(result)
	app_window.ShowAndRun()
}
