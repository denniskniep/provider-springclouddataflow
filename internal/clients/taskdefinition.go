package clients

import (
	"context"
	"encoding/json"
	"errors"

	core "github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	"github.com/denniskniep/spring-cloud-dataflow-sdk-go/v2/client/tasks"
	kiota "github.com/microsoft/kiota-abstractions-go"
)

type TaskDefinitionService interface {
	DescribeTaskDefinition(ctx context.Context, app *core.TaskDefinitionParameters) (*core.TaskDefinitionObservation, error)

	CreateTaskDefinition(ctx context.Context, app *core.TaskDefinitionParameters) error
	UpdateTaskDefinition(ctx context.Context, app *core.TaskDefinitionParameters) error
	DeleteTaskDefinition(ctx context.Context, app *core.TaskDefinitionParameters) error

	MapToTaskDefinitionCompare(app interface{}) (*TaskDefinitionCompare, error)
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

func (s *DataFlowServiceImpl) MapToTaskDefinitionCompare(taskDefinition interface{}) (*TaskDefinitionCompare, error) {
	taskJson, err := json.Marshal(taskDefinition)
	if err != nil {
		return nil, err
	}

	var taskCompare = TaskDefinitionCompare{}
	err = json.Unmarshal(taskJson, &taskCompare)
	if err != nil {
		return nil, err
	}

	return &taskCompare, nil
}

func (s *DataFlowServiceImpl) CreateTaskDefinition(ctx context.Context, task *core.TaskDefinitionParameters) error {
	_, err := s.client.Tasks().Definitions().Post(ctx, &tasks.DefinitionsRequestBuilderPostRequestConfiguration{
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

func (s *DataFlowServiceImpl) UpdateTaskDefinition(ctx context.Context, task *core.TaskDefinitionParameters) error {
	return errors.New("Update not implemented - all properties are immutable!")
}

func (s *DataFlowServiceImpl) DescribeTaskDefinition(ctx context.Context, task *core.TaskDefinitionParameters) (*core.TaskDefinitionObservation, error) {
	result, err := s.client.Tasks().Definitions().ByName(task.Name).Get(ctx, nil)

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

func (s *DataFlowServiceImpl) DeleteTaskDefinition(ctx context.Context, task *core.TaskDefinitionParameters) error {
	_, err := s.client.Tasks().Definitions().ByName(task.Name).Delete(ctx, nil)

	var apiError *kiota.ApiError
	if errors.As(err, &apiError) && apiError.ResponseStatusCode == 404 {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}
