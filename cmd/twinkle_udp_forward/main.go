package main

import (
	"flag"
	"net"

	"github.com/Crtrpt/twinkle"
	"github.com/Crtrpt/twinkle/logger"
)

func processingDown(udp *string, done chan struct{}, forward chan []byte) {
	defer func() {
		done <- struct{}{}
	}()
	udpAddr, err := net.ResolveUDPAddr("udp", *udp)
	if err != nil {
		logger.Errorf("udp 解析错误%s", err.Error())
		return
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Errorf("udp 监听错误%s", err.Error())
		return
	}
	if *v {
		logger.Infof("listen udp:%s", *udp)
	}

	go func() {
		for {
			select {
			case data, ok := <-downStream:
				if ok {
					_, err = conn.WriteToUDP(data, downAddr)
					if err != nil {
						logger.Errorf("写入downclient 出现异常%s", err.Error())
					}
				}
			}
		}
	}()
	for {
		buf := make([]byte, 1024)
		l, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			logger.Errorf("读取udp 数据包出现错误%s", err.Error())
		}

		if l == 4 && buf[0] == (*key)[0] && buf[1] == (*key)[1] && buf[2] == (*key)[2] && buf[3] == (*key)[3] {
			if *v {
				logger.Infof("update downclient:%s", addr)
			}
			downAddr = addr
			// _, err = conn.WriteToUDP([]byte("ok"), addr)
			continue
		}
		upstream <- buf[:l]
	}

}

// 处理
func processingUp(udp *string, done chan struct{}, forward chan []byte) {
	defer func() {
		done <- struct{}{}
	}()
	udpAddr, err := net.ResolveUDPAddr("udp", *udp)
	if err != nil {
		logger.Errorf("udp 解析错误%s", err.Error())
		return
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Errorf("udp 监听错误%s", err.Error())
		return
	}
	defer conn.Close()
	if *v {
		logger.Infof("listen udp:%s", *udp)
	}

	go func() {
		for {
			buf := make([]byte, 1024)
			l, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				logger.Errorf("读取udp 数据包出现错误%s", err.Error())
				break
			}

			downStream <- twinkle.UDPForwardPacket(0, addr.IP, addr.Port, buf[:l])
		}
	}()

	for {
		select {
		case data, ok := <-upstream:
			if ok {
				_, ip, port, payload, err := twinkle.UDPForwardUnPacket(data)
				if err != nil {
					logger.Errorf("udp 解包出现错误%s", err.Error())
					break
				}
				logger.Infof("-> %s:%d", ip, port)
				_, err = conn.WriteTo(payload, &net.UDPAddr{IP: ip, Port: port})
				if err != nil {
					logger.Errorf("写入upclient 出现异常%s", err.Error())
				}
			}
		}
	}

}

var key *string
var v, vv, vvv *bool

var downAddr *net.UDPAddr

var downStream chan []byte
var upstream chan []byte

func main() {
	udp := flag.String("udp", "", "对外暴露的udp地址")
	udp_down := flag.String("udp_down", "", "对外暴露的udp地址")
	key = flag.String("key", "twin", "对外暴露的udp地址")

	v = flag.Bool("v", false, "v")
	vv = flag.Bool("vv", false, "v")
	vvv = flag.Bool("vvv", false, "vvv")

	flag.Parse()
	done := make(chan struct{}, 0)
	forword := make(chan []byte, 1000)

	downStream = make(chan []byte, 0)
	upstream = make(chan []byte, 0)

	go processingDown(udp_down, done, forword)
	go processingUp(udp, done, forword)
	_ = <-done
	println("end===")

}
