package manage_test

import (
	"testing"

	"github.com/go-oauth2/oauth2/v4/manage"

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

func TestDefaultValidateURIRejectsHostSuffixSpoof(t *testing.T) {
	if err := manage.DefaultValidateURI("https://whitelist.com", "https://evil-whitelist.com/callback"); err == nil {
		t.Fatal("expected a host suffix spoof to be rejected")
	}
}

func TestDefaultValidateURIAllowsExactHostAndSubdomain(t *testing.T) {
	for _, redirectURI := range []string{
		"https://whitelist.com/callback",
		"https://login.whitelist.com/callback",
		"https://LOGIN.WHITELIST.COM/callback",
	} {
		if err := manage.DefaultValidateURI("https://whitelist.com", redirectURI); err != nil {
			t.Errorf("expected %q to be accepted: %v", redirectURI, err)
		}
	}
}
