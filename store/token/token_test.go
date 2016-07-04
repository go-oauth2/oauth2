package token

import (
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/oauth2.v2"
	"gopkg.in/oauth2.v2/models"
)

func testAccessStore(store oauth2.TokenStore) {
	info := &models.Token{
		ClientID:        "1",
		UserID:          "1_1",
		RedirectURI:     "http://localhost/",
		Scope:           "all",
		AuthType:        oauth2.Code.String(),
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
}

func testRefreshStore(store oauth2.TokenStore) {
	info := &models.Token{
		ClientID:         "1",
		UserID:           "1_2",
		RedirectURI:      "http://localhost/",
		Scope:            "all",
		AuthType:         oauth2.Code.String(),
		Access:           "1_2_1",
		AccessCreateAt:   time.Now(),
		AccessExpiresIn:  time.Second * 5,
		Refresh:          "1_2_2",
		RefreshCreateAt:  time.Now(),
		RefreshExpiresIn: time.Minute * 1,
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
}
