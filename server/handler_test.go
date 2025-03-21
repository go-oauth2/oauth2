package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-oauth2/oauth2/v4/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRefreshTokenFormResolveHandler(t *testing.T) {
	Convey("Correct Request", t, func() {
		f := url.Values{}
		f.Add("refresh_token", "test_token")

		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		token, err := RefreshTokenFormResolveHandler(r)
		So(err, ShouldBeNil)
		So(token, ShouldEqual, "test_token")
	})

	Convey("Missing Refresh Token", t, func() {
		r := httptest.NewRequest("POST", "/", nil)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		token, err := RefreshTokenFormResolveHandler(r)
		So(err, ShouldBeError, errors.ErrInvalidRequest)
		So(token, ShouldBeEmpty)
	})
}

func TestRefreshTokenCookieResolveHandler(t *testing.T) {
	Convey("Correct Request", t, func() {
		r := httptest.NewRequest(http.MethodPost, "/", nil)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.AddCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    "test_token",
			HttpOnly: true,
			Path:     "/",
			Domain:   ".example.com",
			Expires:  time.Now().Add(time.Hour),
		})

		token, err := RefreshTokenCookieResolveHandler(r)
		So(err, ShouldBeNil)
		So(token, ShouldEqual, "test_token")
	})

	Convey("Missing Refresh Token", t, func() {
		r := httptest.NewRequest("POST", "/", nil)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		token, err := RefreshTokenCookieResolveHandler(r)
		So(err, ShouldBeError, errors.ErrInvalidRequest)
		So(token, ShouldBeEmpty)
	})
}
