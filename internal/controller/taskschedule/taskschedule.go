package taskschedule

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane/crossplane-runtime/pkg/logging"

	"github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	"github.com/denniskniep/provider-springclouddataflow/internal/clients"
	"github.com/denniskniep/provider-springclouddataflow/internal/clients/taskschedule"
	"github.com/denniskniep/provider-springclouddataflow/internal/controllersdk"
)

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	service clients.Service[*v1alpha1.TaskSchedule, v1alpha1.TaskScheduleParameters, v1alpha1.TaskScheduleObservation, taskschedule.TaskScheduleCompare]
	logger  logging.Logger
}

func newExternalClient[R resource.Managed](conn *controllersdk.Connector[R], creds []byte) (managed.ExternalClient, error) {
	service, err := taskschedule.NewTaskScheduleService(creds)
	if err != nil {
		return nil, err
	}

	return &external{
		service: service,
		logger:  conn.Logger,
	}, nil
}

func Setup(mgr ctrl.Manager, o controller.Options) error {
	return controllersdk.Setup(v1alpha1.TaskScheduleGroupVersionKind, &v1alpha1.TaskSchedule{}, mgr, o, newExternalClient[*v1alpha1.TaskSchedule])
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	return controllersdk.Observe(ctx, c.logger, c.service, mg)
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	return controllersdk.Create(ctx, c.logger, c.service, mg)
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	return controllersdk.Update(ctx, c.logger, c.service, mg)
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	return controllersdk.Delete(ctx, c.logger, c.service, mg)
}
