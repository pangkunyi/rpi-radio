package main

import (
	"github.com/gorilla/mux"
	"net/http"
	_ "net/http/pprof"
	"player"
	"strconv"
)

const (
	PORT = 8808
)

func main() {
	go player.StartAndWait()
	r := mux.NewRouter()
	r.HandleFunc("/music/{name}.mp3", musicHandler)
	r.HandleFunc("/list", listHandler)
	r.HandleFunc("/play", playHandler)
	r.HandleFunc("/switch", pauseOrResumeHandler)
	r.HandleFunc("/next", nextHandler)
	r.HandleFunc("/settings", settingsHandler)
	r.HandleFunc("/", indexHandler)
	http.Handle("/", r)
	if err := http.ListenAndServe(":"+strconv.Itoa(PORT), nil); err != nil {
		panic(err)
	}
}
