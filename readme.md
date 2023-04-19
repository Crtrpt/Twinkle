golang proxy server
# 上图
![流程图](./flow.svg "工作流程图")

# 特性
ssh隧道代理


如果需要sshd远程隧道 需要开启 sshd 隧道功能 否则无法监听外部接口
```
GatewayPorts yes
```