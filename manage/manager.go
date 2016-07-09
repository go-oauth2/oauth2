package manage

import (
	"time"

	"github.com/LyricTian/inject"
	"gopkg.in/oauth2.v2"
	"gopkg.in/oauth2.v2/generates"
	"gopkg.in/oauth2.v2/models"
	"gopkg.in/oauth2.v2/store/token"
)

// Config 授权配置参数
type Config struct {
	TokenExp   time.Duration // 令牌有效期
	RefreshExp time.Duration // 更新令牌有效期
}

// NewRedisManager 创建基于redis存储的管理实例
func NewRedisManager(redisCfg *token.RedisConfig) *Manager {
	m := NewManager()
	m.MapClientModel(models.NewClient())
	m.MapTokenModel(models.NewToken())
	m.MapAuthorizeGenerate(generates.NewAuthorizeGenerate())
	m.MapAccessGenerate(generates.NewAccessGenerate())
	m.MustTokenStorage(token.NewRedisStore(redisCfg))

	return m
}

// NewMongoManager 创建基于mongodb存储的管理实例
func NewMongoManager(mongoCfg *token.MongoConfig) *Manager {
	m := NewManager()
	m.MapClientModel(models.NewClient())
	m.MapTokenModel(models.NewToken())
	m.MapAuthorizeGenerate(generates.NewAuthorizeGenerate())
	m.MapAccessGenerate(generates.NewAccessGenerate())
	m.MustTokenStorage(token.NewMongoStore(mongoCfg))

	return m
}

