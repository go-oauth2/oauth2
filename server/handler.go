package server

import (
	"net/http"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
)

type (
	// ClientInfoHandler get client info from request
	// 从请求中获取客户端信息
	ClientInfoHandler func(r *http.Request) (clientID, clientSecret string, err error)

	// ClientAuthorizedHandler check the client allows to use this authorization grant type
	// 检查客户端允许使用此授权授予类型
	ClientAuthorizedHandler func(clientID string, grant oauth2.GrantType) (allowed bool, err error)

	// ClientScopeHandler check the client allows to use scope
	// 检查客户端允许使用范围
	ClientScopeHandler func(clientID, scope string) (allowed bool, err error)

	// UserAuthorizationHandler get user id from request authorization
	// 从请求授权中获取用户ID
	UserAuthorizationHandler func(w http.ResponseWriter, r *http.Request) (userID string, err error)

	// PasswordAuthorizationHandler get user id from username and password
	// 从用户名和密码获取用户ID
	PasswordAuthorizationHandler func(username, password string) (userID string, err error)

	// RefreshingScopeHandler check the scope of the refreshing token
	// 检查刷新令牌的范围
	RefreshingScopeHandler func(newScope, oldScope string) (allowed bool, err error)

	// ResponseErrorHandler response error handing
	// 响应错误处理
	ResponseErrorHandler func(re *errors.Response)

	// InternalErrorHandler internal error handing
	// 内部错误处理程序内部错误处理
	InternalErrorHandler func(err error) (re *errors.Response)

	// AuthorizeScopeHandler set the authorized scope
	// 设置授权范围
	AuthorizeScopeHandler func(w http.ResponseWriter, r *http.Request) (scope string, err error)

	// AccessTokenExpHandler set expiration date for the access token
	// 设置访问令牌的到期日期
	AccessTokenExpHandler func(w http.ResponseWriter, r *http.Request) (exp time.Duration, err error)

	// ExtensionFieldsHandler in response to the access token with the extension of the field
	// 响应带有字段扩展名的访问令牌
	ExtensionFieldsHandler func(ti oauth2.TokenInfo) (fieldsValue map[string]interface{})
)

// ClientFormHandler get client data from form
// 从表单获取客户端数据
func ClientFormHandler(r *http.Request) (string, string, error) {
	clientID := r.Form.Get("client_id")
	if clientID == "" {
		return "", "", errors.ErrInvalidClient
	}
	clientSecret := r.Form.Get("client_secret")
	return clientID, clientSecret, nil
}

// ClientBasicHandler get client data from basic authorization
// 从表单获取客户端数据
func ClientBasicHandler(r *http.Request) (string, string, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return "", "", errors.ErrInvalidClient
	}
	return username, password, nil
}
