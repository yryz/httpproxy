使用Golang实现的HTTP代理转shadowsocks，主要为命令行下`go get`、`docker pull`、`npm install`、`pip install`、`gem install`、`curl`等程序提供HTTP代理服务，解决安装总是失败的问题。这些服务不支持shadowsocks，但对http代理都有支持。

##使用
**安装**

`go get github.com/yryz/httpproxy`

配置文件 ~/.httpproxy/config.json

```
{
        "listen": "127.0.0.1:6666",
        "ss_server": "ip:port",
        "ss_cipher": "aes-128-cfb",
        "ss_password": "your password"
}
```
启动 `httpproxy`

**使用代理**

如果想命令行一直走代理，下面配置加入到 ~/.bash_profile

```
http_proxy=http://127.0.0.1:6666
https_proxy=http://127.0.0.1:6666
```

如果只是想临时使用，可以手动设置http_proxy环境变量或者 使用`httpproxy set` 快速设置（推荐！）。

##特点

* 支持与shadowsocks服务桥接
* 支持CONNECT，支持HTTPS、HTTP2代理

