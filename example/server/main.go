package main

import (
	"log"
	"net/http"

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

	srv := server.NewServer(server.NewConfig(), manager)

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		authReq, err := srv.GetAuthorizeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		authReq.UserID = "000000"
		// TODO: 登录验证、授权处理
		err = srv.HandleAuthorizeRequest(w, authReq)
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

	log.Println("OAuth2 server is running at 9096 port.")
	log.Fatal(http.ListenAndServe(":9096", nil))
}
