package manage

import (
	"time"

	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/models"
)

// NewDefaultManager create to default authorization management instance
func NewDefaultManager() *Manager {
	m := NewManager()
	// default implementation
	m.MapAuthorizeGenerate(generates.NewAuthorizeGenerate())
	m.MapAccessGenerate(generates.NewAccessGenerate())

	return m
}

// NewManager create to authorization management instance
func NewManager() *Manager {
	return &Manager{
		gtcfg:       make(map[oauth2.GrantType]*Config),
		validateURI: DefaultValidateURI,
	}
}

// Manager provide authorization management
type Manager struct {
	codeExp           time.Duration
	gtcfg             map[oauth2.GrantType]*Config
	rcfg              *RefreshingConfig
	validateURI       ValidateURIHandler
	authorizeGenerate oauth2.AuthorizeGenerate
	accessGenerate    oauth2.AccessGenerate
	tokenStore        oauth2.TokenStore
	clientStore       oauth2.ClientStore
}

// get grant type config
func (m *Manager) grantConfig(gt oauth2.GrantType) *Config {
	if c, ok := m.gtcfg[gt]; ok && c != nil {
		return c
	}
	switch gt {
	case oauth2.AuthorizationCode:
		return DefaultAuthorizeCodeTokenCfg
	case oauth2.Implicit:
		return DefaultImplicitTokenCfg
	case oauth2.PasswordCredentials:
		return DefaultPasswordTokenCfg
	case oauth2.ClientCredentials:
		return DefaultClientTokenCfg
	}
	return &Config{}
}

// SetAuthorizeCodeExp set the authorization code expiration time
func (m *Manager) SetAuthorizeCodeExp(exp time.Duration) {
	m.codeExp = exp
}

// SetAuthorizeCodeTokenCfg set the authorization code grant token config
func (m *Manager) SetAuthorizeCodeTokenCfg(cfg *Config) {
	m.gtcfg[oauth2.AuthorizationCode] = cfg
}

// SetImplicitTokenCfg set the implicit grant token config
func (m *Manager) SetImplicitTokenCfg(cfg *Config) {
	m.gtcfg[oauth2.Implicit] = cfg
}

// SetPasswordTokenCfg set the password grant token config
func (m *Manager) SetPasswordTokenCfg(cfg *Config) {
	m.gtcfg[oauth2.PasswordCredentials] = cfg
}

// SetClientTokenCfg set the client grant token config
func (m *Manager) SetClientTokenCfg(cfg *Config) {
	m.gtcfg[oauth2.ClientCredentials] = cfg
}

// SetRefreshTokenCfg set the refreshing token config
func (m *Manager) SetRefreshTokenCfg(cfg *RefreshingConfig) {
	m.rcfg = cfg
}

// SetValidateURIHandler set the validates that RedirectURI is contained in baseURI
func (m *Manager) SetValidateURIHandler(handler ValidateURIHandler) {
	m.validateURI = handler
}

// MapAuthorizeGenerate mapping the authorize code generate interface
func (m *Manager) MapAuthorizeGenerate(gen oauth2.AuthorizeGenerate) {
	m.authorizeGenerate = gen
}

// MapAccessGenerate mapping the access token generate interface
func (m *Manager) MapAccessGenerate(gen oauth2.AccessGenerate) {
	m.accessGenerate = gen
}

// MapClientStorage mapping the client store interface
func (m *Manager) MapClientStorage(stor oauth2.ClientStore) {
	m.clientStore = stor
}

// MustClientStorage mandatory mapping the client store interface
func (m *Manager) MustClientStorage(stor oauth2.ClientStore, err error) {
	if err != nil {
		panic(err.Error())
	}
	m.clientStore = stor
}

// MapTokenStorage mapping the token store interface
func (m *Manager) MapTokenStorage(stor oauth2.TokenStore) {
	m.tokenStore = stor
}

// MustTokenStorage mandatory mapping the token store interface
func (m *Manager) MustTokenStorage(stor oauth2.TokenStore, err error) {
	if err != nil {
		panic(err)
	}
	m.tokenStore = stor
}

// GetClient get the client information
func (m *Manager) GetClient(clientID string) (cli oauth2.ClientInfo, err error) {
	cli, err = m.clientStore.GetByID(clientID)
	if err != nil {
		return
	} else if cli == nil {
		err = errors.ErrInvalidClient
	}
	return
}

