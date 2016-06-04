package oauth2

import (
	"github.com/LyricTian/go.uuid"

	"time"
)

// NewOAuthManager 创建OAuth授权管理实例
// cfg 配置参数
func NewOAuthManager(cfg *OAuthConfig) *OAuthManager {
	if cfg == nil {
		cfg = new(OAuthConfig)
	}
	return &OAuthManager{
		Config: cfg,
	}
}

// NewDefaultOAuthManager 创建默认的OAuth授权管理实例
// cfg 配置参数
// mcfg MongoDB配置参数
// ccName 存储客户端的集合名称(默认为ClientInfo)
// tcName 存储令牌的集合名称(默认为AuthToken)
func NewDefaultOAuthManager(cfg *OAuthConfig, mcfg *MongoConfig, ccName, tcName string) (*OAuthManager, error) {
	oManager := NewOAuthManager(cfg)
	clientStore, err := NewClientMongoStore(mcfg, ccName)
	if err != nil {
		return nil, err
	}
	oManager.SetClientStore(clientStore)
	tokenStore, err := NewTokenMongoStore(mcfg, tcName)
	if err != nil {
		return nil, err
	}
	oManager.SetTokenStore(tokenStore)
	oManager.SetTokenGenerate(NewDefaultTokenGenerate())

	return oManager, nil
}

// OAuthManager OAuth授权管理
type OAuthManager struct {
	Config        *OAuthConfig  // 配置参数
	ACGenerate    ACGenerate    // 授权码生成
	ACStore       ACStore       // 授权码存储
	TokenGenerate TokenGenerate // 令牌生成
	TokenStore    TokenStore    // 令牌存储
	ClientStore   ClientStore   // 客户端存储
}

// SetConfig 设置授权码生成接口
func (om *OAuthManager) SetConfig(cfg *OAuthConfig) {
	om.Config = cfg
}

// SetACGenerate 设置授权码生成接口
func (om *OAuthManager) SetACGenerate(generate ACGenerate) {
	om.ACGenerate = generate
}

// SetACStore 设置授权码存储接口
func (om *OAuthManager) SetACStore(store ACStore) {
	om.ACStore = store
}

// SetTokenGenerate 设置令牌生成接口
func (om *OAuthManager) SetTokenGenerate(generate TokenGenerate) {
	om.TokenGenerate = generate
}

// SetTokenStore 设置令牌存储接口
func (om *OAuthManager) SetTokenStore(store TokenStore) {
	om.TokenStore = store
}

// SetClientStore 设置客户端存储接口
func (om *OAuthManager) SetClientStore(store ClientStore) {
	om.ClientStore = store
}

// GetACManager 获取授权码模式管理实例
func (om *OAuthManager) GetACManager() *ACManager {
	return NewACManager(om, om.Config.ACConfig)
}

// GetImplicitManager 获取简化模式管理实例
func (om *OAuthManager) GetImplicitManager() *ImplicitManager {
	return NewImplicitManager(om, om.Config.ImplicitConfig)
}

// GetPasswordManager 获取密码模式管理实例
func (om *OAuthManager) GetPasswordManager() *PasswordManager {
	return NewPasswordManager(om, om.Config.PasswordConfig)
}

// GetCCManager 获取客户端模式管理实例
func (om *OAuthManager) GetCCManager() *CCManager {
	return NewCCManager(om, om.Config.CCConfig)
}

