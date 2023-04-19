package gps

import (
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/Crtrpt/gps/logger"
	"golang.org/x/crypto/ssh"
)

func (app *App) ListenSSHTunnel(cfg ProxyConfig) {
	sshCfg := cfg.Ssh
	key, err := ioutil.ReadFile(sshCfg.PrivateKey)
	if err != nil {
		panic(err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}

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
		logger.Errorf("ssh.Dial failed: %s", err)
		return
	}

	listener, err := sshClientConn.Listen("tcp", cfg.Ssh.Addr)

	if err != nil {
		logger.Errorf("tcp监听异常%v", err.Error())
		return
	}
	server := &http.Server{Addr: cfg.Ssh.Addr, Handler: app}
	server.Serve(listener)
}
