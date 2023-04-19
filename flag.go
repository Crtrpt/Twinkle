package gps

import (
	"flag"
)

var ConfigFile *string

func InitFlag() {
	ConfigFile = flag.String("f", "./conf/app.toml", "配置文件地址")
	flag.Parse()
}
