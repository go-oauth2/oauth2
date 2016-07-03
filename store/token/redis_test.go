package token

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRedisStore(t *testing.T) {
	Convey("Test redis store", t, func() {
		cfg := &RedisConfig{
			Addr: "192.168.33.70:6379",
		}
		store, err := NewRedisStore(cfg)
		So(err, ShouldBeNil)

		Convey("Test redis store access", func() {
			testAccessStore(store)
		})

		Convey("Test redis store refresh", func() {
			testRefreshStore(store)
		})
	})
}

func TestRedisStoreAccessExpired(t *testing.T) {
	Convey("Test redis store access token expired", t, func() {
		cfg := &RedisConfig{
			Addr: "192.168.33.70:6379",
		}
		store, err := NewRedisStore(cfg)
		So(err, ShouldBeNil)
		testAccessExpired(store)
	})
}

func TestRedisStoreRefreshExpired(t *testing.T) {
	Convey("Test redis store refresh token expired", t, func() {
		cfg := &RedisConfig{
			Addr: "192.168.33.70:6379",
		}
		store, err := NewRedisStore(cfg)
		So(err, ShouldBeNil)
		testRefreshExpired(store)
	})
}
