declare namespace tw {
    interface TWProxy {
        Name: string            //当前代理的名称
        Desc: string           //当前代理的描述信息
        Enable: boolean              //是否启用当前代理
        Url: string            //要代理的地址
        Proxy: string            //代理请求的地址
        Header: { [key: string]: string } //输出的header
        Root: string            //静态文件目录 如果不存在 访问proxy的 后端服务器
        Interrupt: string            //中断文件路径
    }
}