// GenerateAuthToken generate the authorization token(code)
func (m *Manager) GenerateAuthToken(rt oauth2.ResponseType, tgr *oauth2.TokenGenerateRequest) (authToken oauth2.TokenInfo, err error) {
	cli, err := m.GetClient(tgr.ClientID)
	if err != nil {
		return
	} else if tgr.RedirectURI != "" {
		if verr := m.validateURI(cli.GetDomain(), tgr.RedirectURI); verr != nil {
			err = verr
			return
		}
	}

	ti := models.NewToken()
	ti.SetClientID(tgr.ClientID)
	ti.SetUserID(tgr.UserID)
	ti.SetRedirectURI(tgr.RedirectURI)
	ti.SetScope(tgr.Scope)

	createAt := time.Now()
	td := &oauth2.GenerateBasic{
		Client:    cli,
		UserID:    tgr.UserID,
		CreateAt:  createAt,
		TokenInfo: ti,
		Request:   tgr.Request,
	}
	switch rt {
	case oauth2.Code:
		codeExp := m.codeExp
		if codeExp == 0 {
			codeExp = DefaultCodeExp
		}
		ti.SetCodeCreateAt(createAt)
		ti.SetCodeExpiresIn(codeExp)
		if exp := tgr.AccessTokenExp; exp > 0 {
			ti.SetAccessExpiresIn(exp)
		}

		tv, terr := m.authorizeGenerate.Token(td)
		if terr != nil {
			err = terr
			return
		}
		ti.SetCode(tv)
	case oauth2.Token:
		// set access token expires
		icfg := m.grantConfig(oauth2.Implicit)
		aexp := icfg.AccessTokenExp
		if exp := tgr.AccessTokenExp; exp > 0 {
			aexp = exp
		}
		ti.SetAccessCreateAt(createAt)
		ti.SetAccessExpiresIn(aexp)

		if icfg.IsGenerateRefresh {
			ti.SetRefreshCreateAt(createAt)
			ti.SetRefreshExpiresIn(icfg.RefreshTokenExp)
		}

		tv, rv, terr := m.accessGenerate.Token(td, icfg.IsGenerateRefresh)
		if terr != nil {
			err = terr
			return
		}
		ti.SetAccess(tv)

		if rv != "" {
			ti.SetRefresh(rv)
		}
	}

	err = m.tokenStore.Create(ti)
	if err != nil {
		return
	}
	authToken = ti
	return
}

// get authorization code data
func (m *Manager) getAuthorizationCode(code string) (info oauth2.TokenInfo, err error) {
	ti, terr := m.tokenStore.GetByCode(code)
	if terr != nil {
		err = terr
		return
	} else if ti == nil || ti.GetCode() != code || ti.GetCodeCreateAt().Add(ti.GetCodeExpiresIn()).Before(time.Now()) {
		err = errors.ErrInvalidAuthorizeCode
		return
	}
	info = ti
	return
}

// delete authorization code data
func (m *Manager) delAuthorizationCode(code string) (err error) {
	err = m.tokenStore.RemoveByCode(code)
	return
}

// get and delete authorization code data
func (m *Manager) getAndDelAuthorizationCode(tgr *oauth2.TokenGenerateRequest) (info oauth2.TokenInfo, err error) {
	code := tgr.Code
	ti, err := m.getAuthorizationCode(code)
	if err != nil {
		return
	} else if ti.GetClientID() != tgr.ClientID {
		err = errors.ErrInvalidAuthorizeCode
		return
	} else if codeURI := ti.GetRedirectURI(); codeURI != "" && codeURI != tgr.RedirectURI {
		err = errors.ErrInvalidAuthorizeCode
		return
	}

	err = m.delAuthorizationCode(code)
	if err != nil {
		return
	}
	info = ti
	return
}

// GenerateAccessToken generate the access token
func (m *Manager) GenerateAccessToken(gt oauth2.GrantType, tgr *oauth2.TokenGenerateRequest) (accessToken oauth2.TokenInfo, err error) {
	cli, err := m.GetClient(tgr.ClientID)
	if err != nil {
		return
	} else if tgr.ClientSecret != cli.GetSecret() {
		err = errors.ErrInvalidClient
		return
	} else if tgr.RedirectURI != "" {
		if verr := m.validateURI(cli.GetDomain(), tgr.RedirectURI); verr != nil {
			err = verr
			return
		}
	}

	if gt == oauth2.AuthorizationCode {
		ti, verr := m.getAndDelAuthorizationCode(tgr)
		if verr != nil {
			err = verr
			return
		}
		tgr.UserID = ti.GetUserID()
		tgr.Scope = ti.GetScope()
		if exp := ti.GetAccessExpiresIn(); exp > 0 {
			tgr.AccessTokenExp = exp
		}
	}

	ti := models.NewToken()
	ti.SetClientID(tgr.ClientID)
	ti.SetUserID(tgr.UserID)
	ti.SetRedirectURI(tgr.RedirectURI)
	ti.SetScope(tgr.Scope)

	createAt := time.Now()
	ti.SetAccessCreateAt(createAt)

	// set access token expires
	gcfg := m.grantConfig(gt)
	aexp := gcfg.AccessTokenExp
	if exp := tgr.AccessTokenExp; exp > 0 {
		aexp = exp
	}
	ti.SetAccessExpiresIn(aexp)
	if gcfg.IsGenerateRefresh {
		ti.SetRefreshCreateAt(createAt)
		ti.SetRefreshExpiresIn(gcfg.RefreshTokenExp)
	}

	td := &oauth2.GenerateBasic{
		Client:    cli,
		UserID:    tgr.UserID,
		CreateAt:  createAt,
		TokenInfo: ti,
		Request:   tgr.Request,
	}

	av, rv, terr := m.accessGenerate.Token(td, gcfg.IsGenerateRefresh)
	if terr != nil {
		err = terr
		return
	}
	ti.SetAccess(av)

	if rv != "" {
		ti.SetRefresh(rv)
	}

	err = m.tokenStore.Create(ti)
	if err != nil {
		return
	}
	accessToken = ti

	return
}

