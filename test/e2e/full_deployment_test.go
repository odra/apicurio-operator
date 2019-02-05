package e2e

import (
	"context"
	"github.com/integr8ly/apicurio-operator/pkg/apis/integreatly/v1alpha1"
	testapi "github.com/integr8ly/apicurio-operator/test/api"
	"github.com/integr8ly/apicurio-operator/test/api/apicurio"
	"github.com/integr8ly/apicurio-operator/test/api/meta"
	"github.com/integr8ly/apicurio-operator/test/api/watchers"
	"github.com/integr8ly/apicurio-operator/test/mock"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

func TestApicurioStandaloneDeployment(t *testing.T) {
	var err error

	delay := meta.WaitOpts{
		RetryInterval: meta.DefaultRetryInterval,
		Timeout:       meta.DefaultTimeout,
	}

	ctx := testapi.PrepareContext(t, delay)
	defer ctx.Cleanup()

	err = testapi.RegisterTypes(&v1alpha1.ApicurioList{})
	if err != nil {
		t.Fatalf("Failed to register crd scheme: %v", err)
	}

	ns, err := ctx.GetNamespace()
	if err != nil {
		t.Fatalf("Failed to retrieve namespace: %v", err)
	}

	cr := mock.NewApicurioFull("0.2.20.Final", "example.com")
	cr.ObjectMeta.Namespace = ns

	err = framework.Global.Client.Create(context.TODO(), cr, &framework.CleanupOptions{})
	if err != nil {
		t.Fatalf("Failed to create apicurio resource: %v", err)
	}

	reason := new(string)
	*reason = "ApicurioReady"
	checker := watchers.NewAPicurioWatcher(framework.Global.Client.Client)
	checker.Instance.Status.Type = v1alpha1.ApicurioReady
	checker.Instance.Status.Reason = reason

	readInstance := new(v1alpha1.Apicurio)
	key := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Name,
	}
	waitOpts := meta.WaitOpts{
		RetryInterval: meta.DefaultRetryInterval,
		Timeout:       meta.DefaultTimeout,
	}
	waitFn := apicurio.WaitForDeploymentReadiness(readInstance, key, framework.Global.Client)

	err = testapi.WaitForReadiness(checker, waitOpts, waitFn)
	if err != nil {
		t.Fatalf("Error while waiting fot readiness: %v", err)
	}

	err = framework.Global.Client.Delete(context.TODO(), cr)
	if err != nil {
		t.Fatalf("Failed to create apicurio resource: %v", err)
	}

	err = e2eutil.WaitForDeletion(t, framework.Global.Client.Client, readInstance, waitOpts.RetryInterval, waitOpts.Timeout)
	if err != nil {
		t.Fatalf("Error while waiting for apicurio deletion: %v", err)
	}
}
