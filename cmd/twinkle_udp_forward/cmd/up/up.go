package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Crtrpt/twinkle/logger"
)

func main() {
	for {
		udpClient, _ := net.Dial("udp", "127.0.0.1:9001")

		logger.Infof("本地地址:%s 远程地址:%s", udpClient.LocalAddr(), udpClient.RemoteAddr())
		wg := &sync.WaitGroup{}
		wg.Add(2)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				time.Sleep(time.Second * 60)
				data := []byte("from up,time:" + time.Now().Format(time.Layout))
				fmt.Printf("写入数据: %v\n", string(data))
				_, err := udpClient.Write(data)
				if err != nil {
					break
				}
			}
		}(wg)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				buf := make([]byte, 1024)
				l, err := udpClient.Read(buf)
				if err != nil {
					logger.Errorf("出现异常 %v", err)
					time.Sleep(time.Second * 1)
					udpClient.Close()
					break
				}
				fmt.Printf("收到数据: %v\n", string(buf[:l]))
			}
		}(wg)
		wg.Wait()
	}

}
