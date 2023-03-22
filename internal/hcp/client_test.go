package hcp

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shoenig/test/must"

	"github.com/hashicorp/hcp-sdk-go/resource"
)

func Test_New(t *testing.T) {
	testcases := map[string]struct {
		cid     string
		csec    string
		res     resource.Resource
		wantErr bool
	}{
		"Good": {
			cid:  uuid.NewString(),
			csec: uuid.NewString(),
			res: resource.Resource{
				ID:           uuid.NewString(),
				Type:         "type",
				Organization: uuid.NewString(),
				Project:      uuid.NewString(),
			},
		},
		// "NoClientID": {
		// 	cid:  "",
		// 	csec: uuid.NewString(),
		// 	res: resource.Resource{
		// 		ID:           uuid.NewString(),
		// 		Type:         "type",
		// 		Organization: uuid.NewString(),
		// 		Project:      uuid.NewString(),
		// 	},
		// },
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			_, err := New(tc.cid, tc.csec, tc.res.String())
			if tc.wantErr {
				must.Error(t, err)
				return
			}

			must.NoError(t, err)
		})
	}

}
