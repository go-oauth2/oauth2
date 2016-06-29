package models

import "time"

// Token 令牌信息
type Token struct {
	ClientID         string        // 客户端标识
	UserID           string        // 用户标识
	RedirectURI      string        // 重定向URI
	Scope            string        // 权限范围
	AuthType         string        // 令牌授权类型
	Access           string        // 访问令牌
	AccessCreateAt   time.Time     // 访问令牌创建时间
	AccessExpiresIn  time.Duration // 访问令牌有效期
	Refresh          string        // 更新令牌
	RefreshCreateAt  time.Time     // 更新令牌创建时间
	RefreshExpiresIn time.Duration // 更新令牌有效期
}

// GetClientID 客户端ID
func (t *Token) GetClientID() string {
	return t.ClientID
}

// SetClientID 设置客户端ID
func (t *Token) SetClientID(clientID string) {
	t.ClientID = clientID
}

// GetUserID 用户ID
func (t *Token) GetUserID() string {
	return t.UserID
}

// SetUserID 设置用户ID
func (t *Token) SetUserID(userID string) {
	t.UserID = userID
}

// GetRedirectURI 重定向URI
func (t *Token) GetRedirectURI() string {
	return t.RedirectURI
}

// SetRedirectURI 设置重定向URI
func (t *Token) SetRedirectURI(redirectURI string) {
	t.RedirectURI = redirectURI
}

// GetScope 权限范围
func (t *Token) GetScope() string {
	return t.Scope
}

// SetScope 设置权限范围
func (t *Token) SetScope(scope string) {
	t.Scope = scope
}

// GetAuthType 授权类型
func (t *Token) GetAuthType() string {
	return t.AuthType
}

// SetAuthType 设置授权类型
func (t *Token) SetAuthType(authType string) {
	t.AuthType = authType
}

// GetAccess 访问令牌
func (t *Token) GetAccess() string {
	return t.Access
}

// SetAccess 设置访问令牌
func (t *Token) SetAccess(access string) {
	t.Access = access
}

// GetAccessCreateAt 访问令牌创建时间
func (t *Token) GetAccessCreateAt() time.Time {
	return t.AccessCreateAt
}

// SetAccessCreateAt 设置访问令牌创建时间
func (t *Token) SetAccessCreateAt(createAt time.Time) {
	t.AccessCreateAt = createAt
}

// GetAccessExpiresIn 访问令牌有效期
func (t *Token) GetAccessExpiresIn() time.Duration {
	return t.AccessExpiresIn
}

// SetAccessExpiresIn 设置访问令牌有效期
func (t *Token) SetAccessExpiresIn(exp time.Duration) {
	t.AccessExpiresIn = exp
}

// GetRefresh 更新令牌
func (t *Token) GetRefresh() string {
	return t.Refresh
}

// SetRefresh 设置更新令牌
func (t *Token) SetRefresh(refresh string) {
	t.Refresh = refresh
}

// GetRefreshCreateAt 更新令牌创建时间
func (t *Token) GetRefreshCreateAt() time.Time {
	return t.RefreshCreateAt
}

// SetRefreshCreateAt 设置更新令牌创建时间
func (t *Token) SetRefreshCreateAt(createAt time.Time) {
	t.RefreshCreateAt = createAt
}

// GetRefreshExpiresIn 更新令牌有效期
func (t *Token) GetRefreshExpiresIn() time.Duration {
	return t.RefreshExpiresIn
}

// SetRefreshExpiresIn 设置更新令牌有效期
func (t *Token) SetRefreshExpiresIn(exp time.Duration) {
	t.RefreshExpiresIn = exp
}
