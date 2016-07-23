package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
)

// NewServer Create to authorization server instance
func NewServer(cfg *Config, manager oauth2.Manager) *Server {
	srv := &Server{
		Config:            cfg,
		Manager:           manager,
		ClientInfoHandler: ClientFormHandler,
		ResponseErrorHandler: func(re *errors.Response) {
			re.Description = ""
		},
	}
	return srv
}

// Server Provide authorization server
type Server struct {
	Config                       *Config
	Manager                      oauth2.Manager
	ClientInfoHandler            ClientInfoHandler
	ClientAuthorizedHandler      ClientAuthorizedHandler
	ClientScopeHandler           ClientScopeHandler
	UserAuthorizationHandler     UserAuthorizationHandler
	PasswordAuthorizationHandler PasswordAuthorizationHandler
	RefreshingScopeHandler       RefreshingScopeHandler
	ResponseErrorHandler         ResponseErrorHandler
	InternalErrorHandler         InternalErrorHandler
}

// SetAllowedResponseType Allow the authorization types
func (s *Server) SetAllowedResponseType(types ...oauth2.ResponseType) {
	s.Config.AllowedResponseTypes = types
}

// SetAllowedGrantType Allow the grant types
func (s *Server) SetAllowedGrantType(types ...oauth2.GrantType) {
	s.Config.AllowedGrantTypes = types
}

// SetClientInfoHandler Get client info from request
func (s *Server) SetClientInfoHandler(handler ClientInfoHandler) {
	s.ClientInfoHandler = handler
}

// SetClientAuthorizedHandler Check the client allows to use this authorization grant type
func (s *Server) SetClientAuthorizedHandler(handler ClientAuthorizedHandler) {
	s.ClientAuthorizedHandler = handler
}

// SetClientScopeHandler Check the client allows to use scope
func (s *Server) SetClientScopeHandler(handler ClientScopeHandler) {
	s.ClientScopeHandler = handler
}

// SetUserAuthorizationHandler Get user id from request authorization
func (s *Server) SetUserAuthorizationHandler(handler UserAuthorizationHandler) {
	s.UserAuthorizationHandler = handler
}

// SetPasswordAuthorizationHandler Get user id from username and password
func (s *Server) SetPasswordAuthorizationHandler(handler PasswordAuthorizationHandler) {
	s.PasswordAuthorizationHandler = handler
}

// SetRefreshingScopeHandler Check the scope of the refreshing token
func (s *Server) SetRefreshingScopeHandler(handler RefreshingScopeHandler) {
	s.RefreshingScopeHandler = handler
}

// SetResponseErrorHandler Response error handling
func (s *Server) SetResponseErrorHandler(handler ResponseErrorHandler) {
	s.ResponseErrorHandler = handler
}

// SetInternalErrorHandler Internal error handling
func (s *Server) SetInternalErrorHandler(handler InternalErrorHandler) {
	s.InternalErrorHandler = handler
}

// CheckResponseType Check allows response type
func (s *Server) CheckResponseType(rt oauth2.ResponseType) bool {
	for _, art := range s.Config.AllowedResponseTypes {
		if art == rt {
			return true
		}
	}
	return false
}

// CheckGrantType Check allows grant type
func (s *Server) CheckGrantType(gt oauth2.GrantType) bool {
	for _, agt := range s.Config.AllowedGrantTypes {
		if agt == gt {
			return true
		}
	}
	return false
}

// ValidationAuthorizeRequest The authorization request validation
func (s *Server) ValidationAuthorizeRequest(r *http.Request) (req *AuthorizeRequest, rerr, ierr error) {
	redirectURI, err := url.QueryUnescape(r.FormValue("redirect_uri"))
	if err != nil {
		ierr = err
		return
	}
	req = &AuthorizeRequest{
		RedirectURI:  redirectURI,
		ResponseType: oauth2.ResponseType(r.FormValue("response_type")),
		ClientID:     r.FormValue("client_id"),
		State:        r.FormValue("state"),
		Scope:        r.FormValue("scope"),
	}
	if r.Method != "GET" {
		rerr = errors.ErrInvalidRequest
	}
	return
}

