package application

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"

	"github.com/crossplane/crossplane-runtime/pkg/logging"

	"github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	"github.com/denniskniep/provider-springclouddataflow/internal/clients"
	"github.com/denniskniep/provider-springclouddataflow/internal/clients/taskdefinition"
	"github.com/denniskniep/provider-springclouddataflow/internal/controllersdk"
)

const (
	errConnecting = "failed to connect"
)

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	service clients.Service[*v1alpha1.TaskDefinition, v1alpha1.TaskDefinitionParameters, v1alpha1.TaskDefinitionObservation, taskdefinition.TaskDefinitionCompare]
	logger  logging.Logger
}

func newExternalClient[R resource.Managed](conn *controllersdk.Connector[R], creds []byte) (managed.ExternalClient, error) {
	taskDefinitionService, err := taskdefinition.NewTaskDefinitionService(creds)
	if err != nil {
		return nil, errors.Wrap(err, errConnecting)
	}

	return &external{
		service: taskDefinitionService,
		logger:  conn.Logger,
	}, nil
}

// Setup adds a controller that reconciles Application managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.TaskDefinitionGroupKind)
	return controllersdk.Setup[*v1alpha1.Application](name, mgr, o, newExternalClient[*v1alpha1.TaskDefinition])
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	return controllersdk.Observe[*v1alpha1.TaskDefinition, v1alpha1.TaskDefinitionParameters, v1alpha1.TaskDefinitionObservation, taskdefinition.TaskDefinitionCompare](ctx, c.logger, c.service, mg)
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	return controllersdk.Create[*v1alpha1.TaskDefinition, v1alpha1.TaskDefinitionParameters, v1alpha1.TaskDefinitionObservation, taskdefinition.TaskDefinitionCompare](ctx, c.logger, c.service, mg)
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	return controllersdk.Update[*v1alpha1.TaskDefinition, v1alpha1.TaskDefinitionParameters, v1alpha1.TaskDefinitionObservation, taskdefinition.TaskDefinitionCompare](ctx, c.logger, c.service, mg)
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	return controllersdk.Delete[*v1alpha1.TaskDefinition, v1alpha1.TaskDefinitionParameters, v1alpha1.TaskDefinitionObservation, taskdefinition.TaskDefinitionCompare](ctx, c.logger, c.service, mg)
}
