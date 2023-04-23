package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/Crtrpt/twinkle"
	"github.com/Crtrpt/twinkle/logger"
)

// 处理
func processingUdpServer(udp *string, done chan struct{}, forward chan []byte) {
	defer func() {
		done <- struct{}{}
	}()
	udpAddr, err := net.ResolveUDPAddr("udp", *udp)
	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	if err != nil {
		fmt.Println("read from connect failed, err:" + err.Error())
		os.Exit(1)
	}
	logger.Infof("listen:%s", *udp)
	for {
		processingUdp(conn, forward)
	}
}

// 处理udp数据包
func processingUdp(conn *net.UDPConn, forward chan []byte) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	buf := make([]byte, 1024)
	l, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		logger.Errorf("读取udp 数据包出现错误%s", err.Error())
		return
	}
	forward <- twinkle.UDPForwardPacket(0, addr.IP, addr.Port, buf[:l])
}

func processingTcpClient(tcp *string, done chan struct{}, forward chan []byte) {
	defer func() {
		done <- struct{}{}
	}()

	client, err := net.Dial("tcp", *tcp)
	if err != nil {
		os.Exit(1)
	}
	defer client.Close()

	go func(client net.Conn) {
		//读取tcp客户端但会的数据根据pack
		// TODO 发送给指定的udp服务器

	}(client)
	for {
		select {
		case p, ok := <-forward:
			//转发数据包
			if ok {
				// fmt.Printf("转发数据包")
				client.Write(p)
			}
		}
	}

}

func processingTcp() {

}

func main() {
	udp := flag.String("udp", "", "要监听的udp服务器")
	tcp := flag.String("tcp", "", "要转发的tcp服务器")
	flag.Parse()
	done := make(chan struct{}, 0)
	forword := make(chan []byte, 1000)
	go processingTcpClient(tcp, done, forword)
	go processingUdpServer(udp, done, forword)
	_ = <-done
	println("end===")

}