// GetAuthorizeToken Get authorization token(code)
func (s *Server) GetAuthorizeToken(req *AuthorizeRequest) (ti oauth2.TokenInfo, rerr, ierr error) {
	if req.RedirectURI == "" ||
		req.ClientID == "" {
		rerr = errors.ErrInvalidRequest
		return
	} else if req.ResponseType == "" {
		rerr = errors.ErrUnsupportedResponseType
		return
	}
	if allowed := s.CheckResponseType(req.ResponseType); !allowed {
		rerr = errors.ErrUnauthorizedClient
		return
	}
	if fn := s.ClientAuthorizedHandler; fn != nil {
		gt := oauth2.AuthorizationCode
		if req.ResponseType == oauth2.Token {
			gt = oauth2.Implicit
		}
		allowed, err := fn(req.ClientID, gt)
		if err != nil {
			ierr = err
			return
		}
		if !allowed {
			rerr = errors.ErrUnauthorizedClient
			return
		}
	}
	if fn := s.ClientScopeHandler; fn != nil {
		allowed, err := fn(req.ClientID, req.Scope)
		if err != nil {
			ierr = err
			return
		}
		if !allowed {
			rerr = errors.ErrInvalidScope
			return
		}
	}
	tgr := &oauth2.TokenGenerateRequest{
		ClientID:    req.ClientID,
		UserID:      req.UserID,
		RedirectURI: req.RedirectURI,
		Scope:       req.Scope,
	}
	ti, err := s.Manager.GenerateAuthToken(req.ResponseType, tgr)
	if err != nil {
		if err == errors.ErrInvalidClient {
			rerr = err
		} else {
			ierr = err
		}
	}
	return
}

// GetRedirectURI Get redirect uri
func (s *Server) GetRedirectURI(req *AuthorizeRequest, data map[string]interface{}) (uri string, err error) {
	u, err := url.Parse(req.RedirectURI)
	if err != nil {
		return
	}
	q := u.Query()
	if req.State != "" {
		q.Set("state", req.State)
	}
	for k, v := range data {
		q.Set(k, fmt.Sprint(v))
	}
	switch req.ResponseType {
	case oauth2.Code:
		u.RawQuery = q.Encode()
	case oauth2.Token:
		u.RawQuery = ""
		u.Fragment, err = url.QueryUnescape(q.Encode())
		if err != nil {
			return
		}
	}
	uri = u.String()
	return
}

// GetAuthorizeData Get authorization response data
func (s *Server) GetAuthorizeData(rt oauth2.ResponseType, ti oauth2.TokenInfo) (data map[string]interface{}) {
	if rt == oauth2.Code {
		data = map[string]interface{}{
			"code": ti.GetAccess(),
		}
	} else {
		data = s.GetTokenData(ti)
	}
	return
}

// GetErrorData Get error response data
func (s *Server) GetErrorData(rerr, ierr error) (data map[string]interface{}, statusCode int) {
	if ierr != nil {
		rerr = errors.ErrServerError
	}
	re := &errors.Response{
		Error:       rerr,
		Description: errors.Descriptions[rerr],
	}
	if fn := s.ResponseErrorHandler; fn != nil {
		fn(re)
	}
	data = map[string]interface{}{
		"error": re.Error.Error(),
	}
	if v := re.Description; v != "" {
		data["error_description"] = v
	}
	if v := re.URI; v != "" {
		data["error_uri"] = v
	}
	statusCode = re.StatusCode
	return
}

