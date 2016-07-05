package manage

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/oauth2.v2"
	"gopkg.in/oauth2.v2/generates"
	"gopkg.in/oauth2.v2/models"
	"gopkg.in/oauth2.v2/store/client"
	"gopkg.in/oauth2.v2/store/token"
)

func TestManager(t *testing.T) {
	Convey("Manager test", t, func() {
		manager := NewManager()

		manager.MapClientModel(models.NewClient())
		manager.MapTokenModel(models.NewToken())
		manager.MapAuthorizeGenerate(generates.NewAuthorizeGenerate())
		manager.MapAccessGenerate(generates.NewAccessGenerate())
		manager.MapClientStorage(client.NewTempStore())

		Convey("GetClient test", func() {
			cli, err := manager.GetClient("1")
			So(err, ShouldBeNil)
			So(cli.GetSecret(), ShouldEqual, "11")
		})

		Convey("Redis store test", func() {
			manager.MustTokenStorage(token.NewRedisStore(
				&token.RedisConfig{Addr: "192.168.33.70:6379"},
			))
			testManager(manager)
		})

		Convey("MongoDB store test", func() {
			manager.MustTokenStorage(token.NewMongoStore(
				&token.MongoConfig{URL: "mongodb://admin:123456@192.168.33.70:27017"},
			))
			testManager(manager)
		})

	})
}

func testManager(manager oauth2.Manager) {
	reqParams := &oauth2.TokenGenerateRequest{
		ClientID:    "1",
		UserID:      "123456",
		RedirectURI: "http://localhost/oauth2",
		Scope:       "all",
	}
	code, err := manager.GenerateAuthToken(oauth2.Code, reqParams)
	So(err, ShouldBeNil)
	So(code, ShouldNotBeEmpty)

	atParams := &oauth2.TokenGenerateRequest{
		ClientID:          "1",
		RedirectURI:       "http://localhost/oauth2",
		Code:              code,
		IsGenerateRefresh: true,
	}
	accessToken, refreshToken, err := manager.GenerateAccessToken(oauth2.AuthorizationCodeCredentials, atParams)
	So(err, ShouldBeNil)
	So(accessToken, ShouldNotBeEmpty)
	So(refreshToken, ShouldNotBeEmpty)

	_, err = manager.LoadAccessToken(code)
	So(err, ShouldNotBeNil)

	ainfo, err := manager.LoadAccessToken(accessToken)
	So(err, ShouldBeNil)
	So(ainfo.GetClientID(), ShouldEqual, atParams.ClientID)

	rinfo, err := manager.LoadRefreshToken(refreshToken)
	So(err, ShouldBeNil)
	So(rinfo.GetClientID(), ShouldEqual, atParams.ClientID)

	refreshAT, err := manager.RefreshAccessToken(refreshToken, "owner")
	So(err, ShouldBeNil)
	So(refreshAT, ShouldNotBeEmpty)

	_, err = manager.LoadAccessToken(accessToken)
	So(err, ShouldNotBeNil)

	refreshAInfo, err := manager.LoadAccessToken(refreshAT)
	So(err, ShouldBeNil)
	So(refreshAInfo.GetScope(), ShouldEqual, "owner")

	err = manager.RemoveRefreshToken(refreshToken)
	So(err, ShouldBeNil)

	_, err = manager.LoadAccessToken(refreshAT)
	So(err, ShouldNotBeNil)

	_, err = manager.LoadRefreshToken(refreshToken)
	So(err, ShouldNotBeNil)
}
