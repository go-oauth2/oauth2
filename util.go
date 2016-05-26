package oauth2

import (
	"errors"
	"net/url"
)

// ValidateURI 验证基础的Uri与重定向的URI是否一致
func ValidateURI(baseURI string, redirectURI string) error {
	base, err := url.Parse(baseURI)
	if err != nil {
		return err
	}
	redirect, err := url.Parse(redirectURI)
	if err != nil {
		return err
	}
	if base.Fragment != "" || redirect.Fragment != "" {
		return errors.New("Url must not include fragment.")
	}
	if base.Scheme != redirect.Scheme {
		return errors.New("Scheme don't match.")
	}
	if base.Host != redirect.Host {
		return errors.New("Host don't match.")
	}
	return nil
}
