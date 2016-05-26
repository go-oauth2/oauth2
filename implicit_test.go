package oauth2_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/oauth2.v1"

	"testing"
)

func TestImplicitManager(t *testing.T) {
	ClientHandle(func(cli oauth2.Client) {
		userID := "999999"
		oManager, err := oauth2.CreateDefaultOAuthManager(oauth2.NewMongoConfig(MongoURL, DBName), "", "", nil)
		if err != nil {
			t.Error(err)
		}
		Convey("Implicit Manager Test", t, func() {
			manager := oManager.GetImplicitManager()

			token, err := manager.GenerateToken(cli.ID(), userID, "http://www.example.com/cb", "all")
			So(err, ShouldBeNil)

			checkToken, err := oManager.CheckAccessToken(token.AccessToken)
			So(err, ShouldBeNil)
			So(checkToken.ClientID, ShouldEqual, cli.ID())
			So(checkToken.UserID, ShouldEqual, userID)
		})
	})
}
