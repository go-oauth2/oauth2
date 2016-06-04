package oauth2

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestCCManager(t *testing.T) {
	ClientHandle(func(cli Client) {
		oManager, err := NewDefaultOAuthManager(nil, NewMongoConfig(MongoURL, DBName), "", "")
		if err != nil {
			t.Fatal(err)
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
