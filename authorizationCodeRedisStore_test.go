package oauth2

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestACRedisStore(t *testing.T) {
	Convey("Authorization code redis store test", t, func() {
		store, err := NewACRedisStore(&RedisConfig{
			Addr: "192.168.33.70:6379",
			DB:   1,
		}, "")
		So(err, ShouldBeNil)
		item := ACInfo{
			ClientID:  "123456",
			UserID:    "999999",
			Code:      "",
			CreateAt:  time.Now().Unix(),
			ExpiresIn: time.Millisecond * 500,
		}

		Convey("Put Test", func() {
			id, err := store.Put(item)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)
			Convey("Take Test", func() {
				info, err := store.TakeByID(id)
				So(err, ShouldBeNil)
				So(info.ClientID, ShouldEqual, item.ClientID)
				So(info.UserID, ShouldEqual, item.UserID)
			})
		})

		Convey("GC Test", func() {
			id, err := store.Put(item)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)
			Convey("Take GC Test", func() {
				time.Sleep(time.Millisecond * 1500)
				info, err := store.TakeByID(id)
				So(err, ShouldNotBeNil)
				So(info, ShouldBeNil)
			})
		})
	})
}
