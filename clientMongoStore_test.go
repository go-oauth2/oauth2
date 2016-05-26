package oauth2_test

import (
	"testing"

	"gopkg.in/oauth2.v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClientMongoStore(t *testing.T) {
	ClientHandle(func(info oauth2.Client) {
		Convey("Client mongodb store test", t, func() {
			clientStore, err := oauth2.NewClientMongoStore(oauth2.NewMongoConfig(MongoURL, DBName), "")
			So(err, ShouldBeNil)
			client, err := clientStore.GetByID(info.ID())
			So(err, ShouldBeNil)
			So(client.Secret(), ShouldEqual, info.Secret())
		})
	})
}
