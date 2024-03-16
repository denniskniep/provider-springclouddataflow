package stream

import (
	"testing"

	"github.com/denniskniep/provider-springclouddataflow/internal/clients/application"
	"github.com/denniskniep/provider-springclouddataflow/internal/controllersdk"
)

func TestCreateStream(t *testing.T) {
	controllersdk.SkipIfIsShort(t)

	srvApp := application.TestNewApplicationService(t)
	srvStream := TestNewStreamService(t)

	sourceApp := application.TestMakeDefaultApplication("source", "Test020", "v1.0.0")
	controllersdk.TestCreateAndAssert(t, srvApp, sourceApp)

	sinkApp := application.TestMakeDefaultApplication("sink", "Test021", "v1.0.0")
	controllersdk.TestCreateAndAssert(t, srvApp, sinkApp)

	testStream := TestMakeDefaultStream("MyStream01", "MyDesc", "Test020 | Test021", false)
	controllersdk.TestCreateAndAssert(t, srvStream, testStream)

	controllersdk.TestDelete(t, srvStream, testStream)
	controllersdk.TestDelete(t, srvApp, sourceApp)
	controllersdk.TestDelete(t, srvApp, sinkApp)
}
