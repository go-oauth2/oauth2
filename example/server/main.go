package main

import (
	"fmt"
	"log"
	"net/http"

	"net/url"

	"os"

	"github.com/gorilla/sessions"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store/client"
	"gopkg.in/oauth2.v3/store/token"
)

var (
	sessionStore *sessions.CookieStore
)

func main() {
	// Create the session store
	sessionStore = sessions.NewCookieStore([]byte("123456"))

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

	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetErrorHandler(func(err error) {
		fmt.Println("OAuth2 Error:", err.Error())
	})

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/auth", authHandler)

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

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	us, err := sessionStore.Get(r, "user")
	if err != nil {
		return
	}
	if us.IsNew {
		r.ParseForm()
		fs, _ := sessionStore.Get(r, "form")
		fs.Values["Form"] = r.Form
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = us.Values["UserID"].(string)
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		us, _ := sessionStore.Get(r, "user")
		us.Values["UserID"] = "000000"
		us.Save(r, w)
		w.Header().Set("Location", "/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "static/login.html")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	us, _ := sessionStore.Get(r, "user")
	if us.IsNew {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	if r.Method == "POST" {
		fs, _ := sessionStore.Get(r, "form")
		values := fs.Values["Form"].(url.Values)
		w.Header().Set("Location", "/authorize?"+values.Encode())
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "static/auth.html")
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}
