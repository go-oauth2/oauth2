OAuth 2.0
=========
>  [OAuth 2.0](http://oauth.net/2/) is the next evolution of the OAuth protocol which was originally created in late 2006.

[![GoDoc](https://godoc.org/gopkg.in/oauth2.v3?status.svg)](https://godoc.org/gopkg.in/oauth2.v3)
[![Go Report Card](https://goreportcard.com/badge/gopkg.in/oauth2.v3)](https://goreportcard.com/report/gopkg.in/oauth2.v3)

Quick Start
-----------

### Download and install

``` bash
$ go get -u gopkg.in/oauth2.v3/...
```

### Create file `server.go`

``` go
package main

import (
	"net/http"

	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store/token"
)

func main() {
	manager := manage.NewRedisManager(
		&token.RedisConfig{Addr: "192.168.33.70:6379"},
	)
	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		// validation and to get the user id
		userID = "000000"
		return
	})
	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleAuthorizeRequest(w, r)
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
	http.ListenAndServe(":9096", nil)
}
```

### Build and run

``` bash
$ go build server.go
$ ./server
```

Features
--------

* Based on the [RFC 6749](https://tools.ietf.org/html/rfc6749) implementation
* Easy to use
* Modularity
* Flexible
* Elegant

Example
-------

Simulation examples of authorization code model, please check [example](/example)

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
