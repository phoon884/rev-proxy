package application

import (
	"github.com/phoon884/rev-proxy/pkg/config/domain/models"
	"github.com/phoon884/rev-proxy/pkg/config/domain/ports"
)

type ConfigService struct {
	repository    ports.ConfigRepository
	configuration *models.Config
}

var _ ports.ConfigApplication = (*ConfigService)(nil)

func NewConfigService(
	repository ports.ConfigRepository,
) *ConfigService {
	return &ConfigService{
		repository: repository,
	}
}
