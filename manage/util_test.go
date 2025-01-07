package manage_test

import (
	"testing"

	"github.com/daripadabengong/oauth2/v4/manage"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUtil(t *testing.T) {
	Convey("Util Test", t, func() {
		Convey("ValidateURI Test", func() {
			err := manage.DefaultValidateURI("http://www.example.com", "http://www.example.com/cb?code=xxx")
			So(err, ShouldBeNil)
		})
	})
}
