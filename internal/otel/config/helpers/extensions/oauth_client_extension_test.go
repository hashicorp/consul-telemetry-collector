package extensions

import (
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/oauth2clientauthextension"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/confmap"
)

func Test_OauthClient(t *testing.T) {
	cfg := OauthClientCfg("cid", "csec")
	require.NotNil(t, cfg)

	// Marshall the configuration
	conf := confmap.New()
	err := conf.Marshal(cfg)
	require.NoError(t, err)

	// Unmarshall and verify
	unmarshalledCfg := &OauthClientConfig{}
	err = conf.Unmarshal(unmarshalledCfg)
	require.NoError(t, err)

	require.Equal(t, cfg, unmarshalledCfg)
}

// The purpose of this test is so that if the package ever supports marshaling with an opaque
// client secret we might be able to use it. Since our configuration is marshaled we can't
// use it as it gets exported unfortunately.
func Test_OauthClientPkg(t *testing.T) {
	cfg := oauth2clientauthextension.Config{
		ClientSecret: configopaque.String("foo"),
	}

	conf := confmap.New()
	err := conf.Marshal(cfg)
	require.NoError(t, err)

	// Unmarshall and verify
	unmarshalledCfg := &oauth2clientauthextension.Config{}
	err = conf.Unmarshal(unmarshalledCfg)
	require.NoError(t, err)

	require.NotEqual(t, cfg, unmarshalledCfg)
}
