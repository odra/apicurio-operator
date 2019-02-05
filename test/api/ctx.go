package api

import (
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"

	"github.com/integr8ly/apicurio-operator/pkg/apis/integreatly/v1alpha1"
	"github.com/integr8ly/apicurio-operator/test/api/meta"
)

func PrepareContext(t *testing.T, opts meta.WaitOpts) *framework.TestCtx {
	ctx := framework.NewTestCtx(t)
	opt := &framework.CleanupOptions{
		TestContext:   ctx,
		RetryInterval: opts.RetryInterval,
		Timeout:       opts.Timeout,
	}

	err := ctx.InitializeClusterResources(opt)
	if err != nil {
		t.Fatalf("Failed to initialize test context: %v", err)
	}

	ns, err := ctx.GetNamespace()
	if err != nil {
		t.Fatalf("Failed to get context namespace: %v", err)
	}

	globalVars := framework.Global

	err = e2eutil.WaitForDeployment(t, globalVars.KubeClient, ns, "apicurio-operator", 1, opts.RetryInterval, opts.Timeout)
	if err != nil {
		t.Fatalf("Operator deployment failed: %v", err)
	}

	return ctx
}

func RegisterTypes(objs ...runtime.Object) error {
	for _, obj := range objs {
		err := framework.AddToFrameworkScheme(v1alpha1.SchemeBuilder.AddToScheme, obj)
		if err != nil {
			return err
		}
	}

	return nil
}
