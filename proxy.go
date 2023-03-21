package http_bridge

import (
	"http_bridge/logger"
	"io"
	"net/http"
	"net/url"
)

type ProxyServer struct {
	Config ProxyConfig
}

func NewProxyServer(conf ProxyConfig) *ProxyServer {
	return &ProxyServer{Config: conf}
}

func (server *ProxyServer) Close() {
	logger.Info("close:" + server.Config.Name)
}

// ProxyBackend 代理请求到后端地址
func (server *ProxyServer) ProxyHandler(w http.ResponseWriter, r *http.Request) {
	client := http.DefaultClient
	logger.Infof("frontend:%s->backend:%s\r\n", r.URL.RequestURI(), server.Config.Proxy+r.URL.Path)
	req, err := http.NewRequest(r.Method, server.Config.Proxy+r.URL.Path, r.Body)
	if err != nil {
		logger.Errorf("error %+v", err)
	}
	//Header复制
	req.Header = r.Header
	x, err := client.Do(req)
	if err != nil {
		logger.Errorf("backend error %+v", err)
	}
	w.WriteHeader(x.StatusCode)
	data, err := io.ReadAll(x.Body)
	if err != nil {
		logger.Errorf("read backend err %v", err)
	}
	_, err = w.Write(data)
	if err != nil {
		logger.Errorf("write front err  %v", err)
	}
}

func (server *ProxyServer) Run() {
	url, err := url.Parse(server.Config.Url)
	if err != nil {
		logger.Errorf("parse url error %+v", err)
		return
	}

	path := "/"
	if url.Path != "" {
		path = url.Path
	}

	http.HandleFunc(path, server.ProxyHandler)
	err = http.ListenAndServe(":"+url.Port(), nil)
	if err != nil {
		logger.Errorf(" err %f", err)
	}
	logger.Infof("start: %+v", server.Config)
}
