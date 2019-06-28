package collector

import (
	"strconv"
	"testing"
)

func Test_Collector_Env_clusterIDFromPath(t *testing.T) {
	testCases := []struct {
		name       string
		path       string
		expectedID string
	}{
		{
			name:       "case 0",
			path:       "br-ux9ty.env",
			expectedID: "ux9ty",
		},
		{
			name:       "case 1",
			path:       "br-sdfskdux9ty.env",
			expectedID: "sdfskdux9ty",
		},
		{
			name:       "case 2",
			path:       "bridge-ux9ty.env",
			expectedID: "ux9ty",
		},
		{
			name:       "case 3",
			path:       "bridge-sdfskdux9ty.env",
			expectedID: "sdfskdux9ty",
		},
		{
			name:       "case 4",
			path:       "ux9ty.env",
			expectedID: "ux9ty",
		},
		{
			name:       "case 5",
			path:       "sdfskdux9ty.env",
			expectedID: "sdfskdux9ty",
		},
		{
			name:       "case 6",
			path:       "sdfskdux9ty",
			expectedID: "",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			id := clusterIDFromPath(tc.path)

			if id != tc.expectedID {
				t.Fatalf("expected %s to equal %s", tc.expectedID, id)
			}
		})
	}
}
