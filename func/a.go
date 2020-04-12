package main

import (
	"log"
	"net/http"

	"clwen.com/mail/p"
)

func main() {
	// send("hello there", "hoho", "")
	http.HandleFunc("/", p.F)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
