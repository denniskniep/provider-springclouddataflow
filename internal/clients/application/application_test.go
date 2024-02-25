package application

import (
	"testing"

	"github.com/denniskniep/provider-springclouddataflow/internal/controllersdk"
)

func TestCreateApplication(t *testing.T) {
	controllersdk.SkipIfIsShort(t)

	srv := TestNewApplicationService(t)
	testApp := TestMakeDefaultApplication("task", "Test001", "v1.0.0")

	controllersdk.TestCreateAndAssert(t, srv, testApp)
	controllersdk.TestDelete(t, srv, testApp)
}

func TestDeleteNotExisting(t *testing.T) {
	controllersdk.SkipIfIsShort(t)

	srv := TestNewApplicationService(t)
	testApp := TestMakeDefaultApplication("task", "Test002", "v1.0.0")

	controllersdk.TestDelete(t, srv, testApp)
}

func TestUpdateNotExisting(t *testing.T) {
	controllersdk.SkipIfIsShort(t)

	srv := TestNewApplicationService(t)
	testApp := TestMakeDefaultApplication("task", "Test003", "v1.0.0")

	controllersdk.TestUpdate(t, srv, testApp)
}

func TestUpdate(t *testing.T) {
	controllersdk.SkipIfIsShort(t)

	srv := TestNewApplicationService(t)
	testAppV1 := TestMakeDefaultApplication("task", "Test004", "v1.0.0")
	controllersdk.TestCreateAndAssert(t, srv, testAppV1)

	testAppV2 := TestMakeDefaultApplication("task", "Test004", "v2.0.0")
	testAppV2.DefaultVersion = false
	controllersdk.TestCreateAndAssert(t, srv, testAppV2)

	testAppV2.DefaultVersion = true
	controllersdk.TestUpdateAndAssert(t, srv, testAppV2)

	testAppV1.DefaultVersion = false
	controllersdk.TestUpdateAndAssert(t, srv, testAppV1)

	testAppV1.DefaultVersion = true
	controllersdk.TestUpdateAndAssert(t, srv, testAppV1)

	controllersdk.TestDelete(t, srv, testAppV1)
	controllersdk.TestDelete(t, srv, testAppV2)

}
