package clients

import (
	"context"
	"encoding/json"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	client "github.com/denniskniep/spring-cloud-dataflow-sdk-go/v2/client"
	auth "github.com/microsoft/kiota-abstractions-go/authentication"
	http "github.com/microsoft/kiota-http-go"
)

type DataFlowServiceConfig struct {
	Url string `json:"url"`
}

type DataFlowService struct {
	client *client.DataFlowClient
}

func (s *DataFlowService) Client() *client.DataFlowClient {
	return s.client
}

func NewDataFlowService(configData []byte) (*DataFlowService, error) {
	var conf = DataFlowServiceConfig{}
	err := json.Unmarshal(configData, &conf)
	if err != nil {
		return nil, err
	}

	// API requires no authentication, so use the anonymous
	// authentication provider
	authProvider := auth.AnonymousAuthenticationProvider{}

	// Create request adapter using the net/http-based implementation
	adapter, err := http.NewNetHttpRequestAdapter(&authProvider)
	if err != nil {
		return nil, err
	}

	adapter.SetBaseUrl(conf.Url)

	// Create the API client
	client := client.NewDataFlowClient(adapter)

	return &DataFlowService{
		client: client,
	}, err
}

// R=* (i.e Application)
// P=*Parameters (i.e ApplicationParameters)
// O=*Observation (i.e ApplicationObservation)
// C=*Compare (i.e ApplicationCompare)
type Service[R resource.Managed, P any, O any, C any] interface {
	Describe(ctx context.Context, param *P) (*O, error)

	Create(ctx context.Context, param *P) error
	Update(ctx context.Context, param *P) error
	Delete(ctx context.Context, param *P) error

	GetSpec(obj R) *P
	GetStatus(obj R) *O
	SetStatus(obj R, status *O)
	CreateUniqueIdentifier(*P, *O) (*string, error)
}

func GetJsonConfigForTests() string {
	return `{
		"url": "http://localhost:9393"
	}`
}
