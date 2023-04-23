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
	buf := make([]byte, 1024)
	_, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		logger.Errorf("err %v", err)
	}
	logger.Infof("response", addr)
	conn.WriteToUDP([]byte("ok"), addr)
}

func main() {
	addr := "127.0.0.1:9001"
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
	limit <- struct{}{}
	logger.Infof("listen:%s", addr)
	for {
		<-limit
		go processingUdp(conn, limit)
	}

}
