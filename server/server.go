package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gopkg.in/oauth2.v2"
)

// NewServer 创建OAuth2服务实例
func NewServer(cfg *Config, manager oauth2.Manager) *Server {
	srv := &Server{
		cfg:     cfg,
		manager: manager,
	}
	srv.SetClientHandler(ClientFormHandler)
	return srv
}

// Server OAuth2服务处理
type Server struct {
	cfg     *Config
	manager oauth2.Manager
}

// SetTokenType 设置令牌类型
func (s *Server) SetTokenType(tokenType string) {
	s.cfg.TokenType = tokenType
}

// SetAllowedResponseType 设置允许的授权类型
func (s *Server) SetAllowedResponseType(allowedTypes ...oauth2.ResponseType) {
	s.cfg.AllowedResponseType = allowedTypes
}

// SetAllowedGrantType 允许的授权模式
func (s *Server) SetAllowedGrantType(allowedTypes ...oauth2.GrantType) {
	s.cfg.AllowedGrantType = allowedTypes
}

// SetClientHandler 设置客户端处理
func (s *Server) SetClientHandler(handler ClientHandler) {
	s.cfg.Handler.ClientHandler = handler
}

// SetUserHandler 设置用户处理
func (s *Server) SetUserHandler(handler UserHandler) {
	s.cfg.Handler.UserHandler = handler
}

// SetScopeHandler 设置授权范围处理
func (s *Server) SetScopeHandler(handler ScopeHandler) {
	s.cfg.Handler.ScopeHandler = handler
}

// checkResponseType 检查允许的授权类型
func (s *Server) checkResponseType(rt oauth2.ResponseType) bool {
	for _, art := range s.cfg.AllowedResponseType {
		if art == rt {
			return true
		}
	}
	return false
}

// checkGrantType 检查允许的授权模式
func (s *Server) checkGrantType(gt oauth2.GrantType) bool {
	for _, agt := range s.cfg.AllowedGrantType {
		if agt == gt {
			return true
		}
	}
	return false
}

// GetAuthorizeRequest 获取授权请求参数
func (s *Server) GetAuthorizeRequest(r *http.Request) (authReq *AuthorizeRequest, err error) {
	if r.Method != "GET" {
		err = ErrRequestMethodInvalid
		return
	}
	r.ParseForm()
	redirectURI, err := url.QueryUnescape(r.Form.Get("redirect_uri"))
	if err != nil {
		return
	}
	authReq = &AuthorizeRequest{
		Type:        oauth2.ResponseType(r.Form.Get("response_type")),
		RedirectURI: redirectURI,
		State:       r.Form.Get("state"),
		Scope:       r.Form.Get("scope"),
		ClientID:    r.Form.Get("client_id"),
	}
	if authReq.Type == "" || !s.checkResponseType(authReq.Type) {
		err = ErrResponseTypeInvalid
		return
	} else if authReq.ClientID == "" {
		err = ErrClientInvalid
	}
	return
}

// HandleAuthorizeRequest 处理授权请求
func (s *Server) HandleAuthorizeRequest(w http.ResponseWriter, authReq *AuthorizeRequest) (err error) {
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
	ti, err := s.manager.GenerateAuthToken(oauth2.Code, tgr)
	if err != nil {
		return
	}
	redirectURI, err := s.GetRedirectURI(authReq, ti)
	if err != nil {
		return
	}
	w.Header().Set("Location", redirectURI)
	w.WriteHeader(302)
	return
}

// GetRedirectURI 获取重定向URI
func (s *Server) GetRedirectURI(authReq *AuthorizeRequest, ti oauth2.TokenInfo) (uri string, err error) {
	u, err := url.Parse(authReq.RedirectURI)
	if err != nil {
		return
	}
	q := u.Query()
	q.Set("state", authReq.State)
	switch authReq.Type {
	case oauth2.Code:
		q.Set("code", ti.GetAccess())
		u.RawQuery = q.Encode()
	case oauth2.Token:
		q.Set("access_token", ti.GetAccess())
		q.Set("token_type", s.cfg.TokenType)
		q.Set("expires_in", strconv.FormatInt(int64(ti.GetAccessExpiresIn()/time.Second), 10))
		q.Set("scope", ti.GetScope())
		u.RawQuery = ""
		u.Fragment, err = url.QueryUnescape(q.Encode())
		if err != nil {
			return
		}
	}
	uri = u.String()
	return
}

// HandleTokenRequest 处理令牌请求
func (s *Server) HandleTokenRequest(w http.ResponseWriter, r *http.Request) (err error) {
	if r.Method != "POST" {
		err = ErrRequestMethodInvalid
		return
	}
	r.ParseForm()
	gt := oauth2.GrantType(r.Form.Get("grant_type"))
	if gt == "" || !s.checkGrantType(gt) {
		err = ErrGrantTypeInvalid
		return
	}
	var ti oauth2.TokenInfo
	clientID, clientSecret, err := s.cfg.Handler.ClientHandler(r)
	if err != nil {
		return
	}
	tgr := &oauth2.TokenGenerateRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	switch gt {
	case oauth2.AuthorizationCode:
		tgr.RedirectURI = r.Form.Get("redirect_uri")
		tgr.Code = r.Form.Get("code")
		tgr.IsGenerateRefresh = true
		ti, err = s.manager.GenerateAccessToken(oauth2.AuthorizationCode, tgr)
	case oauth2.PasswordCredentials:
		userID, uerr := s.cfg.Handler.UserHandler(r.Form.Get("username"), r.Form.Get("password"))
		if uerr != nil {
			err = uerr
			return
		}
		tgr.UserID = userID
		tgr.Scope = r.Form.Get("scope")
		tgr.IsGenerateRefresh = true
		ti, err = s.manager.GenerateAccessToken(oauth2.PasswordCredentials, tgr)
	case oauth2.ClientCredentials:
		tgr.Scope = r.Form.Get("scope")
		ti, err = s.manager.GenerateAccessToken(oauth2.ClientCredentials, tgr)
	case oauth2.Refreshing:
		tgr.Refresh = r.Form.Get("refresh_token")
		tgr.Scope = r.Form.Get("scope")
		if tgr.Scope != "" { // 检查授权范围
			rti, rerr := s.manager.LoadRefreshToken(tgr.Refresh)
			if rerr != nil {
				err = rerr
				return
			} else if rti.GetClientID() != tgr.ClientID {
				err = ErrRefreshInvalid
				return
			} else if verr := s.cfg.Handler.ScopeHandler(tgr.Scope, rti.GetScope()); verr != nil {
				err = verr
				return
			}
		}
		ti, err = s.manager.RefreshAccessToken(tgr)
		if err == nil {
			ti.SetRefresh("")
		}
	}

	if err != nil {
		return
	}
	err = s.ResJSON(w, ti)
	return
}

// ResJSON 响应Json数据
func (s *Server) ResJSON(w http.ResponseWriter, ti oauth2.TokenInfo) (err error) {
	data := map[string]interface{}{
		"access_token": ti.GetAccess(),
		"token_type":   s.cfg.TokenType,
		"expires_in":   ti.GetAccessExpiresIn() / time.Second,
	}
	if scope := ti.GetScope(); scope != "" {
		data["scope"] = scope
	}
	if refresh := ti.GetRefresh(); refresh != "" {
		data["refresh_token"] = refresh
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "Fri, 01 Jan 1990 00:00:00 GMT")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}
