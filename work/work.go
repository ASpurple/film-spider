package work

import (
	"main/app"
	"main/req"
	"main/tools"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Worker 运行
type Worker struct {
	Info    *app.SiteConf
	curPage int
	list    []string // 详情页URL数组
}

func (w *Worker) finished() bool {
	if len(w.list) < 1 {
		return true
	}
	if app.StopPage > 0 && w.curPage > app.StopPage {
		return true
	}
	return false
}

var tmp = ""

func (w *Worker) readList() bool {
	p := strconv.Itoa(w.curPage)
	tools.Log("\n开始提取：" + w.Info.Name + " 的第" + p + "页")
	addr := strings.Replace(w.Info.ListURL, "*", p, 1)
	html := ""
	if tmp == "" {
		html = req.Get(addr)
	} else {
		html = tmp
		tmp = ""
	}
	if html == "" {
		tools.Log("\n" + addr + ": 请求异常，跳过此页")
		time.Sleep(time.Second * 3)
		return true
	}
	w.list = getList(html, w.Info.HrefSep, w.Info.Host)
	go func() {
		nx := strconv.Itoa(w.curPage + 1)
		ad := strings.Replace(w.Info.ListURL, "*", nx, 1)
		if app.StopPage > 0 && (w.curPage+1) <= app.StopPage {
			tmp = req.Get(ad)
		}
	}()
	w.readDetail()
	w.curPage++
	return w.finished()
}

func (w *Worker) readDetail() {
	var wt sync.WaitGroup
	le := len(w.list)
	for i := 0; i < le; i++ {
		wt.Add(1)
		go func(index int) {
			getDetail(w.list[index], w.Info)
			wt.Done()
		}(i)
	}
	wt.Wait()
}

func newWorker(info *app.SiteConf) Worker {
	w := Worker{
		Info:    info,
		curPage: app.StartPage,
		list:    make([]string, 0),
	}
	if app.StartPage <= 0 {
		w.curPage = 1
	}
	return w
}

// StartAll GO!
func StartAll() {
	for _, s := range app.Sites {
		w := newWorker(s)
		for !w.readList() {
		}
	}
}

// StartOne ..
func StartOne(i int) {
	s := app.Sites[i]
	w := newWorker(s)
	for !w.readList() {
	}
}