// GenerateToken 生成令牌
// cli 客户端信息
// userID 用户标识
// scopes 应用授权标识
// isGenerateRefresh 是否生成更新令牌
func (om *OAuthManager) GenerateToken(cli Client, userID, scopes string, atExpireIn, rtExpireIn int64, isGenerateRefresh bool) (token *Token, err error) {
	createAt := time.Now().Unix()
	atID := uuid.NewV4().String()
	atBI := NewTokenBasicInfo(cli, atID, userID, createAt)
	atValue, err := om.TokenGenerate.AccessToken(atBI)
	if err != nil {
		return
	}
	tokenValue := Token{
		ClientID:    cli.ID(),
		UserID:      userID,
		AccessToken: atValue,
		ATID:        atID,
		ATCreateAt:  createAt,
		ATExpiresIn: time.Duration(atExpireIn) * time.Second,
		Scope:       scopes,
		CreateAt:    createAt,
		Status:      Actived,
	}
	if isGenerateRefresh {
		rtID := uuid.NewV4().String()
		rtBI := NewTokenBasicInfo(cli, rtID, userID, createAt)
		rtValue, rtErr := om.TokenGenerate.RefreshToken(rtBI)
		if rtErr != nil {
			err = rtErr
			return
		}
		tokenValue.RefreshToken = rtValue
		tokenValue.RTID = rtID
		tokenValue.RTCreateAt = createAt
		tokenValue.RTExpiresIn = time.Duration(rtExpireIn) * time.Second
	}
	id, err := om.TokenStore.Create(&tokenValue)
	if err != nil {
		return
	}
	tokenValue.ID = id
	token = &tokenValue
	return
}

// GetClient 根据客户端标识获取客户端信息
// clientID 客户端标识
func (om *OAuthManager) GetClient(clientID string) (cli Client, err error) {
	cli, err = om.ClientStore.GetByID(clientID)
	if err != nil {
		return
	} else if cli == nil {
		err = ErrClientNotFound
	}
	return
}

// ValidateClient 验证客户端的重定向URI
// clientID 客户端标识
// redirectURI 重定向URI
func (om *OAuthManager) ValidateClient(clientID, redirectURI string) (cli Client, err error) {
	cli, err = om.GetClient(clientID)
	if err != nil {
		return
	} else if v := ValidateURI(cli.Domain(), redirectURI); v != nil {
		err = v
	}
	return
}

// CheckAccessToken 检查访问令牌是否可用，同时返回该令牌的相关信息
// accessToken 访问令牌
func (om *OAuthManager) CheckAccessToken(accessToken string) (token *Token, err error) {
	if accessToken == "" {
		err = ErrATNotFound
		return
	}
	tokenValue, err := om.TokenStore.GetByAccessToken(accessToken)
	if err != nil {
		return
	} else if tokenValue == nil {
		err = ErrATNotFound
		return
	} else if tokenValue.Status != Actived {
		err = ErrATInvalid
		return
	} else if v := om.checkRefreshTokenExpire(tokenValue); v != nil {
		err = v
		return
	} else if v := om.checkAccessTokenExpire(tokenValue); v != nil {
		err = v
		return
	} else if v := om.checkAccessTokenValidity(accessToken, tokenValue); v != nil {
		err = v
		return
	}
	token = tokenValue
	return
}

// RevokeAccessToken 废除访问令牌(将该访问令牌的状态更改为删除)
// accessToken 访问令牌
func (om *OAuthManager) RevokeAccessToken(accessToken string) (err error) {
	if accessToken == "" {
		err = ErrATNotFound
		return
	}
	token, err := om.TokenStore.GetByAccessToken(accessToken)
	if err != nil {
		return
	} else if token == nil {
		err = ErrATNotFound
		return
	} else if token.Status != Actived {
		err = ErrATInvalid
		return
	} else if v := om.checkAccessTokenValidity(accessToken, token); v != nil {
		err = v
		return
	}
	info := map[string]interface{}{
		"Status": Deleted,
	}
	err = om.TokenStore.Update(token.ID, info)
	return
}

