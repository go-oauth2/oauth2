package oauth2_test

import (
	"testing"
	"time"

	"gopkg.in/oauth2.v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestACMemoryStore(t *testing.T) {
	Convey("AC memory store test", t, func() {
		store := oauth2.NewACMemoryStore(1)
		item := oauth2.ACInfo{
			ClientID:  "123456",
			UserID:    "999999",
			CreateAt:  time.Now().Unix(),
			ExpiresIn: time.Millisecond * 500,
		}
		Convey("Put Test", func() {
			id, err := store.Put(item)
			So(err, ShouldBeNil)
			So(id, ShouldEqual, 1)
			item.ID = id
			Convey("Take Test", func() {
				info, err := store.TakeByID(id)
				So(err, ShouldBeNil)
				So(info.ClientID, ShouldEqual, item.ClientID)
				So(info.UserID, ShouldEqual, item.UserID)
			})
			Convey("Take GC Test", func() {
				time.Sleep(time.Second * 2)
				info, err := store.TakeByID(id)
				So(err, ShouldNotBeNil)
				So(info, ShouldBeNil)
			})
		})
	})
}
