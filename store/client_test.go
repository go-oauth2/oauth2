package store_test

import (
	"testing"

	"github.com/go-oauth2/oauth2/models"
	"github.com/go-oauth2/oauth2/store"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClientStore(t *testing.T) {
	Convey("Test client store", t, func() {
		clientStore := store.NewClientStore()

		err := clientStore.Set("1", &models.Client{ID: "1", Secret: "2"})
		So(err, ShouldBeNil)

		cli, err := clientStore.GetByID("1")
		So(err, ShouldBeNil)
		So(cli.GetID(), ShouldEqual, "1")
	})
}
