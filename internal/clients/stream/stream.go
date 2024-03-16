package stream

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pkg/errors"

	"github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	core "github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1"
	"github.com/denniskniep/provider-springclouddataflow/internal/clients"
	"github.com/denniskniep/spring-cloud-dataflow-sdk-go/v2/client/streams"
	kiota "github.com/microsoft/kiota-abstractions-go"
)

const (
	errConnecting = "failed to connect"
	errNotStream  = "managed resource is not a Stream custom resource"
)

type StreamService struct {
	clients.DataFlowService
}

func NewStreamService(configData []byte) (clients.Service[*v1alpha1.Stream, v1alpha1.StreamParameters, v1alpha1.StreamObservation, StreamCompare], error) {
	dataFlowService, err := clients.NewDataFlowService(configData)

	if err != nil {
		return nil, errors.Wrap(err, errConnecting)
	}

	return &StreamService{
		*dataFlowService,
	}, nil
}

type StreamCompare struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Definition  string `json:"definition"`
}

type StreamDescribeResponse struct {
	Name              string `json:"name"`
	DslText           string `json:"dslText"`
	OriginalDslText   string `json:"originalDslText"`
	Status            string `json:"status"`
	Description       string `json:"description"`
	StatusDescription string `json:"statusDescription"`
}

func (s *StreamService) GetSpec(app *core.Stream) *core.StreamParameters {
	return &app.Spec.ForProvider
}

func (s *StreamService) GetStatus(app *core.Stream) *core.StreamObservation {
	return &app.Status.AtProvider
}

func (s *StreamService) SetStatus(app *core.Stream, status *core.StreamObservation) {
	app.Status.AtProvider = *status
}

func (s *StreamService) CreateUniqueIdentifier(spec *core.StreamParameters, status *core.StreamObservation) (*string, error) {
	uniqueId := spec.Name
	return &uniqueId, nil
}

func (s *StreamService) Create(ctx context.Context, stream *core.StreamParameters) error {
	err := s.Client().Streams().Definitions().Post(ctx, &streams.DefinitionsRequestBuilderPostRequestConfiguration{
		QueryParameters: &streams.DefinitionsRequestBuilderPostQueryParameters{
			Name:        &stream.Name,
			Description: &stream.Description,
			Definition:  &stream.Definition,
			Deploy:      &stream.Deploy,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *StreamService) Update(ctx context.Context, app *core.StreamParameters) error {
	return errors.New("Update of Stream not implemented - all properties are immutable!")
}

func (s *StreamService) Describe(ctx context.Context, stream *core.StreamParameters) (*core.StreamObservation, error) {
	result, err := s.Client().Streams().Definitions().ByName(stream.Name).Get(ctx, nil)

	var apiError *kiota.ApiError
	if errors.As(err, &apiError) && apiError.ResponseStatusCode == 404 {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var response = StreamDescribeResponse{}
	err = json.Unmarshal(result, &response)
	if err != nil {
		return nil, err
	}

	var observed = core.StreamObservation{
		Name:              response.Name,
		Description:       response.Description,
		Definition:        response.DslText,
		Status:            response.Status,
		StatusDescription: response.StatusDescription,
	}

	return &observed, nil
}

func (s *StreamService) Delete(ctx context.Context, stream *core.StreamParameters) error {
	_, err := s.Client().Streams().Definitions().ByName(stream.Name).Delete(ctx, nil)

	var apiError *kiota.ApiError
	if errors.As(err, &apiError) && apiError.ResponseStatusCode == 404 {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *StreamService) MakeCompare() *StreamCompare {
	return &StreamCompare{}
}

func TestNewStreamService(t *testing.T) clients.Service[*v1alpha1.Stream, v1alpha1.StreamParameters, v1alpha1.StreamObservation, StreamCompare] {
	jsonConfig := clients.GetJsonConfigForTests()

	srv, err := NewStreamService([]byte(jsonConfig))
	if err != nil {
		t.Fatal(err)
	}

	return srv
}

func TestMakeDefaultStream(name string, description string, definition string, deploy bool) *core.StreamParameters {
	return &core.StreamParameters{
		Name:        name,
		Description: description,
		Definition:  definition,
		Deploy:      deploy,
	}
}
