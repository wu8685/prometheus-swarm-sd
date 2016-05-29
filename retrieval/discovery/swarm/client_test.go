package swarm

import (
	"net/http"
	"strings"
	"testing"
)

var (
	masters = []string{"http://example.com:2376"}
)

func TestGetNodeInfo(t *testing.T) {
	client := &swarmClient{}
	mockResp := createMockResponse()

	info, err := client.processNodeInfo(mockResp)
	if err != nil {
		t.Fatalf("Fail to get node info: %s.\n", err.Error())
	}
	if len(info.SystemStatus.Nodes) <= 0 {
		t.Fatalf("Fail to get node info: node number is %d. It is supposed to be greater than 0.\n", len(info.SystemStatus.Nodes))
	}

	node := info.SystemStatus.Nodes[0]
	if strings.HasPrefix(node.Host, " ") || strings.HasSuffix(node.Host, " ") || node.Host == "" {
		t.Fatalf("Fail to get node info: wrong host %s\n", node.Host)
	}
	if strings.HasPrefix(node.Addr, " ") || strings.HasSuffix(node.Addr, " ") || node.Addr == "" {
		t.Fatalf("Fail to get node info: wrong addr %s\n", node.Addr)
	}
	if len(info.SystemStatus.Nodes) <= 0 {
		t.Fatalf("Fail to parse node info: fail to get host: its num is %d", len(info.SystemStatus.Nodes))
	}
}

func createMockResponse() (resp *http.Response) {
	resp = &http.Response{}
	resp.StatusCode = 200
	resp.Body = &mockBody{strings.NewReader(twoNode)}
	return
}

type mockBody struct {
	*strings.Reader
}

func (b *mockBody) Close() error {
	return nil
}
