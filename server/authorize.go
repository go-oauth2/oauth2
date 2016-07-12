package server

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
	"gopkg.in/oauth2.v2"
)

// AuthorizeRequest 授权请求
type AuthorizeRequest struct {
	Type        oauth2.ResponseType // 授权类型
	ClientID    string              // 客户端标识
	Scope       string              // 授权范围
	RedirectURI string              // 重定向URI
	State       string              // 状态
	UserID      string              // 用户标识
}

// TokenRequestHandler 令牌请求处理
type TokenRequestHandler struct {
	// 客户端信息处理
	ClientHandler ClientHandler
	// 客户端信息处理(基于fasthttp)
	ClientFastHandler ClientFastHandler
	// 用户信息处理
	UserHandler UserHandler
	// 授权范围处理
	ScopeHandler ScopeHandler
}

// ClientHandler 获取请求的客户端认证信息
type ClientHandler func(r *http.Request) (clientID, clientSecret string, err error)

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

// ClientFastHandler 基于fasthttp获取客户端认证信息
type ClientFastHandler func(ctx *fasthttp.RequestCtx) (clientID, clientSecret string, err error)

// ClientFormFastHandler 客户端表单信息(基于fasthttp)
func ClientFormFastHandler(ctx *fasthttp.RequestCtx) (clientID, clientSecret string, err error) {
	clientID = string(ctx.FormValue("client_id"))
	clientSecret = string(ctx.FormValue("client_secret"))
	if clientID == "" || clientSecret == "" {
		err = ErrAuthorizationFormInvalid
	}
	return
}

// ClientBasicFastHandler 客户端基础认证信息(基于fasthttp)
func ClientBasicFastHandler(ctx *fasthttp.RequestCtx) (clientID, clientSecret string, err error) {
	auth := string(ctx.Request.Header.Peek("Authorization"))
	const prefix = "Basic "
	if auth == "" || !strings.HasPrefix(auth, prefix) {
		err = ErrAuthorizationHeaderInvalid
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		err = ErrAuthorizationHeaderInvalid
		return
	}
	clientID = cs[:s]
	clientSecret = cs[s+1:]
	return
}

// UserHandler 密码模式下,根据用户名、密码获取用户标识
type UserHandler func(username, password string) (userID string, err error)

// ScopeHandler 更新令牌时的授权范围检查
type ScopeHandler func(new, old string) (err error)
