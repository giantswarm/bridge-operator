package collector

import (
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/exporterkit/collector"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

type SetConfig struct {
	G8sClient versioned.Interface
	Logger    micrologger.Logger
}

// Set is basically only a wrapper for the operator's collector implementations.
// It eases the iniitialization and prevents some weird import mess so we do not
// have to alias packages.
type Set struct {
	*collector.Set
}

func NewSet(config SetConfig) (*Set, error) {
	var err error

	var bridgeCollector *Bridge
	{
		c := BridgeConfig{
			G8sClient: config.G8sClient,
			Logger:    config.Logger,
		}

		bridgeCollector, err = NewBridge(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var envCollector *Env
	{
		c := EnvConfig{
			G8sClient: config.G8sClient,
			Logger:    config.Logger,
		}

		envCollector, err = NewEnv(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var collectorSet *collector.Set
	{
		c := collector.SetConfig{
			Collectors: []collector.Interface{
				bridgeCollector,
				envCollector,
			},
			Logger: config.Logger,
		}

		collectorSet, err = collector.NewSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	s := &Set{
		Set: collectorSet,
	}

	return s, nil
}
