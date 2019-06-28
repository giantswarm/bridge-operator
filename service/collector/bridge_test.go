package collector

import (
	"strconv"
	"testing"
)

func Test_Collector_Bridge_clusterIDFromName(t *testing.T) {
	testCases := []struct {
		name       string
		path       string
		expectedID string
	}{
		{
			name:       "case 0",
			path:       "br-ux9ty",
			expectedID: "ux9ty",
		},
		{
			name:       "case 1",
			path:       "br-sdfskdux9ty",
			expectedID: "sdfskdux9ty",
		},
		{
			name:       "case 2",
			path:       "bridge-ux9ty",
			expectedID: "",
		},
		{
			name:       "case 3",
			path:       "bridge-sdfskdux9ty",
			expectedID: "",
		},
		{
			name:       "case 4",
			path:       "ux9ty",
			expectedID: "",
		},
		{
			name:       "case 5",
			path:       "sdfskdux9ty",
			expectedID: "",
		},
		{
			name:       "case 6",
			path:       "br-sdfskdux9ty.env",
			expectedID: "",
		},
		{
			name:       "case 7",
			path:       "lo",
			expectedID: "",
		},
		{
			name:       "case 8",
			path:       "cali5444ea819af",
			expectedID: "",
		},
		{
			name:       "case 9",
			path:       "veth1a5fde1",
			expectedID: "",
		},
		{
			name:       "case 10",
			path:       "eth0",
			expectedID: "",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			id := clusterIDFromName(tc.path)

			if id != tc.expectedID {
				t.Fatalf("expected %s to equal %s", tc.expectedID, id)
			}
		})
	}
}
