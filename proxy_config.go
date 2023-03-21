package http_bridge

type Config struct {
	Proxy map[string]ProxyConfig
}

type ProxyConfig struct {
	Name  string
	Url   string            //要代理的地址
	Proxy string            //代理请求的地址
	Hook  map[string]string //各种hook的能力
}
