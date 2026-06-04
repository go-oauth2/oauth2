package store_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/store"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTokenStore(t *testing.T) {
	Convey("Test memory store", t, func() {
		store, err := store.NewMemoryTokenStore()
		So(err, ShouldBeNil)
		testToken(store)
	})

	Convey("Test file store", t, func() {
		os.Remove("data.db")

		store, err := store.NewFileTokenStore("data.db")
		So(err, ShouldBeNil)
		testToken(store)
	})
}

func testToken(store oauth2.TokenStore) {
	Convey("Test authorization code store", func() {
		ctx := context.Background()
		info := &models.Token{
			ClientID:      "1",
			UserID:        "1_1",
			RedirectURI:   "http://localhost/",
			Scope:         "all",
			Code:          "11_11_11",
			CodeCreateAt:  time.Now(),
			CodeExpiresIn: time.Second * 5,
		}
		err := store.Create(ctx, info)
		So(err, ShouldBeNil)

		cinfo, err := store.GetByCode(ctx, info.Code)
		So(err, ShouldBeNil)
		So(cinfo.GetUserID(), ShouldEqual, info.UserID)

		err = store.RemoveByCode(ctx, info.Code)
		So(err, ShouldBeNil)

		cinfo, err = store.GetByCode(ctx, info.Code)
		So(err, ShouldBeNil)
		So(cinfo, ShouldBeNil)
	})

	Convey("Test access token store", func() {
		ctx := context.Background()
		info := &models.Token{
			ClientID:        "1",
			UserID:          "1_1",
			RedirectURI:     "http://localhost/",
			Scope:           "all",
			Access:          "1_1_1",
			AccessCreateAt:  time.Now(),
			AccessExpiresIn: time.Second * 5,
		}
		err := store.Create(ctx, info)
		So(err, ShouldBeNil)

		ainfo, err := store.GetByAccess(ctx, info.GetAccess())
		So(err, ShouldBeNil)
		So(ainfo.GetUserID(), ShouldEqual, info.GetUserID())

		err = store.RemoveByAccess(ctx, info.GetAccess())
		So(err, ShouldBeNil)

		ainfo, err = store.GetByAccess(ctx, info.GetAccess())
		So(err, ShouldBeNil)
		So(ainfo, ShouldBeNil)
	})

	Convey("Test refresh token store", func() {
		ctx := context.Background()
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
		err := store.Create(ctx, info)
		So(err, ShouldBeNil)

		rinfo, err := store.GetByRefresh(ctx, info.GetRefresh())
		So(err, ShouldBeNil)
		So(rinfo.GetUserID(), ShouldEqual, info.GetUserID())

		err = store.RemoveByRefresh(ctx, info.GetRefresh())
		So(err, ShouldBeNil)

		rinfo, err = store.GetByRefresh(ctx, info.GetRefresh())
		So(err, ShouldBeNil)
		So(rinfo, ShouldBeNil)
	})

	Convey("Test TTL", func() {
		ctx := context.Background()
		info := &models.Token{
			ClientID:         "1",
			UserID:           "1_1",
			RedirectURI:      "http://localhost/",
			Scope:            "all",
			Access:           "1_3_1",
			AccessCreateAt:   time.Now(),
			AccessExpiresIn:  time.Second * 1,
			Refresh:          "1_3_2",
			RefreshCreateAt:  time.Now(),
			RefreshExpiresIn: time.Second * 1,
		}
		err := store.Create(ctx, info)
		So(err, ShouldBeNil)

		time.Sleep(time.Second * 1)
		ainfo, err := store.GetByAccess(ctx, info.Access)
		So(err, ShouldBeNil)
		So(ainfo, ShouldBeNil)
		rinfo, err := store.GetByRefresh(ctx, info.Refresh)
		So(err, ShouldBeNil)
		So(rinfo, ShouldBeNil)
	})
}
