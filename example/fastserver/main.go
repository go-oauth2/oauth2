package main

import (
	"log"

	"github.com/valyala/fasthttp"
	"gopkg.in/oauth2.v2/manage"
	"gopkg.in/oauth2.v2/models"
	"gopkg.in/oauth2.v2/server"
	"gopkg.in/oauth2.v2/store/client"
	"gopkg.in/oauth2.v2/store/token"
)

func main() {
	// 创建基于redis的oauth2管理实例
	manager := manage.NewRedisManager(
		&token.RedisConfig{Addr: "192.168.33.70:6379"},
	)
	// 使用临时客户端存储
	manager.MapClientStorage(client.NewTempStore(&models.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:9094",
	}))

	srv := server.NewFastServer(server.NewConfig(), manager)

	log.Println("OAuth2 server is running at 9096 port.")
	fasthttp.ListenAndServe(":9096", func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Request.URI().Path()) {
		case "/authorize":
			authReq, err := srv.GetAuthorizeRequest(ctx)
			if err != nil {
				ctx.Error(err.Error(), 400)
				return
			}
			authReq.UserID = "000000"
			// TODO: 登录验证、授权处理
			err = srv.HandleAuthorizeRequest(ctx, authReq)
			if err != nil {
				ctx.Error(err.Error(), 400)
			}
		case "/token":
			err := srv.HandleTokenRequest(ctx)
			if err != nil {
				ctx.Error(err.Error(), 400)
			}
		}
	})
}
