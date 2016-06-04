package oauth2

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClientMongoStore(t *testing.T) {
	ClientHandle(func(info Client) {
		Convey("Client mongodb store test", t, func() {
			clientStore, err := NewClientMongoStore(NewMongoConfig(MongoURL, DBName), "")
			So(err, ShouldBeNil)
			client, err := clientStore.GetByID(info.ID())
			So(err, ShouldBeNil)
			So(client.Secret(), ShouldEqual, info.Secret())
		})
	})
}
