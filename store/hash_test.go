package store_test

import (
	"context"
	"testing"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/store"
	. "github.com/smartystreets/goconvey/convey"
)

func TestClientStoreWithHash(t *testing.T) {
	Convey("Test client store with hash - save", t, func() {
		hasher := &store.BcryptHasher{}
		memory := store.NewClientStore()
		store := store.NewClientStoreWithHash(memory, hasher)
		secret := "123456"
		err := store.Save(context.Background(), &models.Client{
			ID:     "123",
			Secret: secret,
			Domain: "http://localhost",
			Public: false,
			UserID: "123",
		})
		So(err, ShouldBeNil)

		Convey("get by id", func() {
			storedClient, err := store.GetByID(context.Background(), "123")

			So(err, ShouldBeNil)
			So(storedClient.GetID(), ShouldEqual, "123")
			So(storedClient.GetSecret(), ShouldNotEqual, secret)

			verifier := storedClient.(oauth2.ClientPasswordVerifier)

			Convey("verify correct password - success", func() {
				So(verifier.VerifyPassword(secret), ShouldBeTrue)
			})

			Convey("verify incorrect password - fail", func() {
				So(verifier.VerifyPassword("wrong"), ShouldBeFalse)
			})
		})
	})
}

// check interfaces

var _ = (oauth2.ClientStore)((*store.ClientStoreWithHash)(nil))
var _ = (oauth2.SavingClientStore)((*store.ClientStoreWithHash)(nil))

var _ = (oauth2.ClientInfo)((*store.ClientInfoWithHash)(nil))
var _ = (oauth2.ClientPasswordVerifier)((*store.ClientInfoWithHash)(nil))