// RefreshAccessToken refreshing an access token
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

	oldAccess, oldRefresh := ti.GetAccess(), ti.GetRefresh()

	td := &oauth2.GenerateBasic{
		Client:    cli,
		UserID:    ti.GetUserID(),
		CreateAt:  time.Now(),
		TokenInfo: ti,
		Request:   tgr.Request,
	}

	rcfg := DefaultRefreshTokenCfg
	if v := m.rcfg; v != nil {
		rcfg = v
	}

	ti.SetAccessCreateAt(td.CreateAt)
	if v := rcfg.AccessTokenExp; v > 0 {
		ti.SetAccessExpiresIn(v)
	}

	if v := rcfg.RefreshTokenExp; v > 0 {
		ti.SetRefreshExpiresIn(v)
	}

	if rcfg.IsResetRefreshTime {
		ti.SetRefreshCreateAt(td.CreateAt)
	}

	if scope := tgr.Scope; scope != "" {
		ti.SetScope(scope)
	}

	tv, rv, terr := m.accessGenerate.Token(td, rcfg.IsGenerateRefresh)
	if terr != nil {
		err = terr
		return
	}

	ti.SetAccess(tv)
	if rv != "" {
		ti.SetRefresh(rv)
	}

	if verr := m.tokenStore.Create(ti); verr != nil {
		err = verr
		return
	}

	if rcfg.IsRemoveAccess {
		// remove the old access token
		if verr := m.tokenStore.RemoveByAccess(oldAccess); verr != nil {
			err = verr
			return
		}
	}

	if rcfg.IsRemoveRefreshing && rv != "" {
		// remove the old refresh token
		if verr := m.tokenStore.RemoveByRefresh(oldRefresh); verr != nil {
			err = verr
			return
		}
	}

	accessToken = ti
	if rv == "" {
		accessToken.SetRefresh("")
		accessToken.SetRefreshCreateAt(time.Now())
		accessToken.SetRefreshExpiresIn(0)
	}

	return
}

// RemoveAccessToken use the access token to delete the token information
func (m *Manager) RemoveAccessToken(access string) (err error) {
	if access == "" {
		err = errors.ErrInvalidAccessToken
		return
	}
	err = m.tokenStore.RemoveByAccess(access)
	return
}

// RemoveRefreshToken use the refresh token to delete the token information
func (m *Manager) RemoveRefreshToken(refresh string) (err error) {
	if refresh == "" {
		err = errors.ErrInvalidAccessToken
		return
	}
	err = m.tokenStore.RemoveByRefresh(refresh)
	return
}

// LoadAccessToken according to the access token for corresponding token information
func (m *Manager) LoadAccessToken(access string) (info oauth2.TokenInfo, err error) {
	if access == "" {
		err = errors.ErrInvalidAccessToken
		return
	}

	ct := time.Now()
	ti, terr := m.tokenStore.GetByAccess(access)
	if terr != nil {
		err = terr
		return
	} else if ti == nil || ti.GetAccess() != access {
		err = errors.ErrInvalidAccessToken
		return
	} else if ti.GetRefresh() != "" && ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(ct) {
		err = errors.ErrExpiredRefreshToken
		return
	} else if ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).Before(ct) {
		err = errors.ErrExpiredAccessToken
		return
	}
	info = ti
	return
}

// LoadRefreshToken according to the refresh token for corresponding token information
func (m *Manager) LoadRefreshToken(refresh string) (info oauth2.TokenInfo, err error) {
	if refresh == "" {
		err = errors.ErrInvalidRefreshToken
		return
	}

	ti, terr := m.tokenStore.GetByRefresh(refresh)
	if terr != nil {
		err = terr
		return
	} else if ti == nil || ti.GetRefresh() != refresh {
		err = errors.ErrInvalidRefreshToken
		return
	} else if ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(time.Now()) {
		err = errors.ErrExpiredRefreshToken
		return
	}
	info = ti
	return
}
