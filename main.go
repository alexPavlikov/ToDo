package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

var take string
var takes string

func main() {
	fmt.Println("Listen on - " + cfg.ServerHost + ":" + cfg.ServerPort)
	err := connect()
	if err != nil {
		return
	}
	handleRequest()
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir(cfg.Assets))))
	err = http.ListenAndServe(":"+cfg.ServerPort, nil)
	if err != nil {
		log.Fatal("Error - ListenAndServe", err.Error())
	}
}

func handleRequest() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/form/add", addHandler)
	http.HandleFunc("/form/edit", editHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.NotFound(w, r)
	}
	data := selectCase()

	if r.Method == "POST" {
		r.ParseForm()
		take := r.FormValue("new_post_data")
		val, _ := strconv.Atoi(take)
		delCase(val)
	}

	r.ParseForm()
	id := 1
	id += incId()
	text := r.FormValue("text")
	date := r.FormValue("date")
	priority := r.FormValue("priority")
	if text != "" && date != "" {
		fmt.Println(text, date, priority)
		list.Id = id
		list.Text = text
		list.Date = date
		list.Rating, _ = strconv.Atoi(priority)
		addCase(list)
	}

	if r.Method == "GET" {
		r.ParseForm()
		takes = r.FormValue("post_data")
	}

	rows := map[string]interface{}{"Data": data}
	tmpl.Execute(w, rows)
}

var cases string

func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("form.html")
	if err != nil {
		http.NotFound(w, r)
	}

	tmpl.Execute(w, nil)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cases = r.FormValue("textChange")
	val, _ := strconv.Atoi(takes)
	selectById(val, cases)
	fmt.Fprint(w, "Изменения успешно выполненно!")
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := r.FormValue("text")
	date := r.FormValue("date")
	priority := r.FormValue("priority")
	fmt.Println(text, date, priority)
}
