基于Golang的OAuth2服务实现
=======================

> 完全模块化、支持http/fasthttp的服务端处理、令牌存储支持redis/mongodb

[![GoDoc](https://godoc.org/gopkg.in/oauth2.v2?status.svg)](https://godoc.org/gopkg.in/oauth2.v2)
[![Go Report Card](https://goreportcard.com/badge/gopkg.in/oauth2.v2)](https://goreportcard.com/report/gopkg.in/oauth2.v2)

获取
----

``` bash
$ go get -u gopkg.in/oauth2.v2/...
```

HTTP服务端
--------

``` go
package main

import (
	"log"
	"net/http"

	"gopkg.in/oauth2.v2/manage"
	"gopkg.in/oauth2.v2/models"
	"gopkg.in/oauth2.v2/server"
	"gopkg.in/oauth2.v2/store/client"
	"gopkg.in/oauth2.v2/store/token"
)

func main() {
	manager := manage.NewRedisManager(
		&token.RedisConfig{Addr: "192.168.33.70:6379"},
	)
	manager.MapClientStorage(client.NewTempStore())
	srv := server.NewServer(server.NewConfig(), manager)

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		authReq, err := srv.GetAuthorizeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// TODO: 登录验证、授权处理
        authReq.UserID = "000000"

		err = srv.HandleAuthorizeRequest(w, authReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	log.Fatal(http.ListenAndServe(":9096", nil))
}

```

FastHTTP服务端
-------------

``` go
srv := server.NewFastServer(server.NewConfig(), manager)

fasthttp.ListenAndServe(":9096", func(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Request.URI().Path()) {
	case "/authorize":
		authReq, err := srv.GetAuthorizeRequest(ctx)
		if err != nil {
			ctx.Error(err.Error(), 400)
			return
		}
		authReq.UserID = "000000"
		// TODO: 登录验证、授权处理
		err = srv.HandleAuthorizeRequest(ctx, authReq)
		if err != nil {
			ctx.Error(err.Error(), 400)
		}
	case "/token":
		err := srv.HandleTokenRequest(ctx)
		if err != nil {
			ctx.Error(err.Error(), 400)
		}
	}
})
```

测试
----
> [goconvey](https://github.com/smartystreets/goconvey)

``` bash
$ goconvey -port=9092
```

范例
----

模拟授权码模式的测试范例，请查看[example](/example)


License
-------

```
Copyright 2016.All rights reserved.
```

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

```
   http://www.apache.org/licenses/LICENSE-2.0
```

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
