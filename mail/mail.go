package mail

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

// SendMail ..
func SendMail(mailTo []string, subject string, body string) error {
	//定义邮箱服务器连接信息，如果是阿里邮箱 pass填密码，qq邮箱填授权码
	mailConn := map[string]string{
		"user": "809408345@qq.com",
		"pass": "ayyntygbjhdkbege",
		"host": "smtp.qq.com",
		"port": "465",
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()
	m.SetHeader("From", "Film Server Spider"+"<"+mailConn["user"]+">")
	m.SetHeader("To", mailTo...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err

}

// SendLog 发送日志
func SendLog() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	p := filepath.Join(wd, "log.txt")
	file, err := os.Open(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf, err := ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	body := string(buf)
	t := time.Now().Format("2006-01-02 15:04:05")
	err = SendMail([]string{"15255689160@163.com"}, "Spider运行日志-"+t, body)
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Remove(p)
}
