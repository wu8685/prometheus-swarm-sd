package swarm

import (
	neturl "net/url"
	"strings"
	"testing"

	"github.com/prometheus/prometheus/config"
)

var (
	masters = []string{"http://example.com:2376"}
)

func TestGetNodeInfo(t *testing.T) {
	client, err := createSwarmClient()
	if err != nil {
		t.Fatalf("Fail to create swarm client: %s.\n", err.Error())
	}

	info, err := client.getNodeInfo()
	if err != nil {
		t.Fatalf("Fail to get node info: %s.\n", err.Error())
	}
	if len(info.Nodes) <= 0 {
		t.Fatalf("Fail to get node info: node number is %d. It is supposed to be greater than 0.\n", len(info.Nodes))
	}

	node := info.Nodes[0]
	if strings.HasPrefix(node.Host, " ") || strings.HasSuffix(node.Host, " ") || node.Host == "" {
		t.Fatalf("Fail to get node info: wrong host %s\n", node.Host)
	}
	if strings.HasPrefix(node.Addr, " ") || strings.HasSuffix(node.Addr, " ") || node.Addr == "" {
		t.Fatalf("Fail to get node info: wrong addr %s\n", node.Addr)
	}
}

func createSwarmClient() (*swarmClient, error) {
	cfg := &config.DefaultSwarmSDConfig

	masterUrls := []config.URL{}
	for _, m := range masters {
		url, err := neturl.Parse(m)
		if err != nil {
			return nil, err
		}
		masterUrls = append(masterUrls, config.URL{url})
	}
	cfg.Masters = masterUrls
	client := newSwarmClient(cfg)
	client.do = mockOneNodeDo

	return client, nil
}
