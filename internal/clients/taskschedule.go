package clients

import (
	"context"
	"encoding/json"
	"errors"

	core "github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	"github.com/denniskniep/spring-cloud-dataflow-sdk-go/v2/client/tasks"
	kiota "github.com/microsoft/kiota-abstractions-go"
)

type TaskScheduleService interface {
	DescribeTaskSchedule(ctx context.Context, app *core.TaskScheduleParameters) (*core.TaskScheduleObservation, error)

	CreateTaskSchedule(ctx context.Context, app *core.TaskScheduleParameters) error
	UpdateTaskSchedule(ctx context.Context, app *core.TaskScheduleParameters) error
	DeleteTaskSchedule(ctx context.Context, app *core.TaskScheduleParameters) error

	MapToTaskScheduleCompare(app interface{}) (*TaskScheduleCompare, error)
}

type TaskScheduleCompare struct {
	ScheduleName       string  `json:"scheduleName"`
	TaskDefinitionName *string `json:"taskDefinitionName,omitempty"`
}

func (s *DataFlowServiceImpl) MapToTaskScheduleCompare(taskSchedule interface{}) (*TaskScheduleCompare, error) {
	taskJson, err := json.Marshal(taskSchedule)
	if err != nil {
		return nil, err
	}

	var taskCompare = TaskScheduleCompare{}
	err = json.Unmarshal(taskJson, &taskCompare)
	if err != nil {
		return nil, err
	}

	return &taskCompare, nil
}

func (s *DataFlowServiceImpl) CreateTaskSchedule(ctx context.Context, task *core.TaskScheduleParameters) error {

	properties := "scheduler.cron.expression=" + task.CronExpression
	if task.Properties != nil {
		properties = properties + "," + *task.Properties
	}

	err := s.client.Tasks().Schedules().Post(ctx, &tasks.SchedulesRequestBuilderPostRequestConfiguration{
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

func (s *DataFlowServiceImpl) UpdateTaskSchedule(ctx context.Context, task *core.TaskScheduleParameters) error {
	return errors.New("Update not implemented - all properties are immutable!")
}

func (s *DataFlowServiceImpl) DescribeTaskSchedule(ctx context.Context, task *core.TaskScheduleParameters) (*core.TaskScheduleObservation, error) {
	result, err := s.client.Tasks().Schedules().BySchedulesId(task.ScheduleName).Get(ctx, nil)

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

func (s *DataFlowServiceImpl) DeleteTaskSchedule(ctx context.Context, task *core.TaskScheduleParameters) error {
	_, err := s.client.Tasks().Schedules().BySchedulesId(task.ScheduleName).Delete(ctx, nil)

	var apiError *kiota.ApiError
	if errors.As(err, &apiError) && apiError.ResponseStatusCode == 404 {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}
