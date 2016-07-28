package manage

import (
	"reflect"
	"time"

	"github.com/codegangsta/inject"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/models"
)

// Config Configuration parameters
type Config struct {
	AccessTokenExp    time.Duration // Access token expiration time (in seconds)
	RefreshTokenExp   time.Duration // Refresh token expiration time
	IsGenerateRefresh bool          // Whether to generate the refreshing token
}

// NewDefaultManager Create to default authorization management instance
func NewDefaultManager() *Manager {
	m := NewManager()

	// default config
	m.SetAuthorizeCodeExp(time.Minute * 10)
	m.SetImplicitTokenCfg(&Config{AccessTokenExp: time.Hour * 1})
	m.SetClientTokenCfg(&Config{AccessTokenExp: time.Hour * 2})
	m.SetAuthorizeCodeTokenCfg(&Config{IsGenerateRefresh: true, AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 24 * 3})
	m.SetPasswordTokenCfg(&Config{IsGenerateRefresh: true, AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 24 * 7})

	m.MapTokenModel(models.NewToken())
	m.MapAuthorizeGenerate(generates.NewAuthorizeGenerate())
	m.MapAccessGenerate(generates.NewAccessGenerate())

	return m
}

// NewManager Create to authorization management instance
func NewManager() *Manager {
	return &Manager{
		injector: inject.New(),
		gtcfg:    make(map[oauth2.GrantType]*Config),
	}
}

// Manager Provide authorization management
type Manager struct {
	injector inject.Injector              // Dependency injection
	codeExp  time.Duration                // Authorize code expiration time
	gtcfg    map[oauth2.GrantType]*Config // Authorization grant configuration
}

func (m *Manager) newTokenInfo(ti oauth2.TokenInfo) oauth2.TokenInfo {
	in := reflect.ValueOf(ti)
	if in.IsNil() {
		return ti
	}
	out := reflect.New(in.Type().Elem())
	return out.Interface().(oauth2.TokenInfo)
}

// SetAuthorizeCodeExp Set the authorization code expiration time
func (m *Manager) SetAuthorizeCodeExp(exp time.Duration) {
	m.codeExp = exp
}

// SetAuthorizeCodeTokenCfg Set the authorization code grant token config
func (m *Manager) SetAuthorizeCodeTokenCfg(cfg *Config) {
	m.gtcfg[oauth2.AuthorizationCode] = cfg
}

// SetImplicitTokenCfg Set the implicit grant token config
func (m *Manager) SetImplicitTokenCfg(cfg *Config) {
	m.gtcfg[oauth2.Implicit] = cfg
}

// SetPasswordTokenCfg Set the password grant token config
func (m *Manager) SetPasswordTokenCfg(cfg *Config) {
	m.gtcfg[oauth2.PasswordCredentials] = cfg
}

// SetClientTokenCfg Set the client grant token config
func (m *Manager) SetClientTokenCfg(cfg *Config) {
	m.gtcfg[oauth2.ClientCredentials] = cfg
}

// SetRefreshTokenCfg Set the refreshing token config
func (m *Manager) SetRefreshTokenCfg(cfg *Config) {
	m.gtcfg[oauth2.Refreshing] = cfg
}

// MapTokenModel Mapping the token information model
func (m *Manager) MapTokenModel(token oauth2.TokenInfo) error {
	if token == nil {
		return errors.ErrNilValue
	}
	m.injector.Map(token)
	return nil
}

// MapAuthorizeGenerate Mapping the authorize code generate interface
func (m *Manager) MapAuthorizeGenerate(gen oauth2.AuthorizeGenerate) error {
	if gen == nil {
		return errors.ErrNilValue
	}
	m.injector.Map(gen)
	return nil
}

// MapAccessGenerate Mapping the access token generate interface
func (m *Manager) MapAccessGenerate(gen oauth2.AccessGenerate) error {
	if gen == nil {
		return errors.ErrNilValue
	}
	m.injector.Map(gen)
	return nil
}

