package clients

import (
	"testing"

	core "github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/context"
)

func createApplicationService(t *testing.T) ApplicationService {
	jsonConfig := `{
		"Uri": "http://localhost:9393/"
	}`

	srv, err := NewApplicationService([]byte(jsonConfig))
	if err != nil {
		t.Fatal(err)
	}

	return srv
}

func createDefaultApplication(appType string, name string, version string) *core.ApplicationParameters {
	return &core.ApplicationParameters{
		Type:           appType,
		Name:           name,
		Version:        version,
		Uri:            "docker://hello-world:" + version,
		BootVersion:    "2",
		DefaultVersion: true,
	}
}

func TestCreate(t *testing.T) {
	skipIfIsShort(t)

	srv := createApplicationService(t)
	testApp := createDefaultApplication("task", "Test001", "v1.0.0")

	createdApp := create(t, srv, testApp)

	assertApplicationAreEqual(t, srv, createdApp, testApp)

	delete(t, srv, testApp)
}

func TestDeleteNotExisting(t *testing.T) {
	skipIfIsShort(t)
	srv := createApplicationService(t)
	testApp := createDefaultApplication("task", "Test999", "v1.0.0")
	delete(t, srv, testApp)
}

func TestUpdateNotExisting(t *testing.T) {
	skipIfIsShort(t)
	srv := createApplicationService(t)
	testApp := createDefaultApplication("task", "Test888", "v1.0.0")
	update(t, srv, testApp)
}

func TestUpdate(t *testing.T) {
	skipIfIsShort(t)

	srv := createApplicationService(t)
	testAppV1 := createDefaultApplication("task", "Test002", "v1.0.0")
	createdAppV1 := create(t, srv, testAppV1)
	assertApplicationAreEqual(t, srv, createdAppV1, testAppV1)

	testAppV2 := createDefaultApplication("task", "Test002", "v2.0.0")
	testAppV2.DefaultVersion = false
	createdAppV2 := create(t, srv, testAppV2)
	assertApplicationAreEqual(t, srv, createdAppV2, testAppV2)

	testAppV1.DefaultVersion = false
	testAppV2.DefaultVersion = true
	update(t, srv, testAppV1)
	update(t, srv, testAppV2)

	foundAppV1, err := srv.DescribeApplication(context.Background(), testAppV1)
	if err != nil {
		t.Fatal(err)
	}

	assertApplicationAreEqual(t, srv, foundAppV1, testAppV1)

	foundAppV2, err := srv.DescribeApplication(context.Background(), testAppV2)
	if err != nil {
		t.Fatal(err)
	}

	assertApplicationAreEqual(t, srv, foundAppV2, testAppV2)

	delete(t, srv, testAppV1)
	delete(t, srv, testAppV2)
}

func update(t *testing.T, srv ApplicationService, app *core.ApplicationParameters) {
	t.Helper()
	err := srv.UpdateApplication(context.Background(), app)
	if err != nil {
		t.Fatal(err)
	}
}

func create(t *testing.T, srv ApplicationService, app *core.ApplicationParameters) *core.ApplicationObservation {
	t.Helper()
	err := srv.CreateApplication(context.Background(), app)
	if err != nil {
		t.Fatal(err)
	}

	createdApp, err := srv.DescribeApplication(context.Background(), app)
	if err != nil {
		t.Fatal(err)
	}

	if createdApp == nil {
		t.Fatal("App was not found")
	}
	return createdApp
}

func delete(t *testing.T, srv ApplicationService, app *core.ApplicationParameters) {
	t.Helper()
	err := srv.DeleteApplication(context.Background(), app)
	if err != nil {
		t.Fatal(err)
	}

	noApp, err := srv.DescribeApplication(context.Background(), app)
	if err != nil {
		t.Fatal(err)
	}

	if noApp != nil {
		t.Fatal("App was not deleted")
	}
}

func assertApplicationAreEqual(t *testing.T, srv ApplicationService, actual *core.ApplicationObservation, expected *core.ApplicationParameters) {
	t.Helper()
	mappedActual, err := srv.MapToApplicationCompare(actual)
	if err != nil {
		t.Fatal(err)
	}

	mappedExpected, err := srv.MapToApplicationCompare(expected)
	if err != nil {
		t.Fatal(err)
	}

	diff := cmp.Diff(mappedActual, mappedExpected)
	if diff != "" {
		t.Fatal(diff)
	}
}

func skipIfIsShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}
