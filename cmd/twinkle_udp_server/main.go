package main

import (
	"fmt"
	"net"
	"os"

	"github.com/Crtrpt/twinkle/logger"
)

func processingUdp(conn *net.UDPConn, limit chan struct{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("err %v", err)
			return
		}
		limit <- struct{}{}
	}()
	for {
		buf := make([]byte, 1024)
		l, addr, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			logger.Errorf("err %v", err)
			return
		}
		logger.Infof("request  %s %d %s", addr.IP, addr.Port, string(buf[:l]))
		logger.Infof("response %s %d %s", addr.IP, addr.Port, []byte("ok"))
		_, err = conn.WriteToUDP([]byte("ok"), addr)
		if err != nil {
			logger.Infof("err ", err)
			break
		}
	}
}

func main() {
	addr := "127.0.0.1:9003"
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	conn, err := net.ListenUDP("udp", udpAddr)
	limit := make(chan struct{}, 10)
	defer conn.Close()
	if err != nil {
		fmt.Println("read from connect failed, err:" + err.Error())
		os.Exit(1)
	}
	limit <- struct{}{}
	limit <- struct{}{}
	logger.Infof("listen:%s", addr)
	for {
		<-limit
		go processingUdp(conn, limit)
	}

}
