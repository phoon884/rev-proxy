package application

import (
	"time"

	"github.com/phoon884/rev-proxy/internal/httpHdl/ratelimit/domain/ports"
)

type Ratelimiter struct {
	repo     ports.RatelimitRepository
	capacity float64
	leakRate float64
}

func (r *Ratelimiter) Allow(userIP string) (bool, error) {
	if r.capacity <= 0 {
		// options for no ratelimit
		return true, nil
	}
	r.repo.Lock(userIP)
	defer r.repo.Unlock(userIP)
	bucket, err := r.repo.GetBucket(userIP)
	if err != nil {
		return false, err
	}
	var elaspedTime time.Duration
	if bucket.Last_updated.IsZero() {
		elaspedTime = time.Duration(0)
	} else {
		elaspedTime = time.Since(bucket.Last_updated)
	}
	bucket.Last_updated = time.Now()
	bucket.BucketSize = max(bucket.BucketSize-r.leakRate*elaspedTime.Seconds(), 0)

	bucket.BucketSize++
	r.repo.SetBucket(userIP, bucket)
	if bucket.BucketSize < r.capacity {
		return true, nil
	} else {
		return false, nil
	}
}

func NewRatelimiter(
	repo ports.RatelimitRepository,
	capacity float64,
	leakRate float64,
) *Ratelimiter {
	return &Ratelimiter{repo: repo, capacity: capacity, leakRate: leakRate}
}
