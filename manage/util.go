package manage

import (
	"net/url"
	"strings"

	"github.com/go-oauth2/oauth2/v4"

	"github.com/go-oauth2/oauth2/v4/errors"
)

type (
	// ValidateURIHandler validates that redirectURI is contained in baseURI
	ValidateURIHandler      func(redirectURIs []string, redirectURI string) error
	ExtractExtensionHandler func(*oauth2.TokenGenerateRequest, oauth2.ExtendableTokenInfo)
)

// DefaultValidateURI validates that redirectURI is contained in redirectURIs
func DefaultValidateURI(redirectURIs []string, redirectURI string) error {
	redirect, err := url.Parse(redirectURI)
	if err != nil {
		return err
	}

	for _, baseURI := range redirectURIs {
		base, err := url.Parse(baseURI)
		if err != nil {
			continue
		}

		if strings.EqualFold(base.Host, redirect.Host) {
			return nil
		}
	}

	return errors.ErrInvalidRedirectURI
}
