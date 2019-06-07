package collector

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

const (
	subsystemEnv = "env"
)

var (
	envDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemEnv, "inconsistencies"),
		"Information of environment inconsistencies.",
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

	v := &Env{
		g8sClient: config.G8sClient,
		logger:    config.Logger,
	}

	return v, nil
}

func (v *Env) Collect(ch chan<- prometheus.Metric) error {
	return nil
}

func (v *Env) Describe(ch chan<- *prometheus.Desc) error {
	ch <- envDesc
	return nil
}
