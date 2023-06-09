package twinkle

type Config struct {
	Proxy map[string]ProxyConfig `toml:"proxy"`
	Admin AdminConfig
}

type AdminConfig struct {
	UserName string
	Password string
	Host     string
}

type ProxyConfig struct {
	Name      string            //当前代理的名称
	Desc      string            //当前代理的描述信息
	Enable    bool              //是否启用当前代理
	Url       string            //要代理的地址
	Proxy     string            //代理请求的地址
	Header    map[string]string //输出的header
	Hook      map[string]string //各种hook的能力
	Root      string            //静态文件目录 如果不存在 访问proxy的 后端服务器
	Interrupt string            //中断文件路径
	Ssh       struct {
		Auth             string //password key
		Host             string
		UserName         string //用户名
		Password         string //密码
		PrivateKey       string //私钥登陆
		Addr             string //要监听的地址
		Timeout          int    //超时秒数
		RemoteUdpOverTcp string //处理udp的转发
	} `toml:"ssh"`
}
