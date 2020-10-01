package main
 
import (
	"github.com/gotk3/gotk3/gtk"
	"encoding/csv"
	"fmt"
	"os"
)
var Questions []string
var Answers []string
func main() {
	//readFile("data.csv",&s);
	//fmt.Println(s[9][0])
	showWindow();
}
func readFile(filename string){
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comment = '#'
	counter:=0; 
	for {
		record,e := reader.Read()
		if e != nil {
			break
		}
		Questions=append(Questions,record[0]);
		Answers=append(Answers,record[1]);
		fmt.Println(counter);
	}
}
func showWindow(){
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
    obj, err = b.GetObject("open_menubtn")
    if err != nil {
        panic(err)
    }
    openFileBtn := obj.(*gtk.MenuItem)
    openFileBtn.Connect("activate", func() {
	    if err == nil {
    	    filename := openFileDialog(win)
    	    readFile(filename)
    	    fmt.Println(Questions[0])
    	}
	})
    win.ShowAll()
    gtk.Main()
}
func openFileDialog(win *gtk.Window) string{
	openDialog, err := gtk.FileChooserDialogNewWith2Buttons("Select files", win, gtk.FILE_CHOOSER_ACTION_OPEN,
		"Cancel", gtk.RESPONSE_CANCEL, "OK", gtk.RESPONSE_OK)
	if err != nil {
    	panic(err)
	}
	response := openDialog.Run()
	if response != gtk.RESPONSE_OK {
    	//panic(err)
	}
	file := openDialog.GetFilename()
	openDialog.Destroy()
	return file
}
