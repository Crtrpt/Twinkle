package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Crtrpt/twinkle/logger"
	"github.com/creack/pty"
	"github.com/fsnotify/fsnotify"
)

func runCmd(s string, done chan struct{}) {
	cmds := strings.Split(s, " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	// os.StartProcess(cmds[0], cmds[1:], nil)
	f, err := pty.Start(cmd)

	if err != nil {
		panic(err)
	}

	go func() {
		<-done
		defer f.Close()
		if cmd.Process != nil && cmd.ProcessState != nil {
			err := syscall.Kill(cmd.Process.Pid, syscall.SIGKILL)
			if err != nil {
				logger.Error("err kill %v", err)
				return
			}
			_, err = cmd.Process.Wait()
			if err != nil {
				logger.Error("err wait %v", err)
				return
			}
		}

	}()
	io.Copy(os.Stdout, f)

}

func main() {
	path := flag.String("d", "", "要监听的目录或者文件")
	run := flag.String("r", "", "程序变更之后要运行的程序")
	flag.Parse()
	done := make(chan struct{})
	go runCmd(*run, done)
	c := make(chan os.Signal)
	signal.Notify(c, os.Kill, os.Interrupt)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}
				done <- struct{}{}
				go runCmd(*run, done)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(*path)
	if err != nil {
		log.Fatal(err)
	}
	_ = <-c

}
