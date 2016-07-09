OAuth2服务端
===========

> 基于Golang实现的OAuth2协议，具有简单化、模块化的特点

[![GoDoc](https://godoc.org/gopkg.in/oauth2.v2?status.svg)](https://godoc.org/gopkg.in/oauth2.v2)
[![Go Report Card](https://goreportcard.com/badge/gopkg.in/oauth2.v2)](https://goreportcard.com/report/gopkg.in/oauth2.v2)

获取
----

```bash
$ go get -u gopkg.in/oauth2.v2/...
```

使用
----

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

测试
----

``` bash
$ goconvey -port=9092
```

> goconvey使用明细[https://github.com/smartystreets/goconvey](https://github.com/smartystreets/goconvey)

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
