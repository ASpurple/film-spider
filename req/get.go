package req

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// Get get请求
func Get(addr string) (html string) {
	resp, err := http.Get(addr)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return string(buf)
}

// Download ..
func Download(u string) (buf []byte, suf string) {
	resp, err := http.Get(u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	buf = bs
	t := resp.Header.Get("Content-Type")
	i := strings.Index(t, "/")
	if i < 0 || i == len(t)-1 {
		suf = "jpg"
	} else {
		i++
		suf = t[i:]
	}
	return
}
