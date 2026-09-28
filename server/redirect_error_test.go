package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

func TestRejectedRedirectURIIsNotUsedForErrorRedirect(t *testing.T) {
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	clients := store.NewClientStore()
	if err := clients.Set("client", &models.Client{ID: "client", Secret: "secret", Domain: "https://trusted.example"}); err != nil {
		t.Fatal(err)
	}
	manager.MapClientStorage(clients)

	srv := server.NewDefaultServer(manager)
	srv.SetUserAuthorizationHandler(func(http.ResponseWriter, *http.Request) (string, error) {
		return "user", nil
	})
	r := httptest.NewRequest(http.MethodGet, "/authorize?response_type=code&client_id=client&redirect_uri=https%3A%2F%2Fevil.example%2Fcallback", nil)
	w := httptest.NewRecorder()
	if err := srv.HandleAuthorizeRequest(w, r); err != errors.ErrInvalidRedirectURI {
		t.Fatalf("expected invalid redirect URI error, got %v", err)
	}
	if location := w.Header().Get("Location"); strings.Contains(location, "evil.example") {
		t.Fatalf("rejected redirect URI used for error redirect: status=%d location=%q", w.Code, location)
	}
}
