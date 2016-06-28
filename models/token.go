package models

import "time"

// Token 令牌信息
type Token struct {
	ID               int64         `bson:"_id"`              // 唯一标识
	ClientID         string        `bson:"ClientID"`         // 客户端标识
	UserID           string        `bson:"UserID"`           // 用户标识
	RedirectURI      string        `bson:"RedirectURI"`      // 重定向URI
	Scope            string        `bson:"Scope"`            // 权限范围
	Token            string        `bson:"Token"`            // 令牌
	TokenCreateAt    time.Time     `bson:"TokenCreateAt"`    // 令牌创建时间
	TokenExpiresIn   time.Duration `bson:"TokenExpiresIn"`   // 令牌有效期
	Refresh          string        `bson:"Refresh"`          // 更新令牌
	RefreshCreateAt  time.Time     `bson:"RefreshCreateAt"`  // 更新令牌创建时间
	RefreshExpiresIn time.Duration `bson:"RefreshExpiresIn"` // 更新令牌有效期
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

// GetToken 令牌
func (t *Token) GetToken() string {
	return t.Token
}

// SetToken 设置令牌
func (t *Token) SetToken(token string) {
	t.Token = token
}

// GetTokenCreateAt 令牌创建时间
func (t *Token) GetTokenCreateAt() time.Time {
	return t.TokenCreateAt
}

// SetTokenCreateAt 设置令牌创建时间
func (t *Token) SetTokenCreateAt(createAt time.Time) {
	t.TokenCreateAt = createAt
}

// GetTokenExpiresIn 令牌有效期
func (t *Token) GetTokenExpiresIn() time.Duration {
	return t.TokenExpiresIn
}

// SetTokenExpiresIn 设置令牌有效期
func (t *Token) SetTokenExpiresIn(exp time.Duration) {
	t.TokenExpiresIn = exp
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
