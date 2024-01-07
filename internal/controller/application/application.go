package application

import (
	"context"
	"strconv"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane/provider-springclouddataflow/apis/core/v1alpha1"
	apisv1alpha1 "github.com/crossplane/provider-springclouddataflow/apis/v1alpha1"
	clients "github.com/crossplane/provider-springclouddataflow/internal/clients"
	"github.com/crossplane/provider-springclouddataflow/internal/features"
)

const (
	errNotApplication = "managed resource is not a Application custom resource"
	errTrackPCUsage   = "cannot track ProviderConfig usage"
	errGetPC          = "cannot get ProviderConfig"
	errGetCreds       = "cannot get credentials"

	errNewClient = "cannot create new Service"
	errDescribe  = "failed to describe Application resource"
	errCreate    = "failed to create Application resource"
	errUpdate    = "failed to update Application resource"
	errDelete    = "failed to delete Application resource"
	errMapping   = "failed to map Application resource"
)

// Setup adds a controller that reconciles Application managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	o.Logger.Info("Setup Controller: Application")
	name := managed.ControllerName(v1alpha1.ApplicationGroupKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1alpha1.StoreConfigGroupVersionKind))
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.ApplicationGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			kube:         mgr.GetClient(),
			usage:        resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
			newServiceFn: clients.NewApplicationService,
			logger:       o.Logger.WithValues("controller", name)}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithPollInterval(o.PollInterval),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithInitializers(),
		managed.WithReferenceResolver(managed.NewAPISimpleReferenceResolver(mgr.GetClient())),
		managed.WithConnectionPublishers(cps...))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		WithEventFilter(resource.DesiredStateChanged()).
		For(&v1alpha1.Application{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	kube         client.Client
	usage        resource.Tracker
	logger       logging.Logger
	newServiceFn func(creds []byte) (clients.ApplicationService, error)
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	logger := c.logger.WithValues("method", "connect")
	logger.Debug("Start Connect")
	cr, ok := mg.(*v1alpha1.Application)
	if !ok {
		return nil, errors.New(errNotApplication)
	}

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	pc := &apisv1alpha1.ProviderConfig{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: cr.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetPC)
	}

	cd := pc.Spec.Credentials
	data, err := resource.CommonCredentialExtractor(ctx, cd.Source, c.kube, cd.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	svc, err := c.newServiceFn(data)
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}
	logger.Debug("Connected")
	return &external{service: svc, logger: c.logger}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	// A 'client' used to connect to the external resource API. In practice this
	// would be something like an AWS SDK client.
	service clients.ApplicationService
	logger  logging.Logger
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	logger := c.logger.WithValues("method", "observe")
	logger.Debug("Start observe")
	cr, ok := mg.(*v1alpha1.Application)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotApplication)
	}

	c.logger.Debug("ExternalName: '" + meta.GetExternalName(cr) + "'")

	uniqueId := createUniqueIdentifier(&cr.Spec.ForProvider)

	observed, err := c.service.DescribeApplication(ctx, &cr.Spec.ForProvider)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errDescribe)
	}

	if observed == nil {
		c.logger.Debug("Managed resource '" + uniqueId + "' does not exist")
		return managed.ExternalObservation{
			ResourceExists:    false,
			ResourceUpToDate:  false,
			ConnectionDetails: managed.ConnectionDetails{},
		}, nil
	}

	c.logger.Debug("Found '" + uniqueId + "'")

	// Update Status
	cr.Status.AtProvider = *observed
	cr.SetConditions(xpv1.Available().WithMessage("Application exists"))

	observedCompareable, err := c.service.MapToApplicationCompare(observed)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errMapping)
	}

	specCompareable, err := c.service.MapToApplicationCompare(&cr.Spec.ForProvider)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errMapping)
	}

	diff := ""
	resourceUpToDate := cmp.Equal(specCompareable, observedCompareable)

	// Compare Spec with observed
	if !resourceUpToDate {
		diff = cmp.Diff(specCompareable, observedCompareable)
	}
	c.logger.Debug("Managed resource '" + uniqueId + "' upToDate: " + strconv.FormatBool(resourceUpToDate) + "")

	return managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        resourceUpToDate,
		Diff:                    diff,
		ResourceLateInitialized: false,
		ConnectionDetails:       managed.ConnectionDetails{},
	}, nil
}

func createUniqueIdentifier(app *v1alpha1.ApplicationParameters) string {
	return app.Type + "." + app.Name + "." + app.Version
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	logger := c.logger.WithValues("method", "create")
	logger.Debug("Start create")
	cr, ok := mg.(*v1alpha1.Application)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotApplication)
	}

	err := c.service.CreateApplication(ctx, &cr.Spec.ForProvider)

	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
	}

	uniqueId := createUniqueIdentifier(&cr.Spec.ForProvider)
	meta.SetExternalName(cr, uniqueId)
	c.logger.Debug("Managed resource '" + uniqueId + "' created")

	return managed.ExternalCreation{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	logger := c.logger.WithValues("method", "update")
	logger.Debug("Start update")
	cr, ok := mg.(*v1alpha1.Application)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotApplication)
	}

	err := c.service.UpdateApplication(ctx, &cr.Spec.ForProvider)

	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdate)
	}

	uniqueId := createUniqueIdentifier(&cr.Spec.ForProvider)
	c.logger.Debug("Managed resource '" + uniqueId + "' updated")
	return managed.ExternalUpdate{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	logger := c.logger.WithValues("method", "delete")
	logger.Debug("Start delete")
	cr, ok := mg.(*v1alpha1.Application)
	if !ok {
		return errors.New(errNotApplication)
	}

	err := c.service.DeleteApplication(ctx, &cr.Spec.ForProvider)

	if err != nil {
		return errors.Wrap(err, errDelete)
	}

	uniqueId := createUniqueIdentifier(&cr.Spec.ForProvider)
	c.logger.Debug("Managed resource '" + uniqueId + "' deleted")
	return nil
}
