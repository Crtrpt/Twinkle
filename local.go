package twinkle

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/Crtrpt/twinkle/logger"
)

func processTcpConn(conn net.Conn, cfg ProxyConfig) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Info("panic")
		}
	}()
	defer func() {
		logger.Info("exit")
	}()
	defer conn.Close()

	urlBackend, err := url.Parse(cfg.Proxy)

	backendConn, err := net.Dial("tcp", urlBackend.Host)
	if err != nil {
		logger.Infof("tcp backend error", err)
		return
	}
	defer backendConn.Close()
	logger.Infof("dial:%s", cfg.Proxy)
	do := make(chan struct{}, 0)
	sourceExit := make(chan struct{}, 1)
	targetExit := make(chan struct{}, 1)

	go func() {
		defer func() {
			do <- struct{}{}
			logger.Infof("source->target exit")
		}()
		for {
			select {
			case _ = <-sourceExit:
				return
			default:
				conn.SetDeadline(time.Now().Add(10 * time.Second))
				backendConn.SetDeadline(time.Now().Add(10 * time.Second))
				//source->target
				var buf = make([]byte, 1024)
				len, err := conn.Read(buf[:])
				if err != nil {
					return
				}
				if len == 0 {
					continue
				}
				len, err = backendConn.Write(buf[0:len])
				if err != nil {
					return
				}
			}
		}
	}()
	go func() {
		defer func() {
			do <- struct{}{}
			logger.Infof("target->source exit")
		}()
		for {
			select {
			case _ = <-targetExit:
				return
			default:
				conn.SetDeadline(time.Now().Add(10 * time.Second))
				backendConn.SetDeadline(time.Now().Add(10 * time.Second))
				//target->source
				var buf = make([]byte, 1024)
				len, err := backendConn.Read(buf[:])
				if err != nil {
					return
				}
				if len == 0 {
					continue
				}
				len, err = conn.Write(buf[0:len])
				if err != nil {
					return
				}
			}
		}
	}()
	_ = <-do
	sourceExit <- struct{}{}
	targetExit <- struct{}{}
}

func processingTcp(listener net.Listener, cfg ProxyConfig) {
	for {
		conn, err := (listener).Accept()
		if err != nil {
			fmt.Printf("出现错误====")
			break
		}

		go processTcpConn(conn, cfg)
	}
}

// TODO 实现udp协议的透传
func processingUdp(listener net.Listener, cfg ProxyConfig) {
	for {
		conn, err := (listener).Accept()
		if err != nil {
			fmt.Printf("出现错误====")
			break
		}

		go processTcpConn(conn, cfg)
	}
}

// 监听本地端口
func (app *App) ListenHttpPort(scheme string, cfg ProxyConfig, urlp *url.URL) {
	addr := urlp.Host

	logger.Infof("localhost:%s", cfg.Url)
	if scheme == "udp" {
		listener, err := net.Listen("udp", addr)
		if err != nil {
			logger.Errorf("监听服务器异常 %s %v", addr, err)
		}
		go processingUdp(listener, cfg)
		return
	}
	if scheme == "tcp" {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			logger.Errorf("监听服务器异常 %s %v", addr, err)
		}
		go processingTcp(listener, cfg)
		return
	}

	err := http.ListenAndServe(urlp.Host, app)
	if err != nil {
		logger.Errorf("监听服务器异常 %s %v", addr, err)
	}
}
