package models

import "time"

// NewToken Create to token model instance
func NewToken() *Token {
	return &Token{}
}

// Token Token model
type Token struct {
	ClientID         string        `bson:"ClientID"`         // The client id
	UserID           string        `bson:"UserID"`           // The user id
	RedirectURI      string        `bson:"RedirectURI"`      // Redirect URI
	Scope            string        `bson:"Scope"`            // Scope of authorization
	Code             string        `bson:"Code"`             // Authorization code
	CodeCreateAt     time.Time     `bson:"CodeCreateAt"`     // Create Time
	CodeExpiresIn    time.Duration `bson:"CodeExpiresIn"`    // The lifetime in seconds of the authorization code
	Access           string        `bson:"Access"`           // Access Token
	AccessCreateAt   time.Time     `bson:"AccessCreateAt"`   // Create Time
	AccessExpiresIn  time.Duration `bson:"AccessExpiresIn"`  // The lifetime in seconds of the access token
	Refresh          string        `bson:"Refresh"`          // Refresh Token
	RefreshCreateAt  time.Time     `bson:"RefreshCreateAt"`  // Create Time
	RefreshExpiresIn time.Duration `bson:"RefreshExpiresIn"` // The lifetime in seconds of the access token
}

// GetClientID The client id
func (t *Token) GetClientID() string {
	return t.ClientID
}

// SetClientID The client id
func (t *Token) SetClientID(clientID string) {
	t.ClientID = clientID
}

// GetUserID The user id
func (t *Token) GetUserID() string {
	return t.UserID
}

// SetUserID The user id
func (t *Token) SetUserID(userID string) {
	t.UserID = userID
}

// GetRedirectURI Redirect URI
func (t *Token) GetRedirectURI() string {
	return t.RedirectURI
}

// SetRedirectURI Redirect URI
func (t *Token) SetRedirectURI(redirectURI string) {
	t.RedirectURI = redirectURI
}

// GetScope Get Scope of authorization
func (t *Token) GetScope() string {
	return t.Scope
}

// SetScope Get Scope of authorization
func (t *Token) SetScope(scope string) {
	t.Scope = scope
}

// GetCode Authorization code
func (t *Token) GetCode() string {
	return t.Code
}

// SetCode Authorization code
func (t *Token) SetCode(code string) {
	t.Code = code
}

// GetCodeCreateAt Create Time
func (t *Token) GetCodeCreateAt() time.Time {
	return t.CodeCreateAt
}

// SetCodeCreateAt Create Time
func (t *Token) SetCodeCreateAt(createAt time.Time) {
	t.CodeCreateAt = createAt
}

// GetCodeExpiresIn The lifetime in seconds of the authorization code
func (t *Token) GetCodeExpiresIn() time.Duration {
	return t.CodeExpiresIn
}

// SetCodeExpiresIn The lifetime in seconds of the authorization code
func (t *Token) SetCodeExpiresIn(exp time.Duration) {
	t.CodeExpiresIn = exp
}

// GetAccess Access Token
func (t *Token) GetAccess() string {
	return t.Access
}

// SetAccess Access Token
func (t *Token) SetAccess(access string) {
	t.Access = access
}

// GetAccessCreateAt Create Time
func (t *Token) GetAccessCreateAt() time.Time {
	return t.AccessCreateAt
}

// SetAccessCreateAt Create Time
func (t *Token) SetAccessCreateAt(createAt time.Time) {
	t.AccessCreateAt = createAt
}

// GetAccessExpiresIn The lifetime in seconds of the access token
func (t *Token) GetAccessExpiresIn() time.Duration {
	return t.AccessExpiresIn
}

// SetAccessExpiresIn The lifetime in seconds of the access token
func (t *Token) SetAccessExpiresIn(exp time.Duration) {
	t.AccessExpiresIn = exp
}

// GetRefresh Refresh Token
func (t *Token) GetRefresh() string {
	return t.Refresh
}

// SetRefresh Refresh Token
func (t *Token) SetRefresh(refresh string) {
	t.Refresh = refresh
}

// GetRefreshCreateAt Create Time
func (t *Token) GetRefreshCreateAt() time.Time {
	return t.RefreshCreateAt
}

// SetRefreshCreateAt Create Time
func (t *Token) SetRefreshCreateAt(createAt time.Time) {
	t.RefreshCreateAt = createAt
}

// GetRefreshExpiresIn The lifetime in seconds of the access token
func (t *Token) GetRefreshExpiresIn() time.Duration {
	return t.RefreshExpiresIn
}

// SetRefreshExpiresIn The lifetime in seconds of the access token
func (t *Token) SetRefreshExpiresIn(exp time.Duration) {
	t.RefreshExpiresIn = exp
}
