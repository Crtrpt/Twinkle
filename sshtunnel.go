package twinkle

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Crtrpt/twinkle/logger"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// 处理tcp 监听
func processingTunnelTcpClient(conn net.Conn, cfg ProxyConfig) {
	for {
		var buf = make([]byte, 1024)
		len, err := conn.Read(buf[:])
		if err != nil {
			return
		}
		fmt.Printf("收到数据====%d", len)
	}

}

// 处理tcp 监听
func processingTunnelTcpServer(listener net.Listener, cfg ProxyConfig) {
	var downClient net.Conn
	var udpForwardClient net.Conn
	for {
		conn, err := (listener).Accept()
		if err != nil {
			logger.Error(err)
			break
		}
		go TunnelTcpClient(conn, cfg, downClient, udpForwardClient)
	}
}

// 处理tcp 链接
func TunnelTcpClient(client net.Conn, cfg ProxyConfig, downclient net.Conn, udpForwardClient net.Conn) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Info("panic")
		}
	}()
	defer func() {
		logger.Info("exit")
	}()
	defer client.Close()

	do := make(chan struct{}, 0)
	//收到udp 代理来的数据 就发给隧道下行的客户端

	_ = <-do
}

func startForwardServer(sshClientConn *ssh.Client, cfg ProxyConfig) {
	//把foraward 可执行文件复制进去
	sftpClient, err := sftp.NewClient(sshClientConn)
	if err != nil {
		log.Fatal(err)
	}
	path := "/tmp/twinkle_udp_forward"
	_, err = sftpClient.Lstat(path)
	if err != nil {
		f, err := sftpClient.Open(path)
		if err != nil {
			logger.Errorf("打开文件异常%v", err)
		}
		f, err = sftpClient.Create(path)
		defer f.Close()

		srcfile, err := ioutil.ReadFile("./bin/twinkle_udp_forward")
		if err != nil {
			panic(err)
		}
		f.Write(srcfile)
		f.Close()
		sftpClient.Chmod(path, os.ModePerm)
	}
	sftpClient.Close()

	session, err := sshClientConn.NewSession()
	if err != nil {
		logger.Errorf("启动session异常11%v", err)
		return
	}

	RemoteUdpOverTcpUrlParse, err := url.Parse(cfg.Ssh.RemoteUdpOverTcp)
	urlParse, err := url.Parse(cfg.Ssh.Addr)

	cmd := path + " -tcp " + RemoteUdpOverTcpUrlParse.Host + " -udp " + urlParse.Host
	logger.Infof("执行命令%s", cmd)
	err = session.Run(cmd)

	if err != nil {
		logger.Errorf("启动session异常%v", err)
		return
	}
	//写入正常
}

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
		//启动udp over tcp
		if cfg.Ssh.RemoteUdpOverTcp == "" {
			logger.Errorf("udp透传需要配置 udp over tcp 代理")
			return
		}
		RemoteUdpOverTcpUrlParse, err := url.Parse(cfg.Ssh.RemoteUdpOverTcp)
		if err != nil {
			logger.Errorf("udp 代理隧道解析异常: %s", err)
			return
		}
		logger.Infof("启动 隧道服务器端 %s", RemoteUdpOverTcpUrlParse.Host)
		//启动上行监听
		listener, err := sshClientConn.Listen("tcp", RemoteUdpOverTcpUrlParse.Host)
		//暂时不支持的协议
		if err != nil {
			logger.Errorf("udp监听异常%v", err.Error())
			return
		}
		go processingTunnelTcpServer(listener, cfg)

		logger.Infof("启动 隧道客户端 %s", RemoteUdpOverTcpUrlParse.Host)
		client, err := net.Dial("tcp", RemoteUdpOverTcpUrlParse.Host)
		if err != nil {
			logger.Errorf("隧道客户端异常%v", err.Error())
			return
		}

		go processingTunnelTcpClient(client, cfg)

		// startForwardServer(sshClientConn, cfg)

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