// RefreshAccessToken 更新访问令牌(在更新令牌有效期内，更新访问令牌的有效期)，同时返回更新后的令牌信息
// refreshToken 更新令牌
// scopes 申请的权限范围(不可以超出上一次申请的范围，如果省略该参数，则表示与上一次一致)
func (om *OAuthManager) RefreshAccessToken(refreshToken, scopes string) (token *Token, err error) {
	if refreshToken == "" {
		err = ErrRTNotFound
		return
	}
	tokenValue, err := om.TokenStore.GetByRefreshToken(refreshToken)
	if err != nil {
		return
	} else if tokenValue == nil {
		err = ErrRTNotFound
		return
	} else if tokenValue.Status != Actived {
		err = ErrRTInvalid
		return
	} else if v := om.checkRefreshTokenExpire(tokenValue); v != nil {
		err = v
		return
	} else if v := om.checkRefreshTokenValidity(refreshToken, tokenValue); v != nil {
		err = v
		return
	}
	cli, err := om.GetClient(tokenValue.ClientID)
	if err != nil {
		return
	}
	tokenValue.ATCreateAt = time.Now().Unix()
	tokenValue.ATID = uuid.NewV4().String()
	atBI := NewTokenBasicInfo(cli, tokenValue.ATID, tokenValue.UserID, tokenValue.ATCreateAt)
	atValue, err := om.TokenGenerate.AccessToken(atBI)
	if err != nil {
		return
	}
	tokenValue.AccessToken = atValue
	tokenInfo := map[string]interface{}{
		"AccessToken": tokenValue.AccessToken,
		"ATID":        tokenValue.ATID,
		"ATCreateAt":  tokenValue.ATCreateAt,
	}
	if scopes != "" {
		tokenValue.Scope = scopes
		tokenInfo["Scope"] = tokenValue.Scope
	}
	err = om.TokenStore.Update(tokenValue.ID, tokenInfo)
	if err != nil {
		return
	}
	token = tokenValue
	return
}

// checkAccessTokenExpire 检查访问令牌是否过期，
// 如果访问令牌过期同时没有更新令牌的情况下，
// 则将令牌状态更改为过期
func (om *OAuthManager) checkAccessTokenExpire(token *Token) error {
	if token.AccessToken == "" {
		return nil
	}
	nowUnix := time.Now().Unix()
	if (token.ATCreateAt + int64(token.ATExpiresIn/time.Second)) > nowUnix {
		return nil
	}
	var err error
	if token.RefreshToken == "" {
		info := map[string]interface{}{
			"Status": Expired,
		}
		err = om.TokenStore.Update(token.ID, info)
		if err == nil {
			err = ErrATExpire
		}
	}
	return err
}

// checkRefreshTokenExpire 检查更新令牌是否过期，
// 如果更新令牌过期则将令牌状态更改为过期
func (om *OAuthManager) checkRefreshTokenExpire(token *Token) error {
	if token.RefreshToken == "" {
		return nil
	}
	nowUnix := time.Now().Unix()
	if (token.RTCreateAt + int64(token.RTExpiresIn/time.Second)) > nowUnix {
		return nil
	}
	info := map[string]interface{}{
		"Status": Expired,
	}
	err := om.TokenStore.Update(token.ID, info)
	if err == nil {
		err = ErrRTExpire
	}
	return err
}

// checkAccessTokenValidity 检查访问令牌的有效性
func (om *OAuthManager) checkAccessTokenValidity(tv string, token *Token) (err error) {
	cli, err := om.GetClient(token.ClientID)
	if err != nil {
		return
	}
	bi := NewTokenBasicInfo(cli, token.ATID, token.UserID, token.ATCreateAt)
	v, err := om.TokenGenerate.AccessToken(bi)
	if err != nil {
		return
	}
	if tv != v {
		err = ErrATInvalid
	}
	return
}

// checkRefreshTokenValidity 检查刷新令牌的有效性
func (om *OAuthManager) checkRefreshTokenValidity(rv string, token *Token) (err error) {
	cli, err := om.GetClient(token.ClientID)
	if err != nil {
		return
	}
	bi := NewTokenBasicInfo(cli, token.RTID, token.UserID, token.RTCreateAt)
	v, err := om.TokenGenerate.RefreshToken(bi)
	if err != nil {
		return
	}
	if rv != v {
		err = ErrRTInvalid
	}
	return
}
