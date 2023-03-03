package collector

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shoenig/test"
)

func Test_LoadConfig(t *testing.T) {
	testcases := map[string]struct {
		preset Config
		config string
		expect Config
		err    error
	}{
		"Empty": {
			expect: Config{
				Cloud: &Cloud{},
			},
		},
		"PresetAndConfigOverride": {
			preset: Config{
				Cloud: &Cloud{
					ClientID:     "client_id",
					ClientSecret: "",
					ResourceID:   "resource-id",
				},
			},
			config: `cloud { client_id = "secret-client-id" }`,
			expect: Config{
				Cloud: &Cloud{
					ClientID:     "secret-client-id",
					ClientSecret: "",
					ResourceID:   "resource-id",
				},
				HTTPCollectorEndpoint: "",
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
			var testfile string
			if tc.config != "" {
				testfile = writeConfig(t, []byte(tc.config))
			}

			preset := tc.preset.clone()
			test.Eq(t, tc.preset, preset)

			preset.ConfigFile = testfile

			err := LoadConfig(testfile, &preset)
			if tc.err != nil {
				test.Error(t, err)
				return
			}
			test.NoError(t, err)
			test.Eq(t, tc.expect, preset, test.Cmp(cmpopts.IgnoreFields(Config{}, "ConfigFile")))
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
