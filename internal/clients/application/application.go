package application

import (
	"context"
	"encoding/json"
	"errors"

	core "github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	"github.com/denniskniep/provider-springclouddataflow/internal/clients"
	"github.com/denniskniep/spring-cloud-dataflow-sdk-go/v2/client/apps"
	kiota "github.com/microsoft/kiota-abstractions-go"
)

const (
	errNotApplication = "managed resource is not a Application custom resource"
)

type ApplicationService struct {
	clients.DataFlowService
}

func NewApplicationService(configData []byte) (clients.Service[*core.Application, core.ApplicationParameters, core.ApplicationObservation, ApplicationCompare], error) {
	dataFlowService, err := clients.NewDataFlowService(configData)

	if err != nil {
		return nil, err
	}

	return &ApplicationService{
		*dataFlowService,
	}, nil
}

type ApplicationCompare struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Version        string `json:"version"`
	Uri            string `json:"uri"`
	DefaultVersion bool   `json:"defaultVersion"`
	BootVersion    string `json:"bootVersion"`
}

func (s *ApplicationService) GetSpec(app *core.Application) *core.ApplicationParameters {
	return &app.Spec.ForProvider
}

func (s *ApplicationService) GetStatus(app *core.Application) *core.ApplicationObservation {
	return &app.Status.AtProvider
}

func (s *ApplicationService) SetStatus(app *core.Application, status *core.ApplicationObservation) {
	app.Status.AtProvider = *status
}

func (s *ApplicationService) CreateUniqueIdentifier(spec *core.ApplicationParameters, status *core.ApplicationObservation) (*string, error) {
	uniqueId := spec.Type + "." + spec.Name + "." + spec.Version
	return &uniqueId, nil
}

func (s *ApplicationService) Create(ctx context.Context, app *core.ApplicationParameters) error {
	err := s.Client().Apps().ByType(app.Type).ByName(app.Name).ByVersion(app.Version).Post(ctx, &apps.ItemItemWithVersionItemRequestBuilderPostRequestConfiguration{
		QueryParameters: &apps.ItemItemWithVersionItemRequestBuilderPostQueryParameters{
			BootVersion: &app.BootVersion,
			Uri:         &app.Uri,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *ApplicationService) Update(ctx context.Context, app *core.ApplicationParameters) error {
	if app.DefaultVersion {
		err := s.Client().Apps().ByType(app.Type).ByName(app.Name).ByVersion(app.Version).Put(ctx, &apps.ItemItemWithVersionItemRequestBuilderPutRequestConfiguration{})
		var apiError *kiota.ApiError
		if errors.As(err, &apiError) && apiError.ResponseStatusCode == 404 {
			return nil
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ApplicationService) Describe(ctx context.Context, app *core.ApplicationParameters) (*core.ApplicationObservation, error) {
	result, err := s.Client().Apps().ByType(app.Type).ByName(app.Name).ByVersion(app.Version).Get(ctx, nil)

	var apiError *kiota.ApiError
	if errors.As(err, &apiError) && apiError.ResponseStatusCode == 404 {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var observed = core.ApplicationObservation{}
	err = json.Unmarshal(result, &observed)
	if err != nil {
		return nil, err
	}

	return &observed, nil
}

func (s *ApplicationService) Delete(ctx context.Context, app *core.ApplicationParameters) error {
	_, err := s.Client().Apps().ByType(app.Type).ByName(app.Name).ByVersion(app.Version).Delete(ctx, nil)

	var apiError *kiota.ApiError
	if errors.As(err, &apiError) && apiError.ResponseStatusCode == 404 {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}
