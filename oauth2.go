package oauth2

import "time"

// CreateDefaultOAuthManager 创建默认的OAuth授权管理实例
// mongoConfig MongoDB配置参数
// tokenCollectionName 存储令牌的集合名称(默认为AuthToken)
// clientCollectionName 存储客户端的集合名称(默认为ClientInfo)
// oauthConfig 配置参数
func CreateDefaultOAuthManager(mongoConfig *MongoConfig, tokenCollectionName, clientCollectionName string, oauthConfig *OAuthConfig) (*OAuthManager, error) {
	if oauthConfig == nil {
		oauthConfig = new(OAuthConfig)
	}
	oaManager := &OAuthManager{
		Config:        oauthConfig,
		ACGenerate:    NewDefaultACGenerate(),
		ACStore:       NewACMemoryStore(0),
		TokenGenerate: NewDefaultTokenGenerate(),
	}
	tokenStore, err := NewTokenMongoStore(mongoConfig, tokenCollectionName)
	if err != nil {
		return nil, err
	}
	oaManager.TokenStore = tokenStore
	clientStore, err := NewClientMongoStore(mongoConfig, clientCollectionName)
	if err != nil {
		return nil, err
	}
	oaManager.ClientStore = clientStore
	return oaManager, nil
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

// SetACGenerate 设置授权码生成接口
func (om *OAuthManager) SetACGenerate(generate ACGenerate) {
	om.ACGenerate = generate
}

// SetACStore 设置授权码存储接口
func (om *OAuthManager) SetACStore(store ACStore) {
	om.ACStore = store
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
	}
	cli, err := om.GetClient(tokenValue.ClientID)
	if err != nil {
		return
	}
	tokenValue.ATCreateAt = time.Now().Unix()
	atValue, err := om.TokenGenerate.AccessToken(NewTokenBasicInfo(cli, tokenValue.UserID, tokenValue.ATCreateAt))
	if err != nil {
		return
	}
	tokenValue.AccessToken = atValue
	tokenInfo := map[string]interface{}{
		"AccessToken": tokenValue.AccessToken,
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
