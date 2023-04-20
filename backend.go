package twinkle

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Crtrpt/twinkle/logger"
)

// 代理到后端服务器
func (app *App) ProxyBackend(resp http.ResponseWriter, r *http.Request, cfg *ProxyConfig, frontendUrl, RequestURI string) (err error) {
	client := http.DefaultClient
	backendUrl := fmt.Sprintf("%s%s", cfg.Proxy, RequestURI)
	logger.Infof("F:%-50s B:%-50s", frontendUrl, backendUrl)
	req, err := http.NewRequest(r.Method, backendUrl, r.Body)
	if err != nil {
		logger.Errorf("error %+v", err)
		return err
	}
	//Header复制
	req.Header = r.Header

	if r.Header.Get("X-Forwarded-For") == "" {
		req.Header.Set("X-Forwarded-For", r.RemoteAddr)
	} else {
		req.Header.Set("X-Forwarded-For", r.Header.Get("X-Forwarded-For")+","+r.RemoteAddr)
	}

	x, err := client.Do(req)
	if err != nil {
		logger.Errorf("backend error %+v", err)
		return err
	}
	resp.WriteHeader(x.StatusCode)
	data, err := io.ReadAll(x.Body)
	if err != nil {
		logger.Errorf("read backend err %v", err)
		return err
	}
	_, err = resp.Write(data)
	if err != nil {
		logger.Errorf("write front err  %v", err)
		return err
	}
	return nil
}
