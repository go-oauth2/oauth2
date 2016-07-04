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

		Convey("Test access token store", func() {
			testAccessStore(store)
		})

		Convey("Test refresh token store", func() {
			testRefreshStore(store)
		})
	})
}
