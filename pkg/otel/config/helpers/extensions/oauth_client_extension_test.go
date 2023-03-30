package extensions

import (
	"fmt"
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

	//Unmarshall and verify
	unmarshalledCfg := &OauthClientConfig{}
	conf.Unmarshal(unmarshalledCfg)

	require.Equal(t, cfg, unmarshalledCfg)
}

// The purpose of this test is so that if the package ever supports marshalling with an opaque
// client secret we might be able to use it. Since our configuration is marshalled we can't
// use it as it gets exported unfortunately.
func Test_OauthClientPkg(t *testing.T) {
	cfg := oauth2clientauthextension.Config{
		ClientSecret: configopaque.String("foo"),
	}

	fmt.Println(cfg.ClientSecret)

	conf := confmap.New()
	err := conf.Marshal(cfg)
	require.NoError(t, err)

	//Unmarshall and verify
	unmarshalledCfg := &oauth2clientauthextension.Config{}
	conf.Unmarshal(unmarshalledCfg)
	fmt.Println(unmarshalledCfg.ClientSecret)

	require.NotEqual(t, cfg, unmarshalledCfg)

}
