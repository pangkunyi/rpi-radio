package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/music/{name}.mp3", musicHandler)
	r.HandleFunc("/list", listHandler)
	http.Handle("/", r)
	if err := http.ListenAndServe(":8808", nil); err != nil {
		panic(err)
	}
}
