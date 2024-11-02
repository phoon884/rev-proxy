package application

import (
	"gopkg.in/yaml.v3"

	"github.com/phoon884/rev-proxy/pkg/config/domain/models"
)

func (c *ConfigService) Config() error {
	configuration := new(models.Config)
	confBytes, err := c.repository.GetConfig()
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(*confBytes, configuration)
	if err != nil {
		return err
	}

	c.configuration = configuration
	return nil
}
