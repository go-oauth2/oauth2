package oauth2_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/oauth2.v1"

	"testing"
)

func TestCCManager(t *testing.T) {
	ClientHandle(func(cli oauth2.Client) {
		oManager, err := oauth2.CreateDefaultOAuthManager(oauth2.NewMongoConfig(MongoURL, DBName), "", "", nil)
		if err != nil {
			t.Error(err)
		}
		Convey("Client Credentials Manager Test", t, func() {
			manager := oManager.GetCCManager()

			token, err := manager.GenerateToken(cli.ID(), cli.Secret(), "all")
			So(err, ShouldBeNil)

			checkToken, err := oManager.CheckAccessToken(token.AccessToken)
			So(err, ShouldBeNil)
			So(checkToken.ClientID, ShouldEqual, cli.ID())
		})
	})
}
