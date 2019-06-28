package collector

import (
	"reflect"
	"strconv"
	"testing"
)

func Test_Collector_symmetricDifference(t *testing.T) {
	testCases := []struct {
		name string
		a    []string
		b    []string
		l    []string
		r    []string
	}{
		{
			name: "case 0",
			a:    []string{"a", "b", "c", "d"},
			b:    []string{"c", "d", "e", "f"},
			l:    []string{"a", "b"},
			r:    []string{"e", "f"},
		},
		{
			name: "case 1",
			a:    []string{"a", "b", "c", "d"},
			b:    []string{"b", "c", "d", "e"},
			l:    []string{"a"},
			r:    []string{"e"},
		},
		{
			name: "case 2",
			a:    []string{"a", "b", "c", "d"},
			b:    []string{"a", "b", "c", "d", "e"},
			l:    nil,
			r:    []string{"e"},
		},
		{
			name: "case 3",
			a:    []string{"a", "b", "c", "d"},
			b:    []string{"b", "c", "d"},
			l:    []string{"a"},
			r:    nil,
		},
		{
			name: "case 4",
			a:    []string{"b", "c", "d"},
			b:    []string{"b", "c", "d"},
			l:    nil,
			r:    nil,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			l, r := symmetricDifference(tc.a, tc.b)

			if !reflect.DeepEqual(l, tc.l) {
				t.Fatalf("expected %#v to equal %#v", tc.l, l)
			}
			if !reflect.DeepEqual(r, tc.r) {
				t.Fatalf("expected %#v to equal %#v", tc.r, r)
			}
		})
	}
}
