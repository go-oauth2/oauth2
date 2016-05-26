package oauth2_test

import (
	"testing"
	"time"

	"gopkg.in/LyricTian/lib.v2"
	"gopkg.in/oauth2.v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestACGenerate(t *testing.T) {
	Convey("Authorization code generate test", t, func() {
		acGenerate := oauth2.NewDefaultACGenerate()
		info := &oauth2.ACInfo{
			ID:       1,
			ClientID: "123456",
			UserID:   "999999",
			Code:     lib.NewRandom(6).NumberAndLetter(),
			CreateAt: time.Now().Unix(),
		}
		Convey("Generate code", func() {
			code, err := acGenerate.Code(info)
			So(err, ShouldBeNil)
			So(code, ShouldNotBeBlank)
			Convey("Parse code", func() {
				id, err := acGenerate.Parse(code)
				So(err, ShouldBeNil)
				So(id, ShouldEqual, 1)
			})
			Convey("Verify code", func() {
				valid, err := acGenerate.Verify(code, info)
				So(err, ShouldBeNil)
				So(valid, ShouldBeTrue)
			})
		})
	})
}
