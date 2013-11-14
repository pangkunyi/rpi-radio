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
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

const (
	TPL_HEADER = `
<html>
    <head>
        <title>rpi radio</title>
        <style>
            body{
                background:black;
            }
            a, label{
                color: rgb(170, 143, 143);
            }
            #ctrl{
                padding:0;
            }
            #ctrl li{
                width: 100%%;
                text-align: center;
                margin: 4px;
                font-size: 1.3em;
            }               
            #settings_form div{
                width: 100%%;
                text-align: center;
            }
        </style>
    </head>
<body>
`
	TPL_FOOTER = `
</body>
</html>
`
	TPL_CTRL = TPL_HEADER + `
<ul id="ctrl">
    <li><a href="/play">加载播放</a></li>
    <li><a href="/next">下一首歌</a></li>
    <li><a href="/switch">播放暂停</a></li>
    <li><a href="/settings">播放设置</a></li>
</ul>
` + TPL_FOOTER

	TPL_SETTINGS = TPL_HEADER + `
<form id="settings_form" action="/settings" method="POST">
    <div>
        <label>数目：</label>
        <input type="text" name="maxSize" value="%d"/>
    </div>
    <div>
        <label>目录：</label>
        <input type="text" name="folder" value="%s"/>
    </div>
    <div>
        <input type="submit" value="保存"/>
        <a href="/">返回</a>
    </div>
</form>
` + TPL_FOOTER
)

var (
	playing           = make(chan int32, 1)
	idx         int32 = 0
	maxSize           = 20
	BASE_FOLDER       = pcs.OPEN_DIR
	folder            = ""
)

func sendText(w http.ResponseWriter, text string) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, text)
}

func sendErr(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	http.Error(w, "Sever Error: "+err.Error(), 505)
}

func gotoIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", 303)
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sendText(w, fmt.Sprintf(TPL_SETTINGS, maxSize, folder))
	} else if r.Method == "POST" {
		if _maxSize, err := strconv.Atoi(r.FormValue("maxSize")); err == nil {
			maxSize = _maxSize
		}
		folder = r.FormValue("folder")
		gotoIndex(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	sendText(w, TPL_CTRL)
}

func musicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	req := pcs.NewDownloadReq()
	req.AccessToken = pcs.ACCESS_TOKEN
	req.Path = fmt.Sprintf(`%s/%s/%s.mp3`, BASE_FOLDER, folder, vars["name"])

	fmt.Printf("%#v\n", req)
	if err := pcs.Download(req, w); err != nil {
		sendErr(w, err)
	}
}

func nextHandler(w http.ResponseWriter, r *http.Request) {
	if !player.IsStop() {
		if err := player.Next(); err != nil {
			sendErr(w, err)
			return
		}
	}
	gotoIndex(w, r)
}

func pauseOrResumeHandler(w http.ResponseWriter, r *http.Request) {
	if !player.IsStop() {
		if err := player.PauseOrResume(); err != nil {
			sendErr(w, err)
			return
		}
	}
	gotoIndex(w, r)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&idx, int32(1))
	if err := player.Stop(); err != nil {
		sendErr(w, err)
		return
	}

	if list, err := playlist(); err != nil {
		sendErr(w, err)
		return
	} else {
		list = random(list, 5)
		fmt.Printf("ready playlist:\n%#v\n", list)
		go func(list []string, _idx int32) {
			for _, fi := range list {
				player.ReadyToPlay()
				if atomic.LoadInt32(&idx) != _idx {
					fmt.Println("play break with ", _idx)
					player.ResetToPlay()
					break
				}
				fmt.Println("ready to play: ", fi)
				player.Play(fi)
			}
		}(list, atomic.LoadInt32(&idx))
	}
	gotoIndex(w, r)
}

func random(list []string, maxSize int) []string {
	size := len(list)
	count := 1
	for idx, _ := range list {
		tmp := list[idx]
		rIdx := rand.Int() % size
		list[idx] = list[rIdx]
		list[rIdx] = tmp
		if count++; count > maxSize {
			break
		}
	}
	if size > maxSize {
		return list[:maxSize]
	}
	return list
}

func playlist() (list []string, err error) {
	req := pcs.NewListReq()
	req.AccessToken = pcs.ACCESS_TOKEN
	req.Path = fmt.Sprintf("%s/%s", BASE_FOLDER, folder)

	fmt.Printf("%#v\n", req)
	re := regexp.MustCompile(`\+`)
	if resp, err := pcs.List(req); err == nil {
		list = make([]string, 0)
		for _, fi := range resp.Files {
			if strings.HasSuffix(fi.Path, ".mp3") {
				list = append(list, fmt.Sprintf("http://127.0.0.1:%d/music/%s", PORT, re.ReplaceAllString(url.QueryEscape(fi.Path[len(BASE_FOLDER)+len(folder)+1:]), "%20")))
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
		sendErr(w, err)
	}
}
