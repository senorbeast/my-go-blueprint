package spec

import (
	"reflect"
	"testing"
)

func TestResolveAddsDependenciesInStableOrder(t *testing.T) {
	input := DefaultConfig()
	input.Name = "acme"
	input.Module = "example.com/acme"
	input.Features = []Feature{FeatureCRM, FeatureCron}

	got, err := Resolve(input)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	want := []Feature{FeatureAuth, FeatureRBAC, FeatureJobs, FeatureCron, FeatureCRM}
	if !reflect.DeepEqual(got.Features, want) {
		t.Fatalf("features = %#v, want %#v", got.Features, want)
	}
}

func TestResolveRejectsUnknownFeature(t *testing.T) {
	input := DefaultConfig()
	input.Name = "acme"
	input.Module = "example.com/acme"
	input.Features = []Feature{"unknown"}
	if _, err := Resolve(input); err == nil {
		t.Fatal("Resolve() expected an error")
	}
}
