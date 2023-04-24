package twinkle

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/Crtrpt/twinkle/logger"
)

// 处理tcp 链接
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

func processUdpConn(clientConn *net.UDPConn, cfg ProxyConfig, limit chan struct{}) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Info("panic")
		}
		limit <- struct{}{}
	}()

	defer clientConn.Close()
	urlBackend, err := url.Parse(cfg.Proxy)

	backendUdpConn, err := net.Dial("udp", urlBackend.Host)
	defer clientConn.Close()
	if err != nil {
		logger.Infof("udp backend error", err)
		return
	}

	do := make(chan struct{}, 0)
	sourceExit := make(chan struct{}, 1)
	targetExit := make(chan struct{}, 1)

	var buf = make([]byte, 1024)

	len, clientAddr, err := clientConn.ReadFromUDP(buf[:])

	if err != nil {
		logger.Errorf("err client read %s addr %s", err, clientAddr)
		return
	}

	len, err = backendUdpConn.Write(buf[0:len])
	if err != nil {
		logger.Errorf("err backend write %s", err)
		return
	}
	if err == nil {
		logger.Infof("client:%s->local:%s->backend:%s", clientAddr, backendUdpConn.LocalAddr(), backendUdpConn.RemoteAddr())
	}

	go func(clientAddr *net.UDPAddr) {
		defer func() {
			do <- struct{}{}
			logger.Infof("local->backend exit")
		}()
		for {
			select {
			case _ = <-sourceExit:
				return
			default:

				//local->backend
				var buf = make([]byte, 1024)
				//获取客户端IP地址和端口号
				len, addr, err := clientConn.ReadFromUDP(buf[:])
				if err != nil {
					logger.Errorf("err client read %s addr %s", err, addr)
					return
				}
				if len == 0 {
					continue
				}

				len, err = backendUdpConn.Write(buf[0:len])
				if err != nil {
					logger.Errorf("err backend write %s", err)
					return
				}
				logger.Infof("client:%s->local:%s->backend:%s", clientAddr, backendUdpConn.LocalAddr(), backendUdpConn.RemoteAddr())
			}
		}
	}(clientAddr)
	go func(clientAddr *net.UDPAddr) {
		defer func() {
			do <- struct{}{}
			logger.Infof("backend->local exit")
		}()
		for {
			select {
			case _ = <-targetExit:
				return
			default:
				if clientAddr == nil {

					continue
				}

				var buf = make([]byte, 1024)
				len, err := backendUdpConn.Read(buf[:])
				if err != nil {
					logger.Errorf("err backend read %s", err)
					return
				}
				// logger.Errorf("read data  addr %s", backendUdpConn.LocalAddr())
				if len == 0 {
					continue
				}
				_, err = clientConn.WriteToUDP(buf, clientAddr)
				if err != nil {
					logger.Errorf("err client write %s", err)
					return
				}
				logger.Infof("backend:%s->local:%s->client:%s", backendUdpConn.RemoteAddr(), backendUdpConn.LocalAddr(), clientAddr)
			}
		}
	}(clientAddr)
	<-do
	sourceExit <- struct{}{}
	targetExit <- struct{}{}
}

func processingUdp(conn *net.UDPConn, cfg ProxyConfig) {

	//增加并发处理的能力
	limit := make(chan struct{}, 3)
	limit <- struct{}{}
	for {
		<-limit
		processUdpConn(conn, cfg, limit)
	}
}

// 处理tcp 监听
func processingTcp(listener net.Listener, cfg ProxyConfig) {
	for {
		conn, err := (listener).Accept()
		if err != nil {
			logger.Error(err)
			break
		}

		go processTcpConn(conn, cfg)
	}
}

// 监听本地端口
func (app *App) ListenLocalPort(scheme string, cfg ProxyConfig, urlp *url.URL) {
	addr := urlp.Host

	logger.Infof("localhost:%-40s proxy:%-40s", cfg.Url, cfg.Proxy)
	if scheme == "udp" {
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			logger.Errorf("解析udp协议出现错误 %s %v", addr, err)
			return
		}
		udpConn, err := net.ListenUDP("udp", udpAddr)
		if err != nil {
			logger.Errorf("监听服务异常 %s %v", addr, err)
			return
		}
		fmt.Printf("启动udp服务器%s\r\n", addr)
		go processingUdp(udpConn, cfg)
		return
	}
	if scheme == "tcp" {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			logger.Errorf("监听服务器异常 %s %v", addr, err)
			return
		}
		go processingTcp(listener, cfg)
		return
	}

	err := http.ListenAndServe(urlp.Host, app)
	if err != nil {
		logger.Errorf("监听服务器异常 %s %v", addr, err)
	}
}
