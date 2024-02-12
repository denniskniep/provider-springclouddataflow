package application

import (
	"testing"

	core "github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	"github.com/denniskniep/provider-springclouddataflow/internal/clients"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/context"
)

func CreateApplicationService(t *testing.T) *ApplicationService {
	jsonConfig := clients.GetJsonConfigForTests()

	srv, err := NewApplicationService([]byte(jsonConfig))
	if err != nil {
		t.Fatal(err)
	}

	return srv
}

func CreateDefaultApplication(appType string, name string, version string) *core.ApplicationParameters {
	return &core.ApplicationParameters{
		Type:           appType,
		Name:           name,
		Version:        version,
		Uri:            "docker://hello-world:" + version,
		BootVersion:    "2",
		DefaultVersion: true,
	}
}

func TestCreateApplication(t *testing.T) {
	skipIfIsShort(t)

	srv := CreateApplicationService(t)
	testApp := CreateDefaultApplication("task", "Test001", "v1.0.0")

	createdApp := CreateApplication(t, srv, testApp)

	AssertApplicationAreEqual(t, srv, createdApp, testApp)

	DeleteApplication(t, srv, testApp)
}

func TestDeleteNotExisting(t *testing.T) {
	skipIfIsShort(t)
	srv := CreateApplicationService(t)
	testApp := CreateDefaultApplication("task", "Test999", "v1.0.0")
	DeleteApplication(t, srv, testApp)
}

func TestUpdateNotExisting(t *testing.T) {
	skipIfIsShort(t)
	srv := CreateApplicationService(t)
	testApp := CreateDefaultApplication("task", "Test888", "v1.0.0")
	UpdateApplication(t, srv, testApp)
}

func TestUpdate(t *testing.T) {
	skipIfIsShort(t)

	srv := CreateApplicationService(t)
	testAppV1 := CreateDefaultApplication("task", "Test002", "v1.0.0")
	createdAppV1 := CreateApplication(t, srv, testAppV1)
	AssertApplicationAreEqual(t, srv, createdAppV1, testAppV1)

	testAppV2 := CreateDefaultApplication("task", "Test002", "v2.0.0")
	testAppV2.DefaultVersion = false
	createdAppV2 := CreateApplication(t, srv, testAppV2)
	AssertApplicationAreEqual(t, srv, createdAppV2, testAppV2)

	testAppV1.DefaultVersion = false
	testAppV2.DefaultVersion = true
	UpdateApplication(t, srv, testAppV1)
	UpdateApplication(t, srv, testAppV2)

	foundAppV1, err := srv.Describe(context.Background(), testAppV1)
	if err != nil {
		t.Fatal(err)
	}

	AssertApplicationAreEqual(t, srv, foundAppV1, testAppV1)

	foundAppV2, err := srv.Describe(context.Background(), testAppV2)
	if err != nil {
		t.Fatal(err)
	}

	AssertApplicationAreEqual(t, srv, foundAppV2, testAppV2)

	DeleteApplication(t, srv, testAppV1)
	DeleteApplication(t, srv, testAppV2)
}

func UpdateApplication(t *testing.T, srv *ApplicationService, app *core.ApplicationParameters) {
	t.Helper()
	err := srv.Update(context.Background(), app)
	if err != nil {
		t.Fatal(err)
	}
}

func CreateApplication(t *testing.T, srv *ApplicationService, app *core.ApplicationParameters) *core.ApplicationObservation {
	t.Helper()
	err := srv.Create(context.Background(), app)
	if err != nil {
		t.Fatal(err)
	}

	createdApp, err := srv.Describe(context.Background(), app)
	if err != nil {
		t.Fatal(err)
	}

	if createdApp == nil {
		t.Fatal("App was not found")
	}
	return createdApp
}

func DeleteApplication(t *testing.T, srv *ApplicationService, app *core.ApplicationParameters) {
	t.Helper()
	err := srv.Delete(context.Background(), app)
	if err != nil {
		t.Fatal(err)
	}

	noApp, err := srv.Describe(context.Background(), app)
	if err != nil {
		t.Fatal(err)
	}

	if noApp != nil {
		t.Fatal("App was not deleted")
	}
}

func AssertApplicationAreEqual(t *testing.T, srv *ApplicationService, actual *core.ApplicationObservation, expected *core.ApplicationParameters) {
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