// response redirect error
func (s *Server) resRedirectError(w http.ResponseWriter, r *http.Request, req *AuthorizeRequest, rerr, ierr error) (err error) {
	if req == nil {
		err = ierr
		return
	}
	if fn := s.InternalErrorHandler; fn != nil {
		verr := ierr
		if verr == nil {
			verr = rerr
		}
		fn(r, verr)
	}
	data, _ := s.GetErrorData(rerr, ierr)
	err = s.resRedirect(w, req, data)
	return
}

func (s *Server) resRedirect(w http.ResponseWriter, req *AuthorizeRequest, data map[string]interface{}) (err error) {
	uri, verr := s.GetRedirectURI(req, data)
	if verr != nil {
		err = verr
		return
	}
	w.Header().Set("Location", uri)
	w.WriteHeader(302)
	return
}

// HandleAuthorizeRequest The authorization request handling
func (s *Server) HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) (err error) {
	defer func() {
		if verr := recover(); verr != nil {
			err = fmt.Errorf("%v", verr)
		}
	}()
	req, rerr, ierr := s.ValidationAuthorizeRequest(r)
	if rerr != nil || ierr != nil {
		err = s.resRedirectError(w, r, req, rerr, ierr)
		return
	}
	userID, err := s.UserAuthorizationHandler(w, r)
	if err != nil {
		err = s.resRedirectError(w, r, req, nil, err)
		return
	} else if userID == "" {
		return
	}
	req.UserID = userID
	ti, rerr, ierr := s.GetAuthorizeToken(req)
	if rerr != nil || ierr != nil {
		err = s.resRedirectError(w, r, req, rerr, ierr)
		return
	}
	err = s.resRedirect(w, req, s.GetAuthorizeData(req.ResponseType, ti))
	return
}

// ValidationTokenRequest The token request validation
func (s *Server) ValidationTokenRequest(r *http.Request) (gt oauth2.GrantType, tgr *oauth2.TokenGenerateRequest, rerr, ierr error) {
	if r.Method != "POST" {
		rerr = errors.ErrInvalidRequest
		return
	}
	gt = oauth2.GrantType(r.FormValue("grant_type"))
	if gt == "" {
		rerr = errors.ErrUnsupportedGrantType
		return
	}
	clientID, clientSecret, err := s.ClientInfoHandler(r)
	if err != nil {
		ierr = err
		return
	}
	tgr = &oauth2.TokenGenerateRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	switch gt {
	case oauth2.AuthorizationCode:
		tgr.RedirectURI = r.FormValue("redirect_uri")
		tgr.Code = r.FormValue("code")
		if tgr.RedirectURI == "" ||
			tgr.Code == "" {
			rerr = errors.ErrInvalidRequest
		}
	case oauth2.PasswordCredentials:
		tgr.Scope = r.FormValue("scope")
		userID, verr := s.PasswordAuthorizationHandler(r.FormValue("username"), r.FormValue("password"))
		if verr != nil {
			ierr = verr
			return
		} else if userID == "" {
			rerr = errors.ErrInvalidRequest
			return
		}
		tgr.UserID = userID
	case oauth2.ClientCredentials:
		tgr.Scope = r.FormValue("scope")
	case oauth2.Refreshing:
		tgr.Refresh = r.FormValue("refresh_token")
		tgr.Scope = r.FormValue("scope")
		if tgr.Refresh == "" {
			rerr = errors.ErrInvalidRequest
		}
	}
	return
}