// NewManager 创建Manager的实例
func NewManager() *Manager {
	m := &Manager{
		injector: inject.New(),
		rtcfg:    make(map[oauth2.ResponseType]*Config),
		gtcfg:    make(map[oauth2.GrantType]*Config),
	}
	// 设定参数默认值
	// 设定授权码的有效期为10分钟
	m.SetRTConfig(oauth2.Code, &Config{TokenExp: time.Minute * 10})
	// 设定简化模式授权令牌的有效期为1小时
	m.SetRTConfig(oauth2.Token, &Config{TokenExp: time.Hour * 1})

	// 设定授权码模式令牌的有效期为2小时,更新令牌的有效期为3天
	m.SetGTConfig(oauth2.AuthorizationCodeCredentials, &Config{TokenExp: time.Hour * 2, RefreshExp: time.Hour * 24 * 3})
	// 设定密码模式令牌的有效期为2小时,更新令牌的有效期为7天
	m.SetGTConfig(oauth2.PasswordCredentials, &Config{TokenExp: time.Hour * 2, RefreshExp: time.Hour * 24 * 7})
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
func (m *Manager) MapAuthorizeGenerate(gen oauth2.AuthorizeGenerate) {
	if gen == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(gen)
}

// MapAccessGenerate 注入访问令牌生成接口
func (m *Manager) MapAccessGenerate(gen oauth2.AccessGenerate) {
	if gen == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(gen)
}

// MapClientStorage 注入客户端信息存储接口
func (m *Manager) MapClientStorage(stor oauth2.ClientStore) {
	if stor == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(stor)
}

// MustClientStorage 强制注入客户端信息存储接口
func (m *Manager) MustClientStorage(stor oauth2.ClientStore, err error) {
	if err != nil {
		panic(err)
	}
	if stor == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(stor)
}

// MapTokenStorage 注入令牌信息存储接口
func (m *Manager) MapTokenStorage(stor oauth2.TokenStore) {
	if stor == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(stor)
}

// MustTokenStorage 强制注入令牌信息存储接口
func (m *Manager) MustTokenStorage(stor oauth2.TokenStore, err error) {
	if err != nil {
		panic(err)
	}
	if stor == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(stor)
}

// GetClient 获取客户端信息
// clientID 客户端标识
func (m *Manager) GetClient(clientID string) (cli oauth2.ClientInfo, err error) {
	_, ierr := m.injector.Invoke(func(stor oauth2.ClientStore) {
		cli, err = stor.GetByID(clientID)
		if err != nil {
			return
		} else if cli == nil {
			err = ErrClientNotFound
		}
	})
	if err == nil && ierr != nil {
		err = ierr
	}
	return
}

// GenerateAuthToken 生成授权令牌
// rt 授权类型
// tgr 生成令牌的配置参数
func (m *Manager) GenerateAuthToken(rt oauth2.ResponseType, tgr *oauth2.TokenGenerateRequest) (authToken oauth2.TokenInfo, err error) {
	cli, err := m.GetClient(tgr.ClientID)
	if err != nil {
		return
	} else if verr := ValidateURI(cli.GetDomain(), tgr.RedirectURI); verr != nil {
		err = verr
		return
	}
	_, ierr := m.injector.Invoke(func(ti oauth2.TokenInfo, gen oauth2.AuthorizeGenerate, stor oauth2.TokenStore) {
		td := &oauth2.GenerateBasic{
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
		ti.SetAuthType(rt.String())
		ti.SetAccess(tv)
		ti.SetAccessCreateAt(td.CreateAt)
		ti.SetAccessExpiresIn(m.rtcfg[rt].TokenExp)
		err = stor.Create(ti)
		if err != nil {
			return
		}
		authToken = ti
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}

// GenerateAccessToken 生成访问令牌、更新令牌
// gt 授权模式
// tgr 生成令牌的参数
func (m *Manager) GenerateAccessToken(gt oauth2.GrantType, tgr *oauth2.TokenGenerateRequest) (accessToken oauth2.TokenInfo, err error) {
	if gt == oauth2.AuthorizationCodeCredentials { // 授权码模式
		ti, terr := m.LoadAccessToken(tgr.Code)
		if terr != nil {
			err = terr
			return
		} else if ti.GetRedirectURI() != tgr.RedirectURI || ti.GetClientID() != tgr.ClientID {
			err = ErrAuthCodeInvalid
			return
		} else if verr := m.RemoveAccessToken(tgr.Code); verr != nil { // 删除授权码
			err = verr
			return
		}
		tgr.UserID = ti.GetUserID()
		tgr.Scope = ti.GetScope()
	}
	cli, err := m.GetClient(tgr.ClientID)
	if err != nil {
		return
	} else if tgr.ClientSecret != "" && tgr.ClientSecret != cli.GetSecret() {
		err = ErrClientInvalid
		return
	}
	_, ierr := m.injector.Invoke(func(ti oauth2.TokenInfo, gen oauth2.AccessGenerate, stor oauth2.TokenStore) {
		td := &oauth2.GenerateBasic{
			Client:   cli,
			UserID:   tgr.UserID,
			CreateAt: time.Now(),
		}
		av, rv, terr := gen.Token(td, tgr.IsGenerateRefresh)
		if terr != nil {
			err = terr
			return
		}
		ti.SetClientID(tgr.ClientID)
		ti.SetUserID(tgr.UserID)
		ti.SetRedirectURI(tgr.RedirectURI)
		ti.SetScope(tgr.Scope)
		ti.SetAuthType(gt.String())
		ti.SetAccessCreateAt(td.CreateAt)
		ti.SetAccessExpiresIn(m.gtcfg[gt].TokenExp)
		ti.SetAccess(av)
		if rv != "" {
			ti.SetRefreshCreateAt(td.CreateAt)
			ti.SetRefreshExpiresIn(m.gtcfg[gt].RefreshExp)
			ti.SetRefresh(rv)
		}
		err = stor.Create(ti)
		if err != nil {
			return
		}
		accessToken = ti
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}

// RefreshAccessToken 更新访问令牌
func (m *Manager) RefreshAccessToken(tgr *oauth2.TokenGenerateRequest) (accessToken oauth2.TokenInfo, err error) {
	cli, err := m.GetClient(tgr.ClientID)
	if err != nil {
		return
	} else if tgr.ClientSecret != "" && tgr.ClientSecret != cli.GetSecret() {
		err = ErrClientInvalid
		return
	}
	ti, err := m.LoadRefreshToken(tgr.Refresh)
	if err != nil {
		return
	} else if ti.GetClientID() != tgr.ClientID {
		err = ErrRefreshInvalid
		return
	}
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStore, gen oauth2.AccessGenerate) {
		// 移除旧的访问令牌
		if verr := stor.RemoveByAccess(ti.GetAccess()); verr != nil {
			err = verr
			return
		}
		td := &oauth2.GenerateBasic{
			Client:   cli,
			UserID:   ti.GetUserID(),
			CreateAt: time.Now(),
		}
		tv, _, terr := gen.Token(td, false)
		if terr != nil {
			err = terr
			return
		}
		ti.SetAccess(tv)
		ti.SetAccessCreateAt(td.CreateAt)
		if scope := tgr.Scope; scope != "" {
			ti.SetScope(scope)
		}
		if verr := stor.Create(ti); verr != nil {
			err = verr
			return
		}
		accessToken = ti
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}

// RemoveAccessToken 删除访问令牌
func (m *Manager) RemoveAccessToken(access string) (err error) {
	if access == "" {
		err = ErrAccessInvalid
		return
	}
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStore) {
		err = stor.RemoveByAccess(access)
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}

// RemoveRefreshToken 删除更新令牌
func (m *Manager) RemoveRefreshToken(refresh string) (err error) {
	if refresh == "" {
		err = ErrAccessInvalid
		return
	}
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStore) {
		err = stor.RemoveByRefresh(refresh)
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}

// LoadAccessToken 加载访问令牌信息
func (m *Manager) LoadAccessToken(access string) (info oauth2.TokenInfo, err error) {
	if access == "" {
		err = ErrAccessInvalid
		return
	}
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStore) {
		ct := time.Now()
		ti, terr := stor.GetByAccess(access)
		if terr != nil {
			err = terr
			return
		} else if ti == nil {
			err = ErrAccessInvalid
			return
		} else if ti.GetRefresh() != "" && ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(ct) { // 检查更新令牌是否过期
			err = ErrRefreshExpired
		} else if ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).Before(ct) { // 检查访问令牌是否过期
			err = ErrAccessExpired
			return
		}
		info = ti
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}

// LoadRefreshToken 加载更新令牌信息
func (m *Manager) LoadRefreshToken(refresh string) (info oauth2.TokenInfo, err error) {
	if refresh == "" {
		err = ErrRefreshInvalid
		return
	}
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStore) {
		ti, terr := stor.GetByRefresh(refresh)
		if terr != nil {
			err = terr
			return
		} else if ti == nil {
			err = ErrRefreshInvalid
			return
		} else if ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(time.Now()) {
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
