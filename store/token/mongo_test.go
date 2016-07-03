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

		Convey("Test mongo store access", func() {
			testAccessStore(store)
		})

		Convey("Test mongo store refresh", func() {
			testRefreshStore(store)
		})
	})
}

func TestMongoStoreAccessExpired(t *testing.T) {
	Convey("Test mongo store access token expired", t, func() {
		cfg := &MongoConfig{
			URL: mongoURL,
		}
		store, err := NewMongoStore(cfg)
		So(err, ShouldBeNil)
		testAccessExpired(store)
	})
}

func TestMongoStoreRefreshExpired(t *testing.T) {
	Convey("Test mongo store refresh token expired", t, func() {
		cfg := &MongoConfig{
			URL: mongoURL,
		}
		store, err := NewMongoStore(cfg)
		So(err, ShouldBeNil)
		testRefreshExpired(store)
	})
}
