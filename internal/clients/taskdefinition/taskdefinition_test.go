package taskdefinition

import (
	"testing"

	core "github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	"github.com/denniskniep/provider-springclouddataflow/internal/clients"
	"github.com/denniskniep/provider-springclouddataflow/internal/clients/application"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/context"
)

func CreateTaskDefinitionService(t *testing.T) *TaskDefinitionService {
	jsonConfig := clients.GetJsonConfigForTests()

	srv, err := NewTaskDefinitionService([]byte(jsonConfig))
	if err != nil {
		t.Fatal(err)
	}

	return srv
}

func CreateDefaultTaskDefinition(name string, description string, definition string) *core.TaskDefinitionParameters {
	return &core.TaskDefinitionParameters{
		Name:        name,
		Description: description,
		Definition:  definition,
	}
}

func TestCreateTaskDefinition(t *testing.T) {
	skipIfIsShort(t)

	srvApp := application.C
	srvTask := CreateTaskDefinitionService(t)

	testApp := CreateDefaultApplication("task", "Test010", "v1.0.0")
	_ = CreateApplication(t, srvApp, testApp)

	testTask := CreateDefaultTaskDefinition("MyTask01", "MyDesc", "Test010")
	created := CreateTaskDefinition(t, srvTask, testTask)

	AssertTaskDefinitionAreEqual(t, srvTask, created, testTask)

	DeleteTaskDefinition(t, srvTask, testTask)
	DeleteApplication(t, srvApp, testApp)
}

func CreateTaskDefinition(t *testing.T, srv *TaskDefinitionService, task *core.TaskDefinitionParameters) *core.TaskDefinitionObservation {
	t.Helper()
	err := srv.Create(context.Background(), task)
	if err != nil {
		t.Fatal(err)
	}

	createdTask, err := srv.Describe(context.Background(), task)
	if err != nil {
		t.Fatal(err)
	}

	if createdTask == nil {
		t.Fatal("TaskDefinition was not found")
	}
	return createdTask
}

func DeleteTaskDefinition(t *testing.T, srv *TaskDefinitionService, task *core.TaskDefinitionParameters) {
	t.Helper()
	err := srv.Delete(context.Background(), task)
	if err != nil {
		t.Fatal(err)
	}

	noApp, err := srv.Describe(context.Background(), task)
	if err != nil {
		t.Fatal(err)
	}

	if noApp != nil {
		t.Fatal("TaskDefinition was not deleted")
	}
}

func AssertTaskDefinitionAreEqual(t *testing.T, srv *TaskDefinitionService, actual *core.TaskDefinitionObservation, expected *core.TaskDefinitionParameters) {
	t.Helper()
	mappedActual, err := srv.MapToCompare(actual)
	if err != nil {
		t.Fatal(err)
	}

	mappedExpected, err := srv.MapToCompare(expected)
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