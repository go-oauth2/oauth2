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
    "access_token": "143C1A45CFF9E0922F9DC68F7EBC81DC",
    "expires_in": 7200,
    "refresh_token": "5BD7453B8E7C5A3A308166F1675AD57216811391",
    "scope": "all",
    "token_type": "Bearer"
}
```