# wx
微信登录中继和AccessToken缓存

# 安装

## 二进制程序

请到[release][1]页选择适合的平台程序下载

## Docker

```bash
docker pull syutingsong/wx:latest
```

## 从源码编译

```bash
go get -u github.com/syutingsong/wx
```

# 使用

## 启动参数

```
wx [Flags]

  Flags: 
       --version  Displays the program version string.
    -h --help  Displays help with available flag, subcommand, and positional value parameters.
    -d --domains  valid domains for l2 jump
    -a --apps  weixin appId:appSecret pairs
    -l --log-level  log level
    -b --bind  binding ip address default: [::]
    -p --port  binding port number default: 3001
    -c --color  force use color for log output
```

## 例子

作为微信登录中继，支持两个域名
```bash
wx -d example.com -d xts.so -b 127.0.0.1
```

同时作为微信登录中继和AccessToken缓存
```
wx -d example.com -a <appid>:<secret>
```

显示DEBUG级别的日志

```
wx -l debug
```

## 在docker中使用

参数和直接使用一样，但可能需要暴露端口

```bash
docker run -ti --rm -p 80:80 syutingsong/wx -d example.com -p 80
```

# Web API

## 登录中继

### 用途

微信登录回调只能是公众号后台预留的URL地址，无法支持多环境多域名使用。
此登录中继服务可以实现登录回调的二传手：将本服务的URL预留在微信公众号后台，然后本服务可以将浏览器跳转到`l2`参数指定的URL上。

```
GET /login/l2
```

### 参数

name | required | describe
-----|----------|----------
l2   | 是 | 二次跳转URL，必须为启动参数`-d`所指定的域名或其子域
code | 否 | 微信公从号登录返回的`code`
state | 否 | 微信公众号登录返回的`state`
... | 否 | 其它参数

此登录中继会判断`l2`参数中的域名是否在白名单中，如果在的话，会将`code`、`state`及其它参数URL Encode追加在`l2`的末尾。

假设微信登录中继服务部署在`foo.example.com`下，微信登录跳转目标在`bar.example.com`下。

### 例子

浏览器来访：
```
https://example.com/login/l2?l2=https%3A%2F%2Fbar.example.com%2Fwechat%2Fcallback&code=a5472995f0ad2814851f18671830c068&state=abc
```

中继服务会返回：
```
HTTP/1.1 307 Temporary Redirect
Date: Mon, 24 Jun 2019 06:24:56 GMT
Content-Type: text/html; charset=utf-8
Content-Length: 0
Location: https://bar.example.com/wechat/callback?code=a5472995f0ad2814851f18671830c068&state=abc
```

引导浏览器转向实际的服务。

## AccessToken 缓存

```
GET /access_token
```

### 参数

name | required | describe
-----|----------|----------
appid | 是 | 公众号的AppID，此ID应已在`-a`启动参数中配置
force | 否 | 此参数非空字符时，会无视Cache强制从微信服务器获取

### 返回

```json
{
    "access_token": "22_1pQzga6gduIu11N84ph90GM5e9MEyUnUVXWmCrQO4-zgNbVCuS4pTiAtyoa4bz0EtotUeH98vj2pafKd00n51Ywgu5-4JH5ZXUZ5m78zXkajwDN-n-r9u3WeiBfWBrZPY663Q_W9QG9Cp1sMGKEeABATXR",
    "expires_in": 7118
}
```


[1]: https://github.com/syutingsong/wx/releases]
