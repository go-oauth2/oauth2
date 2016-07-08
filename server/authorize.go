package server

import (
	"net/http"

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

// ClientHandler 客户端处理(获取请求的客户端认证信息)
type ClientHandler func(r *http.Request) (clientID, clientSecret string, err error)

// UserHandler 用户处理(密码模式,根据用户名、密码获取用户标识)
type UserHandler func(username, password string) (userID string, err error)

// ScopeHandler 授权范围处理(更新令牌时的授权范围检查)
type ScopeHandler func(new, old string) (err error)

// TokenRequestHandler 令牌请求处理
type TokenRequestHandler struct {
	ClientHandler ClientHandler
	UserHandler   UserHandler
	ScopeHandler  ScopeHandler
}

// ClientFormHandler 客户端表单信息
func ClientFormHandler(r *http.Request) (clientID, clientSecret string, err error) {
	clientID = r.Form.Get("client_id")
	clientSecret = r.Form.Get("client_secret")
	if clientID == "" || clientSecret == "" {
		err = ErrAuthorizationFormInvalid
	}
	return
}

// ClientBasicHandler 客户端基础认证信息
func ClientBasicHandler(r *http.Request) (clientID, clientSecret string, err error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		err = ErrAuthorizationHeaderInvalid
		return
	}
	clientID = username
	clientSecret = password
	return
}
