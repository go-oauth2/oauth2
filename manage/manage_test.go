package manage_test

import (
	"testing"
	"time"

	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/store"

	. "github.com/smartystreets/goconvey/convey"
)

func TestManager(t *testing.T) {
	Convey("Manager test", t, func() {
		manager := manage.NewDefaultManager()

		manager.MustTokenStorage(store.NewMemoryTokenStore())

		clientStore := store.NewClientStore()
		_ = clientStore.Set("1", &models.Client{
			ID:     "1",
			Secret: "11",
			Domain: "http://localhost",
		})
		manager.MapClientStorage(clientStore)

		tgr := &oauth2.TokenGenerateRequest{
			ClientID:    "1",
			UserID:      "123456",
			RedirectURI: "http://localhost/oauth2",
			Scope:       "all",
		}

		Convey("GetClient test", func() {
			cli, err := manager.GetClient("1")
			So(err, ShouldBeNil)
			So(cli.GetSecret(), ShouldEqual, "11")
		})

		Convey("Token test", func() {
			testManager(tgr, manager)
		})

		Convey("zero expiration access token test", func() {
			testZeroAccessExpirationManager(tgr, manager)
			testCannotRequestZeroExpirationAccessTokens(tgr, manager)
		})

		Convey("zero expiration refresh token test", func() {
			testZeroRefreshExpirationManager(tgr, manager)
		})
	})
}

func testManager(tgr *oauth2.TokenGenerateRequest, manager oauth2.Manager) {
	cti, err := manager.GenerateAuthToken(oauth2.Code, tgr)
	So(err, ShouldBeNil)

	code := cti.GetCode()
	So(code, ShouldNotBeEmpty)

	atParams := &oauth2.TokenGenerateRequest{
		ClientID:     tgr.ClientID,
		ClientSecret: "11",
		RedirectURI:  tgr.RedirectURI,
		Code:         code,
	}
	ati, err := manager.GenerateAccessToken(oauth2.AuthorizationCode, atParams)
	So(err, ShouldBeNil)

	accessToken, refreshToken := ati.GetAccess(), ati.GetRefresh()
	So(accessToken, ShouldNotBeEmpty)
	So(refreshToken, ShouldNotBeEmpty)

	ainfo, err := manager.LoadAccessToken(accessToken)
	So(err, ShouldBeNil)
	So(ainfo.GetClientID(), ShouldEqual, atParams.ClientID)

	arinfo, err := manager.LoadRefreshToken(accessToken)
	So(err, ShouldNotBeNil)
	So(arinfo, ShouldBeNil)

	rainfo, err := manager.LoadAccessToken(refreshToken)
	So(err, ShouldNotBeNil)
	So(rainfo, ShouldBeNil)

	rinfo, err := manager.LoadRefreshToken(refreshToken)
	So(err, ShouldBeNil)
	So(rinfo.GetClientID(), ShouldEqual, atParams.ClientID)

	atParams.Refresh = refreshToken
	atParams.Scope = "owner"
	rti, err := manager.RefreshAccessToken(atParams)
	So(err, ShouldBeNil)

	refreshAT := rti.GetAccess()
	So(refreshAT, ShouldNotBeEmpty)

	_, err = manager.LoadAccessToken(accessToken)
	So(err, ShouldNotBeNil)

	refreshAInfo, err := manager.LoadAccessToken(refreshAT)
	So(err, ShouldBeNil)
	So(refreshAInfo.GetScope(), ShouldEqual, "owner")

	err = manager.RemoveAccessToken(refreshAT)
	So(err, ShouldBeNil)

	_, err = manager.LoadAccessToken(refreshAT)
	So(err, ShouldNotBeNil)

	err = manager.RemoveRefreshToken(refreshToken)
	So(err, ShouldBeNil)

	_, err = manager.LoadRefreshToken(refreshToken)
	So(err, ShouldNotBeNil)
}

