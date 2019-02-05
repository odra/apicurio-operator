package meta

import (
	"k8s.io/apimachinery/pkg/runtime"
	"time"
)

const (
	DefaultRetryInterval = time.Second * 5
	DefaultTimeout       = time.Minute * 50
)

type WaitOpts struct {
	RetryInterval time.Duration
	Timeout       time.Duration
}

type ObjectLoader func() (runtime.Object, error)

type ReadinessSpec interface {
	Observe(opts WaitOpts, loader ObjectLoader) error
}
