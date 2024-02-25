package controllersdk

import (
	"context"
	"testing"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/denniskniep/provider-springclouddataflow/internal/clients"
	"github.com/google/go-cmp/cmp"
)

func SkipIfIsShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}

func AssertAreEqual[C any](t *testing.T, actual interface{}, expected interface{}) {
	t.Helper()
	mappedActual, err := MapToCompare[C](actual)
	if err != nil {
		t.Fatal(err)
	}

	mappedExpected, err := MapToCompare[C](expected)
	if err != nil {
		t.Fatal(err)
	}

	diff := cmp.Diff(mappedActual, mappedExpected)
	if diff != "" {
		t.Fatal(diff)
	}
}

func TestCreateAndAssert[R resource.Managed, P any, O any, C any](t *testing.T, srv clients.Service[R, P, O, C], param *P) *O {
	t.Helper()
	err := srv.Create(context.Background(), param)
	if err != nil {
		t.Fatal(err)
	}

	created, err := srv.Describe(context.Background(), param)
	if err != nil {
		t.Fatal(err)
	}

	if created == nil {
		t.Fatal("Created object was not found")
	}

	AssertAreEqual[C](t, created, param)

	return created
}

func TestUpdate[R resource.Managed, P any, O any, C any](t *testing.T, srv clients.Service[R, P, O, C], param *P) {
	t.Helper()
	err := srv.Update(context.Background(), param)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateAndAssert[R resource.Managed, P any, O any, C any](t *testing.T, srv clients.Service[R, P, O, C], param *P) {
	TestUpdate[R, P, O, C](t, srv, param)

	updated, err := srv.Describe(context.Background(), param)
	if err != nil {
		t.Fatal(err)
	}

	AssertAreEqual[C](t, updated, param)
}

func TestDelete[R resource.Managed, P any, O any, C any](t *testing.T, srv clients.Service[R, P, O, C], param *P) {
	t.Helper()
	err := srv.Delete(context.Background(), param)
	if err != nil {
		t.Fatal(err)
	}

	noApp, err := srv.Describe(context.Background(), param)
	if err != nil {
		t.Fatal(err)
	}

	if noApp != nil {
		t.Fatal("App was not deleted")
	}
}
