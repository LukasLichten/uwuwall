package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	User    string
	Message string
	Time    time.Time
}

type Page struct {
	IsRoot      bool
	IsFirstPage bool
	IsLastPage  bool
	PageNumber  int
	PrevPage    int
	NextPage    int
	Messages    []Message
}

type State struct {
	DB              *sql.DB
	Dep             []Message
	ElementsPerPage int
	PageTemlate     *template.Template
}

func main() {
	log.Print("Stawting Sewvew")

	// Sewtting up DB
	CatchAndPanic(os.MkdirAll("data/", 0755))

	db, error := sql.Open("sqlite3", "file:data/database.db")

	CatchAndPanic(error)
	CatchAndPanic(db.Ping())

	query := `
    CREATE TABLE IF NOT EXISTS messages (
        id INT AUTO_INCREMENT,
        user TEXT NOT NULL,
        text TEXT NOT NULL,
        time DATETIME,
        PRIMARY KEY (id)
    );`

	_, error = db.Exec(query)
	CatchAndPanic(error)

	// Sewtting up State
	var state = State{
		DB:              db,
		ElementsPerPage: 4,
		PageTemlate:     template.Must(template.ParseFiles("templates/page.html")),
	}

	r := mux.NewRouter()

	r.HandleFunc("/", state.RootPage)

	r.HandleFunc("/page/{page}", state.SomePage)

	log.Print("Sewvew iws up!!!111")
	log.Fatal(http.ListenAndServe(":6969", r))

}

func CatchAndPanic(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func CatchAndError(e error, w http.ResponseWriter) bool {
	if e != nil {
		log.Print("Error: ", e)
		fmt.Fprintf(w, "Sowwy, thewe was an ewwow in the sewvew")
		return true
	}
	return false
}

func (state State) RootPage(w http.ResponseWriter, r *http.Request) {
	state.RenderPage(w, r, 1, true)
}

func (state State) SomePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageString := vars["page"]

	pageNum, e := strconv.ParseInt(pageString, 10, 64)

	if e != nil {
		fmt.Fprint(w, e.Error())
		return
	}
	pageNumber := int(pageNum)

	state.RenderPage(w, r, pageNumber, false)
}

func (state State) RenderPage(w http.ResponseWriter, r *http.Request, pageNumber int, isRoot bool) {
	var messageCount int
	e := state.DB.QueryRow(`SELECT COUNT(*) FROM messages`).Scan(&messageCount)

	if CatchAndError(e, w) {
		return
	}

	if float64(pageNumber) > (float64(messageCount)/float64(state.ElementsPerPage)+0.5) && pageNumber != 1 {
		fmt.Fprint(w, "Thewe whewe wess messages then expewcted...")
		return
	} else if pageNumber < 1 {
		fmt.Fprint(w, "Page must bwe numbewed 1 ow highew")
		return
	}

	if r.Method == http.MethodPost {
		// New Message
		msg := Message{User: "", Message: r.FormValue("new_message"), Time: time.Now()}

		_, e := state.DB.Exec(`INSERT INTO messages (user, text, time) VALUES (?, ?, ?)`, msg.User, msg.Message, msg.Time)

		if CatchAndError(e, w) {
			return
		}

		messageCount += 1

		// We continuwe nowmawy
		// We shouwd maybe wediwect to get, as bwowsew gets confuwsed with ouw fowm (F5 awsks if we want to resuwbmit fowm)
	}

	rows, e := state.DB.Query(`SELECT user, text, time FROM messages ORDER BY time DESC LIMIT ? OFFSET ?`, state.ElementsPerPage, state.ElementsPerPage*(pageNumber-1))

	defer rows.Close()
	if CatchAndError(e, w) {
		return
	}

	var doit []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.User, &msg.Message, &msg.Time)
		if err == nil {
			doit = append(doit, msg)
		}
	}

	page := Page{
		IsRoot:      isRoot,
		IsFirstPage: pageNumber == 1,
		IsLastPage:  pageNumber*state.ElementsPerPage > messageCount,
		PageNumber:  int(pageNumber),
		PrevPage:    pageNumber - 1,
		NextPage:    pageNumber + 1,
		Messages:    doit,
	}
	state.PageTemlate.Execute(w, page)
}
