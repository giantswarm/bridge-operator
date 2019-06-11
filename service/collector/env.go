package collector

import (
	"io/ioutil"
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

const (
	envDirectory = "/run/flannel/networks/"
)

var (
	envClusterWithoutFileDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemEnv, "cluster_without_env_file"),
		"Clusters without environment files.",
		[]string{
			labelCluster,
		},
		nil,
	)
	envFileWithoutClusterDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemEnv, "env_file_without_cluster"),
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
}

type Env struct {
	g8sClient versioned.Interface
	logger    micrologger.Logger
}

func NewEnv(config EnvConfig) (*Env, error) {
	if config.G8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.G8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	c := &Env{
		g8sClient: config.G8sClient,
		logger:    config.Logger,
	}

	return c, nil
}

func (c *Env) Collect(ch chan<- prometheus.Metric) error {
	var desiredClusterIDs []string
	{
		l, err := c.g8sClient.Core().FlannelConfigs(metav1.NamespaceAll).List(metav1.ListOptions{})
		if err != nil {
			return microerror.Mask(err)
		}

		for _, i := range l.Items {
			desiredClusterIDs = append(desiredClusterIDs, i.Name)
		}
	}

	var envClusterIDs []string
	{
		files, err := ioutil.ReadDir(envDirectory)
		if err != nil {
			return microerror.Mask(err)
		}

		for _, file := range files {
			id := clusterIDFromPath(file.Name())
			if id == "" {
				return microerror.Maskf(executionFailedError, "file %q does not encode a cluster ID", file.Name())
			}

			envClusterIDs = append(envClusterIDs, id)
		}
	}

	{
		l, r := symmetricDifference(desiredClusterIDs, envClusterIDs)

		// Emit metrics for clusters for which we couldn't find any environment
		// file.
		for _, id := range l {
			ch <- prometheus.MustNewConstMetric(
				envClusterWithoutFileDesc,
				prometheus.GaugeValue,
				gaugeValue,
				id,
			)
		}

		// Emit metrics for environment files for which we couldn't find any
		// cluster.
		for _, id := range r {
			ch <- prometheus.MustNewConstMetric(
				envFileWithoutClusterDesc,
				prometheus.GaugeValue,
				gaugeValue,
				id,
			)
		}
	}

	return nil
}

func (c *Env) Describe(ch chan<- *prometheus.Desc) error {
	ch <- envClusterWithoutFileDesc
	ch <- envFileWithoutClusterDesc
	return nil
}

// clusterIDFromPath receives the file name from the result of ioutil.ReadDir.
// The path variable is intended to be the file name of the environment file
// managed by Flannel.
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
