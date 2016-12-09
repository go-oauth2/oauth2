# Golang OAuth 2.0

>  An open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications.

[![License][License-Image]][License-Url] [![ReportCard][ReportCard-Image]][ReportCard-Url] [![Build][Build-Status-Image]][Build-Status-Url] [![Coverage][Coverage-Image]][Coverage-Url] [![GoDoc][GoDoc-Image]][GoDoc-Url] [![Release][Release-Image]][Release-Url] 

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
	"log"
	"net/http"

	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

func main() {
	manager := manage.NewDefaultManager()
	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// client memory store
	clientStore := store.NewClientStore()
	clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "http://localhost",
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		srv.HandleTokenRequest(w, r)
	})

	log.Fatal(http.ListenAndServe(":9096", nil))
}

```

### Build and run

``` bash
$ go build server.go
$ ./server
```

### Open in your web browser

```
http://localhost:9096/token?grant_type=client_credentials&client_id=000000&client_secret=999999&scope=read
```

``` json
{
    "access_token": "J86XVRYSNFCFI233KXDL0Q",
    "expires_in": 7200,
    "scope": "read",
    "token_type": "Bearer"
}
```

## Features

* easy to use
* based on the [RFC 6749](https://tools.ietf.org/html/rfc6749) implementation
* token storage support TTL
* support custom expiration time of the access token
* support custom extension field
* support custom scope

## Example

> A complete example of simulation authorization code model

Simulation examples of authorization code model, please check [example](/example)

## Storage Implements

* [BuntDB](https://github.com/tidwall/buntdb)(The default storage)
* [Redis](https://github.com/go-oauth2/redis)
* [MongoDB](https://github.com/go-oauth2/mongo)

## MIT License

```
Copyright (c) 2016 Lyric
```

[License-Url]: http://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/npm/l/express.svg
[Build-Status-Url]: https://travis-ci.org/go-oauth2/oauth2
[Build-Status-Image]: https://travis-ci.org/go-oauth2/oauth2.svg?branch=master
[Release-Url]: https://github.com/go-oauth2/oauth2/releases/tag/v3.7.0
[Release-image]: http://img.shields.io/badge/release-v3.7.0-1eb0fc.svg
[ReportCard-Url]: https://goreportcard.com/report/gopkg.in/oauth2.v3
[ReportCard-Image]: https://goreportcard.com/badge/gopkg.in/oauth2.v3
[GoDoc-Url]: https://godoc.org/gopkg.in/oauth2.v3
[GoDoc-Image]: https://godoc.org/gopkg.in/oauth2.v3?status.svg
[Coverage-Url]: https://coveralls.io/github/go-oauth2/oauth2?branch=master
[Coverage-Image]: https://coveralls.io/repos/github/go-oauth2/oauth2/badge.svg?branch=master
