package generates_test

import (
	"testing"
	"time"

	"github.com/go-oauth2/oauth2"
	"github.com/go-oauth2/oauth2/generates"
	"github.com/go-oauth2/oauth2/models"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthorize(t *testing.T) {
	Convey("Test Authorize Generate", t, func() {
		data := &oauth2.GenerateBasic{
			Client: &models.Client{
				ID:     "123456",
				Secret: "123456",
			},
			UserID:   "000000",
			CreateAt: time.Now(),
		}
		gen := generates.NewAuthorizeGenerate()
		code, err := gen.Token(data)
		So(err, ShouldBeNil)
		So(code, ShouldNotBeEmpty)
		Println("\nAuthorize Code:" + code)
	})
}
