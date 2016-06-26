package manage

import (
	"errors"
	"net/url"
)

// ValidateURI 校验重定向的URI与域名的一致性
func ValidateURI(domain string, redirectURI string) error {
	base, err := url.Parse(domain)
	if err != nil {
		return err
	}
	redirect, err := url.Parse(redirectURI)
	if err != nil {
		return err
	} else if base.Fragment != "" || redirect.Fragment != "" {
		return errors.New("Url must not include fragment.")
	} else if base.Scheme != redirect.Scheme {
		return errors.New("Scheme don't match.")
	} else if base.Host != redirect.Host {
		return errors.New("Host don't match.")
	}
	return nil
}
