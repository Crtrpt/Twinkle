package twinkle

import (
	"bytes"
	"errors"
	"fmt"
	"net"
)

// 获取bit位
func BitGet(v byte, offset int) bool {
	return v&(1<<offset) != 0
}

// 设置bit位
func BitSet(v byte, offset int) byte {
	if offset == 0 {
		v = v | 0b0000_0001
	}
	if offset == 1 {
		v = v | 0b0000_0010
	}
	if offset == 2 {
		v = v | 0b0000_0100
	}
	if offset == 3 {
		v = v | 0b0000_1000
	}
	if offset == 4 {
		v = v | 0b0001_0000
	}
	if offset == 5 {
		v = v | 0b0010_0000
	}
	if offset == 6 {
		v = v | 0b0100_0000
	}
	if offset == 7 {
		v = v | 0b1000_0000
	}
	return v
}

// 清除bit位
func BitClear(v byte, offset int) byte {
	if offset == 0 {
		v = (^((^v) | 0b0000_0001))
	}
	if offset == 1 {
		v = (^(^v) | 0b0000_0010)
	}
	if offset == 2 {
		v = (^(^v) | 0b0000_0100)
	}
	if offset == 3 {
		v = (^(^v) | 0b0000_1000)
	}
	if offset == 4 {
		v = (^(^v) | 0b0001_0000)
	}
	if offset == 5 {
		v = (^(^v) | 0b0010_0000)
	}
	if offset == 6 {
		v = (^(^v) | 0b0100_0000)
	}
	if offset == 7 {
		v = (^(^v) | 0b1000_0000)
	}
	return v
}

// bytes to int
func Btoi(data []byte) int {
	if len(data) != 4 {
		panic("error btoi")
	}
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
func Itob(x int) []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteByte(byte(x & 0xff))
	x >>= 8
	buf.WriteByte(byte(x & 0xff))
	x >>= 8
	buf.WriteByte(byte(x & 0xff))
	x >>= 8
	buf.WriteByte(byte(x & 0xff))
	x >>= 8
	return buf.Bytes()
}

func PackVarInt(x int) []byte {
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

func UnPackVarInt(data []byte) int {
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
func UDPForwardPacket(role int, ip net.IP, port int, payload []byte) []byte {
	buf := bytes.NewBuffer([]byte("twin"))
	data := bytes.NewBuffer([]byte(""))
	ctrl := byte(0)
	//ip类型
	if len(ip) == 4 {
		ctrl <<= 1
	} else {
		ctrl := ctrl | 0b0000_0001
		ctrl <<= 1
	}
	//udp 客户端发送方角色
	if role == 0 {
		ctrl = ctrl | 0b0000_0000
	}
	//down
	if role == 1 {
		ctrl = ctrl | 0b0000_0010
	}
	fmt.Printf("ctrl: %b", ctrl)
	data.WriteByte(ctrl)
	data.Write(ip)
	data.Write(Itob(port))
	data.Write(payload)

	buf.Write(PackVarInt(data.Len()))
	buf.Write(data.Bytes())
	return buf.Bytes()
}

// 解包
func UDPForwardUnPacket(data []byte) (role int, ip net.IP, port int, payload []byte, err error) {
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
			l = UnPackVarInt(data[start : start+offset+1])
			break
		}
		offset = offset + 1
	}
	// ctrl := data[start+offset+1]
	isIpv4 := true

	if isIpv4 { //ipv4
		ip = data[start+offset+2 : start+offset+6]
		port = Btoi(data[start+offset+6 : start+offset+9])
		payload = data[start+offset+10 : l+4]
	} else { //ipv6
		ip = data[start+offset+2 : start+offset+18]
		port = Btoi(data[start+offset+18 : start+offset+21])
		payload = data[start+offset+22 : l+4]
	}

	return
}
