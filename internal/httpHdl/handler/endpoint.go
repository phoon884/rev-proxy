package handler

import (
	"net"
	"regexp"

	"github.com/phoon884/rev-proxy/internal/httpHdl/healthcheck"
	"github.com/phoon884/rev-proxy/internal/httpHdl/models"
	ratelimit "github.com/phoon884/rev-proxy/internal/httpHdl/ratelimit/application"
	lbPort "github.com/phoon884/rev-proxy/internal/loadbalancing/ports"
)

type Endpoint struct {
	regexPath       string
	loadbalancer    lbPort.Loadbalancer
	proxySetHeaders map[string]string
	healthChecker   *healthcheck.HttpHealthcheck
	DownStreamAddr  []string
	Ratelimter      *ratelimit.Ratelimiter
}

func (e *Endpoint) setHeader(request *models.HTTPReq) {
	for key, value := range e.proxySetHeaders {
		request.ChangeHeader(key, value)
	}
}

func (e *Endpoint) ComparePath(request *models.HTTPReq) (bool, error) {
	return regexp.Match(e.regexPath, []byte(request.Path))
}

func (e *Endpoint) GetDownStreamAddr(conn net.Conn) string {
	if len(e.DownStreamAddr) == 0 {
		return ""
	}
	addrs := make([]string, 0, len(e.DownStreamAddr))
	if e.healthChecker != nil {
		for _, addr := range e.DownStreamAddr {
			if e.healthChecker.GetStatus(addr) {
				addrs = append(addrs, addr)
			}
		}
		if len(addrs) == 0 {
			return ""
		}
		return e.loadbalancer.SelectDownstreamAddr(conn, addrs)
	} else {
		return e.loadbalancer.SelectDownstreamAddr(conn, e.DownStreamAddr)
	}
}

func NewEndpoint(
	regexPath string,
	loadbalancer lbPort.Loadbalancer,
	proxySetHeaders map[string]string,
	downStreamAddr []string,
	healthChecker *healthcheck.HttpHealthcheck,
	ratelimiter *ratelimit.Ratelimiter,
) *Endpoint {
	return &Endpoint{
		regexPath:       regexPath,
		loadbalancer:    loadbalancer,
		proxySetHeaders: proxySetHeaders,
		DownStreamAddr:  downStreamAddr,
		healthChecker:   healthChecker,
		Ratelimter:      ratelimiter,
	}
}
