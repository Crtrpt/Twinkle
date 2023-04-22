package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/Crtrpt/twinkle/logger"
)

// bytes to int
func btoi(data []byte) int {
	y := 0
	l := len(data)
	for i := l - 1; i >= 0; i-- {
		y <<= 8
		lsb := data[i]
		y = y ^ int(lsb)
	}
	return y
}

// int to bytes
func itob(x int) []byte {
	buf := bytes.NewBuffer([]byte{})
	lsb := x & 0xff
	x >>= 8
	buf.WriteByte(byte(lsb))

	lsb = x & 0xff
	x >>= 8
	buf.WriteByte(byte(lsb))

	lsb = x & 0xff
	x >>= 8
	buf.WriteByte(byte(lsb))

	lsb = x & 0xff
	x >>= 8
	buf.WriteByte(byte(lsb))

	return buf.Bytes()
}

func packVarInt(x int) []byte {
	buf := bytes.NewBuffer([]byte{})
	for {
		lsb := x & 0x7f
		x >>= 7
		if x == 0 {
			buf.WriteByte(byte(lsb))
			break
		} else {
			buf.WriteByte(byte(0x80 | lsb))
		}
	}
	return buf.Bytes()
}

func unpackVarInt(data []byte) int {
	y := 0
	l := len(data)
	for i := l - 1; i >= 0; i-- {
		y <<= 7
		lsb := data[i]
		if data[i] > 127 {
			lsb = data[i] ^ 0x80
		}
		y = y ^ int(lsb)
	}
	return y
}

// 组包 ${twin}${payload_len}[${ip version}${ip}${port}${payload}]
func packet(ip net.IP, port int, payload_len int, payload []byte) []byte {
	buf := bytes.NewBuffer([]byte("twin"))
	data := bytes.NewBuffer([]byte(""))
	if len(ip) == 4 {
		data.WriteByte(byte(0))
	} else {
		data.WriteByte(byte(1))
	}
	data.Write(ip)
	data.Write(itob(port))
	data.Write(payload)

	buf.Write(packVarInt(data.Len()))
	buf.Write(data.Bytes())
	return buf.Bytes()
}

// 解包
func unpack(data []byte) (ip net.IP, port int, payload []byte, err error) {
	if data[0] != 't' || data[1] != 'w' || data[2] != 'i' || data[3] != 'n' {
		err = errors.New("数据包解析错误")
		return
	}
	start := 4
	offset := 0
	l := 0
	//TODO 处理数据包过大警告
	for {
		if data[start+offset] < 128 {
			l = unpackVarInt(data[start : start+offset+1])
			break
		}
		offset = offset + 1
	}

	if data[start+offset+1] == 0 { //ipv4
		ip = data[start+offset+2 : start+offset+6]
		port = btoi(data[start+offset+6 : start+offset+9])
		payload = data[start+offset+10 : l+4]
	} else { //ipv6
		ip = data[start+offset+2 : start+offset+18]
		port = btoi(data[start+offset+18 : start+offset+21])
		payload = data[start+offset+22 : l+4]
	}

	return
}

// 处理
func processingUdpServer(udp *string, done chan struct{}) {
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
		processingUdp(conn)
	}
}

// 处理udp数据包
func processingUdp(conn *net.UDPConn) {
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
	fmt.Printf("原始数据包: %s,%d,%s \r\n", addr.IP, addr.Port, buf[:l])
	data := packet(addr.IP, addr.Port, l, buf[:l])
	fmt.Printf("转发数据组包: %0x \r\n", data)
	ip, port, p, err := unpack(data)
	fmt.Printf("解包数据包: %s,%d,%s \r\n", ip.String(), port, p)
}

func processingTcpClient(tcp *string, done chan struct{}) {
	defer func() {
		done <- struct{}{}
	}()

}

func processingTcp() {

}

func main() {
	udp := flag.String("udp", "", "要监听的udp服务器")
	flag.Parse()
	// tcp := flag.String("tcp", "", "要转发的tcp服务器")
	done := make(chan struct{}, 0)
	go processingUdpServer(udp, done)
	// go processingTcpClient(tcp, done)
	_ = <-done
	println("end===")

}
