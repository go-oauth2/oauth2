package server_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

var (
	srv          *server.Server
	tsrv         *httptest.Server
	manager      *manage.Manager
	csrv         *httptest.Server
	clientID     = "111111"
	clientSecret = "11111111"
)

func init() {
	manager = manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())
}

func clientStore(domain string) oauth2.ClientStore {
	clientStore := store.NewClientStore()
	clientStore.Set(clientID, &models.Client{
		ID:     clientID,
		Secret: clientSecret,
		Domain: domain,
	})
	return clientStore
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
			resObj := e.POST("/token").
				WithFormField("redirect_uri", csrv.URL+"/oauth2").
				WithFormField("code", code).
				WithFormField("grant_type", "authorization_code").
				WithFormField("client_id", clientID).
				WithBasicAuth(clientID, clientSecret).
				Expect().
				Status(http.StatusOK).
				JSON().Object()

			t.Logf("%#v\n", resObj.Raw())

			validationAccessToken(t, resObj.Value("access_token").String().Raw())
		}
	}))
	defer csrv.Close()

	manager.MapClientStorage(clientStore(csrv.URL))
	srv = server.NewDefaultServer(manager)
	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		userID = "000000"
		return
	})

	e.GET("/authorize").
		WithQuery("response_type", "code").
		WithQuery("client_id", clientID).
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

	csrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer csrv.Close()

	manager.MapClientStorage(clientStore(csrv.URL))
	srv = server.NewDefaultServer(manager)
	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		userID = "000000"
		return
	})

	e.GET("/authorize").
		WithQuery("response_type", "token").
		WithQuery("client_id", clientID).
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

	manager.MapClientStorage(clientStore(""))
	srv = server.NewDefaultServer(manager)
	srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		if username == "admin" && password == "123456" {
			userID = "000000"
			return
		}
		err = fmt.Errorf("user not found")
		return
	})

	resObj := e.POST("/token").
		WithFormField("grant_type", "password").
		WithFormField("username", "admin").
		WithFormField("password", "123456").
		WithFormField("scope", "all").
		WithBasicAuth(clientID, clientSecret).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	t.Logf("%#v\n", resObj.Raw())

	validationAccessToken(t, resObj.Value("access_token").String().Raw())
}

func TestClientCredentials(t *testing.T) {
	tsrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testServer(t, w, r)
	}))
	defer tsrv.Close()
	e := httpexpect.New(t, tsrv.URL)

	manager.MapClientStorage(clientStore(""))

	srv = server.NewDefaultServer(manager)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		t.Log("OAuth 2.0 Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		t.Log("Response Error:", re.Error)
	})

	srv.SetAllowedGrantType(oauth2.ClientCredentials)
	srv.SetAllowGetAccessRequest(false)
	srv.SetExtensionFieldsHandler(func(ti oauth2.TokenInfo) (fieldsValue map[string]interface{}) {
		fieldsValue = map[string]interface{}{
			"extension": "param",
		}
		return
	})
	srv.SetAuthorizeScopeHandler(func(w http.ResponseWriter, r *http.Request) (scope string, err error) {
		return
	})
	srv.SetClientScopeHandler(func(clientID, scope string) (allowed bool, err error) {
		allowed = true
		return
	})

	resObj := e.POST("/token").
		WithFormField("grant_type", "client_credentials").
		WithFormField("scope", "all").
		WithFormField("client_id", clientID).
		WithFormField("client_secret", clientSecret).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	t.Logf("%#v\n", resObj.Raw())

	validationAccessToken(t, resObj.Value("access_token").String().Raw())
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
			jresObj := e.POST("/token").
				WithFormField("redirect_uri", csrv.URL+"/oauth2").
				WithFormField("code", code).
				WithFormField("grant_type", "authorization_code").
				WithFormField("client_id", clientID).
				WithBasicAuth(clientID, clientSecret).
				Expect().
				Status(http.StatusOK).
				JSON().Object()

			t.Logf("%#v\n", jresObj.Raw())

			validationAccessToken(t, jresObj.Value("access_token").String().Raw())

			resObj := e.POST("/token").
				WithFormField("grant_type", "refresh_token").
				WithFormField("scope", "one").
				WithFormField("refresh_token", jresObj.Value("refresh_token").String().Raw()).
				WithBasicAuth(clientID, clientSecret).
				Expect().
				Status(http.StatusOK).
				JSON().Object()

			t.Logf("%#v\n", resObj.Raw())

			validationAccessToken(t, resObj.Value("access_token").String().Raw())
		}
	}))
	defer csrv.Close()

	manager.MapClientStorage(clientStore(csrv.URL))
	srv = server.NewDefaultServer(manager)
	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		userID = "000000"
		return
	})

	e.GET("/authorize").
		WithQuery("response_type", "code").
		WithQuery("client_id", clientID).
		WithQuery("scope", "all").
		WithQuery("state", "123").
		WithQuery("redirect_uri", url.QueryEscape(csrv.URL+"/oauth2")).
		Expect().Status(http.StatusOK)
}

// validation access token
func validationAccessToken(t *testing.T, accessToken string) {
	req := httptest.NewRequest("GET", "http://example.com", nil)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	ti, err := srv.ValidationBearerToken(req)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if ti.GetClientID() != clientID {
		t.Error("invalid access token")
	}
}
