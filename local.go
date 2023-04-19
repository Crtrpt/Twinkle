package gps

import (
	"net/http"
	"net/url"

	"github.com/Crtrpt/gps/logger"
)

// 监听本地端口
func (app *App) ListenLocalPort(cfg ProxyConfig, urlp *url.URL) {
	addr := urlp.Host
	logger.Infof("listen:%s", addr)
	err := http.ListenAndServe(urlp.Host, app)
	if err != nil {
		logger.Errorf("监听服务器异常 %s %v", addr, err)
	}
}
