package oauth2_test

import (
	"testing"

	"gopkg.in/oauth2.v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUtil(t *testing.T) {
	Convey("ValidateURI Test", t, func() {
		err := oauth2.ValidateURI("http://www.example.com", "http://www.example.com/cb?code=xxx")
		So(err, ShouldBeNil)
	})
}
