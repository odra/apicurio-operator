package apicurio

import (
	"context"
	"github.com/integr8ly/apicurio-operator/pkg/kube/checkers"
	"github.com/openshift/api/apps/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func NewWatcher(client clientSpec) *watcher {
	return &watcher{
		client:           client,
		ResourceCheckers: []*statusChecker{},
	}
}

func (w *watcher) Reload(ns string) error {
	key := &types.NamespacedName{Namespace: ns}

	for _, res := range w.ResourceCheckers {
		key.Name = res.Name

		err, dc := w.getDeploymentConfig(key)
		if err != nil {
			return err
		}

		//info := &statusResource{
		//	CPU: map[string]string {
		//		"MIN": "",
		//		"MAX": "",
		//	},
		//	Memory: map[string]string{
		//		"MIN": "",
		//		"MAX": "",
		//	},
		//	JVM: map[string]string{
		//		"MIN": "",
		//		"MAX": "",
		//	},
		//}

		//dcRes := dc.Spec.Template.Spec.Containers[0].Resources
		//info := &statusResource{
		//	CPU: map[string]string {
		//		"MIN": dcRes.Requests.Cpu().String(),
		//		"MAX": dcRes.Limits.Cpu().String(),
		//	},
		//	Memory: map[string]string{
		//		"MIN": dcRes.Requests.Memory().String(),
		//		"MAX": dcRes.Limits.Memory().String(),
		//	},
		//	JVM: map[string]string{
		//		"MIN": getEnvValue(dc.Spec.Template.Spec.Containers[0].Env, "APICURIO_MIN_HEAP"),
		//		"MAX": getEnvValue(dc.Spec.Template.Spec.Containers[0].Env, "APICURIO_MAX_HEAP"),
		//	},
		//}
		//res.Info = info

		res.checker = checkers.NewDeploymentConfigCheck(dc)
	}

	return nil
}

func (w *watcher) IsReady() bool {
	for _, res := range w.ResourceCheckers {
		res.IsReady = res.checker.IsReady()
		if !res.checker.IsReady() {
			return false
		}
	}

	return true
}

func (w *watcher) AddChecker(name string) {
	for _, checker := range w.ResourceCheckers {
		if checker.Name == name {
			return
		}
	}

	w.ResourceCheckers = append(w.ResourceCheckers, &statusChecker{Name: name})
}

func (w *watcher) getDeploymentConfig(key *types.NamespacedName) (error, *v1.DeploymentConfig) {
	dc := &v1.DeploymentConfig{
		ObjectMeta: v12.ObjectMeta{
			Namespace: key.Namespace,
			Name:      key.Name,
		},
	}

	err := w.client.Get(context.TODO(), *key, dc)
	if err != nil {
		return err, nil
	}

	return nil, dc
}
