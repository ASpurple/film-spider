package app

import (
	"fmt"
	"io/ioutil"
	"main/tools"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// StartPage 开始页
var StartPage = 0

// StopPage 结束页
var StopPage = 0

// Sites 各个站点配置信息
var Sites []*SiteConf

// Connstr 数据库连接信息
var Connstr = ""

// ReadConf 读取配置文件
func ReadConf() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p := filepath.Join(wd, "site.conf")
	file, err := os.Open(p)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conf := string(buf)
	Connstr = tools.GetMid(conf, "<db>", "</db>")
	pg := tools.GetMid(conf, "<range>", "</range>")
	pgarr := strings.Split(pg, ",")
	if len(pgarr) < 2 {
		fmt.Println("\n抓取范围配置错误: 缺少配置项")
		os.Exit(1)
	}
	pgs, e1 := strconv.Atoi(pgarr[0])
	pge, e2 := strconv.Atoi(pgarr[1])
	if e1 != nil || e2 != nil {
		fmt.Println("\n抓取范围配置错误：应为正确的整数")
		os.Exit(1)
	}
	StartPage = pgs
	StopPage = pge
	ss := tools.GetAllMid(conf, "{", "}")
	for i, c := range ss {
		arr := tools.GetAllMid(c, "[[", "]]")
		errStr := "第" + strconv.Itoa(i+1) + "条配置信息错误"
		if len(arr) < 15 {
			fmt.Println(errStr + ": 缺少配置项")
			continue
		}
		if len(arr[0]) < 1 {
			fmt.Println(errStr + ": 名称不能为空")
			continue
		}
		ok := true
		cs := [][]string{}
		for k := 5; k < len(arr); k++ {
			seps := strings.Split(arr[k], "|")
			if len(seps) < 2 {
				ok = false
				break
			}
			cs = append(cs, seps)
		}
		if !ok {
			fmt.Println(errStr + ": 缺少标记")
			continue
		}
		site := SiteConf{
			Name:         arr[0],
			Host:         arr[1],
			ListURL:      arr[2],
			DetailURL:    arr[3],
			HrefSep:      arr[4],
			NameSep1:     cs[0][0],
			NameSep2:     cs[0][1],
			DirectorSep1: cs[1][0],
			DirectorSep2: cs[1][1],
			ActorSep1:    cs[2][0],
			ActorSep2:    cs[2][1],
			TypeSep1:     cs[3][0],
			TypeSep2:     cs[3][1],
			AreaSep1:     cs[4][0],
			AreaSep2:     cs[4][1],
			LangSep1:     cs[5][0],
			LangSep2:     cs[5][1],
			YearSep1:     cs[6][0],
			YearSep2:     cs[6][1],
			SummarySep1:  cs[7][0],
			SummarySep2:  cs[7][1],
			ImgSep1:      cs[8][0],
			ImgSep2:      cs[8][1],
			AddrSep1:     cs[9][0],
			AddrSep2:     cs[9][1],
		}
		Sites = append(Sites, &site)
	}
}
