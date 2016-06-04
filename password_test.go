package oauth2

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPasswordManager(t *testing.T) {
	ClientHandle(func(info Client) {
		userID := "999999"
		oManager, err := NewDefaultOAuthManager(nil, NewMongoConfig(MongoURL, DBName), "", "")
		if err != nil {
			t.Fatal(err)
		}

		Convey("Password Manager Test", t, func() {
			manager := oManager.GetPasswordManager()

			token, err := manager.GenerateToken(info.ID(), userID, info.Secret(), "all", true)
			So(err, ShouldBeNil)

			checkAT, err := oManager.CheckAccessToken(token.AccessToken)
			So(err, ShouldBeNil)
			So(checkAT.ClientID, ShouldEqual, info.ID())
			So(checkAT.UserID, ShouldEqual, userID)

			newAT, err := oManager.RefreshAccessToken(checkAT.RefreshToken, "")
			So(err, ShouldBeNil)
			So(newAT.AccessToken, ShouldNotEqual, checkAT.AccessToken)
		})
	})
}
