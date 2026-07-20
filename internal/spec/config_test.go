package spec

import (
	"reflect"
	"strings"
	"testing"
)

func TestResolveAddsDependenciesInStableOrder(t *testing.T) {
	input := DefaultConfig()
	input.Name = "acme"
	input.Module = "example.com/acme"
	input.Features = []Feature{FeatureSales, FeatureCron}

	got, err := Resolve(input)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	want := []Feature{FeatureAuth, FeatureRBAC, FeatureJobs, FeatureCron, FeatureCustomers, FeatureSales}
	if !reflect.DeepEqual(got.Features, want) {
		t.Fatalf("features = %#v, want %#v", got.Features, want)
	}
}

func TestResolveRejectsRetiredFeatureNamesWithMigrationAdvice(t *testing.T) {
	input := DefaultConfig()
	input.Name = "acme"
	input.Module = "example.com/acme"
	input.Features = []Feature{"crm"}
	_, err := Resolve(input)
	if err == nil || !strings.Contains(err.Error(), "customers") || !strings.Contains(err.Error(), "sales") {
		t.Fatalf("Resolve() error = %v, want CRM migration advice", err)
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

func TestResolveAcceptsCacheWithoutImplicitInfrastructureFeatures(t *testing.T) {
	input := DefaultConfig()
	input.Name = "acme"
	input.Module = "example.com/acme"
	input.Features = []Feature{FeatureCache}

	got, err := Resolve(input)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if !reflect.DeepEqual(got.Features, []Feature{FeatureCache}) {
		t.Fatalf("features = %#v, want cache only", got.Features)
	}
}
