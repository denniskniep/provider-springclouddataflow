package clients

import (
	"context"
	"encoding/json"
	"errors"

	core "github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	"github.com/denniskniep/spring-cloud-dataflow-sdk-go/v2/client/apps"
	kiota "github.com/microsoft/kiota-abstractions-go"
)

type ApplicationService interface {
	DescribeApplication(ctx context.Context, app *core.ApplicationParameters) (*core.ApplicationObservation, error)

	CreateApplication(ctx context.Context, app *core.ApplicationParameters) error
	UpdateApplication(ctx context.Context, app *core.ApplicationParameters) error
	DeleteApplication(ctx context.Context, app *core.ApplicationParameters) error

	MapToApplicationCompare(app interface{}) (*ApplicationCompare, error)
}

type ApplicationCompare struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Version        string `json:"version"`
	Uri            string `json:"uri"`
	DefaultVersion bool   `json:"defaultVersion"`
	BootVersion    string `json:"bootVersion"`
}

func (s *DataFlowServiceImpl) MapToApplicationCompare(app interface{}) (*ApplicationCompare, error) {
	appJson, err := json.Marshal(app)
	if err != nil {
		return nil, err
	}

	var appCompare = ApplicationCompare{}
	err = json.Unmarshal(appJson, &appCompare)
	if err != nil {
		return nil, err
	}

	return &appCompare, nil
}

func (s *DataFlowServiceImpl) CreateApplication(ctx context.Context, app *core.ApplicationParameters) error {

	err := s.client.Apps().ByType(app.Type).ByName(app.Name).ByVersion(app.Version).Post(ctx, &apps.ItemItemWithVersionItemRequestBuilderPostRequestConfiguration{
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

func (s *DataFlowServiceImpl) UpdateApplication(ctx context.Context, app *core.ApplicationParameters) error {
	if app.DefaultVersion {
		err := s.client.Apps().ByType(app.Type).ByName(app.Name).ByVersion(app.Version).Put(ctx, &apps.ItemItemWithVersionItemRequestBuilderPutRequestConfiguration{})
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

func (s *DataFlowServiceImpl) DescribeApplication(ctx context.Context, app *core.ApplicationParameters) (*core.ApplicationObservation, error) {
	result, err := s.client.Apps().ByType(app.Type).ByName(app.Name).ByVersion(app.Version).Get(ctx, nil)

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

func (s *DataFlowServiceImpl) DeleteApplication(ctx context.Context, app *core.ApplicationParameters) error {
	_, err := s.client.Apps().ByType(app.Type).ByName(app.Name).ByVersion(app.Version).Delete(ctx, nil)

	var apiError *kiota.ApiError
	if errors.As(err, &apiError) && apiError.ResponseStatusCode == 404 {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}
