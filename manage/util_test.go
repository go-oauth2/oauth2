package manage

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUtil(t *testing.T) {
	Convey("Util Test", t, func() {
		Convey("ValidateURI Test", func() {
			err := ValidateURI("http://www.example.com", "http://www.example.com/cb?code=xxx")
			So(err, ShouldBeNil)
		})
	})
}
