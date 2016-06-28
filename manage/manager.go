package manage

import (
	"time"

	"github.com/LyricTian/inject"
	"gopkg.in/oauth2.v2"
)

// Config 授权配置参数
type Config struct {
	TokenExp   time.Duration // 令牌有效期
	RefreshExp time.Duration // g令牌有效期
}

// NewManager 创建Manager的实例
func NewManager() *Manager {
	m := &Manager{
		injector: inject.New(),
	}
	// 设定参数默认值

	// 设定授权码的有效期为10分钟
	m.SetRTConfig(oauth2.Code, &Config{TokenExp: time.Minute * 10})
	// 设定简化模式授权令牌的有效期为1小时
	m.SetRTConfig(oauth2.Token, &Config{TokenExp: time.Hour * 1})

	// 设定授权码模式令牌的有效期为2小时,g令牌的有效期为3天
	m.SetGTConfig(oauth2.PasswordCredentials, &Config{TokenExp: time.Hour * 2, RefreshExp: time.Hour * 24 * 3})

	// 设定客户端模式令牌的有效期为1小时
	m.SetGTConfig(oauth2.ClientCredentials, &Config{TokenExp: time.Hour * 2})
	return m
}

// Manager OAuth2授权管理
type Manager struct {
	injector inject.Injector                 // 注入器
	rtcfg    map[oauth2.ResponseType]*Config // 授权类型配置参数
	gtcfg    map[oauth2.GrantType]*Config    // 授权模式配置参数
}

// SetRTConfig 设定授权类型配置参数
// rt 授权类型
// cfg 配置参数
func (m *Manager) SetRTConfig(rt oauth2.ResponseType, cfg *Config) {
	m.rtcfg[rt] = cfg
}

