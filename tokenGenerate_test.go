package oauth2_test

import (
	"testing"
	"time"

	"github.com/LyricTian/go.uuid"

	"gopkg.in/oauth2.v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTokenGenerate(t *testing.T) {
	cli := oauth2.DefaultClient{
		ClientID:     "123456",
		ClientSecret: "654321",
		ClientDomain: "http://www.lyric.name",
	}
	basicInfo := &oauth2.TokenBasicInfo{
		Client:   cli,
		TokenID:  uuid.NewV4().String(),
		UserID:   "999999",
		CreateAt: time.Now().Unix(),
	}
	Convey("Token generate test", t, func() {
		tokenGenerate := oauth2.NewDefaultTokenGenerate()
		Convey("Generate access token", func() {
			token, err := tokenGenerate.AccessToken(basicInfo)
			So(err, ShouldBeNil)
			_, _ = Println("\n [P]Access Token:" + token)
		})
		Convey("Generate refresh token", func() {
			token, err := tokenGenerate.RefreshToken(basicInfo)
			So(err, ShouldBeNil)
			_, _ = Println("\n [P]Refresh Token:" + token)
		})
	})
}
