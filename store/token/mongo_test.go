package token

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	mongoURL = "mongodb://admin:123456@192.168.33.70:27017"
)

func TestMongoStore(t *testing.T) {
	Convey("Test mongo store", t, func() {
		cfg := &MongoConfig{
			URL: mongoURL,
		}
		store, err := NewMongoStore(cfg)
		So(err, ShouldBeNil)

		Convey("Test access token store", func() {
			testAccessStore(store)
		})

		Convey("Test refresh token store", func() {
			testRefreshStore(store)
		})
	})
}
