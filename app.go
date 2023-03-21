package http_bridge

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"http_bridge/logger"
	"net"
	"net/http"
	"net/url"
	"sync"
)

type ProxyHandler interface {
	ProxyHandler(w http.ResponseWriter, r *http.Request)
}

type App struct {
	Config       *Config
	ProxyMapLock *sync.Mutex
	ProxyMap     map[string]any
}

func NewApp(ctx context.Context) *App {
	proxyMap := make(map[string]any, 0)
	app := &App{
		ProxyMapLock: &sync.Mutex{},
		ProxyMap:     proxyMap,
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

type ProxyDispatch struct {
	Port  string
	Proxy map[string]ProxyConfig
}

func (receiver ProxyDispatch) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	fmt.Printf("req: host:%s \r\n", req.Host)
}

// Run 执行
func (app *App) Run(ctx context.Context) (res any, err error) {
	logger.Info("run")
	for _, cfg := range app.Config.Proxy {
		go func(cfg ProxyConfig) {
			defer func() {
				if v := recover(); v != nil {
					logger.Errorf("异常 %v", v)
					return
				}
			}()
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
			fmt.Printf("%v \r\n", names)
			key := names[0].String() + ":" + port
			if app.ProxyMap[key] == nil {
				app.ProxyMapLock.Lock()
				defer app.ProxyMapLock.Unlock()
				app.ProxyMap[key] = true
				go func(cfg ProxyConfig, urlp *url.URL) {
					//fmt.Printf("启动一个端口监听器 %s \r\n", ":"+port)
					disp := ProxyDispatch{
						Port:  port,
						Proxy: app.Config.Proxy,
					}
					addr := urlp.Host
					logger.Infof("listen:%s", addr)
					err := http.ListenAndServe(urlp.Host, disp)
					if err != nil {
						logger.Errorf("监听服务器异常 %s %v", addr, err)
					}
				}(cfg, urlP)
			}
		}(cfg)
	}
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
