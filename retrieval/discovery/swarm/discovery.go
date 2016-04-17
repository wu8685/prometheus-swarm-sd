package swarm

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/prometheus/common/log"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/util/strutil"
)

const (

	// swarmMetaLabelPrefix is the meta prefix used for all meta labels.
	// in this discovery.
	metaLabelPrefix = model.MetaLabelPrefix + "swarm_"
	// nodeLabelPrefix is the prefix for the node labels.
	nodeLabelPrefix = metaLabelPrefix + "node_label_"
	// nodesTargetGroupName is the name given to the target group for nodes.
	nodesTargetGroupName = "swarm_nodes"
	// mastersTargetGroupName is the name given to the target group for swarm masters.
	mastersTargetGroupName = "swarm_masters"
	// roleLabel is the name for the label containing a target's role.
	roleLabel = metaLabelPrefix + "role"
)

type Discovery struct {
	Conf   *config.SwarmSDConfig
	client *swarmClient

	masters []config.URL

	nodes  map[string]*Node
	nodeMu sync.RWMutex
}

// Initialize sets up the discovery for usage.
func (d *Discovery) Initialize() error {
	d.client = newSwarmClient(d.Conf)

	d.masters = d.Conf.Masters
	d.nodes = map[string]*Node{}
	return nil
}

// Sources returns the source identifiers the provider is currently aware of.
func (d *Discovery) Sources() []string {
	sources := []string{}

	for _, masterUrl := range d.masters {
		sources = append(sources, fmt.Sprintf("%s:%s", mastersTargetGroupName, masterUrl.Host))
	}

	info, err := d.client.getNodeInfo()
	if err != nil {
		log.Errorf("Fail to get swarm node info: %s.", err.Error())
		return []string{}
	}

	for _, node := range info.Nodes {
		sources = append(sources, fmt.Sprintf("%s:%s", nodesTargetGroupName, node.Addr))
	}

	return sources
}

// Run hands a channel to the target provider through which it can send
// updated target groups. The channel must be closed by the target provider
// if no more updates will be sent.
// On receiving from done Run must return.
func (d *Discovery) Run(up chan<- config.TargetGroup, done <-chan struct{}) {
	defer close(up)

	if tg := d.masterTargetGroup(); tg != nil {
		select {
		case <-done:
			return
		case up <- *tg:
		}
	}

	retryInterval := time.Duration(d.Conf.RefreshInterval)
	update := make(chan []*Node, 10)

	go d.watchNodes(update, done, retryInterval)

	for {
		select {
		case <-done:
			return
		case nodes := <-update:
			d.updateNodes(nodes)
			tg := d.nodeTargetGroup()
			up <- *tg
		}
	}
}

func (d *Discovery) masterTargetGroup() (tg *config.TargetGroup) {
	tg = &config.TargetGroup{
		Source: mastersTargetGroupName,
		Labels: model.LabelSet{
			roleLabel: model.LabelValue("master"),
		},
	}

	for _, master := range d.Conf.Masters {
		masterAddress := master.Host
		cadvisorAddr := d.useCAdvisorPort(masterAddress)
		t := model.LabelSet{
			model.InstanceLabel: model.LabelValue(masterAddress),
			model.AddressLabel:  model.LabelValue(cadvisorAddr),
			model.SchemeLabel:   model.LabelValue(master.Scheme),
		}
		tg.Targets = append(tg.Targets, t)

	}
	return tg
}

func (d *Discovery) watchNodes(update chan []*Node, done <-chan struct{}, retryInterval time.Duration) {
	info, err := d.client.getNodeInfo()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	update <- info.Nodes

	until(func() {
		info, err = d.client.getNodeInfo()
		if err != nil {
			log.Errorf(err.Error())
			return
		}
		update <- info.Nodes
	}, retryInterval, done)
}

func (d *Discovery) updateNodes(nodes []*Node) {
	d.nodeMu.Lock()
	defer d.nodeMu.Unlock()

	// remove deleted nodes
	deletedNodeNames := []string{}
	for name, _ := range d.nodes {
		exist := false
		for _, node := range nodes {
			if name == node.Host {
				exist = true
				break
			}
		}
		if !exist {
			deletedNodeNames = append(deletedNodeNames, name)
		}
	}
	for _, deleted := range deletedNodeNames {
		delete(d.nodes, deleted)
	}

	// update exist/created nodes
	for _, node := range nodes {
		d.nodes[node.Host] = node
	}
}

func (d *Discovery) nodeTargetGroup() (tg *config.TargetGroup) {
	d.nodeMu.RLock()
	defer d.nodeMu.RUnlock()

	tg = &config.TargetGroup{
		Source: nodesTargetGroupName,
		Labels: model.LabelSet{
			roleLabel: model.LabelValue("node"),
		},
		Targets: []model.LabelSet{},
	}

	for _, node := range d.nodes {
		cadvisorAddr := d.useCAdvisorPort(node.Addr)
		t := model.LabelSet{
			model.AddressLabel:  model.LabelValue(cadvisorAddr),
			model.InstanceLabel: model.LabelValue(node.Addr),
		}

		for k, v := range node.Labels {
			labelName := strutil.SanitizeLabelName(nodeLabelPrefix + k)
			t[model.LabelName(labelName)] = model.LabelValue(v)
		}

		tg.Targets = append(tg.Targets, t)
	}
	return
}

func (d *Discovery) useCAdvisorPort(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	// If error then no port is specified
	if err != nil {
		return fmt.Sprintf("%s:%s", addr, d.Conf.MetricsPort)
	}

	return fmt.Sprintf("%s:%s", host, d.Conf.MetricsPort)
}

// Until loops until stop channel is closed, running f every period.
// f may not be invoked if stop channel is already closed.
func until(f func(), period time.Duration, stopCh <-chan struct{}) {
	select {
	case <-stopCh:
		return
	default:
		f()
	}
	for {
		select {
		case <-stopCh:
			return
		case <-time.After(period):
			f()
		}
	}
}
