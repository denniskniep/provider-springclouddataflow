package taskdefinition

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pkg/errors"

	"github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	core "github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	"github.com/denniskniep/provider-springclouddataflow/internal/clients"
	"github.com/denniskniep/spring-cloud-dataflow-sdk-go/v2/client/tasks"
	kiota "github.com/microsoft/kiota-abstractions-go"
)

const (
	errNotTaskDefinition = "managed resource is not a TaskDefinition custom resource"
	errConnecting        = "failed to connect"
)

type TaskDefinitionService struct {
	clients.DataFlowService
}

func NewTaskDefinitionService(configData []byte) (clients.Service[*v1alpha1.TaskDefinition, v1alpha1.TaskDefinitionParameters, v1alpha1.TaskDefinitionObservation, TaskDefinitionCompare], error) {
	dataFlowService, err := clients.NewDataFlowService(configData)

	if err != nil {
		return nil, errors.Wrap(err, errConnecting)
	}

	return &TaskDefinitionService{
		*dataFlowService,
	}, nil
}

type TaskDefinitionCompare struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Definition  string `json:"definition"`
}

type TaskDefinitionDescribeResponse struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	DslText             string `json:"dslText"`
	Composed            bool   `json:"composed"`
	ComposedTaskElement bool   `json:"composedTaskElement"`
	Status              string `json:"status"`
}

func (s *TaskDefinitionService) GetSpec(taskdef *core.TaskDefinition) *core.TaskDefinitionParameters {
	return &taskdef.Spec.ForProvider
}

func (s *TaskDefinitionService) GetStatus(taskdef *core.TaskDefinition) *core.TaskDefinitionObservation {
	return &taskdef.Status.AtProvider
}

func (s *TaskDefinitionService) SetStatus(taskdef *core.TaskDefinition, status *core.TaskDefinitionObservation) {
	taskdef.Status.AtProvider = *status
}

func (s *TaskDefinitionService) CreateUniqueIdentifier(spec *core.TaskDefinitionParameters, status *core.TaskDefinitionObservation) (*string, error) {
	uniqueId := spec.Name
	return &uniqueId, nil
}

func (s *TaskDefinitionService) Create(ctx context.Context, task *core.TaskDefinitionParameters) error {
	_, err := s.Client().Tasks().Definitions().Post(ctx, &tasks.DefinitionsRequestBuilderPostRequestConfiguration{
		QueryParameters: &tasks.DefinitionsRequestBuilderPostQueryParameters{
			Name:        &task.Name,
			Description: &task.Description,
			Definition:  &task.Definition,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *TaskDefinitionService) Update(ctx context.Context, task *core.TaskDefinitionParameters) error {
	return errors.New("Update not implemented - all properties are immutable!")
}

func (s *TaskDefinitionService) Describe(ctx context.Context, task *core.TaskDefinitionParameters) (*core.TaskDefinitionObservation, error) {
	result, err := s.Client().Tasks().Definitions().ByName(task.Name).Get(ctx, nil)

	var apiError *kiota.ApiError
	if errors.As(err, &apiError) && apiError.ResponseStatusCode == 404 {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var response = TaskDefinitionDescribeResponse{}
	err = json.Unmarshal(result, &response)
	if err != nil {
		return nil, err
	}

	var observed = core.TaskDefinitionObservation{
		Name:                response.Name,
		Description:         response.Description,
		Definition:          response.DslText,
		Composed:            response.Composed,
		ComposedTaskElement: response.ComposedTaskElement,
		Status:              response.Status,
	}

	return &observed, nil
}

func (s *TaskDefinitionService) Delete(ctx context.Context, task *core.TaskDefinitionParameters) error {
	_, err := s.Client().Tasks().Definitions().ByName(task.Name).Delete(ctx, nil)

	var apiError *kiota.ApiError
	if errors.As(err, &apiError) && apiError.ResponseStatusCode == 404 {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *TaskDefinitionService) MakeCompare() *TaskDefinitionCompare {
	return &TaskDefinitionCompare{}
}

func TestNewTaskDefinitionService(t *testing.T) clients.Service[*v1alpha1.TaskDefinition, v1alpha1.TaskDefinitionParameters, v1alpha1.TaskDefinitionObservation, TaskDefinitionCompare] {
	jsonConfig := clients.GetJsonConfigForTests()

	srv, err := NewTaskDefinitionService([]byte(jsonConfig))
	if err != nil {
		t.Fatal(err)
	}

	return srv
}

func TestMakeDefaultTaskDefinition(name string, description string, definition string) *core.TaskDefinitionParameters {
	return &core.TaskDefinitionParameters{
		Name:        name,
		Description: description,
		Definition:  definition,
	}
}
