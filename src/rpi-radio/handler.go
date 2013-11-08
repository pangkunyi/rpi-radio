package main

import (
	"fmt"
	"github.com/gorilla/mux"
	pcs "github.com/pangkunyi/baidu-pcs"
	"math/rand"
	"net/http"
	"net/url"
	"player"
	"regexp"
	"strings"
	"sync/atomic"
)

var (
	playing       = make(chan int32, 1)
	idx     int32 = 0
)

func musicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	req := pcs.NewDownloadReq()
	req.AccessToken = pcs.ACCESS_TOKEN
	req.Path = fmt.Sprintf(`/apps/kunyi/%s.mp3`, vars["name"])

	fmt.Printf("%#v\n", req)
	if err := pcs.Download(req, w); err != nil {
		http.Error(w, "Sever Error", 505)
	}
}

func nextHandler(w http.ResponseWriter, r *http.Request) {
	if err := player.Next(); err != nil {
		http.Error(w, "Sever Error: "+err.Error(), 505)
	} else {
		fmt.Fprint(w, "OK!")
	}
}

func pauseOrResumeHandler(w http.ResponseWriter, r *http.Request) {
	if err := player.PauseOrResume(); err != nil {
		http.Error(w, "Sever Error: "+err.Error(), 505)
	} else {
		fmt.Fprint(w, "OK!")
	}
}
func playHandler(w http.ResponseWriter, r *http.Request) {
	//if cap(playing) == 1 {
	//	player.Next()
	//}
	playing <- atomic.AddInt32(&idx, int32(1))
	<-playing
	if list, err := playlist(); err != nil {
		http.Error(w, "Sever Error: "+err.Error(), 505)
		return
	} else {
		list = random(list)
		go func(list []string, idx int32) {
			for _, fi := range list {
				fmt.Println("ready to play: ", fi)
				playing <- idx
				player.PlayAndWait(fi)
				if idx != <-playing {
					fmt.Println("i am break")
					break
				}
			}
		}(list, idx)
	}
	fmt.Fprint(w, "OK!")
}

func random(list []string) []string {
	size := len(list)
	for idx, _ := range list {
		tmp := list[idx]
		rIdx := rand.Int() % size
		list[idx] = list[rIdx]
		list[rIdx] = tmp
	}
	return list
}

func playlist() (list []string, err error) {
	req := pcs.NewListReq()
	req.AccessToken = pcs.ACCESS_TOKEN
	req.Path = fmt.Sprintf(`/apps/kunyi`)

	fmt.Printf("%#v\n", req)
	re := regexp.MustCompile(`\+`)
	if resp, err := pcs.List(req); err == nil {
		list = make([]string, 0)
		for _, fi := range resp.Files {
			if strings.HasSuffix(fi.Path, ".mp3") {
				list = append(list, "http://127.0.0.1:8808/music/"+re.ReplaceAllString(url.QueryEscape(fi.Path[12:]), "%20"))
			}
		}
	}
	return
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	if list, err := playlist(); err == nil {
		for _, fi := range list {
			fmt.Fprintf(w, "%s ", fi)
		}
		return
	} else {
		http.Error(w, "Sever Error: "+err.Error(), 505)
	}
}
