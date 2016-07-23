OAuth 2.0
=========
>  [OAuth 2.0](http://oauth.net/2/) is the next evolution of the OAuth protocol which was originally created in late 2006.

[![GoDoc](https://godoc.org/gopkg.in/oauth2.v3?status.svg)](https://godoc.org/gopkg.in/oauth2.v3)
[![Go Report Card](https://goreportcard.com/badge/gopkg.in/oauth2.v3)](https://goreportcard.com/report/gopkg.in/oauth2.v3)
[![Build Status](https://travis-ci.org/go-oauth2/oauth2.svg?branch=master)](https://travis-ci.org/go-oauth2/oauth2)

Quick Start
-----------

### Download and install

``` bash
$ go get -u -v gopkg.in/oauth2.v3
```

### Create file `server.go`

``` go
package main

import (
	"net/http"

	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

func main() {
	manager := manage.NewDefaultManager()
	// token memory store
	manager.MapTokenStorage(store.NewMemoryTokenStore(0))
	// client test store
	manager.MapClientStorage(store.NewTestClientStore())

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

Example
-------

Simulation examples of authorization code model, please check [example](/example)

Storage implements
------------------

* [Redis](https://github.com/go-oauth2/redis)
* [MongoDB](https://github.com/go-oauth2/mongo)

License
-------

```
Copyright (c) 2016, OAuth 2.0
All rights reserved.
```