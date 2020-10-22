package work

import (
	"main/app"
	"main/req"
	"main/save"
	"main/tools"
	"strings"
)

func exit(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func getList(html, sep, host string) (res []string) {
	arr := tools.GetAllMid(html, "<a", "</a")
	for _, a := range arr {
		if strings.Index(a, sep) == -1 {
			continue
		}
		i := strings.Index(a, "\"")
		if i < 0 || i == len(a)-1 {
			continue
		}
		i++
		a = a[i:]
		i = strings.Index(a, "\"")
		if i < 0 {
			continue
		}
		s := host + a[:i]
		res = append(res, s)
	}
	return
}

func getDetail(url string, info *app.SiteConf) {
	html := req.Get(url)
	if html == "" {
		tools.Log("\n" + url + ": 请求异常")
		return
	}
	html = tools.ClearTag(html)
	film := app.Film{}
	film.Site = info.Name
	film.Name = tools.GetMid(html, info.NameSep1, info.NameSep2)
	if film.Name == "" {
		tools.Log(url + "：数据异常，电影名为空")
		return
	}
	film.Name = strings.TrimSpace(film.Name)
	film.Director = strings.TrimSpace(tools.GetMid(html, info.DirectorSep1, info.DirectorSep2))
	if film.Director == "" {
		film.Director = "未知"
	}
	film.Actor = strings.TrimSpace(tools.GetMid(html, info.ActorSep1, info.ActorSep2))
	if film.Actor == "" {
		film.Actor = "未知"
	}
	film.Type = strings.TrimSpace(tools.GetMid(html, info.TypeSep1, info.TypeSep2))
	if film.Type == "" {
		film.Type = "未知"
	}
	film.Area = strings.TrimSpace(tools.GetMid(html, info.AreaSep1, info.AreaSep2))
	if film.Area == "" {
		film.Area = "未知"
	}
	film.Lang = strings.TrimSpace(tools.GetMid(html, info.LangSep1, info.LangSep2))
	if film.Lang == "" {
		film.Lang = "未知"
	}
	film.Year = strings.TrimSpace(tools.GetMid(html, info.YearSep1, info.YearSep2))
	if film.Year == "" {
		film.Year = "未知"
	}
	film.Summary = clearTag(tools.GetMid(html, info.SummarySep1, info.SummarySep2))
	if film.Summary == "" {
		film.Summary = "暂无"
	}
	film.Img = strings.TrimSpace(tools.GetMid(html, info.ImgSep1, info.ImgSep2))
	addrs := tools.GetAllMid(html, info.AddrSep1, info.AddrSep2)
	film.Addrs = getAddrs(addrs)

	// 比对保存
	save.Comp(film, info)
}

// 去除文本内的html标签
func clearTag(txt string) string {
	res := ""
	for {
		i := strings.Index(txt, "<")
		if i < 0 {
			res += txt
			break
		}
		if i > 0 {
			res += txt[:i]
		}
		txt = txt[i:]
		i = strings.Index(txt, ">")
		if i < 0 || i == len(txt)-1 {
			break
		}
		i++
		txt = txt[i:]
	}
	return res
}

// 提取地址
func getAddrs(ads []string) (arr []app.Addr) {
	a := []string{}
	for _, v := range ads {
		if strings.Index(v, ".m3u8") != -1 {
			a = append(a, v)
		}
	}
	for _, v := range a {
		if strings.Index(v, "http") == -1 {
			continue
		}
		ad := app.Addr{}
		va := strings.Split(v, "$")
		le := len(va)
		if le == 0 {
			continue
		}
		if le == 1 {
			ad.Remark = ""
			ad.URL = v
			arr = append(arr, ad)
			continue
		}
		ad.Remark = va[0]
		ad.URL = va[1]
		arr = append(arr, ad)
	}
	return
}
