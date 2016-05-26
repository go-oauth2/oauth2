package oauth2_test

import (
	"testing"

	"gopkg.in/oauth2.v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPasswordManager(t *testing.T) {
	ClientHandle(func(info oauth2.Client) {
		userID := "999999"
		oManager, err := oauth2.CreateDefaultOAuthManager(oauth2.NewMongoConfig(MongoURL, DBName), "", "", nil)
		if err != nil {
			t.Error(err)
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