func testZeroAccessExpirationManager(tgr *oauth2.TokenGenerateRequest, manager oauth2.Manager) {
	config := manage.Config{
		AccessTokenExp:    0, // Set explicitly as we're testing 0 (no) expiration
		IsGenerateRefresh: true,
	}
	m, ok := manager.(*manage.Manager)
	So(ok, ShouldBeTrue)
	m.SetAuthorizeCodeTokenCfg(&config)

	cti, err := manager.GenerateAuthToken(oauth2.Code, tgr)
	So(err, ShouldBeNil)

	code := cti.GetCode()
	So(code, ShouldNotBeEmpty)

	atParams := &oauth2.TokenGenerateRequest{
		ClientID:     tgr.ClientID,
		ClientSecret: "11",
		RedirectURI:  tgr.RedirectURI,
		Code:         code,
	}
	ati, err := manager.GenerateAccessToken(oauth2.AuthorizationCode, atParams)
	So(err, ShouldBeNil)

	accessToken, refreshToken := ati.GetAccess(), ati.GetRefresh()
	So(accessToken, ShouldNotBeEmpty)
	So(refreshToken, ShouldNotBeEmpty)

	tokenInfo, err := manager.LoadAccessToken(accessToken)
	So(err, ShouldBeNil)
	So(tokenInfo, ShouldNotBeNil)
	So(tokenInfo.GetAccess(), ShouldEqual, accessToken)
	So(tokenInfo.GetAccessExpiresIn(), ShouldEqual, 0)
}

func testCannotRequestZeroExpirationAccessTokens(tgr *oauth2.TokenGenerateRequest, manager oauth2.Manager) {
	config := manage.Config{
		AccessTokenExp: time.Hour * 5,
	}
	m, ok := manager.(*manage.Manager)
	So(ok, ShouldBeTrue)
	m.SetAuthorizeCodeTokenCfg(&config)

	cti, err := manager.GenerateAuthToken(oauth2.Code, tgr)
	So(err, ShouldBeNil)

	code := cti.GetCode()
	So(code, ShouldNotBeEmpty)

	atParams := &oauth2.TokenGenerateRequest{
		ClientID:       tgr.ClientID,
		ClientSecret:   "11",
		RedirectURI:    tgr.RedirectURI,
		AccessTokenExp: 0, // requesting token without expiration
		Code:           code,
	}
	ati, err := manager.GenerateAccessToken(oauth2.AuthorizationCode, atParams)
	So(err, ShouldBeNil)

	accessToken := ati.GetAccess()
	So(accessToken, ShouldNotBeEmpty)
	So(ati.GetAccessExpiresIn(), ShouldEqual, time.Hour*5)
}

func testZeroRefreshExpirationManager(tgr *oauth2.TokenGenerateRequest, manager oauth2.Manager) {
	config := manage.Config{
		RefreshTokenExp:   0, // Set explicitly as we're testing 0 (no) expiration
		IsGenerateRefresh: true,
	}
	m, ok := manager.(*manage.Manager)
	So(ok, ShouldBeTrue)
	m.SetAuthorizeCodeTokenCfg(&config)

	cti, err := manager.GenerateAuthToken(oauth2.Code, tgr)
	So(err, ShouldBeNil)

	code := cti.GetCode()
	So(code, ShouldNotBeEmpty)

	atParams := &oauth2.TokenGenerateRequest{
		ClientID:       tgr.ClientID,
		ClientSecret:   "11",
		RedirectURI:    tgr.RedirectURI,
		AccessTokenExp: time.Hour,
		Code:           code,
	}
	ati, err := manager.GenerateAccessToken(oauth2.AuthorizationCode, atParams)
	So(err, ShouldBeNil)

	accessToken, refreshToken := ati.GetAccess(), ati.GetRefresh()
	So(accessToken, ShouldNotBeEmpty)
	So(refreshToken, ShouldNotBeEmpty)

	tokenInfo, err := manager.LoadRefreshToken(refreshToken)
	So(err, ShouldBeNil)
	So(tokenInfo, ShouldNotBeNil)
	So(tokenInfo.GetRefresh(), ShouldEqual, refreshToken)
	So(tokenInfo.GetRefreshExpiresIn(), ShouldEqual, 0)

	// LoadAccessToken also checks refresh expiry
	tokenInfo, err = manager.LoadAccessToken(accessToken)
	So(err, ShouldBeNil)
	So(tokenInfo, ShouldNotBeNil)
	So(tokenInfo.GetRefresh(), ShouldEqual, refreshToken)
	So(tokenInfo.GetRefreshExpiresIn(), ShouldEqual, 0)
}
