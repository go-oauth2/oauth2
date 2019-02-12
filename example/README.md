# Use Examples

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

## Authorization Code Grant

### Open the browser

[http://localhost:9094](http://localhost:9094)

```
{
  "access_token": "GIGXO8XWPQSAUGOYQGTV8Q",
  "token_type": "Bearer",
  "refresh_token": "5FBLXQ47XJ2MGTY8YRZQ8W",
  "expiry": "2019-01-08T01:53:45.868194+08:00"
}
```


### Try access token

Open the browser [http://localhost:9094/try](http://localhost:9094/try)

```
{
  "client_id": "222222",
  "expires_in": 7195,
  "user_id": "000000"
}
```

## Refresh token

Open the browser [http://localhost:9094/refresh](http://localhost:9094/refresh)

```
{
  "access_token": "0IIL4_AJN2-SR0JEYZVQWG",
  "token_type": "Bearer",
  "refresh_token": "AG6-63MLXUEFUV2Q_BLYIW",
  "expiry": "2019-01-09T23:03:16.374062+08:00"
}
```

## Password Credentials Grant

Open the browser [http://localhost:9094/pwd](http://localhost:9094/pwd)

```
{
  "access_token": "87JT3N6WOWANXVDNZFHY7Q",
  "token_type": "Bearer",
  "refresh_token": "LDIS6PXAVY-BXHPEDESWNG",
  "expiry": "2019-02-12T10:58:43.734902+08:00"
}
```

## Client Credentials Grant

Open the browser [http://localhost:9094/client](http://localhost:9094/client)

```
{
  "access_token": "OA6ITALNMDOGD58C0SN-MG",
  "token_type": "Bearer",
  "expiry": "2019-02-12T11:10:35.864838+08:00"
}
```

![login](https://raw.githubusercontent.com/go-oauth2/oauth2/master/example/server/static/login.png)
![auth](https://raw.githubusercontent.com/go-oauth2/oauth2/master/example/server/static/auth.png)
![token](https://raw.githubusercontent.com/go-oauth2/oauth2/master/example/server/static/token.png)
