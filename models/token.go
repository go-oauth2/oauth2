package models

import (
	"time"

	"github.com/go-oauth2/oauth2/v4"
)

// NewToken create to token model instance
// 创建令牌模型实例
func NewToken() *Token {
	return &Token{}
}

// Token token model
// 令牌令牌模型
type Token struct {
	ClientID         string        `bson:"ClientID"`
	UserID           string        `bson:"UserID"`
	RedirectURI      string        `bson:"RedirectURI"`
	Scope            string        `bson:"Scope"`
	Code             string        `bson:"Code"`
	CodeCreateAt     time.Time     `bson:"CodeCreateAt"`
	CodeExpiresIn    time.Duration `bson:"CodeExpiresIn"`
	Access           string        `bson:"Access"`
	AccessCreateAt   time.Time     `bson:"AccessCreateAt"`
	AccessExpiresIn  time.Duration `bson:"AccessExpiresIn"`
	Refresh          string        `bson:"Refresh"`
	RefreshCreateAt  time.Time     `bson:"RefreshCreateAt"`
	RefreshExpiresIn time.Duration `bson:"RefreshExpiresIn"`
}

// New create to token model instance
// 新创建令牌模型实例
func (t *Token) New() oauth2.TokenInfo {
	return NewToken()
}

// GetClientID the client id
// 获取客户端ID
func (t *Token) GetClientID() string {
	return t.ClientID
}

// SetClientID the client id
// 设置客户端ID
func (t *Token) SetClientID(clientID string) {
	t.ClientID = clientID
}

// GetUserID the user id
// 获取用户ID
func (t *Token) GetUserID() string {
	return t.UserID
}

// SetUserID the user id
// 设置用户ID
func (t *Token) SetUserID(userID string) {
	t.UserID = userID
}

// GetRedirectURI redirect URI
// 获取重定向URI
func (t *Token) GetRedirectURI() string {
	return t.RedirectURI
}

// SetRedirectURI redirect URI
// 设置重定向URI
func (t *Token) SetRedirectURI(redirectURI string) {
	t.RedirectURI = redirectURI
}

// GetScope get scope of authorization
// 获取授权范围
func (t *Token) GetScope() string {
	return t.Scope
}

// SetScope get scope of authorization
// 设置授权范围
func (t *Token) SetScope(scope string) {
	t.Scope = scope
}

// GetCode authorization code
// 获取code
func (t *Token) GetCode() string {
	return t.Code
}

// SetCode authorization code
// 设置code
func (t *Token) SetCode(code string) {
	t.Code = code
}

// GetCodeCreateAt create Time
// 获取code创建时间
func (t *Token) GetCodeCreateAt() time.Time {
	return t.CodeCreateAt
}

// SetCodeCreateAt create Time
// 设置code创建时间
func (t *Token) SetCodeCreateAt(createAt time.Time) {
	t.CodeCreateAt = createAt
}

// GetCodeExpiresIn the lifetime in seconds of the authorization code
// 获取授权码有效期
func (t *Token) GetCodeExpiresIn() time.Duration {
	return t.CodeExpiresIn
}

// SetCodeExpiresIn the lifetime in seconds of the authorization code
// 设置授权码有效期
func (t *Token) SetCodeExpiresIn(exp time.Duration) {
	t.CodeExpiresIn = exp
}

// GetAccess access Token
// 获取访问令牌
func (t *Token) GetAccess() string {
	return t.Access
}

// SetAccess access Token
// 设置访问令牌
func (t *Token) SetAccess(access string) {
	t.Access = access
}

// GetAccessCreateAt create Time
// 获取令牌创建时间
func (t *Token) GetAccessCreateAt() time.Time {
	return t.AccessCreateAt
}

// SetAccessCreateAt create Time
// 设置令牌创建时间
func (t *Token) SetAccessCreateAt(createAt time.Time) {
	t.AccessCreateAt = createAt
}

// GetAccessExpiresIn the lifetime in seconds of the access token
// 获取令牌有效期
func (t *Token) GetAccessExpiresIn() time.Duration {
	return t.AccessExpiresIn
}

// SetAccessExpiresIn the lifetime in seconds of the access token
// 设置令牌有效期
func (t *Token) SetAccessExpiresIn(exp time.Duration) {
	t.AccessExpiresIn = exp
}

// GetRefresh refresh Token
// 获取刷新令牌
func (t *Token) GetRefresh() string {
	return t.Refresh
}

// SetRefresh refresh Token
// 设置刷新令牌
func (t *Token) SetRefresh(refresh string) {
	t.Refresh = refresh
}

// GetRefreshCreateAt create Time
// 获取刷新令牌创建时间
func (t *Token) GetRefreshCreateAt() time.Time {
	return t.RefreshCreateAt
}

// SetRefreshCreateAt create Time
// 设置刷新令牌创建时间
func (t *Token) SetRefreshCreateAt(createAt time.Time) {
	t.RefreshCreateAt = createAt
}

// GetRefreshExpiresIn the lifetime in seconds of the refresh token
// 获取刷新令牌有效期
func (t *Token) GetRefreshExpiresIn() time.Duration {
	return t.RefreshExpiresIn
}

// SetRefreshExpiresIn the lifetime in seconds of the refresh token
// 设置刷新令牌有效期
func (t *Token) SetRefreshExpiresIn(exp time.Duration) {
	t.RefreshExpiresIn = exp
}
