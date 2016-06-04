package oauth2

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTokenMongoStore(t *testing.T) {
	Convey("Token mongodb store test", t, func() {
		tokenStore, err := NewTokenMongoStore(NewMongoConfig(MongoURL, DBName), "")
		So(err, ShouldBeNil)
		createAt := time.Now().Unix()
		tokenValue := Token{
			ClientID:     "123456",
			UserID:       "999999",
			AccessToken:  "654321",
			ATCreateAt:   createAt,
			ATExpiresIn:  time.Second * 1,
			RefreshToken: "000000",
			RTCreateAt:   createAt,
			RTExpiresIn:  time.Second * 1,
			CreateAt:     createAt,
			Status:       Actived,
		}
		id, err := tokenStore.Create(&tokenValue)
		So(err, ShouldBeNil)
		So(id, ShouldBeGreaterThanOrEqualTo, 1)
		tokenValue.ID = id
		err = tokenStore.Update(id, map[string]interface{}{"Status": Expired})
		So(err, ShouldBeNil)
		at, err := tokenStore.GetByAccessToken("654321")
		So(err, ShouldBeNil)
		So(at.Status, ShouldEqual, Expired)
		rt, err := tokenStore.GetByRefreshToken("000000")
		So(err, ShouldBeNil)
		So(rt.ID, ShouldEqual, id)
	})
}
