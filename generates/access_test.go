package generates

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/oauth2.v2"
	"gopkg.in/oauth2.v2/models"
)

func TestAccess(t *testing.T) {
	Convey("Test Access Generate", t, func() {
		data := &oauth2.GenerateBasic{
			Client: &models.Client{
				ID:     "123456",
				Secret: "123456",
			},
			UserID:   "000000",
			CreateAt: time.Now(),
		}
		gen := NewAccessGenerate()
		access, refresh, err := gen.Token(data, true)
		So(err, ShouldBeNil)
		So(access, ShouldNotBeEmpty)
		So(refresh, ShouldNotBeEmpty)
	})
}
