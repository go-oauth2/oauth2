# OAuth 2.0

>  An open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications.

[![GoDoc](https://godoc.org/gopkg.in/oauth2.v3?status.svg)](https://godoc.org/gopkg.in/oauth2.v3)
[![Go Report Card](https://goreportcard.com/badge/gopkg.in/oauth2.v3)](https://goreportcard.com/report/gopkg.in/oauth2.v3)
[![Build Status](https://travis-ci.org/go-oauth2/oauth2.svg?branch=master)](https://travis-ci.org/go-oauth2/oauth2)

## Protocol Flow

```
     +--------+                               +---------------+
     |        |--(A)- Authorization Request ->|   Resource    |
     |        |                               |     Owner     |
     |        |<-(B)-- Authorization Grant ---|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(C)-- Authorization Grant -->| Authorization |
     | Client |                               |     Server    |
     |        |<-(D)----- Access Token -------|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(E)----- Access Token ------>|    Resource   |
     |        |                               |     Server    |
     |        |<-(F)--- Protected Resource ---|               |
     +--------+                               +---------------+
```

## Quick Start

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
	"gopkg.in/oauth2.v3/store"
)

func main() {
	manager := manage.NewDefaultManager()
	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	// client test store
	manager.MapClientStorage(store.NewTestClientStore())

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)

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

### Open in your web browser

```
http://localhost:9096/token?grant_type=clientcredentials&client_id=1&client_secret=11&scope=all
```

```
{
    "access_token": "ZGF4ARHJPT2Y_QAIOJVL-Q",
    "expires_in": 7200,
    "scope": "all",
    "token_type": "Bearer"
}
```

## Features

* Easy to use
* Based on the [RFC 6749](https://tools.ietf.org/html/rfc6749) implementation
* Token storage support TTL
* Support custom extension field
* Support custom scope
* Support custom expiration time of the access token

## Example

> A complete example of simulation authorization code model

Simulation examples of authorization code model, please check [example](/example)

## Storage implements

* [BuntDB](https://github.com/tidwall/buntdb)(The default storage)
* [Redis](https://github.com/go-oauth2/redis)
* [MongoDB](https://github.com/go-oauth2/mongo)

## License

```
Copyright (c) 2016, OAuth 2.0
All rights reserved.
```