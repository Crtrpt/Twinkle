package twinkle

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Crtrpt/twinkle/logger"
)

// 判断对静态文件的代理
func (app *App) ProxyStatic(resp http.ResponseWriter, r *http.Request, cfg *ProxyConfig, frontUrl, RequestURI string) bool {
	staticUrl := fmt.Sprintf("%s%s", cfg.Root, RequestURI)

	_, err := os.Stat(staticUrl)
	//如果没有错误 就输出文件信息
	if err == nil {
		logger.Infof("F:%50s S:%0s", frontUrl, staticUrl)
		//TODO 判断访问的目录是否是子目录 防止访问父目录文件
		http.ServeFile(resp, r, staticUrl)
		return true
	}
	return false
}
