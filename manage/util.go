package manage

import (
	"net/url"
	"strings"

	"github.com/daripadabengong/oauth2/v4"

	"github.com/daripadabengong/oauth2/v4/errors"
)

type (
	// ValidateURIHandler validates that redirectURI is contained in baseURI
	ValidateURIHandler      func(baseURI, redirectURI string) error
	ExtractExtensionHandler func(*oauth2.TokenGenerateRequest, oauth2.ExtendableTokenInfo)
)

// DefaultValidateURI validates that redirectURI is contained in baseURI
func DefaultValidateURI(baseURI string, redirectURI string) error {
	base, err := url.Parse(baseURI)
	if err != nil {
		return err
	}

	redirect, err := url.Parse(redirectURI)
	if err != nil {
		return err
	}
	if !strings.HasSuffix(redirect.Host, base.Host) {
		return errors.ErrInvalidRedirectURI
	}
	return nil
}
