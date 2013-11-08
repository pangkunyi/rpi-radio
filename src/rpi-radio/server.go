package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"player"
)

func main() {
	go player.StartAndWait()
	r := mux.NewRouter()
	r.HandleFunc("/music/{name}.mp3", musicHandler)
	r.HandleFunc("/list", listHandler)
	r.HandleFunc("/play", playHandler)
	r.HandleFunc("/switch", pauseOrResumeHandler)
	r.HandleFunc("/next", nextHandler)
	http.Handle("/", r)
	if err := http.ListenAndServe(":8808", nil); err != nil {
		panic(err)
	}
}
