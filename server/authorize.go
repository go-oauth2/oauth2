package server

import (
	"encoding/base64"
	"net/http"
	"strings"

	"gopkg.in/oauth2.v2"
)

// AuthorizeRequest 授权请求
type AuthorizeRequest struct {
	Type        oauth2.ResponseType
	ClientID    string
	Scope       string
	RedirectURI string
	State       string
	UserID      string
}

// ClientHandler 获取客户端信息
type ClientHandler func(r *http.Request) (clientID, clientSecret string, err error)

// UserHandler 获取用户信息
type UserHandler func(username, password string) (userID string, err error)

// ClientFormHandler 客户端表单信息
func ClientFormHandler(r *http.Request) (clientID, clientSecret string, err error) {
	clientID = r.Form.Get("client_id")
	clientSecret = r.Form.Get("client_secret")
	return
}

// ClientBasicHandler 客户端基础认证信息
func ClientBasicHandler(r *http.Request) (clientID, clientSecret string, err error) {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		err = ErrAuthorizationHeaderInvalid
		return
	}
	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return
	}
	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		err = ErrAuthorizationHeaderInvalid
		return
	}
	clientID = pair[0]
	clientSecret = pair[1]
	return
}
