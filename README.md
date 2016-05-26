Golang OAuth 2.0
================

[![GoDoc](https://godoc.org/gopkg.in/oauth2.v1?status.svg)](https://godoc.org/gopkg.in/oauth2.v1)
[![Go Report Card](https://goreportcard.com/badge/gopkg.in/oauth2.v1)](https://goreportcard.com/report/gopkg.in/oauth2.v1)

> 基于Golang实现的OAuth 2.0协议相关操作，包括：令牌(或授权码)的生成、存储、验证操作以及更新令牌、废除令牌； 具有简单、灵活的特点； 其中所涉及的相关http请求操作在这里不做处理； 支持授权码模式、简化模式、密码模式、客户端模式； 默认使用MongoDB存储相关信息

获取
----

```bash
$ go get -v gopkg.in/oauth2.v1
```

范例
----

> 数据初始化：初始化相关的客户端信息

```go
package main

import (
	"fmt"

	"gopkg.in/oauth2.v1"
)

func main() {
	mongoConfig := oauth2.NewMongoConfig("mongodb://127.0.0.1:27017", "test")
	// 创建默认的OAuth2管理实例(基于MongoDB)
	manager, err := oauth2.CreateDefaultOAuthManager(mongoConfig, "", "", nil)
	if err != nil {
		panic(err)
	}
	// 模拟授权码模式
	// 使用默认参数，生成授权码
	code, err := manager.GetACManager().
		GenerateCode("clientID_x", "userID_x", "http://www.example.com/cb", "scopes")
	if err != nil {
		panic(err)
	}
	// 生成访问令牌及更新令牌
	genToken, err := manager.GetACManager().
		GenerateToken(code, "http://www.example.com/cb", "clientID_x", "clientSecret_x", true)
	if err != nil {
		panic(err)
	}
	// 检查访问令牌
	checkToken, err := manager.CheckAccessToken(genToken.AccessToken)
	if err != nil {
		panic(err)
	}
	// TODO: 使用用户标识、申请的授权范围响应数据
	fmt.Println(checkToken.UserID, checkToken.Scope)
	// 申请一个新的访问令牌
	newToken, err := manager.RefreshAccessToken(checkToken.RefreshToken, "scopes")
	if err != nil {
		panic(err)
	}
	fmt.Println(newToken.AccessToken, newToken.ATExpiresIn)
	// TODO: 将新的访问令牌响应给客户端
}
```

执行测试
----

```bash
$ go test -v
# 或
$ goconvey --port=9090
```

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
