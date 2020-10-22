package tools

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

var mut sync.Mutex

// F ..
var F *os.File

// OpenF ..
func OpenF() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	p := filepath.Join(wd, "log.txt")
	os.Remove(p)
	f, err := os.OpenFile(p, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	F = f
}

// CloseF ..
func CloseF() {
	F.Close()
	F = nil
}

// Log ..
func Log(msg string) {
	if F == nil {
		return
	}
	msg = "<br>\n" + msg + "<br>\n" + time.Now().Format("15:04:05") + "<br>\n---------------------------------------<br>\n"
	mut.Lock()
	F.Write([]byte(msg))
	mut.Unlock()
}
