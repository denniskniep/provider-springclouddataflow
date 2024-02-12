package taskschedule

import (
	"testing"

	core "github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/context"
)

func CreateTaskScheduleService(t *testing.T) TaskScheduleService {
	jsonConfig := getJsonConfig()

	srv, err := NewTaskScheduleService([]byte(jsonConfig))
	if err != nil {
		t.Fatal(err)
	}

	return srv
}

func CreateDefaultTaskSchedule(scheduleName string, taskDefinitionName string) *core.TaskScheduleParameters {
	return &core.TaskScheduleParameters{
		ScheduleName:       scheduleName,
		TaskDefinitionName: &taskDefinitionName,
		CronExpression:     "0 * * * *",
		Platform:           "default",
	}
}

// Local Env is not able to add Schedules
/*func TestCreateTaskSchedule(t *testing.T) {
	skipIfIsShort(t)

	srvApp := CreateApplicationService(t)
	srvTask := CreateTaskDefinitionService(t)
	srvSchedule := CreateTaskScheduleService(t)

	testApp := CreateDefaultApplication("task", "Test040", "v1.0.0")
	_ = CreateApplication(t, srvApp, testApp)

	testTask := CreateDefaultTaskDefinition("MyTask30", "MyDesc", "Test040")
	_ = CreateTaskDefinition(t, srvTask, testTask)

	testSchedule := CreateDefaultTaskSchedule("schedule-2", "MyTask30")
	created := CreateTaskSchedule(t, srvSchedule, testSchedule)

	AssertTaskScheduleAreEqual(t, srvSchedule, created, testSchedule)

	DeleteTaskSchedule(t,srvSchedule,testSchedule)
	DeleteTaskDefinition(t, srvTask, testTask)
	DeleteApplication(t, srvApp, testApp)
}*/

func CreateTaskSchedule(t *testing.T, srv TaskScheduleService, task *core.TaskScheduleParameters) *core.TaskScheduleObservation {
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
		t.Fatal("TaskSchedule was not found")
	}
	return createdTask
}

func DeleteTaskSchedule(t *testing.T, srv TaskScheduleService, task *core.TaskScheduleParameters) {
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
		t.Fatal("TaskSchedule was not deleted")
	}
}

func AssertTaskScheduleAreEqual(t *testing.T, srv TaskScheduleService, actual *core.TaskScheduleObservation, expected *core.TaskScheduleParameters) {
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
