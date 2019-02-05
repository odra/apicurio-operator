package kube

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestRemoveFinalizer(t *testing.T) {
	cases := []struct {
		Name     string
		Obj      runtime.Object
		Prepare  func(obj runtime.Object)
		Validate func(t *testing.T, obj runtime.Object)
	}{
		{
			Name: "Should remove a finalizer",
			Obj: &v1.Pod{
				ObjectMeta: v12.ObjectMeta{
					Finalizers: []string{"one", "two", "three"},
				},
			},
			Prepare: func(obj runtime.Object) {
				RemoveFinalizer(obj, "two")
			},
			Validate: func(t *testing.T, obj runtime.Object) {
				accessor, err := meta.Accessor(obj)
				if err != nil {
					t.Fatal(err)
				}

				finalizers := accessor.GetFinalizers()
				if len(finalizers) != 2 {
					t.Fatalf("Failed to remove finalizer: %v", finalizers)
				}

				if finalizers[0] != "one" || finalizers[1] != "three" {
					t.Fatalf("Failed to remove finalizer: %v", finalizers)
				}
			},
		},
		{
			Name: "Should not fail when removing a non-existing finalizer",
			Obj: &v1.Pod{
				ObjectMeta: v12.ObjectMeta{
					Finalizers: []string{"one", "three"},
				},
			},
			Prepare: func(obj runtime.Object) {
				RemoveFinalizer(obj, "two")
			},
			Validate: func(t *testing.T, obj runtime.Object) {
				accessor, err := meta.Accessor(obj)
				if err != nil {
					t.Fatal(err)
				}

				finalizers := accessor.GetFinalizers()
				if len(finalizers) != 2 {
					t.Fatalf("Failed to remove finalizer: %v", finalizers)
				}

				if finalizers[0] != "one" || finalizers[1] != "three" {
					t.Fatalf("Failed to remove finalizer: %v", finalizers)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Prepare(tc.Obj)
			tc.Validate(t, tc.Obj)
		})
	}
}

func TestHasFinalizer(t *testing.T) {
	cases := []struct {
		Name     string
		Obj      runtime.Object
		Validate func(t *testing.T, obj runtime.Object)
	}{
		{
			Name: "Should find finalizer",
			Obj: &v1.Pod{
				ObjectMeta: v12.ObjectMeta{
					Finalizers: []string{"one", "two", "three"},
				},
			},
			Validate: func(t *testing.T, obj runtime.Object) {
				finalizers := []string{"one", "two", "three"}

				for _, v := range finalizers {
					found, err := HasFinalizer(obj, v)
					if err != nil {
						t.Fatal(err)
					}

					if !found {
						t.Fatalf("Finalized not found: %s", v)
					}
				}
			},
		},
		{
			Name: "Should not find finalizer",
			Obj: &v1.Pod{
				ObjectMeta: v12.ObjectMeta{
					Finalizers: []string{"one", "two", "three"},
				},
			},
			Validate: func(t *testing.T, obj runtime.Object) {
				found, err := HasFinalizer(obj, "four")
				if err != nil {
					t.Fatal(err)
				}

				if found {
					t.Fatal("Should not find finalizer: four")
				}
			},
		},
		{
			Name: "Should not fail with unser finalizer property",
			Obj:  &v1.Pod{},
			Validate: func(t *testing.T, obj runtime.Object) {
				found, err := HasFinalizer(obj, "four")
				if err != nil {
					t.Fatal(err)
				}

				if found {
					t.Fatal("Should not find finalizer: four")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Validate(t, tc.Obj)
		})
	}
}

func TestAddFinalizer(t *testing.T) {
	cases := []struct {
		Name     string
		Obj      runtime.Object
		Prepare  func(obj runtime.Object)
		Validate func(t *testing.T, obj runtime.Object)
	}{
		{
			Name: "Should add finalizer",
			Obj:  &v1.Pod{},
			Prepare: func(obj runtime.Object) {
				AddFinalizer(obj, "one")
			},
			Validate: func(t *testing.T, obj runtime.Object) {
				accessor, err := meta.Accessor(obj)
				if err != nil {
					t.Fatal(err)
				}

				finalizers := accessor.GetFinalizers()
				if finalizers[0] != "one" {
					t.Fatalf("Failed to add finalizer: %v", finalizers)
				}
			},
		},
		{
			Name: "Should ignore adding existing finalizer",
			Obj: &v1.Pod{
				ObjectMeta: v12.ObjectMeta{
					Finalizers: []string{"one"},
				},
			},
			Prepare: func(obj runtime.Object) {
				AddFinalizer(obj, "one")
			},
			Validate: func(t *testing.T, obj runtime.Object) {
				accessor, err := meta.Accessor(obj)
				if err != nil {
					t.Fatal(err)
				}

				finalizers := accessor.GetFinalizers()
				if len(finalizers) != 1 {
					t.Fatalf("Should not add duplicated finalizer: %v", finalizers)
				}

				if finalizers[0] != "one" {
					t.Fatalf("Failed to add finalizer: %v", finalizers)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Prepare(tc.Obj)
			tc.Validate(t, tc.Obj)
		})
	}
}
