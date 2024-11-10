package ports

import (
	"github.com/phoon884/rev-proxy/internal/httpHdl/ratelimit/domain/models"
)

type RatelimitRepository interface {
	GetBucket(userIP string) (models.UsersBucket, error)
	Lock(userIP string)
	Unlock(userIP string)
	SetBucket(userIP string, value models.UsersBucket) error
}
