package twinkle

import (
	"bytes"
	"errors"
	"fmt"
	"net"
)

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
	//控制位
	ctrl := data[start+offset+1]

	isIpv4 := BitGet(ctrl, 0)
	isRole := BitGet(ctrl, 1)
	if isRole {
		role = 1
	} else {
		role = 0
	}

	if !isIpv4 { //ipv4
		ip = data[start+offset+2 : start+offset+6]
		fmt.Printf("%v", ip)
		port = Btoi(data[start+offset+6 : start+offset+9])
		payload = data[start+offset+10 : l+5]
	} else { //ipv6
		ip = data[start+offset+2 : start+offset+18]
		port = Btoi(data[start+offset+18 : start+offset+21])
		payload = data[start+offset+22 : l+4]
	}

	return
}
