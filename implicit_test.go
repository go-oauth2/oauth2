package oauth2

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestImplicitManager(t *testing.T) {
	ClientHandle(func(cli Client) {
		userID := "999999"
		oManager, err := NewDefaultOAuthManager(nil, NewMongoConfig(MongoURL, DBName), "", "")
		if err != nil {
			t.Fatal(err)
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
