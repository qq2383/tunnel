# Tunnel 
基于 socks5 的 TLS 加密高性能隧道（服务及客户端）

## 生成证书
```
# 生成服务器端的私钥
openssl genrsa -out cert.key 2048

# 生成服务器端证书
openssl req -new -x509 -key cert.key -out cert.pem -days 3650
```

## Server (tunneld)
实例： [cmd/tunneld](./cmd/tunneld) 
1. 修改 tunneld.cnf 文件 
```
// 隧道服务
// port 端口
// 证书路径
server:
    port: 10000
    cert_path: cert
// http 服务器
// port 端口
// 安全密码
http:
    port: 10001
    passwd: "123456"
```

2. 运行服务器
```
cd ./cmd/tunneld
go run .
```

3. 浏览器打开 https://ip:10001
4. 输入密码：123456
5. 进入 User, 新建一个用户：user1, 密码：123456

## Client (tunnel)
实例： [cmd/tunnel](./cmd/tunnel) 
1. 修改 tunnel.cnf 文件 
```
// socks 5 端口
local:
    port: 1080

// 远程隧道 ip 及端口
remote:
    addr: 127.0.0.1
    port: 10001

// 远程隧道的用户名及密码
user:
    -
        name: user
        passwd: 123456
```

2. 运行客户端
```
cd cmd/tunnel
go run .
```

3. 浏览器或系统设置代理
