package collector

import (
	"testing"

	"github.com/shoenig/test"
)

func Test_Validation(t *testing.T) {
	endpoint, cid, csec, crid := "endpoint", "cid", "csec", "crid"
	for name, tc := range map[string]struct {
		input *Config
		err   error
	}{
		"FailNoConfig": {
			err: errNoConfigurationProvided,
		},
		"FailNoCollectorEndpoint": {
			input: &Config{},
			err:   errNoCollectorEndpoint,
		},
		"FailCloudIDOnlySpecified": {
			input: &Config{
				HTTPCollectorEndpoint: &endpoint,
				Cloud: &Cloud{
					ClientID: &cid,
				},
			},
			err: errCloudConfigInvalid,
		},
		"FailCloudSecOnlySpecified": {
			input: &Config{
				HTTPCollectorEndpoint: &endpoint,
				Cloud: &Cloud{
					ClientSecret: &csec,
				},
			},
			err: errCloudConfigInvalid,
		},
		"FailCloudResourceIdOnlySpecified": {
			input: &Config{
				HTTPCollectorEndpoint: &endpoint,
				Cloud: &Cloud{
					ResourceID: &crid,
				},
			},
			err: errCloudConfigInvalid,
		},
		"FailCloudResourceMissingClientID": {
			input: &Config{
				HTTPCollectorEndpoint: &endpoint,
				Cloud: &Cloud{
					ClientSecret: &csec,
					ResourceID:   &crid,
				},
			},
			err: errCloudConfigInvalid,
		},
		"FailCloudResourceMissingResourceID": {
			input: &Config{
				HTTPCollectorEndpoint: &endpoint,
				Cloud: &Cloud{
					ClientSecret: &csec,
					ClientID:     &cid,
				},
			},
			err: errCloudConfigInvalid,
		},
		"FailCloudResourceMissingClientSecret": {
			input: &Config{
				HTTPCollectorEndpoint: &endpoint,
				Cloud: &Cloud{
					ResourceID: &crid,
					ClientID:   &cid,
				},
			},
			err: errCloudConfigInvalid,
		},
		"SuccessfulCloudNotSpecified": {
			input: &Config{
				HTTPCollectorEndpoint: &endpoint,
				Cloud:                 &Cloud{},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {

			err := tc.input.validate()
			if tc.err != nil {
				test.Error(t, err)
				test.ErrorIs(t, err, tc.err)
				return
			}
			test.NoError(t, err)

		})
	}
}
