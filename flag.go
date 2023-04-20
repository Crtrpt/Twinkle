package twinkle

import (
	"flag"
	"fmt"
	"os"
)

var ConfigFile *string
var Help *bool
var Ver *bool

func InitFlag() {
	ConfigFile = flag.String("f", "./conf/app.toml", "配置文件地址")
	Help = flag.Bool("h", false, "帮助文件")
	Ver = flag.Bool("v", false, "当前版本")
	flag.Parse()

	if false != *Help {
		fmt.Printf("-f ./conf/app.toml 指定配置文件 \r\n-h 显示帮助文件\r\n-v 显示当前版本\r\n")
		os.Exit(0)
	}

	if false != *Ver {
		fmt.Printf("当前版本:%s\r\n", VERSION)
		os.Exit(0)
	}

}
