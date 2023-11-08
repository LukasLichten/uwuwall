package main

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

type Message struct {
	User    string
	Message string
	Time    string
}

type Page struct {
	PageNumber int
	Messages   []Message
}

func main() {
	log.Print("Stawting Sewvew")

	var database []Message // Stawte stowage www
	database = append(database, Message{User: "", Message: "Test", Time: "A long time ago"})
	database = append(database, Message{User: "", Message: "Test 2", Time: "A long time ago"})
	database = append(database, Message{User: "", Message: "Test 3", Time: "A long time ago"})

	tmpl := template.Must(template.ParseFiles("templates/page.html"))

	r := mux.NewRouter()

	elements_per_page := 4

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// fiwst page
		counter := elements_per_page
		if len(database) < counter {
			counter = len(database)
		}

		pageStart := len(database) - counter
		if pageStart < 0 {
			pageStart = 0
		}

		doit := slices.Clone[[]Message](database)[pageStart:len(database)]
		page := Page{
			PageNumber: 1,
			Messages:   doit,
		}
		tmpl.Execute(w, page)
	})

	r.HandleFunc("/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		// Some Pawge
		vars := mux.Vars(r)
		pageString := vars["page"]

		pageNum, e := strconv.ParseInt(pageString, 10, 64)
		pageNumber := int(pageNum)

		counter := elements_per_page

		if e != nil {
			fmt.Fprint(w, e.Error())
			return
		} else if float64(pageNum) > (float64(len(database))/float64(counter) + 0.5) {
			fmt.Fprint(w, "Thewe whewe wess messages then expewcted...")
			return
		} else if pageNumber < 1 {
			fmt.Fprint(w, "Page must bwe numbewed 1 ow highew")
			return
		}

		pageStart := counter * (pageNumber - 1)
		pageEnd := counter * pageNumber

		if len(database) < pageEnd {
			pageEnd = len(database)
		}

		doit := slices.Clone[[]Message](database)[pageStart:pageEnd]

		page := Page{
			PageNumber: int(pageNumber),
			Messages:   doit,
		}
		tmpl.Execute(w, page)
	})

	log.Fatal(http.ListenAndServe(":6969", r))

}
