package twinkle

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/Crtrpt/twinkle/logger"

	"github.com/BurntSushi/toml"
)

type App struct {
	Config       *Config
	ProxyMapLock *sync.Mutex
	ProxyList    []*ProxyConfig
	SSHTunnelMap map[string]any //ssh 隧道
	ListenMap    map[string]any //http 代理
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

func (app *App) Do(resp http.ResponseWriter, r *http.Request, cfg *ProxyConfig, url, RequestURI string) (err error) {
	if app.ProxyStatic(resp, r, cfg, url, RequestURI) {
		return
	}
	return app.ProxyBackend(resp, r, cfg, url, RequestURI)
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
			var err error
			if cfg.Interrupt == "" {
				err = app.Do(resp, r, cfg, url, url[len(cfg.Url):])
			} else {
				interrupt, err := NewJavascriptVm(cfg)
				if err == nil {
					err = interrupt.Run(app, resp, r, cfg, url, url[len(cfg.Url):])
				}
			}
			if err != nil {
				resp.Write([]byte(err.Error()))
			}
			return
		}
	}

	resp.Write([]byte("can not find backend serve"))
}

func GetTransportLayer(protocol string) string {
	if protocol == "http" || protocol == "https" || protocol == "tcp" {
		return "tcp"
	}
	if protocol == "udp" {
		return "udp"
	}
	return ""
}

// Run 执行
func (app *App) Run(ctx context.Context) (res any, err error) {
	for _, cfg := range app.Config.Proxy {
		cfg := cfg
		urlP, err := url.Parse(cfg.Url)
		if err != nil {
			logger.Errorf("url解析错误 %v", err)
		}

		port := urlP.Port()
		host := urlP.Hostname()

		names, err := net.LookupIP(host)

		key := cfg.Ssh.Addr + GetTransportLayer(urlP.Scheme) + names[0].String() + ":" + port

		if app.ListenMap[key] == nil {
			//设置为监听状态
			if cfg.Ssh.Auth != "" {
				go app.ListenSSHTunnel(cfg)
			} else {
				app.ListenMap[key] = make(map[string]any, 0)
				go app.ListenHttpPort(urlP.Scheme, cfg, urlP)
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
