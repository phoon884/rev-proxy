package application

import "github.com/phoon884/rev-proxy/pkg/config/domain/models"

func (c *ConfigService) GetConfig() (*models.Config, error) {
	if c.configuration == nil {
		if err := c.Config(); err != nil {
			return nil, err
		}
	}
	if err := c.configuration.Validate(); err != nil {
		return nil, err
	}
	return c.configuration, nil
}
