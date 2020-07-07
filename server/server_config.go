package server

import (
	"github.com/go-oauth2/oauth2/v4"
)

// SetTokenType token type
// 设置Token类型
func (s *Server) SetTokenType(tokenType string) {
	s.Config.TokenType = tokenType
}

// SetAllowGetAccessRequest to allow GET requests for the token
// 设置是否允许对令牌的Get请求
func (s *Server) SetAllowGetAccessRequest(allow bool) {
	s.Config.AllowGetAccessRequest = allow
}

// SetAllowedResponseType allow the authorization types
// 设置允许的响应类型
func (s *Server) SetAllowedResponseType(types ...oauth2.ResponseType) {
	s.Config.AllowedResponseTypes = types
}

// SetAllowedGrantType allow the grant types
// 设置允许的授权类型
func (s *Server) SetAllowedGrantType(types ...oauth2.GrantType) {
	s.Config.AllowedGrantTypes = types
}

// SetClientInfoHandler get client info from request
// 设置客户端信息
func (s *Server) SetClientInfoHandler(handler ClientInfoHandler) {
	s.ClientInfoHandler = handler
}

// SetClientAuthorizedHandler check the client allows to use this authorization grant type
// 设置客户端授权类型
func (s *Server) SetClientAuthorizedHandler(handler ClientAuthorizedHandler) {
	s.ClientAuthorizedHandler = handler
}

// SetClientScopeHandler check the client allows to use scope
// 设置客户端使用范围
func (s *Server) SetClientScopeHandler(handler ClientScopeHandler) {
	s.ClientScopeHandler = handler
}

// SetUserAuthorizationHandler get user id from request authorization
// 设置用户授权
func (s *Server) SetUserAuthorizationHandler(handler UserAuthorizationHandler) {
	s.UserAuthorizationHandler = handler
}

// SetPasswordAuthorizationHandler get user id from username and password
// 设置用户名密码授权
func (s *Server) SetPasswordAuthorizationHandler(handler PasswordAuthorizationHandler) {
	s.PasswordAuthorizationHandler = handler
}

// SetRefreshingScopeHandler check the scope of the refreshing token
// 设置刷新令牌的范围
func (s *Server) SetRefreshingScopeHandler(handler RefreshingScopeHandler) {
	s.RefreshingScopeHandler = handler
}

// SetResponseErrorHandler response error handling
// 设置响应错误方法
func (s *Server) SetResponseErrorHandler(handler ResponseErrorHandler) {
	s.ResponseErrorHandler = handler
}

// SetInternalErrorHandler internal error handling
// 设置内部错误方法
func (s *Server) SetInternalErrorHandler(handler InternalErrorHandler) {
	s.InternalErrorHandler = handler
}

// SetExtensionFieldsHandler in response to the access token with the extension of the field
// 设置扩展字段方法
func (s *Server) SetExtensionFieldsHandler(handler ExtensionFieldsHandler) {
	s.ExtensionFieldsHandler = handler
}

// SetAccessTokenExpHandler set expiration date for the access token
// 设置访问令牌
func (s *Server) SetAccessTokenExpHandler(handler AccessTokenExpHandler) {
	s.AccessTokenExpHandler = handler
}

// SetAuthorizeScopeHandler set scope for the access token
// 设置授权作用域
func (s *Server) SetAuthorizeScopeHandler(handler AuthorizeScopeHandler) {
	s.AuthorizeScopeHandler = handler
}
