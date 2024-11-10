package models

import (
	"time"
)

type UsersBucket struct {
	Last_updated time.Time
	BucketSize   float64
}
