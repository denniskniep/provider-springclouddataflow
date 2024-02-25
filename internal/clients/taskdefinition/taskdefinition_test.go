package taskdefinition

import (
	"testing"

	"github.com/denniskniep/provider-springclouddataflow/internal/clients/application"
	"github.com/denniskniep/provider-springclouddataflow/internal/controllersdk"
)

func TestCreateTaskDefinition(t *testing.T) {
	controllersdk.SkipIfIsShort(t)

	srvApp := application.TestNewApplicationService(t)
	srvTask := TestNewTaskDefinitionService(t)

	testApp := application.TestMakeDefaultApplication("task", "Test010", "v1.0.0")
	controllersdk.TestCreateAndAssert(t, srvApp, testApp)

	testTask := TestMakeDefaultTaskDefinition("MyTask01", "MyDesc", "Test010")
	controllersdk.TestCreateAndAssert(t, srvTask, testTask)

	controllersdk.TestDelete(t, srvTask, testTask)
	controllersdk.TestDelete(t, srvApp, testApp)
}
