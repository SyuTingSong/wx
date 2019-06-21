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

## 参数

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
wx -l debuy
```

## 在docker中使用

参数和直接使用一样，但可能需要暴露端口

```bash
docker run -ti --rm -p 80:80 syutingsong/wx -d example.com -p 80
```


[1]: https://github.com/syutingsong/wx/releases]
