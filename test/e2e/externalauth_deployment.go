package e2e

import (
	"context"
	"github.com/gobuffalo/packr"
	"github.com/integr8ly/apicurio-operator/pkg/apis/integreatly/v1alpha1"
	testapi "github.com/integr8ly/apicurio-operator/test/api"
	"github.com/integr8ly/apicurio-operator/test/api/apicurio"
	"github.com/integr8ly/apicurio-operator/test/api/meta"
	"github.com/integr8ly/apicurio-operator/test/api/watchers"
	"github.com/integr8ly/apicurio-operator/test/mock"
	"github.com/integr8ly/operator-sdk-openshift-utils/pkg/api/kubernetes"
	"github.com/integr8ly/operator-sdk-openshift-utils/pkg/api/template"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"testing"
)

func TestApicurioExternalAuthDeployment(t *testing.T) {
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
	mock.SetExternalAuth(cr, "host", "username", "password", "realm")

	err = framework.Global.Client.Create(context.TODO(), cr, testapi.CleanupOpts(ctx))
	if err != nil {
		t.Fatalf("Failed to create apicurio resource: %v", err)
	}

	//deploy keycloak instance
	box := packr.NewBox("../../res")
	raw, err := box.Find("0.2.20.Final/auth-template.yml")
	if err != nil {
		t.Fatalf("Failed to read keycloak template file content: %v", err)
	}
	jsonData, err := yaml.ToJSON(raw)
	if err != nil {
		t.Fatalf("Failed to parse yaml data to json: %v", err)
	}
	tmpl, err := template.New(framework.Global.KubeConfig, jsonData)
	if err != nil {
		t.Fatalf("Error while initializing template parser: %v", err)
	}
	params := make(map[string]string)
	err = tmpl.Process(params, cr.Namespace)
	if err != nil {
		t.Fatalf("Failed while processing template: %v", err)
	}
	for _, obj := range tmpl.GetObjects(template.NoFilterFn) {
		uo, err := kubernetes.UnstructuredFromRuntimeObject(obj)
		if err != nil {
			t.Fatalf("%v", err)
		}

		uo.SetNamespace(ns)

		err = framework.Global.Client.Create(context.TODO(), uo.DeepCopy(), testapi.CleanupOpts(ctx))
		if err != nil {
			t.Fatalf("Failed to create object: %v", err)
		}
	}

	checker := watchers.NewAPicurioWatcher(framework.Global.Client.Client)
	checker.Instance.Status.Type = v1alpha1.ApicurioReady
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
}