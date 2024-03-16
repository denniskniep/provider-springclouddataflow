package taskschedule

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	core "github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	"github.com/denniskniep/provider-springclouddataflow/internal/clients"
	"github.com/denniskniep/spring-cloud-dataflow-sdk-go/v2/client/tasks"
	kiota "github.com/microsoft/kiota-abstractions-go"
)

const (
	errConnecting      = "failed to connect"
	errNotTaskSchedule = "managed resource is not a TaskSchedule custom resource"
)

type TaskScheduleService struct {
	clients.DataFlowService
}

func NewTaskScheduleService(configData []byte) (*TaskScheduleService, error) {
	dataFlowService, err := clients.NewDataFlowService(configData)

	if err != nil {
		return nil, errors.Wrap(err, errConnecting)
	}

	return &TaskScheduleService{
		*dataFlowService,
	}, nil
}

type TaskScheduleCompare struct {
	ScheduleName       string  `json:"scheduleName"`
	TaskDefinitionName *string `json:"taskDefinitionName,omitempty"`
}

func (s *TaskScheduleService) GetSpec(taskdef *core.TaskSchedule) *core.TaskScheduleParameters {
	return &taskdef.Spec.ForProvider
}

func (s *TaskScheduleService) GetStatus(taskdef *core.TaskSchedule) *core.TaskScheduleObservation {
	return &taskdef.Status.AtProvider
}

func (s *TaskScheduleService) SetStatus(taskdef *core.TaskSchedule, status *core.TaskScheduleObservation) {
	taskdef.Status.AtProvider = *status
}

func (s *TaskScheduleService) CreateUniqueIdentifier(spec *core.TaskScheduleParameters, status *core.TaskScheduleObservation) (*string, error) {
	uniqueId := spec.ScheduleName
	return &uniqueId, nil
}

func (s *TaskScheduleService) Create(ctx context.Context, task *core.TaskScheduleParameters) error {

	properties := "scheduler.cron.expression=" + task.CronExpression
	if task.Properties != nil {
		properties = properties + "," + *task.Properties
	}

	err := s.Client().Tasks().Schedules().Post(ctx, &tasks.SchedulesRequestBuilderPostRequestConfiguration{
		QueryParameters: &tasks.SchedulesRequestBuilderPostQueryParameters{
			ScheduleName:       &task.ScheduleName,
			TaskDefinitionName: task.TaskDefinitionName,
			Platform:           &task.Platform,
			Arguments:          task.Arguments,
			Properties:         &properties,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *TaskScheduleService) Update(ctx context.Context, task *core.TaskScheduleParameters) error {
	return errors.New("Update of TaskSchedule not implemented - all properties are immutable!")
}

func (s *TaskScheduleService) Describe(ctx context.Context, task *core.TaskScheduleParameters) (*core.TaskScheduleObservation, error) {
	result, err := s.Client().Tasks().Schedules().BySchedulesId(task.ScheduleName).Get(ctx, nil)

	var apiError *kiota.ApiError
	if errors.As(err, &apiError) && apiError.ResponseStatusCode == 404 {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var observed = core.TaskScheduleObservation{}
	err = json.Unmarshal(result, &observed)
	if err != nil {
		return nil, err
	}

	return &observed, nil
}

func (s *TaskScheduleService) Delete(ctx context.Context, task *core.TaskScheduleParameters) error {
	_, err := s.Client().Tasks().Schedules().BySchedulesId(task.ScheduleName).Delete(ctx, nil)

	var apiError *kiota.ApiError
	if errors.As(err, &apiError) && apiError.ResponseStatusCode == 404 {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *TaskScheduleService) MakeCompare() *TaskScheduleCompare {
	return &TaskScheduleCompare{}
}
