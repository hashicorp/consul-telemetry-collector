package hcp

import (
	"testing"

	"github.com/go-openapi/errors"
	"github.com/google/uuid"
	"github.com/shoenig/test/must"

	"github.com/hashicorp/hcp-sdk-go/clients/cloud-global-network-manager-service/preview/2022-02-15/client/global_network_manager_service"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-global-network-manager-service/preview/2022-02-15/models"
	"github.com/hashicorp/hcp-sdk-go/resource"
)

func Test_New(t *testing.T) {
	testcases := map[string]struct {
		cid     string
		csec    string
		res     resource.Resource
		wantErr bool
	}{
		"Good": {
			cid:  uuid.NewString(),
			csec: uuid.NewString(),
			res: resource.Resource{
				ID:           uuid.NewString(),
				Type:         "type",
				Organization: uuid.NewString(),
				Project:      uuid.NewString(),
			},
		},
		"NoClientID": {
			cid:  "",
			csec: uuid.NewString(),
			res: resource.Resource{
				ID:           uuid.NewString(),
				Type:         "type",
				Organization: uuid.NewString(),
				Project:      uuid.NewString(),
			},
			wantErr: true,
		},
		"NoClientSecret": {
			cid:  uuid.NewString(),
			csec: "",
			res: resource.Resource{
				ID:           uuid.NewString(),
				Type:         "type",
				Organization: uuid.NewString(),
				Project:      uuid.NewString(),
			},
			wantErr: true,
		},
		"InvalidResource": {
			cid:  uuid.NewString(),
			csec: uuid.NewString(),
			res: resource.Resource{
				ID:           "",
				Type:         "type",
				Organization: "",
				Project:      uuid.NewString(),
			},
			wantErr: true,
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			_, err := New(tc.cid, tc.csec, tc.res.String())
			if tc.wantErr {
				must.Error(t, err)
				return
			}

			must.NoError(t, err)
		})
	}
}

func TestLoadConfig(t *testing.T) {
	testcases := map[string]struct {
		r        resource.Resource
		resp     *global_network_manager_service.AgentTelemetryConfigOK
		endpoint string
		err      error
	}{
		"GoodGlobalEndpoint": {
			resp: &global_network_manager_service.AgentTelemetryConfigOK{
				Payload: &models.
					HashicorpCloudGlobalNetworkManager20220215AgentTelemetryConfigResponse{
					TelemetryConfig: &models.HashicorpCloudGlobalNetworkManager20220215TelemetryConfig{
						Endpoint: "https://global.metrics.com",
						Labels:   nil,
						Metrics: &models.HashicorpCloudGlobalNetworkManager20220215TelemetryMetricsConfig{
							Endpoint:    "",
							IncludeList: []string{"a", "b"},
						},
					},
				},
			},
			endpoint: "https://global.metrics.com",
		},
		"GoodMetricsEndpoint": {
			resp: &global_network_manager_service.AgentTelemetryConfigOK{
				Payload: &models.
					HashicorpCloudGlobalNetworkManager20220215AgentTelemetryConfigResponse{
					TelemetryConfig: &models.HashicorpCloudGlobalNetworkManager20220215TelemetryConfig{
						Endpoint: "https://global.metrics.com",
						Labels:   nil,
						Metrics: &models.HashicorpCloudGlobalNetworkManager20220215TelemetryMetricsConfig{
							Endpoint:    "https://local.metrics.com",
							IncludeList: []string{"a"},
						},
					},
				},
			},
			endpoint: "https://local.metrics.com",
		},
		"BadResponse": {
			r:   resource.Resource{},
			err: errors.New(500, "failed"),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			clientServiceM := &MockClientService{
				MockResponse: tc.resp,
				Err:          tc.err,
			}

			client, err := NewWithDeps(uuid.NewString(), uuid.NewString(), tc.r.String(),
				WithClientService(clientServiceM))
			must.NoError(t, err)

			err = client.ReloadConfig()
			if tc.err != nil {
				must.Error(t, err)
				return
			}

			params := clientServiceM.params
			must.Eq(t, tc.r.ID, params.ClusterID)

			endpoint, err := client.MetricsEndpoint()
			must.NoError(t, err)
			must.Eq(t, tc.endpoint, endpoint)
		})
	}
}
