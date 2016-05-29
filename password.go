package oauth2

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

	token, err = pm.oAuthManager.GenerateToken(cli,
		userID,
		scopes,
		pm.config.ATExpiresIn,
		pm.config.RTExpiresIn,
		isGenerateRefresh)

	return
}
