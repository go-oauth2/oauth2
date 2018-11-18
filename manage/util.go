package manage

import (
	"net/url"
	"strings"

	"github.com/develm/oauth2/errors"
)

type (
	// ValidateURIHandler validates that redirectURI is contained in baseURI
	ValidateURIHandler func(baseURI, redirectURI string) (err error)
)

// DefaultValidateURI validates that redirectURI is contained in baseURI
func DefaultValidateURI(baseURI string, redirectURI string) (err error) {
	base, err := url.Parse(baseURI)
	if err != nil {
		return
	}
	redirect, err := url.Parse(redirectURI)
	if err != nil {
		return
	}
	if !strings.HasSuffix(redirect.Host, base.Host) {
		err = errors.ErrInvalidRedirectURI
	}
	return
}
