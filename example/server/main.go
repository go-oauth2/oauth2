package main

import (
	"log"
	"net/http"

	"fmt"

	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store/client"
	"gopkg.in/oauth2.v3/store/token"
)

func main() {
	manager := manage.NewRedisManager(
		&token.RedisConfig{Addr: "192.168.33.70:6379"},
	)
	// Create the client temporary storage
	manager.MapClientStorage(client.NewTempStore(&models.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:9094",
	}))

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetAllowedResponseType(oauth2.Code)
	srv.SetAllowedGrantType(oauth2.AuthorizationCode)
	srv.SetErrorHandler(func(err error) {
		fmt.Println("OAuth2 Error:", err.Error())
	})
	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
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

	log.Println("Server is running at 9096 port.")
	log.Fatal(http.ListenAndServe(":9096", nil))
}
