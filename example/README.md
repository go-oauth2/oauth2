# Authorization Code Grant

![login](https://raw.githubusercontent.com/go-oauth2/oauth2/master/example/server/static/login.png)
![auth](https://raw.githubusercontent.com/go-oauth2/oauth2/master/example/server/static/auth.png)
![token](https://raw.githubusercontent.com/go-oauth2/oauth2/master/example/server/static/token.png)

## Run Server

``` bash
$ cd example/server
$ go build server.go
$ ./server
```

## Run Client

```
$ cd example/client
$ go build client.go
$ ./client
```

## Open the browser

[http://localhost:9094](http://localhost:9094)

```
{
  "access_token": "GIGXO8XWPQSAUGOYQGTV8Q",
  "token_type": "Bearer",
  "refresh_token": "5FBLXQ47XJ2MGTY8YRZQ8W",
  "expiry": "2019-01-08T01:53:45.868194+08:00"
}
```


## Try access token

Open the browser [http://localhost:9094/try](http://localhost:9094/try)

```
{
  "client_id": "222222",
  "expires_in": 7195,
  "user_id": "000000"
}
```