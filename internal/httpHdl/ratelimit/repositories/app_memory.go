package repositories

import (
	"sync"
	"time"

	"github.com/phoon884/rev-proxy/internal/httpHdl/ratelimit/domain/models"
	"github.com/phoon884/rev-proxy/internal/httpHdl/ratelimit/domain/ports"
)

type AppMemoryRLRepo struct {
	data           map[string]*mapValue
	maxEntry       int
	recordDuration time.Duration // capcity/leak rate
}

type mapValue struct {
	value models.UsersBucket
	mut   sync.Mutex
}

var _ ports.RatelimitRepository = (*AppMemoryRLRepo)(nil)

func NewAppMemoryRLRepo() *AppMemoryRLRepo {
	return &AppMemoryRLRepo{
		data:           make(map[string]*mapValue),
		recordDuration: time.Hour,
	}
}

func (a *AppMemoryRLRepo) Lock(userIP string) {
	if a.data[userIP] == nil {
		a.data[userIP] = &mapValue{}
	}
	a.data[userIP].mut.Lock()
}

func (a *AppMemoryRLRepo) Unlock(userIP string) {
	if a.data[userIP] == nil {
		a.data[userIP] = &mapValue{}
	}
	a.data[userIP].mut.Unlock()
}

func (a *AppMemoryRLRepo) GetBucket(userIP string) (models.UsersBucket, error) {
	if a.data[userIP] == nil {
		a.data[userIP] = &mapValue{}
	}
	return a.data[userIP].value, nil
}

func (a *AppMemoryRLRepo) SetBucket(userIP string, data models.UsersBucket) error {
	if len(a.data) >= a.maxEntry {
		a.makespace()
	}
	a.data[userIP].value = data
	return nil
}

func (a *AppMemoryRLRepo) makespace() {
	for addr, value := range a.data {
		if value.mut.TryLock() {
			if value.value.Last_updated.Before(time.Now().Add(-a.recordDuration)) {
				delete(a.data, addr)
			} else {
				value.mut.Unlock()
			}
		}
	}
}
