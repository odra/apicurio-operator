package apicurio

import (
	"context"
	"encoding/json"
	"fmt"
	integreatlyv1alpha1 "github.com/integr8ly/apicurio-operator/pkg/apis/integreatly/v1alpha1"
	"github.com/integr8ly/apicurio-operator/pkg/kube/template/filters"
	"github.com/integr8ly/operator-sdk-openshift-utils/pkg/api/kubernetes"
	"github.com/integr8ly/operator-sdk-openshift-utils/pkg/api/template"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	kuberr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gobuffalo/packr"
	"github.com/integr8ly/apicurio-operator/pkg/kube"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/rest"

	"github.com/openshift/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var apicurioWatcher *watcher

type ReconcileApiCurio struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	config *rest.Config
	scheme *runtime.Scheme
	tmpl   *template.Tmpl
	box    packr.Box
}

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

	client := mgr.GetClient()
	apicurioWatcher = NewWatcher(client)

	// Watch for changes to primary resource Apicurio
	err = c.Watch(&source.Kind{Type: &integreatlyv1alpha1.Apicurio{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1.DeploymentConfig{}}, &handler.EnqueueRequestForOwner{
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

	switch instance.Status.Type {
	case integreatlyv1alpha1.ApicurioNone:
		err = r.bootstrap(instance)
		if err != nil {
			logrus.Errorf("Error while bootstrapping template: %v", err)
			return reconcile.Result{}, err
		}
	case integreatlyv1alpha1.ApicurioNew:
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

		err = r.reloadCheckers(instance)
		if err != nil {
			logrus.Errorf("Failed to set apicurio checkers: %v", err)
			return reconcile.Result{}, err
		}
	case integreatlyv1alpha1.ApicurioReconcile:
		err = r.reloadCheckers(instance)
		if err != nil {
			return reconcile.Result{}, err
		}

		logrus.Info("Provisioning please wait")

		if !r.isReady(instance) {
			logrus.Info("Deployment not ready yet")
			return reconcile.Result{Requeue: true}, nil
		}

		logrus.Info("Apicurio is ready - finishing reconcile loop")
		err = r.finish(instance)
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	case integreatlyv1alpha1.ApicurioReady:
		if err = r.reloadCheckers(instance); err != nil {
			return reconcile.Result{}, err
		}

		if !apicurioWatcher.isReady() {
			logrus.Info("Readiness checks failed, rolling back reconcile status")
			if err = r.recycle(instance); err != nil {
				return reconcile.Result{}, err
			}

			return reconcile.Result{}, nil
		} else {
			_, err := r.diff(instance)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	default:
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileApiCurio) bootstrap(cr *integreatlyv1alpha1.Apicurio) error {
	cr.Status = integreatlyv1alpha1.ApicurioStatus{
		Type:    integreatlyv1alpha1.ApicurioNew,
		Reason:  new(string),
		Message: new(string),
	}
	*cr.Status.Reason = "NewApicurio"
	*cr.Status.Message = "New apicurio cr detected"
	err := r.client.Status().Update(context.TODO(), cr)
	if err != nil {
		return err
	}

	if r.tmpl != nil {
		return nil
	}

	tmplPath := getTemplatePath(cr)

	yamlData, err := r.box.Find(tmplPath)
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
	filterDcFn := filters.FilterDcs

	if cr.Spec.Auth != nil {
		authCommons := cr.Spec.Auth.ApicurioResource
		updateRoute(cr, &authCommons)
		mapResource(params, "AUTH", &authCommons)
		//external auth instanceAuth
		if cr.Spec.Type == integreatlyv1alpha1.ApicurioExternalAuthType {
			fixAuthUrl(&cr.Spec.Auth.Host)
			params["AUTH_ROUTE"] = cr.Spec.Auth.Host
			params["KC_USER"] = cr.Spec.Auth.Username
			params["KC_PASS"] = cr.Spec.Auth.Password

			if cr.Spec.Auth.Realm != "" {
				params["KC_REALM"] = cr.Spec.Auth.Realm
			}

			filterDcFn = filters.FilterSkipAuthDc
		}
	}

	if cr.Spec.Api != nil {
		updateRoute(cr, cr.Spec.Api)
		mapResource(params, "API", cr.Spec.Api)
	}

	if cr.Spec.WebSocket != nil {
		updateRoute(cr, cr.Spec.WebSocket)
		mapResource(params, "WS", cr.Spec.WebSocket)
	}

	if cr.Spec.Studio != nil {
		updateRoute(cr, cr.Spec.Studio)
		mapResource(params, "UI", cr.Spec.Studio)
	}

	err := r.tmpl.Process(params, cr.Namespace)
	if err != nil {
		return err
	}

	dcs := r.tmpl.GetObjects(filterDcFn)
	for _, dc := range dcs {
		uo, err := kubernetes.UnstructuredFromRuntimeObject(dc)
		if err != nil {
			return err
		}

		apicurioWatcher.addChecker(uo.GetName())
	}

	cr.Status = integreatlyv1alpha1.ApicurioStatus{
		Type:    integreatlyv1alpha1.ApicurioReconcile,
		Reason:  new(string),
		Message: new(string),
	}
	*cr.Status.Reason = "ApicurioReconcile"
	*cr.Status.Message = "Deploying apicurio instance"
	err = r.client.Status().Update(context.TODO(), cr)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReconcileApiCurio) createObjects(ns string, cr *integreatlyv1alpha1.Apicurio) error {
	objects := make([]runtime.Object, 0)
	var filterFn = template.NoFilterFn

	if cr.Spec.Type == integreatlyv1alpha1.ApicurioExternalAuthType {
		filterFn = filters.FilterSkipAuthObjs
	}

	r.tmpl.CopyObjects(filterFn, &objects)

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

func (r *ReconcileApiCurio) reloadCheckers(cr *integreatlyv1alpha1.Apicurio) error {
	return apicurioWatcher.reload(cr.Namespace)
}

func (r *ReconcileApiCurio) deprovision(cr *integreatlyv1alpha1.Apicurio) error {
	ok, err := kube.HasFinalizer(cr, "foregroundDeletion")
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	_, err = kube.RemoveFinalizer(cr, integreatlyv1alpha1.DefaultFinalizer)
	if err != nil {
		return err
	}

	cr.Status = integreatlyv1alpha1.ApicurioStatus{
		Type:    integreatlyv1alpha1.ApicurioDelete,
		Reason:  new(string),
		Message: new(string),
	}
	*cr.Status.Reason = "ApicurioDelete"
	*cr.Status.Message = "Apicurio is being deleted"

	err = r.client.Update(context.TODO(), cr)
	if err != nil {
		return fmt.Errorf("failed to update object: %v", err)
	}

	return nil
}

func (r *ReconcileApiCurio) finish(cr *integreatlyv1alpha1.Apicurio) error {
	b, err := json.Marshal(apicurioWatcher)
	if err != nil {
		return err
	}

	annotations := cr.GetAnnotations()
	if annotations == nil {
		return fmt.Errorf("apicurio annotation is nil: %v", annotations)
	}

	annotations["io.apicurio/state"] = fmt.Sprintf("%s", string(b[:]))
	cr.SetAnnotations(annotations)
	err = r.client.Update(context.TODO(), cr)
	if err != nil {
		return err
	}

	key := types.NamespacedName{
		Name: cr.Name,
		Namespace: cr.Namespace,
	}
	err = r.client.Get(context.TODO(), key, cr)
	if err != nil {
		return err
	}

	cr.Status = integreatlyv1alpha1.ApicurioStatus{
		Type:    integreatlyv1alpha1.ApicurioReady,
		Reason:  new(string),
		Message: new(string),
	}
	*cr.Status.Reason = "ApicurioReady"
	*cr.Status.Message = "Apicurio is ready"

	return r.client.Status().Update(context.TODO(), cr)
}

func (r *ReconcileApiCurio) recycle(cr *integreatlyv1alpha1.Apicurio) error {
	cr.Status = integreatlyv1alpha1.ApicurioStatus{
		Type:    integreatlyv1alpha1.ApicurioNew,
		Reason:  new(string),
		Message: new(string),
	}
	*cr.Status.Reason = "NewApicurio"
	*cr.Status.Message = "New apicurio cr detected"

	return r.client.Status().Update(context.TODO(), cr)
}

func (r *ReconcileApiCurio) diff(cr *integreatlyv1alpha1.Apicurio) ([]string, error) {
	var err error

	resources := make([]string, 0)

	err = apicurioWatcher.reload(cr.Namespace)
	if err != nil {
		return resources, err
	}

	return resources, nil
}


