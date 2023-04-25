package twinkle

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Crtrpt/twinkle/logger"
	"github.com/julienschmidt/httprouter"
)

type Admin struct {
	Ctx      context.Context
	Config   *Config
	User     string
	Password string
}

// 获取首页静态文件
func Doc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "./doc/book/")
}

// 重定向到dashboard入口
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Location", `/index/index.html`)
	w.WriteHeader(301)
}

// 获取配置信息
func (admin *Admin) GetConfig(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cfg, err := json.Marshal(admin.Config.Proxy)
	if err != nil {
		w.Write([]byte("反序列化解析错误"))
	}
	w.Write(cfg)
}

// 基础授权
func (admin *Admin) BaseAuth(h httprouter.Handle) httprouter.Handle {
	return func(resp http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		user, pass, ok := req.BasicAuth()
		if ok && user == admin.User && pass == admin.Password {
			h(resp, req, ps)
			return
		} else {
			resp.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			resp.WriteHeader(401)
			return
		}
	}
}

// 启动admin
func AdminRun(ctx context.Context, cfg *Config) {
	logger.Infof("启动admin 管理端%s", cfg.Admin.Host)
	admin := Admin{
		User:     cfg.Admin.UserName,
		Password: cfg.Admin.Password,
		Ctx:      ctx,
		Config:   cfg,
	}
	router := httprouter.New()

	router.GET("/", Index)
	router.ServeFiles("/index/*filepath", http.Dir("./dist"))
	router.ServeFiles("/doc/*filepath", http.Dir("./doc/book"))
	router.GET("/api", admin.BaseAuth(admin.GetConfig))

	err := http.ListenAndServe(cfg.Admin.Host, router)
	if err != nil {
		logger.Errorf("admin err %v", err)
	}
}
