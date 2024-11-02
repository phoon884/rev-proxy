package ports

import (
	"github.com/phoon884/rev-proxy/pkg/config/domain/models"
)

type ConfigApplication interface {
	Config() error
	GetConfig() (*models.Config, error)
}

type ConfigRepository interface {
	GetConfig() (*[]byte, error)
}
