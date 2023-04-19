package gps

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"golang.org/x/crypto/ssh"
)

func (app *App) ListenSSHTunnel(cfg ProxyConfig) {
	sshCfg := cfg.Ssh
	key, _ := ioutil.ReadFile(sshCfg.PrivateKey)

	signer, _ := ssh.ParsePrivateKey(key)

	config := &ssh.ClientConfig{
		// 服务器用户名
		User: cfg.Ssh.UserName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	sshClientConn, err := ssh.Dial("tcp", sshCfg.Host, config)
	if err != nil {
		fmt.Printf("ssh.Dial failed: %s", err)
		return
	}

	listener, err := sshClientConn.Listen("tcp", cfg.Ssh.Addr)

	if err != nil {
		fmt.Printf("tcp监听异常%v", err.Error())
		return
	}
	server := &http.Server{Addr: cfg.Ssh.Addr, Handler: app}
	server.Serve(listener)
}
