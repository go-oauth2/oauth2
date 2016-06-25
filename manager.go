package oauth2

import (
	"time"

	"github.com/LyricTian/inject"
)

// NewManager 创建Manager的实例
func NewManager() *Manager {
	return nil
}

// Config 授权配置参数
type Config struct {
	CodeExpiresIn    time.Duration // 授权码有效期
	AccessExpiresIn  time.Duration // 访问令牌有效期
	RefreshExpiresIn time.Duration // 刷新令牌有效期
}

// Manager OAuth2授权管理
type Manager struct {
	injector inject.Injector
	configs  map[GrantType]*Config
}

// SetConfig 设定配置参数
func (m *Manager) SetConfig(gt GrantType, cfg *Config) {
	m.configs[gt] = cfg
}

// MapClientModel 注入客户端信息模型
func (m *Manager) MapClientModel(cli ClientInfo) {
	if cli == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(cli)
}

// MapAuthorizeModel 注入授权信息模型
func (m *Manager) MapAuthorizeModel(auth Authorize) {
	if auth == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(auth)
}

// MapTokenModel 注入令牌信息模型
func (m *Manager) MapTokenModel(token TokenInfo) {
	if token == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(token)
}

// MapAuthorizeGenerate 注入授权令牌生成接口
func (m *Manager) MapAuthorizeGenerate(gen AuthorizeGenerate) {
	if gen == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(gen)
}

// MapTokenGenerate 注入访问令牌生成接口
func (m *Manager) MapTokenGenerate(gen TokenGenerate) {
	if gen == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(gen)
}

// MapClientStorage 注入客户端信息存储接口
func (m *Manager) MapClientStorage(stor ClientStorage) {
	if stor == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(stor)
}

// MustClientStorage 注入客户端信息存储接口
func (m *Manager) MustClientStorage(stor ClientStorage, err error) {
	if err != nil {
		panic(err)
	}
	if stor == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(stor)
}

// MapAuthorizeStorage 注入授权码信息存储接口
func (m *Manager) MapAuthorizeStorage(stor AuthorizeStorage) {
	if stor == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(stor)
}

// MustAuthorizeStorage 注入授权码信息存储接口
func (m *Manager) MustAuthorizeStorage(stor AuthorizeStorage, err error) {
	if err != nil {
		panic(err)
	}
	if stor == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(stor)
}

// MapTokenStorage 注入令牌信息存储接口
func (m *Manager) MapTokenStorage(stor TokenStorage) {
	if stor == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(stor)
}

// MustTokenStorage 注入令牌信息存储接口
func (m *Manager) MustTokenStorage(stor TokenStorage, err error) {
	if err != nil {
		panic(err)
	}
	if stor == nil {
		panic(ErrNilValue)
	}
	m.injector.Map(stor)
}
