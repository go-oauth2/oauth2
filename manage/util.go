package manage

import (
	"net/url"
	"strings"

	"github.com/go-errors/errors"
)

// ValidateURI Validates that RedirectURI is contained in domain
func ValidateURI(domain string, redirectURI string) (err error) {
	base, err := url.Parse(domain)
	if err != nil {
		return
	}
	redirect, err := url.Parse(redirectURI)
	if err != nil {
		return
	}
	if !strings.HasSuffix(redirect.Host, base.Host) {
		err = errors.New(ErrInvalidRedirectURI)
	}
	return
}
