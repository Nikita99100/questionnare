package main

import (
	"encoding/csv"
	"github.com/gotk3/gotk3/gtk"
	"math/rand"
	"os"
	"strconv"
)

var QuestionLabel *gtk.Label
var AnswerBox *gtk.Entry
var NextBtn *gtk.Button
var Questions []string
var Answers []string
var QuestionNumber = 0
var IsFileSelected = false
var IsStarted = false
var Name string
var Score = 0

func main() {
	showWindow()
}

func showWindow() {
	gtk.Init(nil)
	b, err := gtk.BuilderNew()
	if err != nil {
		panic(err)
	}
	err = b.AddFromFile("gui.glade")
	if err != nil {
		panic(err)
	}
	obj, err := b.GetObject("window_main")
	if err != nil {
		panic(err)
	}
	win := obj.(*gtk.Window)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	obj, err = b.GetObject("question_text")
	if err != nil {
		panic(err)
	}
	QuestionLabel = obj.(*gtk.Label)
	obj, err = b.GetObject("answer_box")
	if err != nil {
		panic(err)
	}
	AnswerBox = obj.(*gtk.Entry)
	obj, err = b.GetObject("open_menubtn")
	if err != nil {
		panic(err)
	}
	openFileBtn := obj.(*gtk.MenuItem)
	openFileBtn.Connect("activate", func() {
		if err == nil {
			filename := openFileDialog(win)
			readFile(filename)
		}
	})
	obj, err = b.GetObject("next_btn")
	if err != nil {
		panic(err)
	}
	NextBtn = obj.(*gtk.Button)
	NextBtn.Connect("clicked", nextButtonPressed)
	win.ShowAll()
	AnswerBox.Hide()
	NextBtn.SetSensitive(false)
	gtk.Main()
}
func openFileDialog(win *gtk.Window) string {
	openDialog, err := gtk.FileChooserDialogNewWith2Buttons("Select files", win, gtk.FILE_CHOOSER_ACTION_OPEN,
		"Отмена", gtk.RESPONSE_CANCEL, "Выбрать", gtk.RESPONSE_OK)
	if err != nil {
		panic(err)
	}
	var filter *gtk.FileFilter
	filter, err = gtk.FileFilterNew()
	filter.SetName("csv таблица")
	filter.AddPattern("*.csv")

	openDialog.AddFilter(filter)
	response := openDialog.Run()
	if response != gtk.RESPONSE_OK {
		openDialog.Destroy()
		IsFileSelected = false
		return ""

	}
	IsFileSelected = true
	NextBtn.SetSensitive(true)
	file := openDialog.GetFilename()
	openDialog.Destroy()
	return file
}
func readFile(filename string) {
	if !IsFileSelected {
		return
	}
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comment = '#'
	for {
		record, e := reader.Read()
		if e != nil {
			break
		}
		Questions = append(Questions, record[0])
		Answers = append(Answers, record[1])
		//Questions = shuffleQuestions(Questions)
	}
}
func shuffleQuestions(src []string) []string {
	dest := make([]string, len(src))
	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}
	return dest
}
func nextButtonPressed() {

	if !IsStarted {
		if IsFileSelected {
			AnswerBox.SetVisible(true)
			IsStarted = true
		}
	} else {
		QuestionLabel.SetText(Questions[0])
		answer, _ := AnswerBox.GetText()
		if Answers[QuestionNumber-1] == answer {
			Score++
		}
	}

	if QuestionNumber > len(Questions)-1 {
		end()
	} else {
		QuestionLabel.SetText(Questions[QuestionNumber])
		QuestionNumber++
	}
}
func end() {
	NextBtn.Hide()
	QuestionLabel.SetText(strconv.Itoa(Score))

}
