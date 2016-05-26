package oauth2_test

import (
	"testing"

	"gopkg.in/oauth2.v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestACManager(t *testing.T) {
	ClientHandle(func(info oauth2.Client) {
		userID := "999999"
		oManager, err := oauth2.CreateDefaultOAuthManager(oauth2.NewMongoConfig(MongoURL, DBName), "", "", nil)
		if err != nil {
			t.Error(err)
		}
		Convey("Authorization Code Manager Test", t, func() {
			manager := oManager.GetACManager()

			redirectURI := "http://www.example.com/cb"
			code, err := manager.GenerateCode(info.ID(), userID, redirectURI, "all")
			So(err, ShouldBeNil)

			accessToken, err := manager.GenerateToken(code, redirectURI, info.ID(), info.Secret(), true)
			So(err, ShouldBeNil)
			So(accessToken.UserID, ShouldEqual, userID)

			checkAT, err := oManager.CheckAccessToken(accessToken.AccessToken)
			So(err, ShouldBeNil)
			So(checkAT.ClientID, ShouldEqual, info.ID())
			So(checkAT.UserID, ShouldEqual, userID)

			newAT, err := oManager.RefreshAccessToken(checkAT.RefreshToken, "")
			So(err, ShouldBeNil)
			So(newAT.AccessToken, ShouldNotEqual, checkAT.AccessToken)

			err = oManager.RevokeAccessToken(newAT.AccessToken)
			So(err, ShouldBeNil)

			checkAT, err = oManager.CheckAccessToken(newAT.AccessToken)
			So(err, ShouldNotBeNil)
			So(checkAT, ShouldBeNil)
		})
	})
}
