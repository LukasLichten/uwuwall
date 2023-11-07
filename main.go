package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Print("Stawting Sewvew")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Hewwo " + r.RemoteAddr + " you want " + r.RequestURI + " uwu?")
		fmt.Fprint(w, "Wewcome to my website!")
	})

	var error = http.ListenAndServe(":6969", nil)

	log.Print("Sewvew Stopped!!!11")

	if error != nil {
		log.Fatal(error.Error())
	}

}
