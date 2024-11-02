package models

import "errors"

type ProxyPass struct {
	Addresses         []string `yaml:"addresses"`
	DownstreamTimeout int      `yaml:"downstream_timeout"`
	Loadbalancing     string   `yaml:"loadbalancing"`
	Healthcheck       struct {
		Enabled  bool `yaml:"enabled"`
		Interval int  `yaml:"interval"`
		Fails    int  `yaml:"fails"`
		Passes   int  `yaml:"passes"`
	} `yaml:"health_check"`
}

func (p *ProxyPass) Validate() error {
	if len(p.Addresses) == 0 {
		return errors.New("Proxy pass contains no Addresses")
	}
	if p.Loadbalancing != "random" && p.Loadbalancing != "hash" &&
		p.Loadbalancing != "least-connection" && p.Loadbalancing != "round-robin" {
		return errors.New("Loadbalancing method is unknown")
	}
	if p.Healthcheck.Enabled {
		if p.Healthcheck.Interval == 0 {
			return errors.New("health_check.Interval field is empty")
		}
		if p.Healthcheck.Fails == 0 {
			return errors.New("health_check.fails field is empty")
		}
		if p.Healthcheck.Passes == 0 {
			return errors.New("health_check.passes field is empty")
		}
	}

	return nil
}
