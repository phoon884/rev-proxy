package factories

import (
	"sync"

	httpHdl "github.com/phoon884/rev-proxy/internal/httpHdl/handler"
	"github.com/phoon884/rev-proxy/internal/httpHdl/healthcheck"
	ratelimitApp "github.com/phoon884/rev-proxy/internal/httpHdl/ratelimit/application"
	ratelimitRepo "github.com/phoon884/rev-proxy/internal/httpHdl/ratelimit/repositories"
	"github.com/phoon884/rev-proxy/internal/loadbalancing/loadbalancer"
	lbPort "github.com/phoon884/rev-proxy/internal/loadbalancing/ports"
	hdlPort "github.com/phoon884/rev-proxy/internal/ports"
	configApp "github.com/phoon884/rev-proxy/pkg/config/application"
	configModel "github.com/phoon884/rev-proxy/pkg/config/domain/models"
	configRepo "github.com/phoon884/rev-proxy/pkg/config/repositories"
	loggerApp "github.com/phoon884/rev-proxy/pkg/logger/application"
	"github.com/phoon884/rev-proxy/pkg/server"
)

type Factory struct {
	configFilePath string

	logger       *loggerApp.Logger
	configurator *configApp.ConfigService
}

func NewFactory(configFilePath string) *Factory {
	return &Factory{
		configFilePath: configFilePath,
	}
}

func (f *Factory) InitializeConfigurator() *configApp.ConfigService {
	if f.configurator == nil {
		path := f.configFilePath

		repo := configRepo.NewFileRepository(path)
		app := configApp.NewConfigService(repo)
		err := app.Config()
		if err != nil {
			panic(err)
		}
		f.configurator = app
		return app
	}
	return f.configurator
}

func (f *Factory) InitializeLogger(logLevel string) *loggerApp.Logger {
	if f.logger == nil {
		app := loggerApp.NewLogger(logLevel)
		f.logger = app
	}
	return f.logger
}

func (f *Factory) ServerBuilder(wg *sync.WaitGroup, serverCfg configModel.Server) *server.Server {
	logger := f.InitializeLogger("")
	downstreamConnTable := new(sync.Map)
	var hdl hdlPort.Handler
	if serverCfg.Http {
		endpoints := make([]*httpHdl.Endpoint, len(serverCfg.Endpoints))
		for idx, endpointsCfg := range serverCfg.Endpoints {
			var downStreamAddr []string
			var healthchecker *healthcheck.HttpHealthcheck
			if endpointsCfg.ProxyPass != nil {
				downStreamAddr = endpointsCfg.ProxyPass.Addresses
				changeHdrMap := make(map[string]string)
				var lb lbPort.Loadbalancer
				switch endpointsCfg.ProxyPass.Loadbalancing {
				case "random":
					lb = loadbalancer.NewRandomLoadbalancer()
					break
				case "ip-hashing":
					lb = loadbalancer.NewIPHashingLoadbalancer(downStreamAddr)
					break
				case "least-connection":
					lb = loadbalancer.NewLeastConnectionBalancer(downstreamConnTable)
					break
				case "round-robin":
					lb = loadbalancer.NewRoundRobinLoadBalancer()
				default:
					panic("wrong loadbalancing config")
				}
				for _, headerPairs := range endpointsCfg.ProxySetHeaders {
					changeHdrMap[headerPairs.HeaderName] = headerPairs.HeaderValue
				}
				if endpointsCfg.ProxyPass.Healthcheck.Enabled {
					healthchecker = healthcheck.NewHttpHealthcheck(
						endpointsCfg.ProxyPass.Healthcheck.Interval,
						endpointsCfg.ProxyPass.Healthcheck.Fails,
						endpointsCfg.ProxyPass.Healthcheck.Passes,
						downStreamAddr,
						logger,
					)
				}
				ratelimiterRepo := ratelimitRepo.NewAppMemoryRLRepo()
				ratelimiter := ratelimitApp.NewRatelimiter(
					ratelimiterRepo,
					endpointsCfg.RateLimit,
					1,
				)
				endpoints[idx] = httpHdl.NewEndpoint(
					endpointsCfg.RegexPath,
					lb,
					changeHdrMap,
					downStreamAddr,
					healthchecker,
					ratelimiter,
				)
			}
		}

		hdl = httpHdl.NewHTTPHdl(endpoints, logger, downstreamConnTable)

	} else {
		panic("not implemented")
	}
	createdServer := server.NewServer(
		serverCfg.Name,
		hdl,
		logger,
		serverCfg.ProxyTimeout,
		wg,
	)
	return createdServer
}
