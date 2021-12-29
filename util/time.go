package util

import (
	"time"
)

func Duration(t time.Time) float64 {
	return time.Now().Sub(t).Seconds()
}
