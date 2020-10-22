package tools

import (
	"strings"
	"unicode"
)

// GetMid 获取字符串指定片段，s或e传空代表从开头或结尾选取，只会取第一个
func GetMid(str string, s string, e string) string {
	r, _ := pickup(str, s, e)
	return r
}

// 提取两个子串中间部分，可传空，可传通配符
func pickup(txt string, s string, e string) (result, after string) {
	so := txt
	if s != "" && s != "*" {
		sa := strings.Split(s, "*")
		for _, str := range sa {
			i := strings.Index(txt, str)
			if i < 0 {
				return "", so
			}
			i += len(str)
			if i == len(txt) {
				return "", so
			}
			txt = txt[i:]
		}
	}
	if e == "" || e == "*" {
		result = txt
		after = txt
		return
	}
	ea := strings.Split(e, "*")
	if len(ea) < 1 {
		return "", so
	}
	k := strings.Index(txt, ea[0])
	if k < 0 {
		return "", so
	}
	result = txt[:k]
	for index, str := range ea {
		i := strings.Index(txt, str)
		if i < 0 {
			result = ""
			after = so
			break
		}
		i += len(str)
		if i == len(txt) {
			if index == len(ea)-1 {
				after = ""
				break
			} else {
				result = ""
				after = so
				break
			}
		}
		after = txt[i:]
	}
	return
}

// GetAllMid 根据开头子串和结尾子串获取字符串的所有子串
func GetAllMid(str string, s string, e string) []string {
	var list []string
	for {
		r, af := pickup(str, s, e)
		if r == "" {
			break
		}
		list = append(list, r)
		str = af
	}
	return list
}

func clearSep(txt string) string {
	rs := []rune(txt)
	result := []rune{}
	for _, r := range rs {
		if r == 10 || r == 13 || r == 32 || r == 9 {
			continue
		}
		result = append(result, r)
	}
	return string(result)
}

// 判断是否有中文或英文
func hasTxt(str string) bool {
	b := false
	rs := []rune(str)
	for _, r := range rs {
		if unicode.IsLetter(r) {
			b = true
			break
		}
	}
	return b
}

// ClearTag 去除所有标签内部和标签之间的空格和换行
func ClearTag(html string) (h string) {
	for {
		i := strings.Index(html, "<")
		if i < 0 {
			h += html
			break
		}
		if i > 0 {
			html = html[i:]
		}
		i = strings.Index(html, ">")
		if i < 0 || i == len(html)-1 {
			h += html
			break
		}
		i++
		h += clearSep(html[:i])
		html = html[i:]
		i = strings.Index(html, "<")
		if i < 0 {
			h += html
			break
		}
		txt := html[:i]
		if hasTxt(txt) {
			h += txt
		} else {
			h += clearSep(txt)
		}
		html = html[i:]
	}
	return
}
