package server_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gavv/httpexpect"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store/client"
	"gopkg.in/oauth2.v3/store/token"
)

var (
	srv     *server.Server
	tsrv    *httptest.Server
	manager *manage.Manager
	csrv    *httptest.Server
)

func init() {
	manager = manage.NewRedisManager(
		&token.RedisConfig{Addr: "192.168.33.70:6379"},
	)
}

func testServer(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/authorize":
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			t.Error(err)
		}
	case "/token":
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestAuthorizeCode(t *testing.T) {
	tsrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testServer(t, w, r)
	}))
	defer tsrv.Close()
	e := httpexpect.New(t, tsrv.URL)

	csrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2":
			r.ParseForm()
			code, state := r.Form.Get("code"), r.Form.Get("state")
			if state != "123" {
				t.Error("unrecognized state:", state)
				return
			}
			val := e.POST("/token").
				WithFormField("redirect_uri", csrv.URL+"/oauth2").
				WithFormField("code", code).
				WithFormField("grant_type", "authorization_code").
				WithFormField("client_id", "333333").
				WithFormField("client_secret", "33333333").
				Expect().
				Status(http.StatusOK).
				JSON().Raw()

			t.Log(val)
		}
	}))
	defer csrv.Close()

	manager.MapClientStorage(client.NewTempStore(&models.Client{
		ID:     "333333",
		Secret: "33333333",
		Domain: csrv.URL,
	}))

	srv = server.NewServer(server.NewConfig(), manager)
	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		userID = "111111"
		return
	})

	e.GET("/authorize").
		WithQuery("response_type", "code").
		WithQuery("client_id", "333333").
		WithQuery("scope", "all").
		WithQuery("state", "123").
		WithQuery("redirect_uri", url.QueryEscape(csrv.URL+"/oauth2")).
		Expect().Status(http.StatusOK)
}

func TestImplicit(t *testing.T) {
	tsrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testServer(t, w, r)
	}))
	defer tsrv.Close()
	e := httpexpect.New(t, tsrv.URL)

	csrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2":
			t.Log(r.RequestURI)
		}
	}))
	defer csrv.Close()

	manager.MapClientStorage(client.NewTempStore(&models.Client{
		ID:     "55555",
		Secret: "5555555",
		Domain: csrv.URL,
	}))

	srv = server.NewServer(server.NewConfig(), manager)
	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		userID = "222222"
		return
	})

	e.GET("/authorize").
		WithQuery("response_type", "token").
		WithQuery("client_id", "55555").
		WithQuery("scope", "all").
		WithQuery("state", "123").
		WithQuery("redirect_uri", url.QueryEscape(csrv.URL+"/oauth2")).
		Expect().Status(http.StatusOK)
}

func TestPasswordCredentials(t *testing.T) {
	tsrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testServer(t, w, r)
	}))
	defer tsrv.Close()
	e := httpexpect.New(t, tsrv.URL)

	manager.MapClientStorage(client.NewTempStore(&models.Client{
		ID:     "666666",
		Secret: "66666666",
		Domain: csrv.URL,
	}))

	srv = server.NewServer(server.NewConfig(), manager)
	srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		if username == "admin" && password == "123456" {
			userID = "666666"
			return
		}
		err = errors.New("user not found")
		return
	})

	val := e.POST("/token").
		WithFormField("grant_type", "password").
		WithFormField("client_id", "666666").
		WithFormField("client_secret", "66666666").
		WithFormField("username", "admin").
		WithFormField("password", "123456").
		WithFormField("scope", "all").
		Expect().
		Status(http.StatusOK).
		JSON().Raw()

	t.Log(val)
}

func TestClientCredentials(t *testing.T) {
	tsrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testServer(t, w, r)
	}))
	defer tsrv.Close()
	e := httpexpect.New(t, tsrv.URL)

	manager.MapClientStorage(client.NewTempStore(&models.Client{
		ID:     "777777",
		Secret: "77777777",
		Domain: csrv.URL,
	}))

	srv = server.NewServer(server.NewConfig(), manager)

	val := e.POST("/token").
		WithFormField("grant_type", "clientcredentials").
		WithFormField("client_id", "777777").
		WithFormField("client_secret", "77777777").
		WithFormField("scope", "all").
		Expect().
		Status(http.StatusOK).
		JSON().Raw()

	t.Log(val)
}

func TestRefreshing(t *testing.T) {
	tsrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testServer(t, w, r)
	}))
	defer tsrv.Close()
	e := httpexpect.New(t, tsrv.URL)

	csrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2":
			r.ParseForm()
			code, state := r.Form.Get("code"), r.Form.Get("state")
			if state != "123" {
				t.Error("unrecognized state:", state)
				return
			}
			jval := e.POST("/token").
				WithFormField("redirect_uri", csrv.URL+"/oauth2").
				WithFormField("code", code).
				WithFormField("grant_type", "authorization_code").
				WithFormField("client_id", "888888").
				WithFormField("client_secret", "88888888").
				Expect().
				Status(http.StatusOK).
				JSON()

			refresh := jval.Object().Value("refresh_token").String().Raw()

			rval := e.POST("/token").
				WithFormField("grant_type", "refreshtoken").
				WithFormField("client_id", "888888").
				WithFormField("client_secret", "88888888").
				WithFormField("scope", "one").
				WithFormField("refresh_token", refresh).
				Expect().
				Status(http.StatusOK).
				JSON().Raw()

			t.Log(rval)
		}
	}))
	defer csrv.Close()

	manager.MapClientStorage(client.NewTempStore(&models.Client{
		ID:     "888888",
		Secret: "88888888",
		Domain: csrv.URL,
	}))

	srv = server.NewServer(server.NewConfig(), manager)
	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		userID = "888888"
		return
	})

	e.GET("/authorize").
		WithQuery("response_type", "code").
		WithQuery("client_id", "888888").
		WithQuery("scope", "all").
		WithQuery("state", "123").
		WithQuery("redirect_uri", url.QueryEscape(csrv.URL+"/oauth2")).
		Expect().Status(http.StatusOK)
}
