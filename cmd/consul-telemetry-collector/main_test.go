package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shoenig/test"

	"github.com/hashicorp/consul-telemetry-collector/pkg/collector"
)

func Test_LoadConfig(t *testing.T) {
	testcases := map[string]struct {
		flags  []string
		env    map[string]string
		config string
		expect collector.Config
		err    error
	}{
		"Empty": {
			expect: collector.Config{
				Cloud: &collector.Cloud{},
			},
		},
		"EmptyConfig": {
			expect: collector.Config{
				Cloud: &collector.Cloud{},
			},
			config: `http_collector_endpoint=""`,
		},
		"FlagsOnly": {
			flags: []string{"-hcp-client-id", "ID"},
			expect: collector.Config{
				Cloud: &collector.Cloud{
					ClientID: "ID",
				},
			},
		},
		"EnvOnly": {
			env: map[string]string{
				"HCP_CLIENT_ID":     "ID",
				"HCP_CLIENT_SECRET": "SECRET",
			},
			expect: collector.Config{
				Cloud: &collector.Cloud{
					ClientID:     "ID",
					ClientSecret: "SECRET",
				},
				HTTPCollectorEndpoint: "",
				ConfigFile:            "",
			},
		},
		"FlagsAndEnv": {
			env: map[string]string{
				"HCP_CLIENT_ID":     "ID",
				"HCP_CLIENT_SECRET": "SECRET",
			},
			flags: []string{"-hcp-resource-id", "resource-id", "-hcp-client-id", "client_id"},
			expect: collector.Config{
				Cloud: &collector.Cloud{
					ClientID:     "client_id",
					ClientSecret: "SECRET",
					ResourceID:   "resource-id",
				},
				HTTPCollectorEndpoint: "",
				ConfigFile:            "",
			},
		},
		"FlagsAndConfig": {
			flags:  []string{"-hcp-client-id", "client_id", "-hcp-resource-id", "resource-id"},
			config: `cloud { client_id = "secret-client-id" }`,
			expect: collector.Config{
				Cloud: &collector.Cloud{
					ClientID:     "secret-client-id",
					ClientSecret: "",
					ResourceID:   "resource-id",
				},
				HTTPCollectorEndpoint: "",
				ConfigFile:            "",
			},
		},
		"EnvAndConfig": {
			env: map[string]string{
				"HCP_RESOURCE_ID": "resource_id",
				"HCP_CLIENT_ID":   "foo",
			},
			config: `
cloud { 
	client_id = "hcp-client-id" 
	client_secret ="hcp-client-secret" 
}`,
			expect: collector.Config{
				Cloud: &collector.Cloud{
					ClientID:     "hcp-client-id",
					ClientSecret: "hcp-client-secret",
					ResourceID:   "resource_id",
				},
				HTTPCollectorEndpoint: "",
				ConfigFile:            "",
			},
		},
		"FlagsEnvAndConfig": {
			env: map[string]string{
				"HCP_RESOURCE_ID": "resource_id",
				"HCP_CLIENT_ID":   "foo",
			},
			config: `
cloud { 
	client_id = "hcp-client-id" 
	client_secret ="hcp-client-secret" 
}`,
			flags: []string{"-http-collector-endpoint", "http://localhost:5000"},
			expect: collector.Config{
				Cloud: &collector.Cloud{
					ClientID:     "hcp-client-id",
					ClientSecret: "hcp-client-secret",
					ResourceID:   "resource_id",
				},
				HTTPCollectorEndpoint: "http://localhost:5000",
				ConfigFile:            "",
			},
		},
		"InvalidConfigFile": {
			config: `cloud = {}`,
			err:    errors.New("unsupported argument"),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			setupEnv(t, tc.env)

			// This lets us parse the flags from our test instead of the default global one
			fs := flag.NewFlagSet(name, flag.ExitOnError)
			flag.CommandLine = fs
			flags()
			test.NoError(t, fs.Parse(tc.flags))

			var testfile string
			if tc.config != "" {
				testfile = writeConfig(t, []byte(tc.config))
			}
			configFile = testfile

			cfg, err := loadConfig()
			if tc.err != nil {
				test.Error(t, err)
				return
			}
			test.NoError(t, err)
			test.Eq(t, tc.expect, cfg, test.Cmp(cmpopts.IgnoreFields(collector.Config{}, "ConfigFile")))
		})
	}
}

func writeConfig(t *testing.T, contents []byte) string {

	t.Helper()
	td := t.TempDir()
	f := "config.hcl"
	testfile := filepath.Join(td, f)
	test.NoError(t, os.WriteFile(testfile, contents, 0644))
	return testfile
}
