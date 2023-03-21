package http_bridge

import (
	"flag"
)

var ConfigFile *string

func InitFlag() {
	ConfigFile = flag.String("f", "./conf/default.toml", "配置文件地址")
	flag.Parse()
}
