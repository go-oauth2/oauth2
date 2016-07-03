package token

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/oauth2.v2"
	"gopkg.in/oauth2.v2/models"
)

func TestRedisStore(t *testing.T) {
	Convey("Test redis store", t, func() {
		cfg := &RedisConfig{
			Addr: "192.168.33.70:6379",
		}
		store, err := NewRedisStore(cfg)
		So(err, ShouldBeNil)

		info := &models.Token{
			ClientID:         "1",
			UserID:           "1_1",
			RedirectURI:      "http://localhost/",
			Scope:            "all",
			AuthType:         oauth2.Code.String(),
			Access:           "1_1_1",
			AccessCreateAt:   time.Now(),
			AccessExpiresIn:  time.Second * 10,
			Refresh:          "1_1_2",
			RefreshCreateAt:  time.Now(),
			RefreshExpiresIn: time.Minute * 1,
		}
		err = store.Create(info)
		So(err, ShouldBeNil)

		ainfo, err := store.GetByAccess(info.GetAccess())
		So(err, ShouldBeNil)
		So(ainfo.GetRefresh(), ShouldEqual, info.GetRefresh())

		err = store.RemoveByAccess(info.GetAccess())
		So(err, ShouldBeNil)

		ainfo, err = store.GetByAccess(info.GetAccess())
		So(err, ShouldBeNil)
		So(ainfo, ShouldBeNil)

		rinfo, err := store.GetByRefresh(info.GetRefresh())
		So(err, ShouldBeNil)
		So(rinfo.GetAccess(), ShouldEqual, info.GetAccess())

		err = store.RemoveByRefresh(info.GetRefresh())
		So(err, ShouldBeNil)

		rinfo, err = store.GetByRefresh(info.GetRefresh())
		So(err, ShouldBeNil)
		So(rinfo, ShouldBeNil)
	})
}

func TestRedisStoreAccessExpired(t *testing.T) {
	Convey("Test redis store access token expired", t, func() {
		cfg := &RedisConfig{
			Addr: "192.168.33.70:6379",
		}
		store, err := NewRedisStore(cfg)
		So(err, ShouldBeNil)
		info := &models.Token{
			ClientID:         "1",
			UserID:           "1_2",
			RedirectURI:      "http://localhost/",
			Scope:            "all",
			AuthType:         oauth2.Code.String(),
			Access:           "1_2_1",
			AccessCreateAt:   time.Now(),
			AccessExpiresIn:  time.Second * 1,
			Refresh:          "1_2_2",
			RefreshCreateAt:  time.Now(),
			RefreshExpiresIn: time.Second * 5,
		}
		err = store.Create(info)
		So(err, ShouldBeNil)

		time.Sleep(time.Second * 1)

		ainfo, err := store.GetByAccess(info.GetAccess())
		So(err, ShouldBeNil)
		So(ainfo, ShouldBeNil)

		rinfo, err := store.GetByRefresh(info.GetRefresh())
		So(err, ShouldBeNil)
		So(rinfo, ShouldNotBeNil)
	})
}

func TestRedisStoreRefreshExpired(t *testing.T) {
	Convey("Test redis store refresh token expired", t, func() {
		cfg := &RedisConfig{
			Addr: "192.168.33.70:6379",
		}
		store, err := NewRedisStore(cfg)
		So(err, ShouldBeNil)
		info := &models.Token{
			ClientID:         "1",
			UserID:           "1_3",
			RedirectURI:      "http://localhost/",
			Scope:            "all",
			AuthType:         oauth2.Code.String(),
			Access:           "1_3_1",
			AccessCreateAt:   time.Now(),
			AccessExpiresIn:  time.Second * 2,
			Refresh:          "1_3_2",
			RefreshCreateAt:  time.Now(),
			RefreshExpiresIn: time.Second * 1,
		}
		err = store.Create(info)
		So(err, ShouldBeNil)

		time.Sleep(time.Second * 1)

		ainfo, err := store.GetByAccess(info.GetAccess())
		So(err, ShouldBeNil)
		So(ainfo, ShouldBeNil)

		rinfo, err := store.GetByRefresh(info.GetRefresh())
		So(err, ShouldBeNil)
		So(rinfo, ShouldBeNil)
	})
}
