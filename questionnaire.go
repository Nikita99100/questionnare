package main

import (
	"encoding/csv"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	COLUMN_NAME = iota
	COLUMN_SCORE
)

var QuestionBox *gtk.TextView
var AnswerBox *gtk.Entry
var StartBtn *gtk.Button
var Notebook *gtk.Notebook
var TreeViewListStore *gtk.ListStore
var ScoreLabel *gtk.Label
var TimeBar *gtk.ProgressBar
var Questions []string
var Answers []string
var Finished = false
var QuestionNumber = 0
var Filename = ""
var Name string
var Score = 0
var TimeRemaining = 20.0
var IsTimeUp = false

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
	obj, err = b.GetObject("score_label")
	if err != nil {
		panic(err)
	}
	ScoreLabel = obj.(*gtk.Label)
	obj, err = b.GetObject("question_box")
	if err != nil {
		panic(err)
	}
	QuestionBox = obj.(*gtk.TextView)
	obj, err = b.GetObject("notebook")
	if err != nil {
		panic(err)
	}
	Notebook = obj.(*gtk.Notebook)
	obj, err = b.GetObject("open_btn")
	if err != nil {
		panic(err)
	}
	openBtn := obj.(*gtk.MenuItem)
	openBtn.Connect("activate", func() {
		Filename = openFileDialog(win)
		startBtnCheck()
	})
	obj, err = b.GetObject("start_btn")
	if err != nil {
		panic(err)
	}
	StartBtn = obj.(*gtk.Button)
	StartBtn.Connect("clicked", startBtnPressed)
	obj, err = b.GetObject("nxt_btn")
	if err != nil {
		panic(err)
	}
	nextBtn := obj.(*gtk.Button)
	nextBtn.Connect("clicked", nextBtnPressed)
	obj, err = b.GetObject("name_text")
	if err != nil {
		panic(err)
	}
	nameText := obj.(*gtk.Entry)
	obj, err = b.GetObject("time_bar")
	if err != nil {
		panic(err)
	}
	TimeBar = obj.(*gtk.ProgressBar)
	obj, err = b.GetObject("answer_box")
	if err != nil {
		panic(err)
	}
	AnswerBox = obj.(*gtk.Entry)
	AnswerBox.Connect("activate", nextBtnPressed)
	nameText.Connect("changed", func() {
		Name, _ = nameText.GetText()
		startBtnCheck()
	})
	obj, err = b.GetObject("score_tree")
	if err != nil {
		panic(err)
	}
	scoreTree := obj.(*gtk.TreeView)
	setupTreeView(scoreTree)
	win.ShowAll()
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
		return ""

	}
	file := openDialog.GetFilename()
	openDialog.Destroy()
	return file
}
func readQuestions() {
	file, err := os.Open(Filename)
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
		shuffleQuestions()
	}
}
func shuffleQuestions() {
	destq := make([]string, len(Questions))
	desta := make([]string, len(Answers))
	perm := rand.Perm(len(Questions))
	for i, v := range perm {
		destq[v] = Questions[i]
		desta[v] = Answers[i]
	}
	Questions = destq
	Answers = desta
}
func readScores() map[string]int {
	file, err := os.Open("results.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comment = '#'
	scores := make(map[string]int)
	for {
		record, e := reader.Read()
		if e != nil {
			break
		}
		scores[record[0]], _ = strconv.Atoi(record[1])
	}
	return scores
}
func writeScores(scores map[string]int) {
	file, err := os.Create("results.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for key, value := range scores {
		writer.Write([]string{key, strconv.Itoa(value)})
	}
}

func startBtnPressed() {
	Notebook.NextPage()
	readQuestions()
	QuestionNumber = 0
	setQuestionText(Questions[0])
	go countTime()
}

func nextBtnPressed() {
	QuestionNumber++
	answer, _ := AnswerBox.GetText()
	AnswerBox.SetText("")
	if Answers[QuestionNumber-1] == answer {
		Score++
	}
	if QuestionNumber > len(Questions)-1 || IsTimeUp {
		end()
	} else {
		setQuestionText(Questions[QuestionNumber])
	}
}
func countTime() {
	curTime := TimeRemaining
	for curTime > 0 {
		TimeBar.SetFraction(curTime / TimeRemaining)
		time.Sleep(time.Second)
		curTime--
	}
	TimeBar.SetFraction(0)
	IsTimeUp = true
}
func end() {
	Notebook.NextPage()
	ScoreLabel.SetText(strconv.Itoa(Score) + " правильных ответов из " + strconv.Itoa(len(Questions)))
	score := readScores()
	score[Name] = Score
	fillTable(score)
	writeScores(score)
	Finished = true
}
func startBtnCheck() {
	if Name != "" && Filename != "" {
		StartBtn.SetSensitive(true)
	} else {
		StartBtn.SetSensitive(false)
	}
}
func fillTable(scores map[string]int) {
	for key, value := range scores {
		addRow(TreeViewListStore, key, strconv.Itoa(value))
	}
}
func setQuestionText(text string) {
	buffer, _ := QuestionBox.GetBuffer()
	buffer.SetText(text)
	QuestionBox.SetBuffer(buffer)
}
func setupTreeView(treeView *gtk.TreeView) {
	treeView.AppendColumn(createColumn("Имя", COLUMN_NAME))
	treeView.AppendColumn(createColumn("Баллы", COLUMN_SCORE))
	TreeViewListStore, _ = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
	treeView.SetModel(TreeViewListStore)
}
func addRow(listStore *gtk.ListStore, name, score string) {
	iter := listStore.Append()
	err := listStore.Set(iter,
		[]int{COLUMN_NAME, COLUMN_SCORE},
		[]interface{}{name, score})
	if err != nil {
		panic(err)
	}
}
func createColumn(title string, id int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		panic(err)
	}
	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		panic(err)
	}
	return column
}
