package models

import "errors"

type Endpoint struct {
	RegexPath       string `yaml:"regex_path"`
	ProxySetHeaders []struct {
		HeaderName  string `yaml:"header_name"`
		HeaderValue string `yaml:"header_value"`
	} `yaml:"proxy_set_headers"`
	ProxyPass *ProxyPass `yaml:"proxy_pass"`
	RateLimit float64    `yaml:"rate_limit"`
}

func (e *Endpoint) Validate() error {
	if e.RegexPath == "" {
		return errors.New("No servers[*].endpoint.regex_path found")
	}
	for _, proxy_set_headers := range e.ProxySetHeaders {
		if proxy_set_headers.HeaderName == "" {
			return errors.New("No header_name in proxy_set_headers")
		}
	}
	if e.ProxyPass != nil {
		return e.ProxyPass.Validate()
	}
	return nil
}
