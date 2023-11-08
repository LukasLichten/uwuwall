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

type State struct {
	Database        []Message
	ElementsPerPage int
	PageTemlate     *template.Template
}

func main() {
	log.Print("Stawting Sewvew")

	var state = State{
		ElementsPerPage: 2,
		PageTemlate:     template.Must(template.ParseFiles("templates/page.html")),
	}

	state.Database = append(state.Database, Message{User: "", Message: "Test", Time: "A long time ago"})
	state.Database = append(state.Database, Message{User: "", Message: "Test 2", Time: "A long time ago"})
	state.Database = append(state.Database, Message{User: "", Message: "Test 3", Time: "A long time ago"})

	r := mux.NewRouter()

	r.HandleFunc("/", state.root_page)

	r.HandleFunc("/page/{page}", state.some_page)

	log.Fatal(http.ListenAndServe(":6969", r))

}

func (state State) root_page(w http.ResponseWriter, r *http.Request) {
	state.render_page(w, r, 1)
}

func (state State) some_page(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageString := vars["page"]

	pageNum, e := strconv.ParseInt(pageString, 10, 64)

	if e != nil {
		fmt.Fprint(w, e.Error())
		return
	}
	pageNumber := int(pageNum)

	state.render_page(w, r, pageNumber)
}

func (state State) render_page(w http.ResponseWriter, r *http.Request, pageNumber int) {
	log.Print("Page call " + strconv.FormatInt(int64(len(state.Database)), 10))

	if float64(pageNumber) > (float64(len(state.Database))/float64(state.ElementsPerPage) + 0.5) {
		fmt.Fprint(w, "Thewe whewe wess messages then expewcted...")
		return
	} else if pageNumber < 1 {
		fmt.Fprint(w, "Page must bwe numbewed 1 ow highew")
		return
	}

	if r.Method == http.MethodPost {
		// New Message
		msg := r.FormValue("new_message")
		state.Database = append(state.Database, Message{User: "", Message: msg, Time: "A long time ago"})
		// We continuwe nowmawy
		// We shouwd maybe wediwect to get, as bwowsew gets confuwsed with ouw fowm

		// Messages awe wost anyway wight now *sad uwu*
	}

	pageStart := len(state.Database) - state.ElementsPerPage*pageNumber
	pageEnd := len(state.Database) - state.ElementsPerPage*(pageNumber-1)

	if 0 > pageStart {
		pageStart = 0
	}

	doit := slices.Clone[[]Message](state.Database)[pageStart:pageEnd]
	slices.Reverse[[]Message](doit)

	page := Page{
		PageNumber: int(pageNumber),
		Messages:   doit,
	}
	state.PageTemlate.Execute(w, page)
}
