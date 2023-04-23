package twinkle

import (
	"bytes"
	"fmt"
)

func BitGetRange(v byte, start, end int) int {
	panic("暂时没有实现")
}

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
	if len(data) != 3 {
		fmt.Printf("%v===============", len(data))
		panic("error btoi ")
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

// 动态int 写入协议
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

// 反解析 动态byte 到 int
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
