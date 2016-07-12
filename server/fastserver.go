package server

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/valyala/fasthttp"
	"gopkg.in/oauth2.v2"
)

// NewFastServer 创建基于fasthttp的OAuth2服务实例
func NewFastServer(cfg *Config, manager oauth2.Manager) *FastServer {
	srv := &FastServer{}
	srv.cfg = cfg
	srv.manager = manager
	srv.SetClientHandler(ClientFormFastHandler)
	return srv
}

// FastServer 基于fasthttp(https://github.com/valyala/fasthttp)的OAuth2服务处理
type FastServer struct {
	Server
}

// SetClientHandler 设置客户端处理
func (fs *FastServer) SetClientHandler(handler ClientFastHandler) {
	fs.cfg.Handler.ClientFastHandler = handler
}

// GetAuthorizeRequest 获取授权请求参数
func (fs *FastServer) GetAuthorizeRequest(ctx *fasthttp.RequestCtx) (authReq *AuthorizeRequest, err error) {
	if !ctx.IsGet() {
		err = ErrRequestMethodInvalid
		return
	}
	redirectURI, err := url.QueryUnescape(string(ctx.FormValue("redirect_uri")))
	if err != nil {
		return
	}
	authReq = &AuthorizeRequest{
		Type:        oauth2.ResponseType(string(ctx.FormValue("response_type"))),
		RedirectURI: redirectURI,
		State:       string(ctx.FormValue("state")),
		Scope:       string(ctx.FormValue("scope")),
		ClientID:    string(ctx.FormValue("client_id")),
	}
	if authReq.Type == "" || !fs.checkResponseType(authReq.Type) {
		err = ErrResponseTypeInvalid
	} else if authReq.ClientID == "" {
		err = ErrClientInvalid
	}
	return
}

// HandleAuthorizeRequest 处理授权请求
func (fs *FastServer) HandleAuthorizeRequest(ctx *fasthttp.RequestCtx, authReq *AuthorizeRequest) (err error) {
	if authReq.UserID == "" {
		err = ErrUserInvalid
		return
	}
	tgr := &oauth2.TokenGenerateRequest{
		ClientID:    authReq.ClientID,
		UserID:      authReq.UserID,
		RedirectURI: authReq.RedirectURI,
		Scope:       authReq.Scope,
	}
	ti, terr := fs.manager.GenerateAuthToken(oauth2.Code, tgr)
	if terr != nil {
		err = terr
		return
	}
	redirectURI, err := fs.GetRedirectURI(authReq, ti)
	if err != nil {
		return
	}
	ctx.Redirect(redirectURI, 302)
	return
}

// HandleTokenRequest 处理令牌请求
func (fs *FastServer) HandleTokenRequest(ctx *fasthttp.RequestCtx) (err error) {
	if !ctx.IsPost() {
		err = ErrRequestMethodInvalid
		return
	}
	gt := oauth2.GrantType(string(ctx.FormValue("grant_type")))
	if gt == "" || !fs.checkGrantType(gt) {
		err = ErrGrantTypeInvalid
		return
	}

	var ti oauth2.TokenInfo
	clientID, clientSecret, err := fs.cfg.Handler.ClientFastHandler(ctx)
	if err != nil {
		return
	}
	tgr := &oauth2.TokenGenerateRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	switch gt {
	case oauth2.AuthorizationCodeCredentials:
		tgr.RedirectURI = string(ctx.FormValue("redirect_uri"))
		tgr.Code = string(ctx.FormValue("code"))
		tgr.IsGenerateRefresh = true
		ti, err = fs.manager.GenerateAccessToken(oauth2.AuthorizationCodeCredentials, tgr)
	case oauth2.PasswordCredentials:
		userID, uerr := fs.cfg.Handler.UserHandler(string(ctx.FormValue("username")), string(ctx.FormValue("password")))
		if uerr != nil {
			err = uerr
			return
		}
		tgr.UserID = userID
		tgr.Scope = string(ctx.FormValue("scope"))
		tgr.IsGenerateRefresh = true
		ti, err = fs.manager.GenerateAccessToken(oauth2.PasswordCredentials, tgr)
	case oauth2.ClientCredentials:
		tgr.Scope = string(ctx.FormValue("scope"))
		ti, err = fs.manager.GenerateAccessToken(oauth2.ClientCredentials, tgr)
	case oauth2.RefreshCredentials:
		tgr.Refresh = string(ctx.FormValue("refresh_token"))
		tgr.Scope = string(ctx.FormValue("scope"))
		if tgr.Scope != "" { // 检查授权范围
			rti, rerr := fs.manager.LoadRefreshToken(tgr.Refresh)
			if rerr != nil {
				err = rerr
				return
			} else if rti.GetClientID() != tgr.ClientID {
				err = ErrRefreshInvalid
				return
			} else if verr := fs.cfg.Handler.ScopeHandler(tgr.Scope, rti.GetScope()); verr != nil {
				err = verr
				return
			}
		}
		ti, err = fs.manager.RefreshAccessToken(tgr)
		if err == nil {
			ti.SetRefresh("")
		}
	}

	if err != nil {
		return
	}
	err = fs.ResJSON(ctx, ti)
	return
}

// ResJSON 响应Json数据
func (fs *FastServer) ResJSON(ctx *fasthttp.RequestCtx, ti oauth2.TokenInfo) (err error) {
	data := map[string]interface{}{
		"access_token": ti.GetAccess(),
		"token_type":   fs.cfg.TokenType,
		"expires_in":   ti.GetAccessExpiresIn() / time.Second,
	}
	if scope := ti.GetScope(); scope != "" {
		data["scope"] = scope
	}
	if refresh := ti.GetRefresh(); refresh != "" {
		data["refresh_token"] = refresh
	}
	buf, err := json.Marshal(data)
	if err != nil {
		return
	}
	ctx.Response.Header.Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	ctx.Response.Header.Set("Pragma", "no-cache")
	ctx.Response.Header.Set("Expires", "Fri, 01 Jan 1990 00:00:00 GMT")
	ctx.SetContentType("application/json;charset=UTF-8")
	ctx.SetStatusCode(200)
	_, err = ctx.Write(buf)
	return nil
}
