package gps

import (
	"context"
	"fmt"
	"gps/logger"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
)

type App struct {
	Config       *Config
	ProxyMapLock *sync.Mutex
	ProxyList    []*ProxyConfig
	ListenMap    map[string]any
}

func NewApp(ctx context.Context) *App {
	ProxyList := make([]*ProxyConfig, 0)
	ListenMap := make(map[string]any, 0)
	app := &App{
		ProxyMapLock: &sync.Mutex{},
		ProxyList:    ProxyList,
		ListenMap:    ListenMap,
	}
	app.InitConfig(ctx)
	return app
}

func (app *App) InitConfig(ctx context.Context) (res any, err error) {
	app.Config = &Config{}
	res, err = toml.DecodeFile(*ConfigFile, app.Config)
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	return
}

// 判断对静态文件的代理
func (app *App) ProxyStatic(resp http.ResponseWriter, r *http.Request, cfg *ProxyConfig, RequestURI string) bool {
	staticUrl := fmt.Sprintf("%s%s", cfg.Root, RequestURI)
	logger.Info("check static file  %+v", staticUrl)
	_, err := os.Stat(staticUrl)
	//如果没有错误 就输出文件信息
	if err == nil {
		//TODO 判断访问的目录是否是子目录 防止访问父目录文件
		http.ServeFile(resp, r, staticUrl)
		return true
	}
	return false
}

// 判断对后端文件的代理
func (app *App) ProxyBackend(resp http.ResponseWriter, r *http.Request, cfg *ProxyConfig, RequestURI string) (err error) {

	client := http.DefaultClient
	backendUrl := fmt.Sprintf("%s%s", cfg.Proxy, RequestURI)
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

func (app *App) ServeHTTP(resp http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Infof("backend error:  %s", err)
		}
		return
	}()

	url := fmt.Sprintf("%s://%s", "http", r.Host+r.RequestURI)
	logger.Infof("frontend:  %+v", url)

	for _, cfg := range app.ProxyList {
		if strings.HasPrefix(url, cfg.Url) {
			for k, v := range cfg.Header {
				resp.Header().Add(k, v)
			}
			if app.ProxyStatic(resp, r, cfg, url[len(cfg.Url):]) {
				return
			}
			err := app.ProxyBackend(resp, r, cfg, url[len(cfg.Url):])
			if err != nil {
				resp.Write([]byte(err.Error()))
			}
			return
		}
	}

	resp.Write([]byte("can not find backend serve"))
}

func (app *App) ListenNewPort(cfg ProxyConfig, urlp *url.URL) {
	addr := urlp.Host
	logger.Infof("listen:%s", addr)
	err := http.ListenAndServe(urlp.Host, app)
	if err != nil {
		logger.Errorf("监听服务器异常 %s %v", addr, err)
	}
}

// Run 执行
func (app *App) Run(ctx context.Context) (res any, err error) {
	for _, cfg := range app.Config.Proxy {
		cfg := cfg
		urlP, err := url.Parse(cfg.Url)
		if err != nil {
			logger.Errorf("url解析错误 %v", err)
		}
		port := "80"
		if urlP.Port() != "" {
			port = urlP.Port()
		}
		host := "127.0.0.1"
		if urlP.Host != "" {
			host = urlP.Hostname()
		}

		names, err := net.LookupIP(host)

		key := names[0].String() + ":" + port

		if app.ListenMap[key] == nil {
			//设置为监听状态
			app.ListenMap[key] = make(map[string]any, 0)
			go app.ListenNewPort(cfg, urlP)
		}
		app.ProxyList = append(app.ProxyList, &cfg)
	}

	sort.SliceStable(app.ProxyList, func(i, j int) bool {
		return len(app.ProxyList[i].Url) > len(app.ProxyList[j].Url)
	})
	return
}

// Stop 停止执行
func (app *App) Stop(ctx context.Context) (res any, err error) {
	logger.Info("stop")
	return
}

// Reload 重启
func (app *App) Reload(ctx context.Context) (res any, err error) {
	logger.Info("reload")
	return
}
