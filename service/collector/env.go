package collector

import (
	"context"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	subsystemEnv = "env"
)

const (
	envDirectory = "/run/flannel/networks/"
)

var (
	envClusterWithoutFlannelNetworkEnvFileDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemEnv, "cluster_without_flannel_network_env_file"),
		"Clusters without environment files.",
		[]string{
			labelCluster,
		},
		nil,
	)
	envFlannelNetworkEnvFileWithoutClusterDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemEnv, "flannel_network_env_file_without_cluster"),
		"Environment files without cluster.",
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
	ctx := context.Background()

	var desiredClusterIDs []string
	{
		c.logger.LogCtx(ctx, "level", "debug", "message", "finding all cluster IDs in FlannelConfigs")

		l, err := c.g8sClient.Core().FlannelConfigs(metav1.NamespaceAll).List(metav1.ListOptions{})
		if err != nil {
			return microerror.Mask(err)
		}

		for _, i := range l.Items {
			desiredClusterIDs = append(desiredClusterIDs, i.Name)
		}

		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("found %d cluster IDs: %#v", len(desiredClusterIDs), desiredClusterIDs))
	}

	var envClusterIDs []string
	{
		c.logger.LogCtx(ctx, "level", "debug", "message", "finding flannel environment files")

		files, err := ioutil.ReadDir(envDirectory)
		if err != nil {
			return microerror.Mask(err)
		}

		for _, file := range files {
			id := clusterIDFromPath(file.Name())
			if id == "" {
				return microerror.Maskf(executionFailedError, "file %#q does not encode a cluster ID", file.Name())
			}

			envClusterIDs = append(envClusterIDs, id)
		}

		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("found %d flannel environment files for cluster IDs: %#v", len(envClusterIDs), envClusterIDs))
	}

	{
		l, r := symmetricDifference(desiredClusterIDs, envClusterIDs)

		if len(l) == 0 {
			c.logger.LogCtx(ctx, "level", "debug", "message", "no orphaned cluster IDs found")
		}
		if len(r) == 0 {
			c.logger.LogCtx(ctx, "level", "debug", "message", "no orphaned flannel environment files found")
		}

		// Emit metrics for clusters for which we couldn't find any environment
		// file.
		for _, id := range l {
			ch <- prometheus.MustNewConstMetric(
				envClusterWithoutFlannelNetworkEnvFileDesc,
				prometheus.GaugeValue,
				gaugeValue,
				id,
			)
		}

		// Emit metrics for environment files for which we couldn't find any
		// cluster.
		for _, id := range r {
			ch <- prometheus.MustNewConstMetric(
				envFlannelNetworkEnvFileWithoutClusterDesc,
				prometheus.GaugeValue,
				gaugeValue,
				id,
			)
		}
	}

	return nil
}

func (c *Env) Describe(ch chan<- *prometheus.Desc) error {
	ch <- envClusterWithoutFlannelNetworkEnvFileDesc
	ch <- envFlannelNetworkEnvFileWithoutClusterDesc
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
