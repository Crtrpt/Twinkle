package main

import (
	"fmt"
	"http_bridge/logger"
	"net/http"
	"time"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	logger.Infof("new request %+v", r)
	w.Write([]byte(fmt.Sprintf("url:%v\r\n", r.RequestURI)))
	w.Write([]byte("time:" + time.Now().Format(time.Layout)))
}

func main() {
	http.HandleFunc("/", Handle)
	http.ListenAndServe(":8088", nil)
}
