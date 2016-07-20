Authorization code simulation
=============================

Run Server
---------

``` bash
$ cd example/server
$ go run main.go
```

Run Client
----------

```
$ cd example/client
$ go run main.go
```

Open the browser
----------------

[http://localhost:9094](http://localhost:9094)

``` json
{
    "access_token": "BIX-RYRPMHYY4L7O4QTP3Q",
    "expires_in": 7200,
    "refresh_token": "JRITD106WU6YNRE4UUEV_A",
    "scope": "all",
    "token_type": "Bearer"
}
```