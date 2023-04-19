package gps

type Config struct {
	Proxy map[string]ProxyConfig
}

type ProxyConfig struct {
	Name   string
	Root   string            //静态文件目录 如果不存在 访问proxy的 后端服务器
	Url    string            //要代理的地址
	Proxy  string            //代理请求的地址
	Header map[string]string //输出的header
	Hook   map[string]string //各种hook的能力
}
