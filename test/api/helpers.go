package api

import (
	metatest "github.com/integr8ly/apicurio-operator/test/api/meta"
	"github.com/operator-framework/operator-sdk/pkg/test"
	"os"
	"time"
)

func CleanupOpts(ctx *test.TestCtx) *test.CleanupOptions {
	return &test.CleanupOptions{
		TestContext:   ctx,
		Timeout:       time.Minute * 2,
		RetryInterval: time.Second * 5,
	}
}

func WaitForReadiness(checker metatest.ReadinessSpec, opts metatest.WaitOpts, loader metatest.ObjectLoader) error {
	return checker.Observe(opts, loader)
}

func GetVar(name string) string {
	return os.Getenv(name)
}
