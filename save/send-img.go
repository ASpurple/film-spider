package save

import (
	"main/tools"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

// Cli ..
type Cli struct {
	user       string
	pwd        string
	addr       string
	client     *ssh.Client
	session    *ssh.Session
	LastResult string
}

// Connect ..
func (c *Cli) connect() error {
	config := &ssh.ClientConfig{}
	config.SetDefaults()
	config.User = c.user
	config.Auth = []ssh.AuthMethod{ssh.Password(c.pwd)}
	config.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }
	client, err := ssh.Dial("tcp", c.addr, config)
	if nil != err {
		return err
	}
	c.client = client
	return nil
}

// Run ..
func (c *Cli) Run(shell string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	c.session = session
	buf, err := c.session.CombinedOutput(shell)
	c.LastResult = string(buf)
	c.session.Close()
	return c.LastResult, err
}

// Close ..
func (c *Cli) Close() {
	c.session.Close()
	c.client.Close()
}

// SendFile ..
func (c *Cli) SendFile(path string) {
	arr := strings.Split(path, "/")
	if len(arr) < 2 {
		tools.Log("文件路径异常：" + path)
		return
	}
	i := len(arr) - 1
	name := arr[i]
	rarr := arr[:i]
	remoteDir := "."
	pp := strings.Join(rarr, "/")
	if len(rarr) > 0 {
		remoteDir += "/" + pp
	}
	remotePath := "/usr/web/file/" + pp + "/" + name
	wd, err := os.Getwd()
	if err != nil {
		tools.Log("发送文件时未获取到当前工作目录：" + err.Error())
		return
	}
	ps := []string{wd}
	ps = append(ps, arr...)
	localPath := filepath.Join(ps...)
	_, err = c.Run("cd /usr/web/file && mkdir -vp " + remoteDir)
	if err != nil {
		tools.Log("远程服务器创建文件夹出错：" + err.Error())
		return
	}
	ss, err := c.client.NewSession()
	if err != nil {
		tools.Log("上传文件时未能创建Session：" + err.Error())
		return
	}
	err = scp.CopyPath(localPath, remotePath, ss)
	if err != nil {
		tools.Log("上传文件时出错：" + err.Error())
	}
}

// NewSSH ..
func NewSSH() *Cli {
	conn := Cli{
		user: "root",
		pwd:  "Qq809408345",
		addr: "103.210.236.64:22",
	}
	err := conn.connect()
	if err != nil {
		panic(err)
	}
	return &conn
}
