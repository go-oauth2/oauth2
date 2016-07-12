OAuth2授权码模式模拟
=================

运行服务端
--------
> 运行fasthttp服务端，请使用`cd example/fastserver`

```
$ cd example/server
$ go run main.go
```

运行客户端
--------

```
$ cd example/client
$ go run main.go
```

打开浏览器
--------

[http://localhost:9094](http://localhost:9094)

``` json
{
    "access_token": "143C1A45CFF9E0922F9DC68F7EBC81DC",
    "expires_in": 7200,
    "refresh_token": "5BD7453B8E7C5A3A308166F1675AD57216811391",
    "scope": "all",
    "token_type": "Bearer"
}
```