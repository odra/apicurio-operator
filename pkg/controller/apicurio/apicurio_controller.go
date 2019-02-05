package apicuriodeployment

import (
	"context"
	"log"
	integreatlyv1alpha1 "github.com/integr8ly/apicurio-operator/pkg/apis/integreatly/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"fmt"
	"github.com/integr8ly/operator-sdk-openshift-utils/pkg/api/template"
	"github.com/integr8ly/operator-sdk-openshift-utils/pkg/api/kubernetes"
	"github.com/sirupsen/logrus"
	kuberr "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"

	"github.com/gobuffalo/packr"
	"k8s.io/apimachinery/pkg/util/yaml"
	"github.com/integr8ly/apicurio-operator/pkg/kube"
)

// Add creates a new Apicurio Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileApiCurio{
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		config: mgr.GetConfig(),
		box:    packr.NewBox("../../../res"),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	opts := controller.Options{Reconciler: r}
	c, err := controller.New("apicuriodeployment-controller", mgr, opts)
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Apicurio
	err = c.Watch(&source.Kind{Type: &integreatlyv1alpha1.Apicurio{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &integreatlyv1alpha1.Apicurio{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileApiCurio{}

func (r *ReconcileApiCurio) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Printf("Reconciling Apicurio %s/%s\n", request.Namespace, request.Name)

	// Fetch the Apicurio instance
	instance := &integreatlyv1alpha1.Apicurio{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	if instance.GetDeletionTimestamp() != nil {
		err = r.deprovision(instance)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("deprovisioning failed: %v", err)
		}
		return reconcile.Result{}, nil
	}

	ok, err := kube.HasFinalizer(instance, integreatlyv1alpha1.DefaultFinalizer)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not read CR finalizer: %v", err)
	}

	if !ok {
		err = kube.AddFinalizer(instance, integreatlyv1alpha1.DefaultFinalizer)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to set finalizer in object: %v", err)
		}
		err = r.client.Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("failed update in object: %v", err)
		}
	}

	err = r.bootstrap(request, instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.processTemplate(instance)
	if err != nil {
		logrus.Errorf("Error while processing template: %v", err)
		return reconcile.Result{}, err
	}

	err = r.createObjects(request.Namespace, instance)
	if err != nil {
		logrus.Errorf("Error creating runtime objects: %v", err)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileApiCurio) bootstrap(request reconcile.Request, cr *integreatlyv1alpha1.Apicurio) error {
	if r.tmpl != nil {
		return nil
	}

	tmplPath := cr.Spec.Template
	if tmplPath == "" {
		return fmt.Errorf("Spec.Template.Path property is not defined")
	}

	yamlData, err := r.box.Find(cr.Spec.Template)
	if err != nil {
		return err
	}

	jsonData, err := yaml.ToJSON(yamlData)
	if err != nil {
		return err
	}

	tmpl, err := template.New(r.config, jsonData)
	if err != nil {
		return err
	}

	r.tmpl = tmpl

	return nil
}

func (r *ReconcileApiCurio) processTemplate(cr *integreatlyv1alpha1.Apicurio) error {
	params := make(map[string]string)
	for k, v := range routeParams {
		if k == "AUTH_ROUTE" && cr.Spec.ExternalAuthUrl != "" {
			params[k] = cr.Spec.ExternalAuthUrl
			continue
		}
		params[k] = v + "." + cr.Spec.AppDomain
	}

	if cr.Spec.AuthRealm != "" {
		params["KC_REALM"] = cr.Spec.AuthRealm
	}

	if len(cr.Spec.JvmHeap) == 2 {
		params["API_JVM_MIN"] = cr.Spec.JvmHeap[0]
		params["API_JVM_MAX"] = cr.Spec.JvmHeap[1]
	}
	if len(cr.Spec.MemLimit) == 2 {
		params["API_MEM_REQUEST"] = cr.Spec.MemLimit[0]
		params["API_MEM_LIMIT"] = cr.Spec.MemLimit[1]
	}

	err := r.tmpl.Process(params, cr.Namespace)
	if err != nil {
		logrus.Infof("Error: %v", err)
		return err
	}

	for _, ro := range r.tmpl.Source.Objects {
		b, _ := ro.MarshalJSON()
		logrus.Info("Object: %s", string(b[:]))
	}

	logrus.Infof("Template Objects: %+v", r.tmpl.Source.Objects)

	return nil
}

func (r *ReconcileApiCurio) createObjects(ns string, cr *integreatlyv1alpha1.Apicurio) error {

	objects := r.tmpl.GetObjects(filterAuthFn(cr))
	for _, o := range objects {
		uo, err := kubernetes.UnstructuredFromRuntimeObject(o)
		if err != nil {
			return err
		}

		uo.SetNamespace(ns)

		err = controllerutil.SetControllerReference(cr, uo, r.scheme)
		if err != nil {
			return err
		}

		err = r.client.Create(context.TODO(), uo.DeepCopyObject())
		if err != nil && !kuberr.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}

func (r *ReconcileApiCurio) deprovision(cr *integreatlyv1alpha1.Apicurio) error {
	ok, err := integreatlyv1alpha1.HasFinalizer(cr, "foregroundDeletion")
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	_, err = integreatlyv1alpha1.RemoveFinalizer(cr, integreatlyv1alpha1.ApicurioFinalizer)
	if err != nil {
		return err
	}

	err = r.client.Update(context.TODO(), cr)
	if err != nil {
		return fmt.Errorf("failed to update object: %v", err)
	}

	return nil
}

func filterAuthFn(cr *integreatlyv1alpha1.Apicurio) template.FilterFn {
	return func(obj *runtime.Object) error {
		uo, err := kubernetes.UnstructuredFromRuntimeObject(*obj)
		if err != nil {
			return err
		}

		isAuthObj := strings.Contains(uo.GetName(), "auth")
		if cr.Spec.ExternalAuthUrl != "" && isAuthObj {
			return fmt.Errorf("auth object should not be created: %v", obj)
		}

		return nil
	}
}
