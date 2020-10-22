package main

import (
	"main/app"
	"main/mail"
	"main/save"
	"main/tools"
	"main/work"
	"time"
)

func main() {
	app.ReadConf()
	save.SetDatabase()
	task()
}

func task() {
	work.StartAll()
	ticker := time.NewTicker(time.Minute * 15)
	for range ticker.C {
		h, _, _ := time.Now().Clock()
		if h == 5 {
			tools.OpenF()
			work.StartAll()
			tools.CloseF()
			mail.SendLog()
			time.Sleep(time.Hour * 1)
		}
	}
}