// SetGTConfig 设定授权模式配置参数
// gt 授权模式
// cfg 配置参数
func (m *Manager) SetGTConfig(gt oauth2.GrantType, cfg *Config) {
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
func (m *Manager) MapAuthorizeGenerate(gen oauth2.AuthorizeTokenGenerate) {
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
// tgr 生成令牌的配置参数
func (m *Manager) GenerateAuthToken(rt oauth2.ResponseType, tgr *oauth2.TokenGenerateRequest) (token string, err error) {
	cli, err := m.GetClient(tgr.ClientID)
	if err != nil {
		return
	} else if verr := ValidateURI(cli.GetDomain(), tgr.RedirectURI); verr != nil {
		err = verr
		return
	}
	_, ierr := m.injector.Invoke(func(ti oauth2.TokenInfo, gen oauth2.AuthorizeTokenGenerate, stor oauth2.TokenStorage) {
		td := &oauth2.TokenGenerateBasic{
			Client:   cli,
			UserID:   tgr.UserID,
			CreateAt: time.Now(),
		}
		tv, terr := gen.Token(td)
		if terr != nil {
			err = terr
			return
		}
		ti.SetClientID(tgr.ClientID)
		ti.SetUserID(tgr.UserID)
		ti.SetRedirectURI(tgr.RedirectURI)
		ti.SetScope(tgr.Scope)
		ti.SetTokenCreateAt(td.CreateAt)
		ti.SetTokenExpiresIn(m.rtcfg[rt].TokenExp)
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
func (m *Manager) checkAuthToken(tgr *oauth2.TokenGenerateRequest) (err error) {
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStorage) {
		ti, terr := stor.TakeByToken(tgr.Code)
		if terr != nil {
			err = terr
			return
		} else if ti.GetRedirectURI() != tgr.RedirectURI || ti.GetClientID() != tgr.ClientID {
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
// tgr 生成令牌的参数
func (m *Manager) GenerateToken(gt oauth2.GrantType, tgr *oauth2.TokenGenerateRequest) (token, refresh string, err error) {
	if gt == oauth2.AuthorizationCodeCredentials {
		err = m.checkAuthToken(tgr)
		if err != nil {
			return
		}
	}
	cli, err := m.GetClient(tgr.ClientID)
	if err != nil {
		return
	} else if tgr.ClientSecret != "" && tgr.ClientSecret != cli.GetSecret() {
		err = ErrClientInvalid
		return
	}
	_, ierr := m.injector.Invoke(func(ti oauth2.TokenInfo, gen oauth2.TokenGenerate, stor oauth2.TokenStorage) {
		td := &oauth2.TokenGenerateBasic{
			Client:   cli,
			UserID:   tgr.UserID,
			CreateAt: time.Now(),
		}
		tv, rv, terr := gen.Token(td, tgr.IsGenerateRefresh)
		if terr != nil {
			err = terr
			return
		}
		ti.SetClientID(tgr.ClientID)
		ti.SetUserID(tgr.UserID)
		ti.SetRedirectURI(tgr.RedirectURI)
		ti.SetScope(tgr.Scope)
		ti.SetTokenCreateAt(td.CreateAt)
		ti.SetTokenExpiresIn(m.gtcfg[gt].TokenExp)
		ti.SetToken(tv)
		if rv != "" {
			ti.SetRefreshCreateAt(td.CreateAt)
			ti.SetRefreshExpiresIn(m.gtcfg[gt].RefreshExp)
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

// RefreshToken 更新访问令牌
func (m *Manager) RefreshToken(refresh, scope string) (token string, err error) {
	ti, err := m.CheckRefreshToken(refresh)
	if err != nil {
		return
	}
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStorage, gen oauth2.TokenGenerate) {
		cli, cerr := m.GetClient(ti.GetClientID())
		if cerr != nil {
			err = cerr
			return
		}
		td := &oauth2.TokenGenerateBasic{
			Client:   cli,
			UserID:   ti.GetUserID(),
			CreateAt: time.Now(),
		}
		tv, _, terr := gen.Token(td, false)
		if terr != nil {
			err = terr
			return
		}
		ti.SetToken(tv)
		ti.SetTokenCreateAt(td.CreateAt)
		if scope != "" {
			ti.SetScope(scope)
		}
		err = stor.UpdateByRefresh(refresh, ti)
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

// RevokeToken 废除令牌
func (m *Manager) RevokeToken(token string) (err error) {
	if token == "" {
		err = ErrTokenInvalid
		return
	}
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStorage) {
		err = stor.DeleteByToken(token)
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}

// CheckToken 令牌检查
func (m *Manager) CheckToken(token string) (info oauth2.TokenInfo, err error) {
	if token == "" {
		err = ErrTokenInvalid
		return
	}
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStorage) {
		ct := time.Now()
		ti, terr := stor.GetByToken(token)
		if terr != nil {
			err = terr
			return
		} else if ti == nil {
			err = ErrTokenInvalid
			return
		} else if ti.GetRefresh() != "" && ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(ct) { // 检查g令牌是否过期
			if verr := stor.ExpiredByRefresh(ti.GetRefresh()); verr != nil {
				err = verr
				return
			}
			err = ErrRefreshExpired
		} else if ti.GetTokenCreateAt().Add(ti.GetTokenExpiresIn()).Before(ct) { // 检查令牌是否过期
			if verr := stor.ExpiredByToken(token); verr != nil {
				err = verr
				return
			}
			err = ErrTokenExpired
			return
		}
		info = ti
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}

// CheckRefreshToken 更新令牌检查
func (m *Manager) CheckRefreshToken(refresh string) (info oauth2.TokenInfo, err error) {
	if refresh == "" {
		err = ErrRefreshInvalid
		return
	}
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStorage) {
		ti, terr := stor.GetByRefresh(refresh)
		if terr != nil {
			err = terr
			return
		} else if ti == nil {
			err = ErrRefreshInvalid
			return
		} else if ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(time.Now()) {
			// 废除过期的令牌
			if verr := stor.ExpiredByRefresh(refresh); verr != nil {
				err = verr
				return
			}
			err = ErrRefreshExpired
			return
		}
		info = ti
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}
