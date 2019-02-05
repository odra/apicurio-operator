package apicurio

import (
	"fmt"
	"github.com/integr8ly/apicurio-operator/pkg/apis/integreatly/v1alpha1"
	"testing"
)

func TestGetTemplatePath(t *testing.T) {
	cases := []struct {
		Name        string
		TmplVersion string
		TmplType    v1alpha1.ApicurioType
		Apicurio    func(v string, t v1alpha1.ApicurioType) *v1alpha1.Apicurio
		Validate    func(t *testing.T, cr *v1alpha1.Apicurio, v string)
	}{
		{
			Name:        "Should create correct template path string",
			TmplVersion: "0.2.18.Final",
			TmplType:    "apicurio-full",
			Apicurio: func(v string, t v1alpha1.ApicurioType) *v1alpha1.Apicurio {
				return &v1alpha1.Apicurio{
					Spec: v1alpha1.ApicurioSpec{
						Version: v,
						Type:    t,
					},
				}
			},
			Validate: func(t *testing.T, cr *v1alpha1.Apicurio, v string) {
				expectedPath := fmt.Sprintf("%s/template.yml", v)
				actualPath := getTemplatePath(cr)
				if actualPath != expectedPath {
					t.Fatalf("expected Path: %s but got: %s", expectedPath, actualPath)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			instance := tc.Apicurio(tc.TmplVersion, tc.TmplType)
			tc.Validate(t, instance, tc.TmplVersion)
		})
	}
}

func TestMapResource(t *testing.T) {
	cases := []struct {
		Name     string
		Prefix   string
		Params   map[string]string
		Resource *v1alpha1.ApicurioResource
		Validate func(t *testing.T, instance *v1alpha1.ApicurioResource, params map[string]string)
	}{
		{
			Name:   "Should map all fields",
			Prefix: "UI",
			Params: make(map[string]string),
			Resource: &v1alpha1.ApicurioResource{
				Route: "STUDIO",
				CPU: &v1alpha1.MinMaxCfg{
					Min: "1",
					Max: "2",
				},
				Memory: &v1alpha1.MinMaxCfg{
					Min: "1",
					Max: "2",
				},
				JVM: &v1alpha1.MinMaxCfg{
					Min: "1",
					Max: "2",
				},
			},
			Validate: func(t *testing.T, instance *v1alpha1.ApicurioResource, params map[string]string) {
				//route
				if params["UI_ROUTE"] != instance.Route {
					t.Fatalf("Expected %s but found %s", params["UI_ROUTE"], instance.Route)
				}

				//cpu
				if params["UI_CPU_REQUEST"] != instance.CPU.Min {
					t.Fatalf("Expected %s but found %s", params["UI_CPU_REQUEST"], instance.CPU.Min)
				}

				if params["UI_CPU_LIMIT"] != instance.CPU.Max {
					t.Fatalf("Expected %s but found %s", params["UI_CPU_LIMIT"], instance.CPU.Max)
				}

				//memory
				if params["UI_MEM_REQUEST"] != instance.CPU.Min {
					t.Fatalf("Expected %s but found %s", params["UI_MEM_REQUEST"], instance.Memory.Min)
				}

				if params["UI_MEM_LIMIT"] != instance.CPU.Max {
					t.Fatalf("Expected %s but found %s", params["UI_MEM_LIMIT"], instance.Memory.Max)
				}

				//jvm
				if params["UI_JVM_MIN"] != instance.JVM.Min {
					t.Fatalf("Expected %s but found %s", params["API_JVM_MIN"], instance.JVM.Min)
				}
				if params["UI_JVM_MAX"] != instance.JVM.Max {
					t.Fatalf("Expected %s but found %s", params["API_JVM_MAX"], instance.JVM.Max)
				}
			},
		},
		{
			Name:   "Should map memory field",
			Prefix: "UI",
			Params: make(map[string]string),
			Resource: &v1alpha1.ApicurioResource{
				Memory: &v1alpha1.MinMaxCfg{
					Min: "1",
					Max: "2",
				},
			},
			Validate: func(t *testing.T, instance *v1alpha1.ApicurioResource, params map[string]string) {
				fields := []string{"UI_ROUTE", "UI_JVM_MIN", "UI_JVM_MAX", "UI_CPU_REQUEST", "UI_CPU_LIMIT"}
				for _, v := range fields {
					if _, ok := params[v]; ok {
						t.Fatalf("Field should not exist: %s", v)
					}
				}

				//memory
				if params["UI_MEM_REQUEST"] != instance.Memory.Min {
					t.Fatalf("Expected %s but found %s", params["UI_MEM_REQUEST"], instance.Memory.Min)
				}

				if params["UI_MEM_LIMIT"] != instance.Memory.Max {
					t.Fatalf("Expected %s but found %s", params["UI_MEM_LIMIT"], instance.Memory.Max)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mapResource(tc.Params, tc.Prefix, tc.Resource)
			tc.Validate(t, tc.Resource, tc.Params)
		})
	}
}

func TestPrefixKey(t *testing.T) {
	cases := []struct {
		Name     string
		Prefix   string
		Key      string
		Validate func(t *testing.T, key string)
	}{
		{
			Name:   "Should add prefix",
			Prefix: "PRIVATE",
			Key:    "KEY",
			Validate: func(t *testing.T, key string) {
				if key != "PRIVATE_KEY" {
					t.Fatalf("expected key PRIVATE_KEY but found %s", key)
				}
			},
		},
		{
			Name: "Should not add prefix",
			Key:  "KEY",
			Validate: func(t *testing.T, key string) {
				if key != "KEY" {
					t.Fatalf("expected key KEY but found %s", key)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Validate(t, prefixKey(tc.Prefix, tc.Key))
		})
	}
}

func TestUpdateRoute(t *testing.T) {
	cases := []struct {
		Name     string
		Apicurio *v1alpha1.Apicurio
		Validate func(t *testing.T, instance *v1alpha1.Apicurio)
	}{
		{
			Name: "Should update route",
			Apicurio: &v1alpha1.Apicurio{
				Spec: v1alpha1.ApicurioSpec{
					AppDomain: "apps.local",
					Api: &v1alpha1.ApicurioResource{
						Route: "api",
					},
				},
			},
			Validate: func(t *testing.T, instance *v1alpha1.Apicurio) {
				route := instance.Spec.Api.Route
				if route != "api.apps.local" {
					t.Fatalf("expected api.apps.local but found %s", route)
				}
			},
		},
		{
			Name: "Should not update route",
			Apicurio: &v1alpha1.Apicurio{
				Spec: v1alpha1.ApicurioSpec{
					Api: &v1alpha1.ApicurioResource{
						Route: "api",
					},
				},
			},
			Validate: func(t *testing.T, instance *v1alpha1.Apicurio) {
				route := instance.Spec.Api.Route
				if route != "api" {
					t.Fatalf("expected api but found %s", route)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			updateRoute(tc.Apicurio, tc.Apicurio.Spec.Api)
			tc.Validate(t, tc.Apicurio)
		})
	}
}

func TestFixAuthUrl(t *testing.T) {
	cases := []struct {
		Name        string
		Url         string
		ExpectedUrl string
		Validate    func(t *testing.T, expected string, actual string)
	}{
		{
			Name:        "Should fix url",
			Url:         "https://myurl.com/",
			ExpectedUrl: "myurl.com",
			Validate: func(t *testing.T, expected string, actual string) {
				if actual != "myurl.com" {
					t.Fatalf("Expected url \"%s\" but got \"%s\"", expected, actual)
				}
			},
		},
		{
			Name:        "Should ignore correct url",
			Url:         "myurl.com",
			ExpectedUrl: "myurl.com",
			Validate: func(t *testing.T, expected string, actual string) {
				if actual != "myurl.com" {
					t.Fatalf("Expected url \"%s\" but got \"%s\"", expected, actual)
				}
			},
		},
		{
			Name:        "Should fix prefix only",
			Url:         "https://myurl.com",
			ExpectedUrl: "myurl.com",
			Validate: func(t *testing.T, expected string, actual string) {
				if actual != "myurl.com" {
					t.Fatalf("Expected url \"%s\" but got \"%s\"", expected, actual)
				}
			},
		},
		{
			Name:        "Should fix sufix only",
			Url:         "myurl.com/",
			ExpectedUrl: "myurl.com",
			Validate: func(t *testing.T, expected string, actual string) {
				if actual != "myurl.com" {
					t.Fatalf("Expected url \"%s\" but got \"%s\"", expected, actual)
				}
			},
		},
		{
			Name:        "Should fix http url",
			Url:         "http://myurl.com/",
			ExpectedUrl: "myurl.com",
			Validate: func(t *testing.T, expected string, actual string) {
				if actual != "myurl.com" {
					t.Fatalf("Expected url \"%s\" but got \"%s\"", expected, actual)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			fixAuthUrl(&tc.Url)
			tc.Validate(t, tc.ExpectedUrl, tc.Url)
		})
	}
}
