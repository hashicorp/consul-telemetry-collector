package hcp

import (
	"errors"
	"testing"

	oErrors "github.com/go-openapi/errors"
	"github.com/google/uuid"
	"github.com/shoenig/test/must"

	"github.com/hashicorp/hcp-sdk-go/clients/cloud-global-network-manager-service/preview/2022-02-15/client/global_network_manager_service"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-global-network-manager-service/preview/2022-02-15/models"
	"github.com/hashicorp/hcp-sdk-go/resource"
)

func testResource() *resource.Resource {
	return &resource.Resource{
		ID:           uuid.NewString(),
		Type:         "type",
		Organization: uuid.NewString(),
		Project:      uuid.NewString(),
	}
}

func Test_ParseResource(t *testing.T) {
	for testname, tc := range map[string]struct {
		resourceString string
		expectedError  error
	}{
		"success": {resourceString: testResource().String()},
		"invalidURL": {
			resourceString: "foobar",
			expectedError:  errors.New("failed to parse resource_url could not parse resource: unexpected number of tokens 1"),
		},
	} {
		t.Run(testname, func(t *testing.T) {
			res, err := parseResource(tc.resourceString)
			if tc.expectedError != nil {
				must.Error(t, err)
				must.EqError(t, err, tc.expectedError.Error())
				return
			}
			must.NoError(t, err)
			must.NotNil(t, res)

		})
	}
}

func Test_ParseConfig(t *testing.T) {
	for testname, tc := range map[string]struct {
		p             *Params
		expectedError error
	}{
		"success": {p: &Params{uuid.NewString(), uuid.NewString(), testResource().String()}},
		"emptyclientid": {
			p:             &Params{"", uuid.NewString(), testResource().String()},
			expectedError: errors.New("client credentials are empty"),
		},
		"emptyclientsec": {
			p:             &Params{uuid.NewString(), "", testResource().String()},
			expectedError: errors.New("client credentials are empty"),
		},
	} {
		t.Run(testname, func(t *testing.T) {
			res, err := parseResource(tc.p.ResourceURL)
			must.NoError(t, err)
			config, err := parseConfig(tc.p, res)

			if tc.expectedError != nil {
				must.Error(t, err)
				must.EqError(t, err, tc.expectedError.Error())
				return
			}
			must.NoError(t, err)
			must.NotNil(t, config)

		})
	}
}

func Test_New(t *testing.T) {
	testcases := map[string]struct {
		cid     string
		csec    string
		res     *resource.Resource
		wantErr bool
	}{
		"Good": {
			cid:  uuid.NewString(),
			csec: uuid.NewString(),
			res:  testResource(),
		},
		"NoClientID": {
			cid:     "",
			csec:    uuid.NewString(),
			res:     testResource(),
			wantErr: true,
		},
		"NoClientSecret": {
			cid:     uuid.NewString(),
			csec:    "",
			res:     testResource(),
			wantErr: true,
		},
		"InvalidResource": {
			cid:     uuid.NewString(),
			csec:    uuid.NewString(),
			res:     testResource(),
			wantErr: true,
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			_, err := New(&Params{tc.cid, tc.csec, tc.res.String()})
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
		filters  []string
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
			filters:  []string{"a", "b"},
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
			filters:  []string{"a"},
		},
		"BadResponse": {
			r:   resource.Resource{},
			err: oErrors.New(500, "failed"),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			clientServiceM := &MockClientService{
				MockResponse: tc.resp,
				Err:          tc.err,
			}
			p := &Params{uuid.NewString(), uuid.NewString(), tc.r.String()}

			client, err := newClient(p, clientServiceM)

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

			filters, err := client.MetricFilters()
			must.NoError(t, err)
			must.Eq(t, tc.filters, filters)
		})
	}
}
