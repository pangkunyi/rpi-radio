package main

import (
	"fmt"
	"github.com/gorilla/mux"
	pcs "github.com/pangkunyi/baidu-pcs"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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

func listHandler(w http.ResponseWriter, r *http.Request) {
	req := pcs.NewListReq()
	req.AccessToken = pcs.ACCESS_TOKEN
	req.Path = fmt.Sprintf(`/apps/kunyi`)

	fmt.Printf("%#v\n", req)
	re := regexp.MustCompile(`\+`)
	if resp, err := pcs.List(req); err == nil {
		for _, fi := range resp.Files {
			if strings.HasSuffix(fi.Path, ".mp3") {
				fmt.Fprintf(w, "%s ", "http://127.0.0.1:8808/music/"+re.ReplaceAllString(url.QueryEscape(fi.Path[12:]), "%20"))
			}
		}
		return
	} else {
		panic(err)
	}
	http.Error(w, "Sever Error", 505)
}
