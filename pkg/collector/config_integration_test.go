//go:build integration
// +build integration

package collector

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/shoenig/test"
)

func Test_parseFile(t *testing.T) {
	var (
		clientid     = "id"
		clientsecret = "secret"
		endpoint     = "endpoint"
	)
	testcases := map[string]struct {
		config string
		expect *Config
		err    error
	}{
		"Empty": {
			config: `cloud {}`,
			// preset should always have non-nil values
			expect: &Config{
				Cloud: &Cloud{},
			},
		},
		"ClientId": {
			config: fmt.Sprintf(`cloud { client_id = "%s" }`, clientid),
			expect: &Config{
				Cloud: &Cloud{
					ClientID: &clientid,
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
					ClientID:     &clientid,
					ClientSecret: &clientsecret,
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
					ClientID:     &clientid,
					ClientSecret: &clientsecret,
				},
				HTTPCollectorEndpoint: &endpoint,
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
		"InvalidConfigFile": {
			config: `cloud = {}`,
			err:    errors.New("Unsupported argument"),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {

			testfile := writeConfig(t, []byte(tc.config))

			outputConfig, err := parseFile(&testfile)

			if tc.err != nil {
				test.Error(t, err)
				test.ErrorContains(t, err, tc.err.Error())
				return
			}
			test.NoError(t, err)
			test.Eq(t, outputConfig, tc.expect)

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
