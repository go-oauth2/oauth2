package manage

import (
	"time"

	"github.com/LyricTian/errors"
	"github.com/LyricTian/inject"

	"gopkg.in/oauth2.v2"
	"gopkg.in/oauth2.v2/generates"
	"gopkg.in/oauth2.v2/models"
	"gopkg.in/oauth2.v2/store/token"
)

// Config Configuration parameters
type Config struct {
	AccessTokenExp  time.Duration // Access token expiration time (in seconds)
	RefreshTokenExp time.Duration // Refresh token expiration time
}

// NewRedisManager Create to based on redis store authorization management instance
func NewRedisManager(redisCfg *token.RedisConfig) *Manager {
	m := NewManager()
	m.MapClientModel(models.NewClient())
	m.MapTokenModel(models.NewToken())
	m.MapAuthorizeGenerate(generates.NewAuthorizeGenerate())
	m.MapAccessGenerate(generates.NewAccessGenerate())
	m.MustTokenStorage(token.NewRedisStore(redisCfg))

	return m
}

// NewManager Create to authorization management instance
func NewManager() *Manager {
	m := &Manager{
		injector: inject.New(),
		gtcfg:    make(map[oauth2.GrantType]*Config),
	}
	m.SetAuthorizeCodeExp(time.Minute * 10)
	m.SetAuthorizeCodeTokenExp(&Config{AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 24 * 3})
	m.SetImplicitTokenExp(&Config{AccessTokenExp: time.Hour * 1})
	m.SetPasswordTokenExp(&Config{AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 24 * 7})
	m.SetClientTokenExp(&Config{AccessTokenExp: time.Hour * 2})

	return m
}

// Manager Provide authorization management
type Manager struct {
	injector inject.Injector              // Dependency injection
	codeExp  time.Duration                // Authorize code expiration time
	gtcfg    map[oauth2.GrantType]*Config // Authorization grant configuration
}

// SetAuthorizeCodeExp Set the authorization code expiration time
func (m *Manager) SetAuthorizeCodeExp(exp time.Duration) {
	m.codeExp = exp
}

// SetAuthorizeCodeTokenExp Set the authorization code grant token expiration time
func (m *Manager) SetAuthorizeCodeTokenExp(cfg *Config) {
	m.gtcfg[oauth2.AuthorizationCode] = cfg
}

// SetImplicitTokenExp Set the implicit grant token expiration time
func (m *Manager) SetImplicitTokenExp(cfg *Config) {
	m.gtcfg[oauth2.Implicit] = cfg
}

// SetPasswordTokenExp Set the password grant token expiration time
func (m *Manager) SetPasswordTokenExp(cfg *Config) {
	m.gtcfg[oauth2.PasswordCredentials] = cfg
}

// SetClientTokenExp Set the client grant token expiration time
func (m *Manager) SetClientTokenExp(cfg *Config) {
	m.gtcfg[oauth2.ClientCredentials] = cfg
}

// MapClientModel Mapping the client information model
func (m *Manager) MapClientModel(cli oauth2.ClientInfo) error {
	if cli == nil {
		return errors.New(ErrNilValue)
	}
	m.injector.Map(cli)
	return nil
}

// MapTokenModel Mapping the token information model
func (m *Manager) MapTokenModel(token oauth2.TokenInfo) error {
	if token == nil {
		return errors.New(ErrNilValue)
	}
	m.injector.Map(token)
	return nil
}

// MapAuthorizeGenerate Mapping the authorize code generate interface
func (m *Manager) MapAuthorizeGenerate(gen oauth2.AuthorizeGenerate) error {
	if gen == nil {
		return errors.New(ErrNilValue)
	}
	m.injector.Map(gen)
	return nil
}

// MapAccessGenerate Mapping the access token generate interface
func (m *Manager) MapAccessGenerate(gen oauth2.AccessGenerate) error {
	if gen == nil {
		return errors.New(ErrNilValue)
	}
	m.injector.Map(gen)
	return nil
}

// MapClientStorage Mapping the client store interface
func (m *Manager) MapClientStorage(stor oauth2.ClientStore) error {
	if stor == nil {
		return errors.New(ErrNilValue)
	}
	m.injector.Map(stor)
	return nil
}

// MustClientStorage Mandatory mapping the client store interface
func (m *Manager) MustClientStorage(stor oauth2.ClientStore, err error) {
	if err != nil {
		panic(err.Error())
	}
	if stor == nil {
		panic("client store can't be nil value")
	}
	m.injector.Map(stor)
}

// MapTokenStorage Mapping the token store interface
func (m *Manager) MapTokenStorage(stor oauth2.TokenStore) error {
	if stor == nil {
		return errors.New(ErrNilValue)
	}
	m.injector.Map(stor)
	return nil
}

// MustTokenStorage Mandatory mapping the token store interface
func (m *Manager) MustTokenStorage(stor oauth2.TokenStore, err error) {
	if err != nil {
		panic(err)
	}
	if stor == nil {
		panic("token store can't be nil value")
	}
	m.injector.Map(stor)
}

// GetClient Get the client information
func (m *Manager) GetClient(clientID string) (cli oauth2.ClientInfo, err error) {
	_, ierr := m.injector.Invoke(func(stor oauth2.ClientStore) {
		cli, err = stor.GetByID(clientID)
		if err != nil {
			return
		} else if cli == nil {
			err = errors.New(ErrInvalidClient)
		}
	})
	if err == nil && ierr != nil {
		err = ierr
	}
	return
}

