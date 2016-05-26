package oauth2

import "time"

// NewPasswordManager 创建默认的密码模式管理实例
// oaManager OAuth授权管理
// config 配置参数(nil则使用默认值)
func NewPasswordManager(oaManager *OAuthManager, config *PasswordConfig) *PasswordManager {
	if config == nil {
		config = new(PasswordConfig)
	}
	if config.ATExpiresIn == 0 {
		config.ATExpiresIn = DefaultATExpiresIn
	}
	if config.RTExpiresIn == 0 {
		config.RTExpiresIn = DefaultRTExpiresIn
	}
	pManager := &PasswordManager{
		oAuthManager: oaManager,
		config:       config,
	}
	return pManager
}

// PasswordManager 密码模式管理
type PasswordManager struct {
	oAuthManager *OAuthManager   // 授权管理
	config       *PasswordConfig // 配置参数
}

// GenerateToken 生成令牌(只生成访问令牌)
// clientID 客户端标识
// userID 用户标识
// clientSecret 客户端秘钥
// scopes 应用授权标识
func (pm *PasswordManager) GenerateToken(clientID, userID, clientSecret, scopes string, isGenerateRefresh bool) (token *Token, err error) {
	cli, err := pm.oAuthManager.GetClient(clientID)
	if err != nil {
		return
	} else if cli.Secret() != clientSecret {
		err = ErrCSInvalid
		return
	}
	createAt := time.Now().Unix()
	basicInfo := NewTokenBasicInfo(cli, userID, createAt)
	atValue, err := pm.oAuthManager.TokenGenerate.AccessToken(basicInfo)
	if err != nil {
		return
	}
	tokenValue := Token{
		ClientID:    clientID,
		UserID:      userID,
		AccessToken: atValue,
		ATCreateAt:  createAt,
		ATExpiresIn: time.Duration(pm.config.ATExpiresIn) * time.Second,
		Scope:       scopes,
		CreateAt:    createAt,
		Status:      Actived,
	}
	if isGenerateRefresh {
		rtValue, rtErr := pm.oAuthManager.TokenGenerate.RefreshToken(basicInfo)
		if rtErr != nil {
			err = rtErr
			return
		}
		tokenValue.RefreshToken = rtValue
		tokenValue.RTCreateAt = createAt
		tokenValue.RTExpiresIn = time.Duration(pm.config.RTExpiresIn) * time.Second
	}
	id, err := pm.oAuthManager.TokenStore.Create(tokenValue)
	if err != nil {
		return
	}
	tokenValue.ID = id
	token = &tokenValue
	return
}
