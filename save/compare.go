package save

import (
	"crypto/md5"
	"fmt"
	"main/app"
	"main/req"
	"main/tools"
	"os"
	"path/filepath"
	"strings"
)

func filter(f *app.Film) bool {
	if len(f.Addrs) == 0 {
		return true
	}
	t := f.Type
	n := f.Name
	if strings.Index(t, "福利") != -1 || strings.Index(n, "说电影") != -1 || strings.Index(t, "音乐MV") != -1 {
		return true
	}
	return false
}

func amd(f *app.Film) {
	a := ""
	for _, v := range f.Addrs {
		a += v.URL
	}
	f.Amd = fmt.Sprintf("%X", md5.Sum([]byte(a)))
}

// 保存图片
func saveImg(info *app.SiteConf, film *app.Film) string {
	u := film.Img
	if strings.Index(u, "http") == -1 {
		u = info.Host + u
	}
	buf, suf := req.Download(u)
	if len(buf) == 0 {
		tools.Log("\n" + u + ": 下载异常")
		return ""
	}
	ex := hasImg(buf)
	if ex != "" {
		return ex
	}
	wd, err := os.Getwd()
	if err != nil {
		tools.Log("\nos.Getwd 函数运行错误: " + err.Error())
		return ""
	}
	t := tools.Md(film.Type)
	p := filepath.Join(wd, "img", info.Name, t)
	os.MkdirAll(p, os.ModePerm)
	name := tools.Md(film.Name) + "." + suf
	p = filepath.Join(p, name)
	f, err := os.Create(p)
	if err != nil {
		tools.Log("\n下载图片时创建文件错误：" + err.Error())
		return ""
	}
	defer f.Close()
	f.Write(buf)
	lp := "/img/" + info.Name + "/" + t + "/" + name
	imgMD5(buf, lp)
	return lp
}

// Comp 比对保存
func Comp(film app.Film, info *app.SiteConf) {
	if filter(&film) {
		return
	}
	amd(&film)
	id, site, md := SomeInfo(film.Name)
	if id == 0 {
		film.Stat = 1
		film.Img = saveImg(info, &film)
		if InsertFilm(&film) {
			tools.Log("\n新增数据：" + film.Name)
		} else {
			tools.Log("\n新增数据失败：" + film.Name)
		}
		return
	}
	if site != film.Site {
		return
	}
	if md != film.Amd {
		film.ID = id
		if UpdateAddr(&film) {
			SetAmd(film.ID, film.Amd)
			tools.Log("\n更新播放地址：" + film.Name)
		} else {
			tools.Log("\n更新播放地址失败：" + film.Name)
		}
	}
}
