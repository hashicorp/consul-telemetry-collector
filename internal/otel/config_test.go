// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package otel

import (
	"context"
	"fmt"
	"testing"

	"github.com/shoenig/test"
	"github.com/shoenig/test/must"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/otelcol"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config/helpers/exporters"
	"github.com/hashicorp/hcp-sdk-go/resource"
)

func Test_newConfigProvider(t *testing.T) {
	testcases := map[string]struct {
		testfile    string
		exporter    *config.ExporterConfig
		hcpResource *resource.Resource
	}{
		"stock": {
			testfile: "stock.yaml",
		},
		"stock-with-forwarder": {
			testfile: "stock-with-forwarder.yaml",
			exporter: &config.ExporterConfig{
				ID: exporters.BaseOtlpExporterID,
				Exporter: &exporters.ExporterConfig{
					Endpoint: "https://test-forwarder-endpoint:4138",
					Headers: map[string]string{
						"authorization": "abc123",
					},
				},
			},
		},
		"stock-with-forwarder-grpc": {
			testfile: "stock-with-forwarder-grpc.yaml",
			exporter: &config.ExporterConfig{
				ID: exporters.GRPCOtlpExporterID,
				Exporter: &exporters.ExporterConfig{
					Endpoint: "https://test-forwarder-endpoint:4138",
					Headers: map[string]string{
						"authorization": "abc123",
					},
				},
			},
		},
		"hcp": {
			testfile: "hcp.yaml",
			hcpResource: &resource.Resource{
				ID:           "otel-cluster",
				Type:         "hashicorp.consul.cluster",
				Organization: "00000000-0000-0000-0000-000000000000",
				Project:      "00000000-0000-0000-0000-000000000001",
			},
		},
		"hcp-with-forwarder": {
			testfile: "hcp-with-forwarder.yaml",
			hcpResource: &resource.Resource{
				ID:           "otel-with-cluster",
				Type:         "hashicorp.consul.cluster",
				Organization: "00000000-0000-0000-0000-000000000003",
				Project:      "00000000-0000-0000-0000-000000000004",
			},
			exporter: &config.ExporterConfig{
				ID: exporters.BaseOtlpExporterID,
				Exporter: &exporters.ExporterConfig{
					Endpoint: "https://test-forwarder-endpoint:4138",
				},
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			var resourceURL string
			var mockClient *hcp.MockClient
			if tc.hcpResource != nil {
				resourceURL = tc.hcpResource.String()
				mockClient = &hcp.MockClient{
					MockMetricsEndpoint: "https://hcp-metrics-endpoint",
					MockMetricFilters:   []string{"^a", "b$"},
					MockMetricAttributes: map[string]string{
						"cluster": "name",
					},
				}
			}
			c := DefaultConfig()
			c.ClientID = "cid"
			c.ClientSecret = "csec"
			c.Client = mockClient
			c.ResourceID = resourceURL
			c.ExporterConfig = tc.exporter

			provider, err := newProvider(c)
			test.NoError(t, err)

			ctx := context.Background()

			// This provider.Get call will perform a configuration retrieval and ensure that it can be unmarshal'd in the
			// expected component config. To perform that Unmarshal we need the actual component code to unmarshal the map
			// [string]interface{} into the receiver/exporter/etc Config struct.
			factories, err := components()

			test.NoError(t, err)
			cfg, err := provider.Get(ctx, factories)

			must.NoError(t, err, must.Sprintf("%#v", cfg))
			must.NoError(t, cfg.Validate(), must.Sprint("provider configuration is invalid"))

			testprovider := testConfigProvider(t, []string{tc.testfile})
			golden, err := testprovider.Get(ctx, factories)
			must.NoError(t, err)
			must.NoError(t, golden.Validate(), must.Sprint("golden configuration is invalid"))

			compareComponents(t, golden.Receivers, cfg.Receivers, test.Sprint("receivers do not match"))
			compareComponents(t, golden.Exporters, cfg.Exporters, test.Sprint("exporters do not match"))
			compareComponents(t, golden.Processors, cfg.Processors, test.Sprint("processors do not match"))
			compareComponents(t, golden.Extensions, cfg.Extensions, test.Sprint("extensions do not match"))
			compareComponents(t, golden.Connectors, cfg.Connectors, test.Sprint("connectors do not match"))
			test.Eq(t, golden.Service.Telemetry, cfg.Service.Telemetry, test.Sprint("telemetry does not match"))
			test.Eq(t, golden.Service.Extensions, cfg.Service.Extensions, test.Sprint("extensions do not match"))
			for name, goldenPipelineConfig := range golden.Service.Pipelines {
				cfgPipelineConfig, ok := cfg.Service.Pipelines[name]
				test.True(t, ok, test.Sprintf("%s golden pipeline does not exist", name))
				test.Eq(t, goldenPipelineConfig, cfgPipelineConfig,
					test.Sprintf("%s golden pipeline does not match configuration", name))
			}
		})
	}
}

func compareComponents(t *testing.T, golden, components map[component.ID]component.Config, settings ...test.Setting) {
	t.Helper()
	for id, goldenCfg := range golden {
		cfg := components[id]
		test.Eq(t, goldenCfg, cfg, settings...)
	}
	for id := range components {
		_, ok := golden[id]
		test.True(t, ok, test.Sprintf("component %s does not exist in the golden config", id))
	}
	test.Eq(t, golden, components)
}

func testConfigProvider(t *testing.T, uris []string) otelcol.ConfigProvider {
	t.Helper()
	provider, err := otelcol.NewConfigProvider(otelcol.ConfigProviderSettings{
		ResolverSettings: confmap.ResolverSettings{
			URIs:       makeURIs(t, uris),
			Providers:  makeMapProvidersMap(fileprovider.New()),
			Converters: nil,
		},
	})
	must.NoError(t, err)
	return provider
}

func makeURIs(t *testing.T, files []string) []string {
	t.Helper()
	uris := make([]string, 0, len(files))
	for _, f := range files {
		uris = append(uris, fmt.Sprintf("file:testdata/%s", f))
	}
	return uris
}
