package main

import (
	"fmt"
	id3 "github.com/mikkyang/id3-go"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const MUSICDIR string = "/music/"

type Page struct {
	Title string
	File  File
	Music []File
}

type File struct {
	Name string
	Path string
	Size int64
	*ID3
}

type ID3 struct {
	Artist string
	Title  string
}

func getID3(p string) *ID3 {
	mp3File, err := id3.Open(p)
	if err != nil {
		fmt.Println(err)
	}
	return &ID3{mp3File.Artist(), mp3File.Title()}
}

func getFileInfo(p string) File {
	var file File

	f, err := os.Stat(p)
	if err != nil {
		fmt.Println(err)
	}

	id3 := getID3(p)

	file.Path = p
	file.Name = f.Name()
	file.Size = f.Size()
	file.ID3 = id3

	return file
}

func findMusic(p string) []File {
	var files []File

	find := func(p string, f os.FileInfo, err error) error {
		if !f.IsDir() && path.Ext(p) == ".mp3" {
			files = append(files, getFileInfo(p))
		}
		return nil
	}

	filepath.Walk(p, find)

	return files
}

func indexPage(title string, files []File) *Page {
	return &Page{Title: title, Music: files}
}

func playPage(title string, file File) *Page {
	return &Page{Title: title, File: file}
}

func musicHandler(w http.ResponseWriter, r *http.Request) {
	files := findMusic(MUSICDIR)
	p := indexPage("Music", files)
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, p)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path[len("/play"):]
	file := getFileInfo(filename)
	p := playPage("Play", file)
	t, _ := template.ParseFiles("templates/play.html")
	t.Execute(w, p)
}

func serve(port string) {
	http.Handle("/music/", http.StripPrefix("/music"+MUSICDIR, http.FileServer(http.Dir(MUSICDIR))))
	http.HandleFunc("/", musicHandler)
	http.HandleFunc("/play/", playHandler)
	http.ListenAndServe(":"+port, nil)
}

func main() {
	serve("8181")
}
