// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package agent

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shoenig/test"
)

func Test_Validation(t *testing.T) {
	endpoint, cid, csec, crid := "endpoint", "cid", "csec", "crid"

	for name, tc := range map[string]struct {
		input       *Config
		err         error
		errContains string
	}{
		"FailNoConfig": {
			err: errNoConfigurationProvided,
		},
		"FailCloudIDOnlySpecified": {
			input: &Config{
				Cloud: &Cloud{
					ClientID: cid,
				},
			},
			err: errCloudConfigInvalid,
		},
		"FailCloudSecOnlySpecified": {
			input: &Config{
				Cloud: &Cloud{
					ClientSecret: csec,
				},
			},
			err: errCloudConfigInvalid,
		},
		"FailCloudResourceIdOnlySpecified": {
			input: &Config{
				Cloud: &Cloud{
					ResourceID: crid,
				},
			},
			err:         errCloudConfigInvalid,
			errContains: "missing client_id, client_secret",
		},
		"FailCloudResourceMissingClientID": {
			input: &Config{
				Cloud: &Cloud{
					ClientSecret: csec,
					ResourceID:   crid,
				},
			},
			err:         errCloudConfigInvalid,
			errContains: "missing client_id",
		},
		"FailCloudResourceMissingResourceID": {
			input: &Config{
				Cloud: &Cloud{
					ClientSecret: csec,
					ClientID:     cid,
				},
			},
			err:         errCloudConfigInvalid,
			errContains: "missing resource_id",
		},
		"FailCloudResourceMissingClientSecret": {
			input: &Config{
				Cloud: &Cloud{
					ResourceID: crid,
					ClientID:   cid,
				},
			},
			err: errCloudConfigInvalid,
		},
		"SuccessfulCloudNotSpecified": {
			input: &Config{
				Cloud: &Cloud{},
			},
		},
		"SuccessfulCloudNotSpecifiedAndOptionalOtel": {
			input: &Config{
				Cloud:                 &Cloud{},
				HTTPCollectorEndpoint: endpoint,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			err := tc.input.validate()
			if tc.err != nil {
				test.Error(t, err)
				test.ErrorIs(t, err, tc.err)
				test.ErrorContains(t, err, tc.errContains)
				return
			}
			test.NoError(t, err)
		})
	}
}

func Test_ReadFile(t *testing.T) {
	var (
		clientid     = "id"
		clientsecret = "secret"
		endpoint     = "endpoint"
	)
	testcases := map[string]struct {
		config string
		expect *Config
		err    error
		json   bool
	}{
		"EmptyHCL": {
			config: `cloud {}`,
			// preset should always have non-nil values
			expect: &Config{
				Cloud: &Cloud{},
			},
		},
		"EmptyJSON": {
			json:   true,
			config: `{"cloud":{}}`,
			// preset should always have non-nil values
			expect: &Config{
				Cloud: &Cloud{},
			},
		},
		"ClientIDJson": {
			json:   true,
			config: fmt.Sprintf(`{"cloud":{"client_id":%q}}`, clientid),
			expect: &Config{
				Cloud: &Cloud{
					ClientID: clientid,
				},
			},
		},
		"ClientId": {
			config: fmt.Sprintf(`cloud {client_id = %q }`, clientid),
			expect: &Config{
				Cloud: &Cloud{
					ClientID: clientid,
				},
			},
		},
		"ClientIDSecretJson": {
			json: true,
			config: fmt.Sprintf(`{
				"cloud": {
					"client_id": "%s",
					"client_secret": "%s"
				}
			}`, clientid, clientsecret),
			expect: &Config{
				Cloud: &Cloud{
					ClientID:     clientid,
					ClientSecret: clientsecret,
				},
			},
		},
		"ClientIdAndSecretId": {
			config: fmt.Sprintf(`cloud {
				client_id = "%s"
				client_secret = "%s"
			}`, clientid, clientsecret),
			expect: &Config{
				Cloud: &Cloud{
					ClientID:     clientid,
					ClientSecret: clientsecret,
				},
			},
		},
		"MinimalExporterConfig": {
			config: `
				exporter_config "otelgrpc" {
					endpoint = "http://otel:3749"
				}
			`,
			expect: &Config{
				ExporterConfig: &ExporterConfig{
					Type:     "otelgrpc",
					Endpoint: "http://otel:3749",
				},
			},
		},
		"FullExporterConfig": {
			config: `
				exporter_config "otelgrpc" {
					endpoint = "http://otel:3749"
					timeout = "10s"
					headers = {
						a = "b"
					}
				}
			`,
			expect: &Config{
				ExporterConfig: &ExporterConfig{
					Type:     "otelgrpc",
					Endpoint: "http://otel:3749",
					Timeout:  "10s",
					Headers: map[string]string{
						"a": "b",
					},
				},
			},
		},
		"AllFieldsJson": {
			json: true,
			config: fmt.Sprintf(`{
			"cloud": {
				"client_id": "%s",
				"client_secret": "%s"
			},
			"http_collector_endpoint": "%s",
			"exporter_config": {
				"otelhttp": {
					"endpoint": "%s",
					"headers": {
						"a": "b"
					},
					"timeout": "10s"
				}
			}
			}`, clientid, clientsecret, endpoint, endpoint),
			expect: &Config{
				Cloud: &Cloud{
					ClientID:     clientid,
					ClientSecret: clientsecret,
				},
				HTTPCollectorEndpoint: endpoint,
				ExporterConfig: &ExporterConfig{
					Type:     "otelhttp",
					Endpoint: endpoint,
					Headers:  map[string]string{"a": "b"},
					Timeout:  "10s",
				},
			},
		},
		"AllFields": {
			config: fmt.Sprintf(`
			cloud {
				client_id = "%s"
				client_secret = "%s"
			}
			http_collector_endpoint = "%s"
			`, clientid, clientsecret, endpoint),
			expect: &Config{
				Cloud: &Cloud{
					ClientID:     clientid,
					ClientSecret: clientsecret,
				},
				HTTPCollectorEndpoint: endpoint,
			},
		},
		"InvalidHCL": {
			config: fmt.Sprintf(`
			cloud {
				client_id = "%s"
				client_secret = "%s"
			}http_collector_endpoint = "%s"
			`, clientid, clientsecret, endpoint),
			err: errors.New("Missing newline after block definition; A block definition must end with a newline."),
		},
		"InvalidJson": {
			json: true,
			config: fmt.Sprintf(`{
			"http_collector_endpoint" = "%s"
			}`, endpoint),
			err: errors.New("Missing property value colon; JSON uses a colon as its name/value delimiter, not an equals sign."),
		},
		"InvalidConfigFile": {
			config: `cloud = {}`,
			err:    errors.New(`Unsupported argument; An argument named "cloud" is not expected here. Did you mean to define a block of type "cloud"?`),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			r := strings.NewReader(tc.config)
			filename := "config.hcl"
			if tc.json {
				filename = "config.json"
			}
			outputConfig, err := readConfiguration(r, filename)

			if tc.err != nil {
				test.Error(t, err)
				test.ErrorContains(t, err, tc.err.Error())
				return
			}

			diff := cmp.Diff(outputConfig, tc.expect, cmp.AllowUnexported(ExporterConfig{}))
			test.NoError(t, err)
			test.Eq(t, outputConfig, tc.expect, test.Sprint(diff))
		})
	}
}
