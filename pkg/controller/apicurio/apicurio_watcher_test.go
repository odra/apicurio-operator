package apicurio

import (
	"github.com/integr8ly/operator-sdk-openshift-utils/pkg/api/schemes"
	"github.com/openshift/api/apps/v1"
	v13 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestNewWatcher(t *testing.T) {
	cases := []struct {
		Name     string
		client   func() client.Client
		Validate func(t *testing.T, watcher *watcher)
	}{
		{
			Name: "Should create client",
			client: func() client.Client {
				c := fake.NewFakeClient()
				return c
			},
			Validate: func(t *testing.T, watcher *watcher) {
				if watcher.client == nil {
					t.Fatalf("Failed to create client: %+v", watcher)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			w := NewWatcher(tc.client())
			tc.Validate(t, w)
		})
	}
}

func TestWatcher_AddChecker(t *testing.T) {
	cases := []struct{
		Name string
		CheckerName string
		Client   func() client.Client
		Watcher func(client client.Client) *watcher
		Validate func(t *testing.T, w *watcher, name string)
	}{
		{
			Name: "Should add new checker",
			CheckerName: "apicurio-studio-ui",
			Client: func() client.Client {
				c := fake.NewFakeClient()
				return c
			},
			Watcher: func(client client.Client) *watcher {
				return NewWatcher(client)
			},
			Validate: func(t *testing.T, w *watcher, name string) {
				if len(w.ResourceCheckers) < 1 {
					t.Fatalf("ResourceCheckers slice is empty: %v", w.ResourceCheckers)
				}

				if w.ResourceCheckers[0].Name != name {
					t.Fatalf("Expected name: %s but got %s", name, w.ResourceCheckers[0].Name)
				}
			},
		},
		{
			Name: "Should not add duplicated checkers",
			CheckerName: "apicurio-studio-ui",
			Client: func() client.Client {
				c := fake.NewFakeClient()
				return c
			},
			Watcher: func(client client.Client) *watcher {
				w := NewWatcher(client)
				w.addChecker( "apicurio-studio-ui")
				return w
			},
			Validate: func(t *testing.T, w *watcher, name string) {
				if len(w.ResourceCheckers) < 1 {
					t.Fatalf("ResourceCheckers slice is empty: %v", w.ResourceCheckers)
				}

				if w.ResourceCheckers[0].Name != name {
					t.Fatalf("Expected name: %s but got %s", name, w.ResourceCheckers[0].Name)
				}
			},
		},
	}

	for _, tc := range cases {
		watcher := tc.Watcher(tc.Client())
		watcher.addChecker(tc.CheckerName)
		tc.Validate(t, watcher, tc.CheckerName)
	}
}

func TestNewWatcher_GetDeploymentConfig(t *testing.T) {
	cases := []struct{
		Name string
		Key types.NamespacedName
		DC *v1.DeploymentConfig
		Client func() client.Client
		Validate func(t *testing.T, dc *v1.DeploymentConfig, key types.NamespacedName)
		ExpectError bool
	}{
		{
			Name: "Should retrieve a valid deployment config",
			DC: &v1.DeploymentConfig{},
			Key:types.NamespacedName{
				Name:      "mydc",
				Namespace: "default",
			},
			Client: func() client.Client {
				scheme := runtime.NewScheme()
				builder := runtime.NewSchemeBuilder(schemes.AddToScheme)
				builder.AddToScheme(scheme)

				objs := []v1.DeploymentConfig{
					{
						TypeMeta: v12.TypeMeta{
							APIVersion: "apps.openshift.io/v1",
							Kind:       "DeploymentConfig",
						},
						ObjectMeta: v12.ObjectMeta{
							Name:      "mydc",
							Namespace: "default",
						},
					},
				}
				ros := make([]runtime.Object, 0)
				for _, obj := range objs {
					ros = append(ros, obj.DeepCopyObject())
				}


				return fake.NewFakeClientWithScheme(scheme, ros...)

			},
			Validate: func(t *testing.T, dc *v1.DeploymentConfig, key types.NamespacedName) {
				if dc.Name != key.Name {
					t.Fatalf("Deployment config name does not match: %+v", dc)
				}

				if dc.Namespace != key.Namespace {
					t.Fatalf("Deployment config namespace does not match: %+v", dc)
				}
			},
			ExpectError: false,
		},
		{
			Name: "Should fail to retrieve a valid deployment config",
			DC: &v1.DeploymentConfig{},
			Key:types.NamespacedName{
				Name:      "mydc",
				Namespace: "default",
			},
			Client: func() client.Client {
				scheme := runtime.NewScheme()
				builder := runtime.NewSchemeBuilder(schemes.AddToScheme)
				builder.AddToScheme(scheme)

				return fake.NewFakeClientWithScheme(scheme)
			},
			Validate: func(t *testing.T, dc *v1.DeploymentConfig, key types.NamespacedName) {
				if dc != nil {
					t.Fatalf("Deployment config should be nil: %v", dc)
				}
			},
			ExpectError: true,
		},
	}
	
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			client := tc.Client()
			watcher := NewWatcher(client)
			err, dc := watcher.getDeploymentConfig(&tc.Key)

			if tc.ExpectError && err == nil {
				t.Fatal("Expected error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			tc.Validate(t, dc, tc.Key)
		})
	}
}

func TestNewWatcher_Reload(t *testing.T) {
	cases := []struct {
		Name string
		Key types.NamespacedName
		Client func() client.Client
		Validate func(t *testing.T, w *watcher)
		ExpectError bool
	}{
		{
			Name: "Should not fail when reloading empty checker list",
			Key:types.NamespacedName{
				Name: "mydc",
				Namespace: "default",
			},
			Client: func() client.Client {
				scheme := runtime.NewScheme()
				builder := runtime.NewSchemeBuilder(schemes.AddToScheme)
				builder.AddToScheme(scheme)

				return fake.NewFakeClientWithScheme(scheme)
			},
			Validate: func(t *testing.T, w *watcher) {
				if len(w.ResourceCheckers) > 0 {
					t.Fatalf("Should not have any checkers: %v", w.ResourceCheckers)
				}
			},
			ExpectError: false,
		},
		{
			Name: "Should not fail when reloading empty checker list",
			Key:types.NamespacedName{
				Name: "mydc",
				Namespace: "default",
			},
			Client: func() client.Client {
				objs := []v1.DeploymentConfig{
					{
						TypeMeta: v12.TypeMeta{
							APIVersion: "apps.openshift.io/v1",
							Kind:       "DeploymentConfig",
						},
						ObjectMeta: v12.ObjectMeta{
							Name:      "mydc",
							Namespace: "default",
						},
					},
				}
				ros := make([]runtime.Object, 0)
				for _, obj := range objs {
					ros = append(ros, obj.DeepCopyObject())
				}

				scheme := runtime.NewScheme()
				builder := runtime.NewSchemeBuilder(schemes.AddToScheme)
				builder.AddToScheme(scheme)

				return fake.NewFakeClientWithScheme(scheme, ros...)
			},
			Validate: func(t *testing.T, w *watcher) {
				w.addChecker("mydc")
				if len(w.ResourceCheckers)  != 1 {
					t.Fatalf("Should have one checker: %v", w.ResourceCheckers)
				}

				if w.ResourceCheckers[0].Name != "mydc" {
					t.Fatalf("Name does not match: %s", w.ResourceCheckers[0].Name)
				}
			},
			ExpectError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			client := tc.Client()
			watcher := NewWatcher(client)
			err := watcher.reload(tc.Key.Namespace)

			if tc.ExpectError && err == nil {
				t.Fatal("Expected error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			tc.Validate(t, watcher)
		})
	}
}

func TestNewWatcher_IsReady(t *testing.T) {
	cases := []struct{
		Name string
		Client func() client.Client
		Prepare func(w *watcher)
		Validate func(t *testing.T, w *watcher)
	}{
		{
			Name: "Should validate as ready using empty list",
			Client: func() client.Client {
				scheme := runtime.NewScheme()
				builder := runtime.NewSchemeBuilder(schemes.AddToScheme)
				builder.AddToScheme(scheme)

				return fake.NewFakeClientWithScheme(scheme)
			},
			Prepare: func(w *watcher) {

			},
			Validate: func(t *testing.T, w *watcher) {
				if !w.isReady() {
					t.Fatalf("Watcher should be ready: %v", w)
				}
			},
		},
		{
			Name: "Should validate as ready",
			Client: func() client.Client {
				scheme := runtime.NewScheme()
				builder := runtime.NewSchemeBuilder(schemes.AddToScheme)
				builder.AddToScheme(scheme)

				objs := []v1.DeploymentConfig{
					{
						TypeMeta: v12.TypeMeta{
							APIVersion: "apps.openshift.io/v1",
							Kind:       "DeploymentConfig",
						},
						ObjectMeta: v12.ObjectMeta{
							Name:      "mydc",
							Namespace: "default",
						},
						Status: v1.DeploymentConfigStatus{
							ReadyReplicas:     1,
							Replicas:          1,
							AvailableReplicas: 1,
							Conditions: []v1.DeploymentCondition{
								{
									Type:   v1.DeploymentAvailable,
									Status: v13.ConditionTrue,
								},
							},
						},
					},
				}
				ros := make([]runtime.Object, 0)
				for _, obj := range objs {
					ros = append(ros, obj.DeepCopyObject())
				}

				return fake.NewFakeClientWithScheme(scheme, ros...)
			},
			Prepare: func(w *watcher) {
				w.addChecker("mydc")
				w.reload("default")
			},
			Validate: func(t *testing.T, w *watcher) {
				if !w.isReady() {
					t.Fatalf("Watcher should be ready: %v", w)
				}
			},
		},
		//{
		//	Name: "Should not validate as ready",
		//	Client: func() client.Client {
		//		scheme := runtime.NewScheme()
		//		builder := runtime.NewSchemeBuilder(schemes.AddToScheme)
		//		builder.AddToScheme(scheme)
		//
		//		objs := []v1.DeploymentConfig{
		//			{
		//				TypeMeta: v12.TypeMeta{
		//					APIVersion: "apps.openshift.io/v1",
		//					Kind:       "DeploymentConfig",
		//				},
		//				ObjectMeta: v12.ObjectMeta{
		//					Name:      "mydc",
		//					Namespace: "default",
		//				},
		//				Status: v1.DeploymentConfigStatus{
		//					ReadyReplicas:     1,
		//					Replicas:          1,
		//					AvailableReplicas: 1,
		//					Conditions: []v1.DeploymentCondition{
		//						{
		//							Type:   v1.DeploymentAvailable,
		//							Status: v13.ConditionFalse,
		//						},
		//					},
		//				},
		//			},
		//		}
		//		ros := make([]runtime.Object, 0)
		//		for _, obj := range objs {
		//			ros = append(ros, obj.DeepCopyObject())
		//		}
		//
		//		return fake.NewFakeClientWithScheme(scheme, ros...)
		//	},
		//	Prepare: func(w *watcher) {
		//		w.addChecker("mydc")
		//		w.reload("default")
		//	},
		//	Validate: func(t *testing.T, w *watcher) {
		//		if w.isReady() {
		//			t.Fatalf("Watcher should not be ready: %v", w)
		//		}
		//	},
		//},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			client := tc.Client()
			watcher := NewWatcher(client)
			tc.Prepare(watcher)
			tc.Validate(t, watcher)
		})
	}
}