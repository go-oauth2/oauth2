package oauth2

// NewCCManager 创建默认的客户端模式管理实例
// oaManager OAuth授权管理
// config 配置参数(nil则使用默认值)
func NewCCManager(oaManager *OAuthManager, config *CCConfig) *CCManager {
	if config == nil {
		config = new(CCConfig)
	}
	if config.ATExpiresIn == 0 {
		config.ATExpiresIn = DefaultCCATExpiresIn
	}
	ccManager := &CCManager{
		oAuthManager: oaManager,
		config:       config,
	}
	return ccManager
}

// CCManager 客户端模式管理(Client Credentials Manager)
type CCManager struct {
	oAuthManager *OAuthManager // 授权管理
	config       *CCConfig     // 配置参数
}

// GenerateToken 生成令牌(只生成访问令牌)
// clientID 客户端标识
// clientSecret 客户端秘钥
// scopes 应用授权标识
func (cm *CCManager) GenerateToken(clientID, clientSecret, scopes string) (token *Token, err error) {
	cli, err := cm.oAuthManager.GetClient(clientID)
	if err != nil {
		return
	} else if cli.Secret() != clientSecret {
		err = ErrCSInvalid
		return
	}
	token, err = cm.oAuthManager.GenerateToken(cli, "", scopes, cm.config.ATExpiresIn, 0, false)
	return
}
