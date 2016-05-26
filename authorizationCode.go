package oauth2

import (
	"time"

	"gopkg.in/LyricTian/lib.v2"
)

// NewACManager 创建授权码模式管理实例
// oaManager OAuth授权管理
// config 配置参数(nil则使用默认值)
func NewACManager(oaManager *OAuthManager, config *ACConfig) *ACManager {
	if config == nil {
		config = new(ACConfig)
	}
	if config.RandomCodeLen == 0 {
		config.RandomCodeLen = DefaultRandomCodeLen
	}
	if config.ACExpiresIn == 0 {
		config.ACExpiresIn = DefaultACExpiresIn
	}
	if config.ATExpiresIn == 0 {
		config.ATExpiresIn = DefaultATExpiresIn
	}
	if config.RTExpiresIn == 0 {
		config.RTExpiresIn = DefaultRTExpiresIn
	}
	acManager := &ACManager{
		oAuthManager: oaManager,
		config:       config,
	}
	return acManager
}

// ACManager 授权码模式管理(Authorization Code Manager)
type ACManager struct {
	oAuthManager *OAuthManager // 授权管理
	config       *ACConfig     // 配置参数
}

// GenerateCode 生成授权码
// clientID 客户端标识
// userID 用户标识
// redirectURI 重定向URI
// scopes 应用授权标识
func (am *ACManager) GenerateCode(clientID, userID, redirectURI, scopes string) (code string, err error) {
	cli, err := am.oAuthManager.ValidateClient(clientID, redirectURI)
	if err != nil {
		return
	}
	acInfo := ACInfo{
		ClientID:    cli.ID(),
		UserID:      userID,
		RedirectURI: redirectURI,
		Scope:       scopes,
		Code:        lib.NewRandom(am.config.RandomCodeLen).NumberAndLetter(),
		CreateAt:    time.Now().Unix(),
		ExpiresIn:   time.Duration(am.config.ACExpiresIn) * time.Second,
	}
	id, err := am.oAuthManager.ACStore.Put(acInfo)
	if err != nil {
		return
	}
	acInfo.ID = id
	code, err = am.oAuthManager.ACGenerate.Code(&acInfo)
	return
}

// GenerateToken 生成令牌
// code 授权码
// redirectURI 重定向URI
// clientID 客户端标识
// clientSecret 客户端秘钥
// isGenerateRefresh 是否生成更新令牌
func (am *ACManager) GenerateToken(code, redirectURI, clientID, clientSecret string, isGenerateRefresh bool) (token *Token, err error) {
	acInfo, err := am.getACInfo(code)
	if err != nil {
		return
	} else if acInfo.RedirectURI != redirectURI {
		err = ErrACInvalid
		return
	} else if acInfo.ClientID != clientID {
		err = ErrACInvalid
		return
	}
	cli, err := am.oAuthManager.ClientStore.GetByID(acInfo.ClientID)
	if err != nil {
		return
	} else if clientSecret != cli.Secret() {
		err = ErrCSInvalid
		return
	}
	createAt := time.Now().Unix()
	basicInfo := NewTokenBasicInfo(cli, acInfo.UserID, createAt)
	atValue, err := am.oAuthManager.TokenGenerate.AccessToken(basicInfo)
	if err != nil {
		return
	}
	tokenValue := Token{
		ClientID:    acInfo.ClientID,
		UserID:      acInfo.UserID,
		AccessToken: atValue,
		ATCreateAt:  createAt,
		ATExpiresIn: time.Duration(am.config.ATExpiresIn) * time.Second,
		Scope:       acInfo.Scope,
		CreateAt:    createAt,
		Status:      Actived,
	}
	if isGenerateRefresh {
		rtValue, rtErr := am.oAuthManager.TokenGenerate.RefreshToken(basicInfo)
		if rtErr != nil {
			err = rtErr
			return
		}
		tokenValue.RefreshToken = rtValue
		tokenValue.RTCreateAt = createAt
		tokenValue.RTExpiresIn = time.Duration(am.config.RTExpiresIn) * time.Second
	}
	id, err := am.oAuthManager.TokenStore.Create(tokenValue)
	if err != nil {
		return
	}
	tokenValue.ID = id
	token = &tokenValue
	return
}

// getACInfo 根据授权码获取授权信息
func (am *ACManager) getACInfo(code string) (info *ACInfo, err error) {
	if code == "" {
		err = ErrACNotFound
		return
	}
	acID, err := am.oAuthManager.ACGenerate.Parse(code)
	if err != nil {
		return
	}
	acInfo, err := am.oAuthManager.ACStore.TakeByID(acID)
	if err != nil {
		return
	}
	acValid, err := am.oAuthManager.ACGenerate.Verify(code, acInfo)
	if err != nil {
		return
	}
	if !acValid ||
		(acInfo.CreateAt+int64(acInfo.ExpiresIn/time.Second)) < time.Now().Unix() {
		err = ErrACInvalid
		return
	}
	info = acInfo
	return
}
