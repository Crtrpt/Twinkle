package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Crtrpt/twinkle"
	"github.com/Crtrpt/twinkle/logger"
)

func main() {
	uptobackend := make(map[string]any, 0)
	for {
		udpClient, _ := net.Dial("udp", "127.0.0.1:9002")
		udpClient.Write([]byte("twin"))
		logger.Infof("本地地址:%s 远程地址:%s", udpClient.LocalAddr(), udpClient.RemoteAddr())
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			defer func() {
				//程序挂掉了？
				uptobackend = make(map[string]any, 0)
			}()
			for {
				buf := make([]byte, 1024)
				l, err := udpClient.Read(buf)
				if err != nil {
					logger.Errorf("出现异常 %v", err)
					// time.Sleep(time.Second * 1)
					udpClient.Close()
					break
				}
				_, ip, port, payload, err := twinkle.UDPForwardUnPacket(buf[:l])
				if err != nil {
					logger.Errorf("数据包解错误 %v %s", err, hex.Dump(buf[:l]))

					return
				}
				fmt.Printf("收到数据:addr:%s port:%d  payload:%s\n", ip, port, string(payload))

				//发起udp请求
				var backendClient *net.Conn
				key := fmt.Sprintf("%s:%d", ip, port)
				if uptobackend[key] == nil {
					conn, err := net.Dial("udp", "127.0.0.1:9003")
					if err != nil {
						logger.Errorf("创建链接出现异常  %v", err)
						break
					}
					uptobackend[key] = &conn
					go func(ip net.IP, port int) {
						defer func() {
							//从后端读取出错删除映射关系
							delete(uptobackend, fmt.Sprintf("%s:%d", ip, port))
						}()

						for {
							buf := make([]byte, 1024)
							l, err := (*backendClient).Read(buf)
							if err != nil {
								logger.Errorf("出现异常 %v", err)
								time.Sleep(time.Second * 1)
								udpClient.Close()
								break
							}
							l, err = udpClient.Write(twinkle.UDPForwardPacket(0, ip, port, buf[:l]))
							if err != nil {
								logger.Errorf("写入forward异常")
								break
							}
							logger.Infof("forward->  %s:%d", ip, port)
						}
					}(ip, port)
				}
				backendClient = uptobackend[key].(*net.Conn)

				//写入后端出错页删除掉映射关系
				_, err = (*backendClient).Write(payload)
				if err != nil {
					delete(uptobackend, fmt.Sprintf("%s:%d", ip, port))
				}

			}
		}(wg)
		wg.Wait()
	}

}