// GenerateAuthToken Generate the authorization token(code)
func (m *Manager) GenerateAuthToken(rt oauth2.ResponseType, tgr *oauth2.TokenGenerateRequest) (authToken oauth2.TokenInfo, err error) {
	cli, err := m.GetClient(tgr.ClientID)
	if err != nil {
		return
	} else if verr := ValidateURI(cli.GetDomain(), tgr.RedirectURI); verr != nil {
		err = verr
		return
	}
	_, ierr := m.injector.Invoke(func(ti oauth2.TokenInfo, gen oauth2.AuthorizeGenerate, tgen oauth2.AccessGenerate, stor oauth2.TokenStore) {
		var (
			tv   string
			terr error
		)
		td := &oauth2.GenerateBasic{
			Client:   cli,
			UserID:   tgr.UserID,
			CreateAt: time.Now(),
		}
		if rt == oauth2.Code {
			ti.SetAccessExpiresIn(m.codeExp)
			tv, terr = gen.Token(td)
		} else {
			ti.SetAccessExpiresIn(m.gtcfg[oauth2.Implicit].AccessTokenExp)
			tv, _, terr = tgen.Token(td, false)
		}
		if terr != nil {
			err = terr
			return
		}
		ti.SetAccess(tv)
		ti.SetAccessCreateAt(td.CreateAt)
		ti.SetClientID(tgr.ClientID)
		ti.SetUserID(tgr.UserID)
		ti.SetRedirectURI(tgr.RedirectURI)
		ti.SetScope(tgr.Scope)
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

// GenerateAccessToken Generate the access token
func (m *Manager) GenerateAccessToken(gt oauth2.GrantType, tgr *oauth2.TokenGenerateRequest) (accessToken oauth2.TokenInfo, err error) {
	if gt == oauth2.AuthorizationCode {
		ti, terr := m.LoadAccessToken(tgr.Code)
		if terr != nil {
			err = terr
			return
		} else if ti.GetRedirectURI() != tgr.RedirectURI || ti.GetClientID() != tgr.ClientID {
			err = errors.New(ErrInvalidAuthorizeCode)
			return
		} else if verr := m.RemoveAccessToken(tgr.Code); verr != nil { // remove authorize code
			err = verr
			return
		}
		tgr.UserID = ti.GetUserID()
		tgr.Scope = ti.GetScope()
	}
	cli, err := m.GetClient(tgr.ClientID)
	if err != nil {
		return
	} else if tgr.ClientSecret != cli.GetSecret() {
		err = errors.New(ErrInvalidClient)
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
		ti.SetAccessCreateAt(td.CreateAt)
		ti.SetAccessExpiresIn(m.gtcfg[gt].AccessTokenExp)
		ti.SetAccess(av)
		if tgr.IsGenerateRefresh && rv != "" {
			ti.SetRefreshCreateAt(td.CreateAt)
			ti.SetRefreshExpiresIn(m.gtcfg[gt].RefreshTokenExp)
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

// RefreshAccessToken Refreshing an access token
func (m *Manager) RefreshAccessToken(tgr *oauth2.TokenGenerateRequest) (accessToken oauth2.TokenInfo, err error) {
	cli, err := m.GetClient(tgr.ClientID)
	if err != nil {
		return
	} else if tgr.ClientSecret != cli.GetSecret() {
		err = errors.New(ErrInvalidClient)
		return
	}
	ti, err := m.LoadRefreshToken(tgr.Refresh)
	if err != nil {
		return
	} else if ti.GetClientID() != tgr.ClientID {
		err = errors.New(ErrInvalidRefreshToken)
		return
	}
	oldAccess := ti.GetAccess()
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStore, gen oauth2.AccessGenerate) {
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
		// remove the old access token
		if verr := stor.RemoveByAccess(oldAccess); verr != nil {
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

// RemoveAccessToken Use the access token to delete the token information
func (m *Manager) RemoveAccessToken(access string) (err error) {
	if access == "" {
		err = errors.New(ErrInvalidAccessToken)
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

// RemoveRefreshToken Use the refresh token to delete the token information
func (m *Manager) RemoveRefreshToken(refresh string) (err error) {
	if refresh == "" {
		err = errors.New(ErrInvalidAccessToken)
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

// LoadAccessToken According to the access token for corresponding token information
func (m *Manager) LoadAccessToken(access string) (info oauth2.TokenInfo, err error) {
	if access == "" {
		err = errors.New(ErrInvalidAccessToken)
		return
	}
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStore) {
		ct := time.Now()
		ti, terr := stor.GetByAccess(access)
		if terr != nil {
			err = terr
			return
		} else if ti == nil {
			err = errors.New(ErrInvalidAccessToken)
			return
		} else if ti.GetRefresh() != "" && ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(ct) {
			err = errors.New(ErrExpiredRefreshToken)
		} else if ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).Before(ct) {
			err = errors.New(ErrExpiredAccessToken)
			return
		}
		info = ti
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}

// LoadRefreshToken According to the refresh token for corresponding token information
func (m *Manager) LoadRefreshToken(refresh string) (info oauth2.TokenInfo, err error) {
	if refresh == "" {
		err = errors.New(ErrInvalidRefreshToken)
		return
	}
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStore) {
		ti, terr := stor.GetByRefresh(refresh)
		if terr != nil {
			err = terr
			return
		} else if ti == nil {
			err = errors.New(ErrInvalidRefreshToken)
			return
		} else if ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(time.Now()) {
			err = errors.New(ErrExpiredRefreshToken)
			return
		}
		info = ti
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}
