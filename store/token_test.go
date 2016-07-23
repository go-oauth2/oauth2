package store_test

import (
	"testing"
	"time"

	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/store"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTokenStore(t *testing.T) {
	Convey("Test memory store", t, func() {
		store := store.NewMemoryTokenStore(time.Second * 1)

		Convey("Test access token store", func() {
			info := &models.Token{
				ClientID:        "1",
				UserID:          "1_1",
				RedirectURI:     "http://localhost/",
				Scope:           "all",
				Access:          "1_1_1",
				AccessCreateAt:  time.Now(),
				AccessExpiresIn: time.Second * 5,
			}
			err := store.Create(info)
			So(err, ShouldBeNil)

			ainfo, err := store.GetByAccess(info.GetAccess())
			So(err, ShouldBeNil)
			So(ainfo.GetUserID(), ShouldEqual, info.GetUserID())

			err = store.RemoveByAccess(info.GetAccess())
			So(err, ShouldBeNil)

			ainfo, err = store.GetByAccess(info.GetAccess())
			So(err, ShouldBeNil)
			So(ainfo, ShouldBeNil)
		})

		Convey("Test refresh token store", func() {
			info := &models.Token{
				ClientID:         "1",
				UserID:           "1_2",
				RedirectURI:      "http://localhost/",
				Scope:            "all",
				Access:           "1_2_1",
				AccessCreateAt:   time.Now(),
				AccessExpiresIn:  time.Second * 5,
				Refresh:          "1_2_2",
				RefreshCreateAt:  time.Now(),
				RefreshExpiresIn: time.Second * 15,
			}
			err := store.Create(info)
			So(err, ShouldBeNil)

			rinfo, err := store.GetByRefresh(info.GetRefresh())
			So(err, ShouldBeNil)
			So(rinfo.GetUserID(), ShouldEqual, info.GetUserID())

			err = store.RemoveByRefresh(info.GetRefresh())
			So(err, ShouldBeNil)

			rinfo, err = store.GetByRefresh(info.GetRefresh())
			So(err, ShouldBeNil)
			So(rinfo, ShouldBeNil)
		})

		Convey("Test gc", func() {
			info := &models.Token{
				ClientID:        "1",
				UserID:          "1_3",
				RedirectURI:     "http://localhost/",
				Scope:           "all",
				Access:          "1_3_1",
				AccessCreateAt:  time.Now(),
				AccessExpiresIn: time.Second * 1,
			}
			err := store.Create(info)
			So(err, ShouldBeNil)

			time.Sleep(time.Second * 1)
			ainfo, err := store.GetByRefresh(info.GetAccess())
			So(err, ShouldBeNil)
			So(ainfo, ShouldBeNil)
		})
	})
}
