package agent

import (
	"context"
	"testing"
	"time"

	"github.com/shoenig/test/must"
)

func Test_runSvc(t *testing.T) {
	testcases := map[string]struct {
		cfg Config
	}{
		"MinimalConfig": {
			cfg: Config{
				HTTPCollectorEndpoint: "https://otel-http-endpoint",
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			s, err := NewService(&tc.cfg)
			must.NoError(t, err)
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			errCh := make(chan error)
			go func() {
				err := s.Run(ctx)
				errCh <- err
			}()

			time.After(time.Second)
			cancel()
			<-ctx.Done()
			err = <-errCh
			must.NoError(t, err)
		})
	}

}
