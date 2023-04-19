package gps

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/Crtrpt/gps/logger"

	"github.com/BurntSushi/toml"
)

type App struct {
	Config       *Config
	ProxyMapLock *sync.Mutex
	ProxyList    []*ProxyConfig
	SSHTunnelMap map[string]any
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

func (app *App) ServeHTTP(resp http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Infof("backend error:  %s", err)
		}
		return
	}()

	url := fmt.Sprintf("%s://%s", "http", r.Host+r.RequestURI)

	for _, cfg := range app.ProxyList {
		if strings.HasPrefix(url, cfg.Url) {
			for k, v := range cfg.Header {
				resp.Header().Add(k, v)
			}
			if app.ProxyStatic(resp, r, cfg, url, url[len(cfg.Url):]) {
				return
			}
			err := app.ProxyBackend(resp, r, cfg, url, url[len(cfg.Url):])
			if err != nil {
				resp.Write([]byte(err.Error()))
			}
			return
		}
	}

	resp.Write([]byte("can not find backend serve"))
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

			if cfg.Ssh.Auth != "" {
				go app.ListenSSHTunnel(cfg)
			} else {
				app.ListenMap[key] = make(map[string]any, 0)
				go app.ListenLocalPort(cfg, urlP)
			}

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
