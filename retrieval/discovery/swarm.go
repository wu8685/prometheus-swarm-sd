package discovery

import (
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/retrieval/discovery/swarm"
)

// NewSwarmDiscovery creates a Swarm service discovery based on the passed-in configuration.
func NewSwarmDiscovery(conf *config.SwarmSDConfig) (*swarm.Discovery, error) {
	sd := &swarm.Discovery{
		Conf: conf,
	}
	err := sd.Initialize()
	if err != nil {
		return nil, err
	}
	return sd, nil
}
