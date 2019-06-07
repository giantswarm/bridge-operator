package collector

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

const (
	subsystemEnv = "env"
)

var (
	envClusterWithoutFileDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemEnv, "cluster_without_file"),
		"Clusters without environment files.",
		[]string{
			labelCluster,
		},
		nil,
	)
	envFileWithoutClusterDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemEnv, "file_without_cluster"),
		"Environment files without associated cluster.",
		[]string{
			labelCluster,
		},
		nil,
	)
)

type EnvConfig struct {
	G8sClient versioned.Interface
	Logger    micrologger.Logger

	Directory string
}

type Env struct {
	g8sClient versioned.Interface
	logger    micrologger.Logger

	directory string
}

func NewEnv(config EnvConfig) (*Env, error) {
	if config.G8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.G8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.Directory == "" {
		config.Directory = "/run/flannel/networks/"
	}

	c := &Env{
		g8sClient: config.G8sClient,
		logger:    config.Logger,

		directory: config.Directory,
	}

	return c, nil
}

func (c *Env) Collect(ch chan<- prometheus.Metric) error {
	desiredClusterIDs := map[string]struct{}{}
	{
		l, err := c.g8sClient.Core().FlannelConfigs(metav1.NamespaceAll).List(metav1.ListOptions{})
		if err != nil {
			return microerror.Mask(err)
		}

		for _, i := range l.Items {
			desiredClusterIDs[i.Name] = struct{}{}
		}
	}

	envClusterIDs := map[string]struct{}{}
	{
		f := func(path string, info os.FileInfo, err error) error {
			envClusterIDs[clusterIDFromPath(path)] = struct{}{}
			return nil
		}

		err := filepath.Walk(c.directory, f)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	{
		l, r := symmetricDifference(mapToSlice(desiredClusterIDs), mapToSlice(envClusterIDs))

		// Emit metrics for clusters for which we couldn't find any environment
		// file.
		for _, id := range l {
			//
		}

		// Emit metrics for environment files for which we couldn't find any
		// cluster.
		for _, id := range r {
			//
		}
	}

	return nil
}

func (c *Env) Describe(ch chan<- *prometheus.Desc) error {
	ch <- envClusterWithoutFileDesc
	ch <- envFileWithoutClusterDesc
	return nil
}

// clusterIDFromPath received the path variabled from the executed walk
// function. The path variable is intended to be the file name of the
// environment file managed by Flannel.
//
//     br-ux9ty.env
//
func clusterIDFromPath(path string) string {
	r := regexp.MustCompile(`([a-z0-9]+).env$`)
	l := r.FindStringSubmatch(path)

	if len(l) == 2 {
		return l[1]
	}

	return ""
}

func containsString(l []string, s string) bool {
	for _, i := range l {
		if i == s {
			return true
		}
	}

	return false
}

func mapToSlice(m map[string]struct{}) []string {
	var l []string

	for k, _ := range m {
		l = append(l, k)
	}

	return l
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
