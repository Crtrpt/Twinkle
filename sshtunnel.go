package twinkle

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/Crtrpt/twinkle/logger"
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
		User:    cfg.Ssh.UserName,
		Auth:    []ssh.AuthMethod{},
		Timeout: time.Duration(cfg.Ssh.Timeout) * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	if cfg.Ssh.Auth == "key" {
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	if cfg.Ssh.Auth == "password" {
		config.Auth = append(config.Auth, ssh.Password(cfg.Ssh.Password))
	}

	sshClientConn, err := ssh.Dial("tcp", sshCfg.Host, config)
	if err != nil {
		logger.Errorf("ssh.Dial failed: %s", err)
		return
	}

	logger.Infof("remote:%s %s", cfg.Ssh.Host, cfg.Ssh.Addr)
	urlParse, err := url.Parse(cfg.Ssh.Addr)
	if err != nil {
		logger.Errorf("url parse  failed: %s", err)
		return
	}

	if urlParse.Scheme == "udp" {
		listener, err := sshClientConn.Listen("udp", urlParse.Host)
		//暂时不支持的协议
		if err != nil {
			logger.Errorf("udp监听异常%v", err.Error())
			return
		}
		processingUdp(listener, cfg)
		return
	}

	if urlParse.Scheme == "tcp" {

		listener, err := sshClientConn.Listen("tcp", urlParse.Host)

		if err != nil {
			logger.Errorf("tcp监听异常%v", err.Error())
			return
		}
		processingTcp(listener, cfg)
		return
	}

	if urlParse.Scheme == "http" {
		listener, err := sshClientConn.Listen("tcp", urlParse.Host)

		if err != nil {
			logger.Errorf("tcp监听异常%v", err.Error())
			return
		}
		//增加应用层协议判断
		server := &http.Server{Addr: cfg.Ssh.Addr, Handler: app}
		server.Serve(listener)
		return
	}

}
