package main

import (
	"sync"

	"github.com/phoon884/rev-proxy/pkg/factories"
)

func main() {
	factory := factories.NewFactory("./config.yaml")
	config, err := factory.InitializeConfigurator().GetConfig()
	logger := factory.InitializeLogger(config.LogLevel)
	logger.Info("Service starting...")
	if err != nil {
		logger.Error("Configuration error:", err.Error())
		return
	}
	var wg sync.WaitGroup
	for _, serverCfg := range config.Servers {
		wg.Add(1)
		createdServer := factory.ServerBuilder(&wg, serverCfg)
		go createdServer.StartServer(serverCfg.Port)
	}
	wg.Wait()
}
