package gps

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Crtrpt/gps/logger"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
)

type Interrupt interface {
	Run() error
	SetHeader(key string, val any)
	Do(resp http.ResponseWriter, r *http.Request, cfg *ProxyConfig, url, RequestURI string)
}

type JavaScriptInterrupt struct {
	Script *otto.Script
	Vm     *otto.Otto
}

func NewJavascriptVm(cfg *ProxyConfig) (*JavaScriptInterrupt, error) {
	vm := otto.New()
	source, err := ioutil.ReadFile(cfg.Interrupt)
	if err != nil {
		return nil, err
	}
	script, err := vm.Compile(cfg.Interrupt, string(source))
	if err != nil {
		return nil, err
	}
	return &JavaScriptInterrupt{
		Script: script,
		Vm:     vm,
	}, nil
}

func (vm *JavaScriptInterrupt) Run(app *App, resp http.ResponseWriter, req *http.Request, cfg *ProxyConfig, url, RequestURI string) error {

	//获取url
	vm.Vm.Set("GetUrl", func() otto.Value {
		v, _ := otto.ToValue(url)
		return v
	})

	//获取请求的方法
	vm.Vm.Set("GetMethod", func() otto.Value {
		v, err := otto.ToValue(req.Method)
		if err != nil {
			logger.Error(err)
		}
		return v
	})

	//获取请求的header
	vm.Vm.Set("SetRequestHeader", func(key otto.Value, val otto.Value) otto.Value {
		req.Header.Set(key.String(), val.String())
		return otto.Value{}
	})

	//获取请求的header
	vm.Vm.Set("GetRequestHeader", func(key otto.Value) otto.Value {
		v, err := otto.ToValue(req.Header.Get(key.String()))
		if err != nil {
			logger.Error(err)
		}
		return v
	})

	//设置返回的header
	vm.Vm.Set("SetResponseHeader", func(key otto.Value, val otto.Value) otto.Value {
		resp.Header().Set(key.String(), val.String())
		return otto.Value{}
	})

	//获取返回的header
	vm.Vm.Set("GetResponseHeader", func(key otto.Value) otto.Value {
		v, err := otto.ToValue(resp.Header().Get(key.String()))
		if err != nil {
			logger.Error(err)
		}
		return v
	})

	//获取请求的body
	vm.Vm.Set("GetRequestBody", func() otto.Value {
		buf := new(bytes.Buffer)
		buf.ReadFrom(req.Body)
		newStr := buf.String()
		v, err := otto.ToValue(newStr)
		if err != nil {
			logger.Error(err)
		}
		return v
	})

	//获取返回的body
	vm.Vm.Set("SetCode", func(body otto.Value) {
		codeStr, err := body.ToString()
		if err != nil {
			logger.Error(err)
		}
		code, err := strconv.Atoi(codeStr)
		if err != nil {
			logger.Error(err)
		}
		resp.WriteHeader(code)
	})

	//设置body
	vm.Vm.Set("SetBody", func(body otto.Value) otto.Value {
		resp.Write([]byte(body.String()))
		return otto.Value{}
	})

	//执行默认全球
	vm.Vm.Set("Run", func() otto.Value {
		app.Do(resp, req, cfg, url, RequestURI)
		return otto.Value{}
	})

	//每次执行都会加载 TODO 修改为判断更新
	source, err := ioutil.ReadFile(cfg.Interrupt)
	if err != nil {
		return err
	}
	vm.Script, err = vm.Vm.Compile(cfg.Interrupt, string(source))
	if err != nil {
		return err
	}

	_, err = vm.Vm.Run(vm.Script)

	if err != nil {
		logger.Errorf("vm error", err)
		return err
	}

	return nil
}
