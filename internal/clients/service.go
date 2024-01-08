package clients

import (
	"encoding/json"

	client "github.com/denniskniep/spring-cloud-dataflow-sdk-go/v2/client"
	auth "github.com/microsoft/kiota-abstractions-go/authentication"
	http "github.com/microsoft/kiota-http-go"
)

type DataFlowServiceConfig struct {
	Url string `json:"url"`
}

type DataFlowServiceImpl struct {
	client *client.DataFlowClient
}

func NewDataFlowService(configData []byte) (*DataFlowServiceImpl, error) {
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

	return &DataFlowServiceImpl{
		client: client,
	}, err
}

func NewApplicationService(configData []byte) (ApplicationService, error) {
	return NewDataFlowService(configData)
}

func NewTaskDefinitionService(configData []byte) (TaskDefinitionService, error) {
	return NewDataFlowService(configData)
}
