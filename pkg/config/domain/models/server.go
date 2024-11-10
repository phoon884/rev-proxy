package models

import "errors"

type Server struct {
	Name            string     `yaml:"name"`
	Port            int        `yaml:"port"`
	Http            bool       `yaml:"http"`
	ProxyTimeout    int        `yaml:"proxy_timeout"`
	ProxyBufferSize string     `yaml:"proxy_buffer_size"`
	Endpoints       []Endpoint `yaml:"endpoints"`
	TcpProxyPass    *ProxyPass `yaml:"tcp_proxy_pass"`
}

func (s *Server) Validate() error {
	if s.Name == "" {
		return errors.New("No servers[*].name found")
	}
	if s.Port == 0 {
		return errors.New("No servers[*].port found")
	}
	if s.Http == true && s.TcpProxyPass != nil {
		return errors.New("conflicting configuration[HTTP/tcp_proxy_pass]")
	}
	if s.Http == false && s.Endpoints != nil {
		return errors.New("conflicting configuration[HTTP/endpoint]")
	}
	if len(s.Endpoints) != 0 {
		for _, endpoint := range s.Endpoints {
			err := endpoint.Validate()
			if err != nil {
				return err
			}
		}
	}
	if s.TcpProxyPass != nil && s.TcpProxyPass.Healthcheck.Enabled {
		return errors.New("conflicting configuration[TcpProxyPass/Healthcheck]")
	}
	return nil
}
