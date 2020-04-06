package collector

const (
	gaugeValue float64 = 1
	namespace  string  = "bridge_operator"
)

const (
	labelCluster = "cluster_id"
)

func containsString(l []string, s string) bool {
	for _, i := range l {
		if i == s {
			return true
		}
	}

	return false
}

// symmetricDifference implements the selection of a relative complement of two
// lists. See also https://en.wikipedia.org/wiki/Set_(mathematics)#Complements.
// Given input arguments a and b, return value l contains only values that are
// exclusively in a and r contains only values that are exclusively in b.
//
//     a = [1, 2, 3, 4]
//     b = [3, 4, 5, 6]
//     l = [1, 2]
//     r = [5, 6]
//
func symmetricDifference(a, b []string) (l []string, r []string) {
	for _, s := range a {
		if !containsString(b, s) {
			l = append(l, s)
		}
	}

	for _, s := range b {
		if !containsString(a, s) {
			r = append(r, s)
		}
	}

	return l, r
}
