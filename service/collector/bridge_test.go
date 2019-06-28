package collector

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Collector_Bridge_clusterIDFromName(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		expectedID   string
		errorMatcher func(err error) bool
	}{
		{
			name:         "case 0",
			path:         "br-ux9ty",
			expectedID:   "ux9ty",
			errorMatcher: nil,
		},
		{
			name:         "case 1",
			path:         "br-sdfskdux9ty",
			expectedID:   "sdfskdux9ty",
			errorMatcher: nil,
		},
		{
			name:         "case 2",
			path:         "bridge-ux9ty",
			expectedID:   "",
			errorMatcher: IsExecutionFailed,
		},
		{
			name:         "case 3",
			path:         "bridge-sdfskdux9ty",
			expectedID:   "",
			errorMatcher: IsExecutionFailed,
		},
		{
			name:         "case 4",
			path:         "ux9ty",
			expectedID:   "",
			errorMatcher: IsExecutionFailed,
		},
		{
			name:         "case 5",
			path:         "sdfskdux9ty",
			expectedID:   "",
			errorMatcher: IsExecutionFailed,
		},
		{
			name:         "case 6",
			path:         "br-sdfskdux9ty.env",
			expectedID:   "",
			errorMatcher: IsExecutionFailed,
		},
		{
			name:         "case 7",
			path:         "lo",
			expectedID:   "",
			errorMatcher: IsExecutionFailed,
		},
		{
			name:         "case 8",
			path:         "cali5444ea819af",
			expectedID:   "",
			errorMatcher: IsExecutionFailed,
		},
		{
			name:         "case 9",
			path:         "veth1a5fde1",
			expectedID:   "",
			errorMatcher: IsExecutionFailed,
		},
		{
			name:         "case 10",
			path:         "eth0",
			expectedID:   "",
			errorMatcher: IsExecutionFailed,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			id, err := clusterIDFromName(tc.path)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

			if id != tc.expectedID {
				t.Fatalf("\n\n%s\n", cmp.Diff(id, tc.expectedID))
			}
		})
	}
}
