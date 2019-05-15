package api

import (
	"github.com/integr8ly/apicurio-operator/test/api/meta"
	appsv1 "github.com/openshift/api/apps"
	v12 "github.com/openshift/api/apps/v1"
	imagev1 "github.com/openshift/api/image/v1"
	routev1 "github.com/openshift/api/route/v1"
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

func WaitForReadiness(watcher meta.ResourceWatcherSpec, opts meta.WaitOpts) error {
	return watcher.WaitForReadiness(opts.RetryInterval, opts.Timeout)
}

func WaitForDeletion(watcher meta.ResourceWatcherSpec, opts meta.WaitOpts) error {
	return watcher.WaitForDeletion(opts.RetryInterval, opts.Timeout)
}

func GetVar(name string) string {
	return os.Getenv(name)
}

func RegisterOpenshiftSchemes() error {
	if err := RegisterTypes(appsv1.Install, &v12.DeploymentConfigList{}); err !=nil {
		return  err
	}

	if err := RegisterTypes(routev1.Install, &routev1.Route{}); err !=nil {
		return  err
	}

	if err := RegisterTypes(imagev1.Install, &imagev1.ImageStream{}); err !=nil {
		return  err
	}

	return nil
}