// GetAccessToken Get access token
func (s *Server) GetAccessToken(gt oauth2.GrantType, tgr *oauth2.TokenGenerateRequest) (ti oauth2.TokenInfo, rerr, ierr error) {
	if allowed := s.CheckGrantType(gt); !allowed {
		rerr = errors.ErrUnauthorizedClient
		return
	}
	if fn := s.ClientAuthorizedHandler; fn != nil {
		allowed, err := fn(tgr.ClientID, gt)
		if err != nil {
			ierr = err
			return
		}
		if !allowed {
			rerr = errors.ErrUnauthorizedClient
			return
		}
	}
	switch gt {
	case oauth2.AuthorizationCode:
		ti, ierr = s.Manager.GenerateAccessToken(gt, tgr)
		if ierr != nil {
			if ierr == errors.ErrInvalidAuthorizeCode {
				rerr = errors.ErrInvalidGrant
				ierr = nil
			} else if ierr == errors.ErrInvalidClient {
				rerr = errors.ErrInvalidClient
				ierr = nil
			}
		}
	case oauth2.PasswordCredentials, oauth2.ClientCredentials:
		if fn := s.ClientScopeHandler; fn != nil {
			allowed, err := fn(tgr.ClientID, tgr.Scope)
			if err != nil {
				ierr = err
				return
			}
			if !allowed {
				rerr = errors.ErrInvalidScope
				return
			}
		}
		ti, ierr = s.Manager.GenerateAccessToken(gt, tgr)
		if ierr != nil {
			if ierr == errors.ErrInvalidClient {
				rerr = errors.ErrInvalidClient
				ierr = nil
			}
		}
	case oauth2.Refreshing:
		if scope := tgr.Scope; scope != "" {
			rti, err := s.Manager.LoadRefreshToken(tgr.Refresh)
			if err != nil {
				if err == errors.ErrInvalidRefreshToken {
					rerr = err
					return
				}
				ierr = err
				return
			}
			if fn := s.RefreshingScopeHandler; fn != nil && !fn(scope, rti.GetScope()) {
				rerr = errors.ErrInvalidScope
				return
			}
		}
		ti, ierr = s.Manager.RefreshAccessToken(tgr)
		if ierr != nil {
			if ierr == errors.ErrInvalidClient {
				rerr = errors.ErrInvalidClient
				ierr = nil
			} else if ierr == errors.ErrInvalidRefreshToken {
				rerr = errors.ErrInvalidRefreshToken
				ierr = nil
			}
		} else {
			ti.SetRefresh("")
		}
	}

	return
}

// GetTokenData Get token data
func (s *Server) GetTokenData(ti oauth2.TokenInfo) (data map[string]interface{}) {
	data = map[string]interface{}{
		"access_token": ti.GetAccess(),
		"token_type":   s.Config.TokenType,
		"expires_in":   int64(ti.GetAccessExpiresIn() / time.Second),
	}
	if scope := ti.GetScope(); scope != "" {
		data["scope"] = scope
	}
	if refresh := ti.GetRefresh(); refresh != "" {
		data["refresh_token"] = refresh
	}
	return
}

// HandleTokenRequest The token request handling
func (s *Server) HandleTokenRequest(w http.ResponseWriter, r *http.Request) (err error) {
	defer func() {
		if verr := recover(); verr != nil {
			err = fmt.Errorf("%v", verr)
		}
	}()
	gt, tgr, rerr, ierr := s.ValidationTokenRequest(r)
	if rerr != nil || ierr != nil {
		err = s.resTokenError(w, r, rerr, ierr)
		return
	}
	ti, rerr, ierr := s.GetAccessToken(gt, tgr)
	if rerr != nil || ierr != nil {
		err = s.resTokenError(w, r, rerr, ierr)
		return
	}
	err = s.resToken(w, s.GetTokenData(ti))
	return
}

func (s *Server) resTokenError(w http.ResponseWriter, r *http.Request, rerr, ierr error) (err error) {
	if fn := s.InternalErrorHandler; fn != nil {
		verr := ierr
		if verr == nil {
			verr = rerr
		}
		fn(r, verr)
	}
	data, statusCode := s.GetErrorData(rerr, ierr)
	s.resToken(w, data, statusCode)
	return
}

func (s *Server) resToken(w http.ResponseWriter, data map[string]interface{}, statusCode ...int) (err error) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	status := http.StatusOK
	if len(statusCode) > 0 && statusCode[0] > 0 {
		status = statusCode[0]
	}
	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(data)
	return
}
