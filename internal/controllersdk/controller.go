package controllersdk

import (
	"context"
	"strconv"

	"encoding/json"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	apisv1alpha1 "github.com/denniskniep/provider-springclouddataflow/apis/v1alpha1"
	"github.com/denniskniep/provider-springclouddataflow/internal/clients"
	"github.com/denniskniep/provider-springclouddataflow/internal/features"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	errNotType      = "managed resource is not of spcific custom resource type"
	errTrackPCUsage = "cannot track ProviderConfig usage"
	errGetPC        = "cannot get ProviderConfig"
	errGetCreds     = "cannot get credentials"

	errNewClient       = "cannot create new Service"
	errDescribe        = "failed to describe resource"
	errCreate          = "failed to create resource"
	errUpdate          = "failed to update resource"
	errDelete          = "failed to delete  resource"
	errExtract         = "failed to extract Spec and Status from resource"
	errMapping         = "failed to map resource"
	errMappingObserved = "failed to map observed resource to compareable"
	errMappingSpec     = "failed to map spec resource to compareable"
	errCreateUniqueId  = "failed to create unique identifier for resource"
)

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type Connector[R resource.Managed] struct {
	Kube                client.Client
	Usage               resource.Tracker
	Logger              logging.Logger
	NewExternalClientFn func(conn *Connector[R], creds []byte) (managed.ExternalClient, error)
}

// Setup adds a controller that reconciles Application managed resources.
func Setup[R resource.Managed](name string, mgr ctrl.Manager, o controller.Options, newExternalClientFn func(conn *Connector[R], creds []byte) (managed.ExternalClient, error)) error {
	o.Logger.Info("Setup Controller: " + name)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1alpha1.StoreConfigGroupVersionKind))
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.ApplicationGroupVersionKind),
		managed.WithExternalConnecter(&Connector[R]{
			Kube:                mgr.GetClient(),
			Usage:               resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
			NewExternalClientFn: newExternalClientFn,
			Logger:              o.Logger.WithValues("controller", name)}),
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

func (c *Connector[R]) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	logger := c.Logger.WithValues("method", "connect")
	logger.Debug("Start Connect")

	cr, ok := mg.(R)
	if !ok {
		return nil, errors.New(errNotType)
	}

	if err := c.Usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	pc := &apisv1alpha1.ProviderConfig{}
	if err := c.Kube.Get(ctx, types.NamespacedName{Name: cr.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetPC)
	}

	cd := pc.Spec.Credentials
	data, err := resource.CommonCredentialExtractor(ctx, cd.Source, c.Kube, cd.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	externalClient, err := c.NewExternalClientFn(c, data)
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}

	logger.Debug("Connected")
	return externalClient, nil
}

func MapToCompare[C any](obj interface{}) (*C, error) {
	objJson, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var objCompare = *new(C)
	err = json.Unmarshal(objJson, &objCompare)
	if err != nil {
		return nil, err
	}

	return &objCompare, nil
}

func Observe[R resource.Managed, P any, O any, C any](ctx context.Context, logger logging.Logger, srv clients.Service[R, P, O, C], mg resource.Managed) (managed.ExternalObservation, error) {
	logger = logger.WithValues("method", "observe")
	logger.Debug("Start observe")

	cr, spec, status, err := cast(srv, mg)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errExtract)
	}
	crWithAssert := any(mg).(R)

	logger.Debug("ExternalName: '" + meta.GetExternalName(cr) + "'")

	uniqueId, err := srv.CreateUniqueIdentifier(spec, status)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errCreateUniqueId)
	}

	observed, err := srv.Describe(ctx, spec)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errDescribe)
	}

	if observed == nil {
		logger.Debug("Managed resource '" + *uniqueId + "' does not exist")
		return managed.ExternalObservation{
			ResourceExists:    false,
			ResourceUpToDate:  false,
			ConnectionDetails: managed.ConnectionDetails{},
		}, nil
	}

	logger.Debug("Found '" + *uniqueId + "'")

	// Update Status
	srv.SetStatus(crWithAssert, observed)
	cr.SetConditions(xpv1.Available().WithMessage("Application exists"))

	observedCompareable, err := MapToCompare[C](observed)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errMappingObserved)
	}

	specCompareable, err := MapToCompare[C](spec)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errMappingSpec)
	}

	diff := ""
	resourceUpToDate := cmp.Equal(specCompareable, observedCompareable)

	// Compare Spec with observed
	if !resourceUpToDate {
		diff = cmp.Diff(specCompareable, observedCompareable)
	}
	logger.Debug("Managed resource '" + *uniqueId + "' upToDate: " + strconv.FormatBool(resourceUpToDate) + "")

	return managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        resourceUpToDate,
		Diff:                    diff,
		ResourceLateInitialized: false,
		ConnectionDetails:       managed.ConnectionDetails{},
	}, nil
}

func Create[R resource.Managed, P any, O any, C any](ctx context.Context, logger logging.Logger, srv clients.Service[R, P, O, C], mg resource.Managed) (managed.ExternalCreation, error) {
	logger = logger.WithValues("method", "create")
	logger.Debug("Start create")
	cr, spec, status, err := cast(srv, mg)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errExtract)
	}

	err = srv.Create(ctx, spec)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
	}

	uniqueId, err := srv.CreateUniqueIdentifier(spec, status)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateUniqueId)
	}

	meta.SetExternalName(cr, *uniqueId)
	logger.Debug("Managed resource '" + *uniqueId + "' created")

	return managed.ExternalCreation{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func Update[R resource.Managed, P any, O any, C any](ctx context.Context, logger logging.Logger, srv clients.Service[R, P, O, C], mg resource.Managed) (managed.ExternalUpdate, error) {
	logger = logger.WithValues("method", "update")
	logger.Debug("Start update")

	_, spec, status, err := cast(srv, mg)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errExtract)
	}

	err = srv.Update(ctx, spec)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdate)
	}

	uniqueId, err := srv.CreateUniqueIdentifier(spec, status)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errCreateUniqueId)
	}

	logger.Debug("Managed resource '" + *uniqueId + "' updated")

	return managed.ExternalUpdate{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func Delete[R resource.Managed, P any, O any, C any](ctx context.Context, logger logging.Logger, srv clients.Service[R, P, O, C], mg resource.Managed) error {
	logger = logger.WithValues("method", "delete")
	logger.Debug("Start delete")

	_, spec, status, err := cast(srv, mg)
	if err != nil {
		return err
	}

	err = srv.Delete(ctx, spec)

	if err != nil {
		return errors.Wrap(err, errDelete)
	}

	uniqueId, err := srv.CreateUniqueIdentifier(spec, status)
	if err != nil {
		return errors.Wrap(err, errCreateUniqueId)
	}

	logger.Debug("Managed resource '" + *uniqueId + "' deleted")
	return nil
}

func cast[R resource.Managed, P any, O any, C any](srv clients.Service[R, P, O, C], mg resource.Managed) (resource.Managed, *P, *O, error) {
	cr, ok := mg.(R)
	if !ok {
		return nil, nil, nil, errors.New(errNotType)
	}

	return cr, srv.GetSpec(cr), srv.GetStatus(cr), nil
}
