package manage

import (
	"time"

	"github.com/LyricTian/inject"
	"gopkg.in/oauth2.v2"
)

// NewManager 创建Manager的实例
func NewManager() *Manager {
	return nil
}

// Manager OAuth2授权管理
type Manager struct {
	injector inject.Injector          // 注入器
	rtcfg    map[ResponseType]*Config // 授权类型配置参数
	gtcfg    map[GrantType]*Config    // 授权模式配置参数
}

// SetRTConfig 设定授权类型配置参数
// rt 授权类型
// cfg 配置参数
func (m *Manager) SetRTConfig(rt ResponseType, cfg *Config) {
	m.rtcfg[rt] = cfg
}

// SetGTConfig 设定授权模式配置参数
// gt 授权模式
// cfg 配置参数
func (m *Manager) SetGTConfig(gt GrantType, cfg *Config) {
	m.gtcfg[gt] = cfg
}

// MapClientModel 注入客户端信息模型
func (m *Manager) MapClientModel(cli oauth2.ClientInfo) {
	if cli == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(cli)
}

// MapTokenModel 注入令牌信息模型
func (m *Manager) MapTokenModel(token oauth2.TokenInfo) {
	if token == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(token)
}

// MapAuthorizeGenerate 注入授权令牌生成接口
func (m *Manager) MapAuthorizeGenerate(gen oauth2.AuthorizeGenerate) {
	if gen == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(gen)
}

// MapTokenGenerate 注入访问令牌生成接口
func (m *Manager) MapTokenGenerate(gen oauth2.TokenGenerate) {
	if gen == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(gen)
}

// MapClientStorage 注入客户端信息存储接口
func (m *Manager) MapClientStorage(stor oauth2.ClientStorage) {
	if stor == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(stor)
}

// MustClientStorage 注入客户端信息存储接口
func (m *Manager) MustClientStorage(stor oauth2.ClientStorage, err error) {
	if err != nil {
		panic(err)
	}
	if stor == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(stor)
}

// MapTokenStorage 注入令牌信息存储接口
func (m *Manager) MapTokenStorage(stor oauth2.TokenStorage) {
	if stor == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(stor)
}

// MustTokenStorage 注入令牌信息存储接口
func (m *Manager) MustTokenStorage(stor oauth2.TokenStorage, err error) {
	if err != nil {
		panic(err)
	}
	if stor == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(stor)
}

// GetClient 获取客户端信息
func (m *Manager) GetClient(clientID string) (cli oauth2.ClientInfo, err error) {
	err = m.injector.Apply(func(stor oauth2.ClientStorage) {
		cli, err = stor.GetByID(clientID)
		if err != nil {
			return
		} else if cli == nil {
			err = ErrClientNotFound
		}
	})
	return
}

// GenerateAuthToken 生成授权令牌
// rt 授权类型
// config 生成令牌的配置参数
func (m *Manager) GenerateAuthToken(rt ResponseType, config *TokenGenerateData) (token string, err error) {
	cli, err := m.GetClient(config.ClientID)
	if err != nil {
		return
	} else if verr := ValidateURI(cli.GetDomain(), config.RedirectURI); verr != nil {
		err = verr
		return
	}
	_, ierr := m.injector.Invoke(func(ti oauth2.TokenInfo, gen oauth2.AuthorizeGenerate, stor oauth2.TokenStorage) {
		td := &oauth2.TokenGenerateData{
			Client:   cli,
			UserID:   config.UserID,
			Scope:    config.Scope,
			CreateAt: time.Now(),
		}
		tv, terr := gen.Token(td)
		if terr != nil {
			err = terr
			return
		}
		ti.SetClientID(config.ClientID)
		ti.SetUserID(config.UserID)
		ti.SetRedirectURI(config.RedirectURI)
		ti.SetScope(config.Scope)
		ti.SetTokenCreateAt(td.CreateAt)
		ti.SetTokenExpiresIn(m.rtcfg[rt].TokenExpiresIn)
		ti.SetToken(tv)
		err = stor.Create(ti)
		if err != nil {
			return
		}
		token = tv
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}

// checkAuthToken 检查授权令牌
func (m *Manager) checkAuthToken(config *TokenGenerateData) (err error) {
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStorage) {
		ti, terr := stor.TakeByToken(config.Code)
		if terr != nil {
			err = terr
			return
		} else if ti.GetRedirectURI() != config.RedirectURI || ti.GetClientID() != config.ClientID {
			err = ErrAuthTokenInvalid
			return
		} else if ti.GetTokenCreateAt().Add(ti.GetTokenExpiresIn()).Before(time.Now()) {
			err = ErrAuthTokenInvalid
			return
		}
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}

// GenerateToken 生成令牌
// gt 授权模式
// config 生成令牌的参数
func (m *Manager) GenerateToken(gt GrantType, config *TokenGenerateData) (token, refresh string, err error) {
	if gt == AuthorizationCode {
		err = m.checkAuthToken(config)
		if err != nil {
			return
		}
	}
	cli, err := m.GetClient(config.ClientID)
	if err != nil {
		return
	} else if config.ClientSecret != "" && config.ClientSecret != cli.GetSecret() {
		err = ErrClientInvalid
		return
	}
	_, ierr := m.injector.Invoke(func(ti oauth2.TokenInfo, gen oauth2.TokenGenerate, stor oauth2.TokenStorage) {
		td := &oauth2.TokenGenerateData{
			Client:   cli,
			UserID:   config.UserID,
			Scope:    config.Scope,
			CreateAt: time.Now(),
		}
		tv, rv, terr := gen.Token(td, config.IsGenerateRefresh)
		if terr != nil {
			err = terr
			return
		}
		ti.SetClientID(config.ClientID)
		ti.SetUserID(config.UserID)
		ti.SetRedirectURI(config.RedirectURI)
		ti.SetScope(config.Scope)
		ti.SetTokenCreateAt(td.CreateAt)
		ti.SetTokenExpiresIn(m.gtcfg[gt].TokenExpiresIn)
		ti.SetToken(tv)
		if rv != "" {
			ti.SetRefreshCreateAt(td.CreateAt)
			ti.SetRefreshExpiresIn(m.gtcfg[gt].RefreshExpiresIn)
			ti.SetRefresh(rv)
		}
		err = stor.Create(ti)
		if err != nil {
			return
		}
		token = tv
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}
