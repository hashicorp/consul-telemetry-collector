package inmem

import (
	"context"
	"testing"

	"github.com/shoenig/test"
)

func Test_InMem(t *testing.T) {
	provider := NewProvider()
	_, err := provider.Retrieve(context.Background(), "", nil)
	test.NoError(t, err)
}
