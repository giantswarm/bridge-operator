package collector

import (
	"context"
	"fmt"
	"net"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

const (
	subsystemBridge = "bridge"
)

var (
	bridgeClusterWithoutNetworkInterfaceDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemBridge, "cluster_without_network_interface"),
		"Clusters without network interface.",
		[]string{
			labelCluster,
		},
		nil,
	)
	bridgeNetworkInterfaceWithoutClusterDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemBridge, "network_interface_without_cluster"),
		"Network interfaces without cluster.",
		[]string{
			labelCluster,
		},
		nil,
	)
)

type BridgeConfig struct {
	G8sClient versioned.Interface
	Logger    micrologger.Logger
}

type Bridge struct {
	g8sClient versioned.Interface
	logger    micrologger.Logger
}

func NewBridge(config BridgeConfig) (*Bridge, error) {
	if config.G8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.G8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	c := &Bridge{
		g8sClient: config.G8sClient,
		logger:    config.Logger,
	}

	return c, nil
}

func (c *Bridge) Collect(ch chan<- prometheus.Metric) error {
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

	var bridgeClusterIDs []string
	{
		c.logger.LogCtx(ctx, "level", "debug", "message", "finding network interfaces")

		interfaces, err := net.Interfaces()
		if err != nil {
			return microerror.Mask(err)
		}

		for _, i := range interfaces {
			id := clusterIDFromName(i.Name)
			if id == "" {
				// There are many different network interfaces and we cannot parse all
				// of them. Thus we continue and go ahead with the next one we found.
				continue
			}

			bridgeClusterIDs = append(bridgeClusterIDs, id)
		}

		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("found %d network interfaces for cluster IDs: %#v", len(bridgeClusterIDs), bridgeClusterIDs))
	}

	{
		l, r := symmetricDifference(desiredClusterIDs, bridgeClusterIDs)

		if len(l) == 0 {
			c.logger.LogCtx(ctx, "level", "debug", "message", "no orphaned cluster IDs found")
		}
		if len(r) == 0 {
			c.logger.LogCtx(ctx, "level", "debug", "message", "no orphaned network interfaces found")
		}

		// Emit metrics for clusters for which we couldn't find any environment
		// file.
		for _, id := range l {
			ch <- prometheus.MustNewConstMetric(
				bridgeClusterWithoutNetworkInterfaceDesc,
				prometheus.GaugeValue,
				gaugeValue,
				id,
			)
		}

		// Emit metrics for environment files for which we couldn't find any
		// cluster.
		for _, id := range r {
			ch <- prometheus.MustNewConstMetric(
				bridgeNetworkInterfaceWithoutClusterDesc,
				prometheus.GaugeValue,
				gaugeValue,
				id,
			)
		}
	}

	return nil
}

func (c *Bridge) Describe(ch chan<- *prometheus.Desc) error {
	ch <- bridgeClusterWithoutNetworkInterfaceDesc
	ch <- bridgeNetworkInterfaceWithoutClusterDesc
	return nil
}

// clusterIDFromName receives the interface name from the result of
// net.Interfaces. The name variable is intended to be the interface name of the
// KVM node's network bridge.
//
//     br-6m5o8
//
func clusterIDFromName(name string) string {
	r := regexp.MustCompile(`br-([a-z0-9]+)$`)
	l := r.FindStringSubmatch(name)

	if len(l) == 2 {
		return l[1]
	}

	return ""
}
