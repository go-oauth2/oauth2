package oauth2

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUtil(t *testing.T) {
	Convey("ValidateURI Test", t, func() {
		err := ValidateURI("http://www.example.com", "http://www.example.com/cb?code=xxx")
		So(err, ShouldBeNil)
	})
}