// MapClientStorage Mapping the client store interface
func (m *Manager) MapClientStorage(stor oauth2.ClientStore) error {
	if stor == nil {
		return errors.ErrNilValue
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
		return errors.ErrNilValue
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
			err = errors.ErrInvalidClient
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
		ti = m.newTokenInfo(ti)

		td := &oauth2.GenerateBasic{
			Client:   cli,
			UserID:   tgr.UserID,
			CreateAt: time.Now(),
		}
		switch rt {
		case oauth2.Code:
			tv, terr := gen.Token(td)
			if terr != nil {
				err = terr
				return
			}
			ti.SetCode(tv)
			ti.SetCodeExpiresIn(m.codeExp)
			ti.SetCodeCreateAt(td.CreateAt)
			if exp := tgr.AccessTokenExp; exp > 0 {
				ti.SetAccessExpiresIn(exp)
			}
		case oauth2.Token:
			tv, rv, terr := tgen.Token(td, m.gtcfg[oauth2.Implicit].IsGenerateRefresh)
			if terr != nil {
				err = terr
				return
			}
			ti.SetAccess(tv)
			ti.SetAccessCreateAt(td.CreateAt)
			aexp := m.gtcfg[oauth2.Implicit].AccessTokenExp
			if exp := tgr.AccessTokenExp; exp > 0 {
				aexp = exp
			}
			ti.SetAccessExpiresIn(aexp)
			if rv != "" && m.gtcfg[oauth2.Implicit].IsGenerateRefresh {
				ti.SetRefresh(rv)
				ti.SetRefreshCreateAt(td.CreateAt)
				ti.SetRefreshExpiresIn(m.gtcfg[oauth2.Implicit].RefreshTokenExp)
			}
		}
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

// get authorization code data
func (m *Manager) getAuthorizationCode(code string) (info oauth2.TokenInfo, err error) {
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStore) {
		ti, terr := stor.GetByCode(code)
		if terr != nil {
			err = terr
			return
		} else if ti == nil {
			err = errors.ErrInvalidAuthorizeCode
			return
		} else if ti.GetCodeCreateAt().Add(ti.GetCodeExpiresIn()).Before(time.Now()) {
			err = errors.ErrInvalidAuthorizeCode
			return
		}
		info = ti
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}

// delete authorization code data
func (m *Manager) delAuthorizationCode(code string) (err error) {
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStore) {
		err = stor.RemoveByCode(code)
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}

// GenerateAccessToken Generate the access token
func (m *Manager) GenerateAccessToken(gt oauth2.GrantType, tgr *oauth2.TokenGenerateRequest) (accessToken oauth2.TokenInfo, err error) {
	if gt == oauth2.AuthorizationCode {
		ti, terr := m.getAuthorizationCode(tgr.Code)
		if terr != nil {
			err = terr
			return
		} else if ti.GetRedirectURI() != tgr.RedirectURI || ti.GetClientID() != tgr.ClientID {
			err = errors.ErrInvalidAuthorizeCode
			return
		} else if verr := m.delAuthorizationCode(tgr.Code); verr != nil {
			err = verr
			return
		}
		tgr.UserID = ti.GetUserID()
		tgr.Scope = ti.GetScope()
		if exp := ti.GetAccessExpiresIn(); exp > 0 {
			tgr.AccessTokenExp = exp
		}
	}
	cli, err := m.GetClient(tgr.ClientID)
	if err != nil {
		return
	} else if tgr.ClientSecret != cli.GetSecret() {
		err = errors.ErrInvalidClient
		return
	}
	_, ierr := m.injector.Invoke(func(ti oauth2.TokenInfo, gen oauth2.AccessGenerate, stor oauth2.TokenStore) {
		ti = m.newTokenInfo(ti)
		td := &oauth2.GenerateBasic{
			Client:   cli,
			UserID:   tgr.UserID,
			CreateAt: time.Now(),
		}
		av, rv, terr := gen.Token(td, m.gtcfg[gt].IsGenerateRefresh)
		if terr != nil {
			err = terr
			return
		}
		ti.SetClientID(tgr.ClientID)
		ti.SetUserID(tgr.UserID)
		ti.SetRedirectURI(tgr.RedirectURI)
		ti.SetScope(tgr.Scope)
		ti.SetAccessCreateAt(td.CreateAt)
		ti.SetAccess(av)

		aexp := m.gtcfg[gt].AccessTokenExp
		if exp := tgr.AccessTokenExp; exp > 0 {
			aexp = exp
		}
		ti.SetAccessExpiresIn(aexp)
		if rv != "" && m.gtcfg[gt].IsGenerateRefresh {
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
		err = errors.ErrInvalidClient
		return
	}
	ti, err := m.LoadRefreshToken(tgr.Refresh)
	if err != nil {
		return
	} else if ti.GetClientID() != tgr.ClientID {
		err = errors.ErrInvalidRefreshToken
		return
	}
	oldAccess := ti.GetAccess()
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStore, gen oauth2.AccessGenerate) {
		td := &oauth2.GenerateBasic{
			Client:   cli,
			UserID:   ti.GetUserID(),
			CreateAt: time.Now(),
		}
		isGenRefresh := false
		if rcfg, ok := m.gtcfg[oauth2.Refreshing]; ok {
			isGenRefresh = rcfg.IsGenerateRefresh
		}
		tv, rv, terr := gen.Token(td, isGenRefresh)
		if terr != nil {
			err = terr
			return
		}
		ti.SetAccess(tv)
		ti.SetAccessCreateAt(td.CreateAt)
		if scope := tgr.Scope; scope != "" {
			ti.SetScope(scope)
		}
		if rv != "" {
			ti.SetRefresh(rv)
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
		err = errors.ErrInvalidAccessToken
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
		err = errors.ErrInvalidAccessToken
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
		err = errors.ErrInvalidAccessToken
		return
	}
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStore) {
		ct := time.Now()
		ti, terr := stor.GetByAccess(access)
		if terr != nil {
			err = terr
			return
		} else if ti == nil {
			err = errors.ErrInvalidAccessToken
			return
		} else if ti.GetRefresh() != "" && ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(ct) {
			err = errors.ErrExpiredRefreshToken
		} else if ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).Before(ct) {
			err = errors.ErrExpiredAccessToken
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
		err = errors.ErrInvalidRefreshToken
		return
	}
	_, ierr := m.injector.Invoke(func(stor oauth2.TokenStore) {
		ti, terr := stor.GetByRefresh(refresh)
		if terr != nil {
			err = terr
			return
		} else if ti == nil {
			err = errors.ErrInvalidRefreshToken
			return
		} else if ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(time.Now()) {
			err = errors.ErrExpiredRefreshToken
			return
		}
		info = ti
	})
	if ierr != nil && err == nil {
		err = ierr
	}
	return
}
