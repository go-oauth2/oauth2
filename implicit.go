package oauth2

// NewImplicitManager 创建默认的简化模式管理实例
// oaManager OAuth授权管理
// config 配置参数(nil则使用默认值)
func NewImplicitManager(oaManager *OAuthManager, config *ImplicitConfig) *ImplicitManager {
	if config == nil {
		config = new(ImplicitConfig)
	}
	if config.ATExpiresIn == 0 {
		config.ATExpiresIn = DefaultIATExpiresIn
	}
	iManager := &ImplicitManager{
		oAuthManager: oaManager,
		config:       config,
	}
	return iManager
}

// ImplicitManager 简化模式管理
type ImplicitManager struct {
	oAuthManager *OAuthManager   // 授权管理
	config       *ImplicitConfig // 配置参数
}

// GenerateToken 生成令牌(只生成访问令牌)
// clientID 客户端标识
// userID 用户标识
// redirectURI 重定向URI
// scopes 应用授权标识
func (im *ImplicitManager) GenerateToken(clientID, userID, redirectURI, scopes string) (token *Token, err error) {
	cli, err := im.oAuthManager.ValidateClient(clientID, redirectURI)
	if err != nil {
		return
	}
	token, err = im.oAuthManager.GenerateToken(cli, userID, scopes, im.config.ATExpiresIn, 0, false)
	return
}
